package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions"
	changesdescriber "github.com/WinPooh32/reviewpls/cmd/reviewpls/functions/changes-describer"
	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions/reviewer"
	"github.com/WinPooh32/reviewpls/internal/gitrepo"
	"github.com/WinPooh32/reviewpls/internal/slogutil"
	"github.com/openai/openai-go/v2"
)

const maxRetries = 10

const (
	errorCodeUnknown     = 1
	errorCodeHasComments = 2
)

const (
	defaultTemperature = 0.1
)

var errHasReviewComments = errors.New("has review comments")

func main() {
	exitCode := 0
	defer func(code *int) {
		os.Exit(*code)
	}(&exitCode)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		AddSource:   false,
		ReplaceAttr: nil,
	}))

	ctx = slogutil.WithContext(ctx, logger)

	deps, err := assembleDependencies(ctx)
	if err != nil {
		logger.Error("assemble dependencies", slog.String("error", err.Error()))

		exitCode = errorCodeUnknown

		return
	}

	if err := run(ctx, deps); err != nil {
		logger.Error("run", slog.String("error", err.Error()))

		if errors.Is(err, errHasReviewComments) {
			exitCode = errorCodeHasComments
		} else {
			exitCode = errorCodeUnknown
		}

		return
	}
}

type Dependencies struct {
	Mr MergeRequest
	Pf PipelineFunctions
}

func assembleDependencies(ctx context.Context) (*Dependencies, error) {
	gitBaseBranch := flag.String("git-branch-base", "master", "base git repositroy branch name")
	gitHeadBaranch := flag.String("git-branch-head", "feature-1234", "head branch name")
	gitRootDir := flag.String("git-root-dir", ".", "root path to a git repository directory")
	language := flag.String("language", "english", "commentary language")
	model := flag.String("model", "gpt-4", "model name")

	flag.Parse()

	repo, err := gitrepo.OpenRepository(ctx, *gitRootDir)
	if err != nil {
		return nil, fmt.Errorf("gitrepo: open repository: %w", err)
	}

	openaiClient := openai.NewClient()

	mr := MergeRequest{
		Repo:       repo,
		BaseBranch: *gitBaseBranch,
		HeadBranch: *gitHeadBaranch,
	}

	cd, err := changesdescriber.New(
		changesdescriber.ChangesDescriberConfig{
			Model:            *model,
			Temperature:      defaultTemperature,
			RetryMaxAttempts: maxRetries,
			RetryDelay:       0,
		},
		openaiClient,
	)

	rv, err := reviewer.New(
		reviewer.ChangesDescriberConfig{
			Model:              *model,
			RetryMaxAttempts:   maxRetries,
			RetryDelay:         0,
			CommentaryLanguage: strings.ToLower(*language),
		},
		openaiClient, *model,
	)

	pf := PipelineFunctions{
		ChangesDescriber: cd,
		Reviewer:         rv,
	}

	return &Dependencies{
		Mr: mr,
		Pf: pf,
	}, nil
}

type MergeRequest struct {
	Repo       *gitrepo.Repository
	BaseBranch string
	HeadBranch string
}

type PipelineFunctions struct {
	ChangesDescriber functions.ChangesDescriber
	Reviewer         functions.Reviewer
}

func run(ctx context.Context, dp *Dependencies) error {
	files, err := dp.Mr.Repo.DiffFiles(ctx, dp.Mr.BaseBranch, dp.Mr.HeadBranch)
	if err != nil {
		return fmt.Errorf("get repo merge request diff files; %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	diffCommits, err := dp.Mr.Repo.DiffCommits(ctx, dp.Mr.BaseBranch, dp.Mr.HeadBranch)
	if err != nil {
		return fmt.Errorf("get repo diff commits; %w", err)
	}

	summaries := make([]*functions.ChangesSummary, 0, len(files))

	for _, file := range files {
		if strings.Contains(file, "pb.go") || strings.Contains(file, "gen.go") {
			continue
		}

		patch, err := dp.Mr.Repo.BlamePatch(ctx, dp.Mr.HeadBranch, diffCommits, file)
		if err != nil {
			return fmt.Errorf("get repo blame patch for the file %q: %w", file, err)
		}

		if len(patch.Commits) == 0 {
			continue
		}

		patchStr := patch.String()

		slogutil.Ctx(ctx).Debug("blame patch", slog.String("file", file), slog.String("patch", patchStr))

		summ, err := dp.Pf.ChangesDescriber.DescribeFileChanges(ctx, file, patchStr)
		if err != nil {
			return fmt.Errorf("describe file changes of %q: %w", file, err)
		}

		summaries = append(summaries, summ)
	}

	slogutil.Ctx(ctx).Info("describe changes", slog.Any("summary", summaries))

	var comments []*functions.ReviewComment

	for _, summ := range summaries {
		if len(summ.Hypotheses) == 0 {
			continue
		}

		comment, err := dp.Pf.Reviewer.AnalyzeChangesSummary(ctx, summ)
		if err != nil {
			return fmt.Errorf("reviewer: analyze summaries: %w", err)
		}

		if len(comment.Text) == 0 {
			slogutil.Ctx(ctx).Warn("empty commentary text, then skip it", slog.String("file", summ.File))
			continue
		}
	}

	if len(comments) == 0 {
		return nil
	}

	bs, err := json.Marshal(&comments)
	if err != nil {
		return fmt.Errorf("marshal json of comments: %w", err)
	}

	if _, err := os.Stdout.Write(bs); err != nil {
		return fmt.Errorf("print comments: %w", err)
	}

	return errHasReviewComments
}
