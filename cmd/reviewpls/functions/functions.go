package functions

import (
	"context"
)

type ChangeCommentary struct {
	Line int
	Text string
}

type ChangesSummary struct {
	File     string
	Summary  string
	Comments []ChangeCommentary
}

type ChangesDescriber interface {
	DescribeFileChanges(ctx context.Context, file, filePatch string) (ChangesSummary, error)
}

type ReviewComment struct {
	File string
	Line int
	Text string
}

type Reviewer interface {
	AnalyzeChangesSummary(ctx context.Context, summaries []ChangesSummary) ([]ReviewComment, error)
}
