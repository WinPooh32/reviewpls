package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions"
	changesdescriber "github.com/WinPooh32/reviewpls/cmd/reviewpls/functions/changes-describer"
	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions/reviewer"
	"github.com/WinPooh32/reviewpls/internal/gitrepo"
	"github.com/openai/openai-go/v2"
)

const maxRetries = 10

const (
	errorCodeUnknown     = 1
	errorCodeHasComments = 2
)

var errHasReviewComments = errors.New("has review comments")

func main() {
	exitCode := 0
	defer func(code *int) {
		os.Exit(*code)
	}(&exitCode)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	gitBaseBranch := flag.String("git-branch-base", "master", "base git repositroy branch name")
	gitHeadBaranch := flag.String("git-branch-head", "feature-1234", "head branch name")
	gitRootDir := flag.String("git-root-dir", ".", "root path to a git repository directory")
	model := flag.String("model", "gpt-4", "model name")
	flag.Parse()

	repo, err := gitrepo.OpenRepository(ctx, *gitRootDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())

		exitCode = errorCodeUnknown

		return
	}

	openaiClient := openai.NewClient()

	mr := MergeRequest{
		Repo:       repo,
		BaseBranch: *gitBaseBranch,
		HeadBranch: *gitHeadBaranch,
	}

	cd, err := changesdescriber.NewChangesDescriberOpenAI(
		changesdescriber.ChangesDescriberOpenAIConfig{
			Model:            *model,
			RetryMaxAttempts: maxRetries,
			RetryDelay:       0,
		},
		openaiClient,
	)

	pf := PipelineFunctions{
		ChangesDescriber: cd,
		Reviewer:         reviewer.NewReviewerOpenAI(openaiClient, *model),
	}

	if err := run(ctx, mr, pf); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())

		if errors.Is(err, errHasReviewComments) {
			exitCode = errorCodeHasComments
		} else {
			exitCode = errorCodeUnknown
		}

		return
	}
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

func run(ctx context.Context, mr MergeRequest, pf PipelineFunctions) error {
	files, err := mr.Repo.DiffFiles(ctx, mr.BaseBranch, mr.HeadBranch)
	if err != nil {
		return fmt.Errorf("get repo merge request diff files; %w", err)
	}

	if len(files) == 0 {
		return nil
	}

	diffCommits, err := mr.Repo.DiffCommits(ctx, mr.BaseBranch, mr.HeadBranch)
	if err != nil {
		return fmt.Errorf("get repo diff commits; %w", err)
	}

	summaries := make([]*functions.ChangesSummary, 0, len(files))

	for _, file := range files {
		if strings.Contains(file, "pb.go") || strings.Contains(file, "gen.go") {
			continue
		}

		patch, err := mr.Repo.BlamePatch(ctx, mr.HeadBranch, diffCommits, file)
		if err != nil {
			return fmt.Errorf("get repo blame patch for the file %q: %w", file, err)
		}

		if len(patch.Commits) == 0 {
			continue
		}

		summ, err := pf.ChangesDescriber.DescribeFileChanges(ctx, file, patch.String())
		if err != nil {
			return fmt.Errorf("describe file changes of %q: %w", file, err)
		}

		summaries = append(summaries, summ)
	}

	comments, err := pf.Reviewer.AnalyzeChangesSummary(ctx, summaries)
	if err != nil {
		return fmt.Errorf("reviewer: analyze summaries: %w", err)
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
