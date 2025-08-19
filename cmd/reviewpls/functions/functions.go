package functions

import (
	"context"
)

type ChangeCommentary struct {
	Line string
	Text string
}

type ChangesSummary struct {
	File     string
	Summary  string
	Comments []ChangeCommentary
}

type ChangesDescriber interface {
	DescribeFileChanges(ctx context.Context, file string) (ChangesSummary, error)
}

type ReviewComment struct {
	File string
	Line string
	Text string
}

type Reviewer interface {
	AnalyzeChangesSummary(summaries []ChangesSummary) ([]ReviewComment, error)
}
