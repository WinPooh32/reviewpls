package changesdescriber

import (
	"context"

	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions"
	"github.com/openai/openai-go/v2"
)

type ChangesDescriberOpenAI struct {
	cli   openai.Client
	model string
}

func NewChangesDescriberOpenAI(cli openai.Client, model string) *ChangesDescriberOpenAI {
	return &ChangesDescriberOpenAI{
		cli:   cli,
		model: model,
	}
}

func (c *ChangesDescriberOpenAI) DescribeFileChanges(ctx context.Context, file string) (functions.ChangesSummary, error) {
	panic("TODO: Implement")
}
