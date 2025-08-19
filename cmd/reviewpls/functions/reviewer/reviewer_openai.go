package reviewer

import (
	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions"
	"github.com/openai/openai-go/v2"
)

type ReviewerOpenAI struct {
	cli   openai.Client
	model string
}

func NewReviewerOpenAI(cli openai.Client, model string) *ReviewerOpenAI {
	return &ReviewerOpenAI{
		cli:   cli,
		model: model,
	}
}

func (r *ReviewerOpenAI) AnalyzeChangesSummary(summaries []functions.ChangesSummary) ([]functions.ReviewComment, error) {
	panic("TODO: Implement")
}
