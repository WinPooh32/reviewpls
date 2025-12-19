package reviewer

import (
	"context"

	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions"
	"github.com/openai/openai-go/v2"
)

type Reviewer struct {
	cli   openai.Client
	model string
}

func New(cli openai.Client, model string) *Reviewer {
	return &Reviewer{
		cli:   cli,
		model: model,
	}
}

func (r *Reviewer) AnalyzeChangesSummary(ctx context.Context, summaries []*functions.ChangesSummary) ([]*functions.ReviewComment, error) {
	panic("TODO: Implement")
}
