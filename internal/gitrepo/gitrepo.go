package gitrepo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/WinPooh32/git"
)

type DiffOperation int

const (
	Equal DiffOperation = iota
	Add
)

type BalamePart struct {
	Operation DiffOperation
	Commit    *git.Commit
	StartLine int
	Lines     []string
}

type BlamePatch struct {
	File       string
	Parts      []BalamePart
	Commits    []*git.Commit
	Size       int64
	LinesCount int
}

func (bp *BlamePatch) String() string {
	var sb strings.Builder

	const (
		fileSymbs    = 9
		numberSymbs  = 3
		commitsSymbs = 18
	)

	sb.Grow(int(bp.Size) + fileSymbs + len(bp.File) + (numberSymbs+commitsSymbs)*bp.LinesCount)

	sb.WriteString("<|file|>")
	sb.WriteString(bp.File)
	sb.WriteByte('\n')

	sb.WriteString("<|commits|>")
	sb.WriteByte('\n')

	for _, c := range bp.Commits {
		sb.WriteString("<|commit|>")
		sb.WriteString(c.ID.String()[:7])
		sb.WriteString("<|commit_message|>")
		sb.WriteString(strings.TrimSuffix(c.CommitMessage, "\n"))
		sb.WriteByte('\n')
	}

	sb.WriteString("<|patch_content|>")
	sb.WriteByte('\n')

	var prevOpIsEqual bool

	for _, part := range bp.Parts {
		switch part.Operation {
		case Equal:
			if !prevOpIsEqual {
				sb.WriteString("<|equal|>")
				sb.WriteByte('\n')
			}
		case Add:
			sb.WriteString("<|commit|>")
			sb.WriteString(part.Commit.ID.String()[:7])
			sb.WriteString("<|add|>")
			sb.WriteByte('\n')
		}

		bp.writeLines(&sb, part.StartLine, part.Lines)

		prevOpIsEqual = part.Operation == Equal
	}

	return sb.String()
}

func (*BlamePatch) writeLines(sb *strings.Builder, startLine int, lines []string) {
	for i, l := range lines {
		n := startLine + i

		sb.WriteString(strconv.Itoa(n))
		sb.WriteString(":")
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
}

type Repository struct {
	rootDir      string
	objectFormat git.ObjectFormat
	repo         *git.Repository
}

func OpenRepository(ctx context.Context, rootDir string) (*Repository, error) {
	repo, err := git.OpenRepository(ctx, rootDir)
	if err != nil {
		return nil, fmt.Errorf("git: open repository: %w", err)
	}

	objectFormat, err := repo.GetObjectFormat()
	if err != nil {
		return nil, fmt.Errorf("git repo: get object format %w", err)
	}

	return &Repository{
		rootDir:      rootDir,
		repo:         repo,
		objectFormat: objectFormat,
	}, nil
}

func (r *Repository) Close() error {
	if err := r.repo.Close(); err != nil {
		return fmt.Errorf("close git repo: %w", err)
	}

	return nil
}

func (r *Repository) DiffFiles(_ context.Context, baseBranch, headBranch string) ([]string, error) {
	commit, _, err := r.repo.GetMergeBase("", baseBranch, headBranch)
	if err != nil {
		return nil, fmt.Errorf("git repo: get merge base: %w", err)
	}

	files, err := r.repo.GetFilesChangedBetween(commit, headBranch)
	if err != nil {
		return nil, fmt.Errorf("git repo: get files changed: %w", err)
	}

	return files, nil
}

func (r *Repository) DiffCommits(_ context.Context, baseBranch, headBranch string) (map[string]*git.Commit, error) {
	compInfo, err := r.repo.GetCompareInfo(r.rootDir, baseBranch, headBranch, true, false)
	if err != nil {
		return nil, fmt.Errorf("git: get compare info of base branch=%q, head branch=%q: %w", baseBranch, headBranch, err)
	}

	diffCommits := make(map[string]*git.Commit, len(compInfo.Commits))

	for _, c := range compInfo.Commits {
		diffCommits[c.ID.String()] = c
	}

	return diffCommits, nil
}

func (r *Repository) BlamePatch(ctx context.Context, headBranch string, diffCommits map[string]*git.Commit, fileRelPath string) (*BlamePatch, error) {
	headCommit, err := r.repo.GetBranchCommit(headBranch)
	if err != nil {
		return nil, fmt.Errorf("git repo: get branch commit for %q: %w", headBranch, err)
	}

	fileEntry, err := headCommit.GetTreeEntryByPath(fileRelPath)
	if err != nil {
		return nil, fmt.Errorf("head commit: get file entry: %w", err)
	}

	brd, err := git.CreateBlameReader(ctx, r.objectFormat, r.rootDir, headCommit, fileRelPath, false)
	if err != nil {
		return nil, fmt.Errorf("git: create blame reader for %q at the commit %q: %w", fileRelPath, headCommit.ID, err)
	}
	defer brd.Close()

	metCommits := map[string]struct{}{}

	bp := &BlamePatch{
		File:       fileRelPath,
		Parts:      []BalamePart{},
		Commits:    []*git.Commit{},
		Size:       fileEntry.Size(),
		LinesCount: 0,
	}

reading:
	for i := 0; ; i++ {
		blamePart, err := brd.NextPart()
		if errors.Is(err, io.EOF) {
			break reading
		}
		if err != nil {
			return nil, fmt.Errorf("blame reader: next part: %w", err)
		}

		if blamePart == nil {
			break reading
		}

		blameCommit, err := r.repo.GetCommit(blamePart.Sha)
		if err != nil {
			return nil, fmt.Errorf("get commit of the blame part: %w", err)
		}

		part := BalamePart{
			StartLine: bp.LinesCount + 1,
			Lines:     blamePart.Lines,
		}

		if _, ok := diffCommits[blameCommit.ID.String()]; ok {
			part.Operation = Add
			part.Commit = blameCommit

			if _, ok := metCommits[blameCommit.ID.String()]; !ok {
				metCommits[blameCommit.ID.String()] = struct{}{}
				bp.Commits = append(bp.Commits, blameCommit)
			}
		} else {
			part.Operation = Equal
		}

		bp.Parts = append(bp.Parts, part)
		bp.LinesCount += len(part.Lines)
	}

	return bp, nil
}
