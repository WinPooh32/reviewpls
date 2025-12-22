package reviewer

import (
	"context"
	"fmt"
	"time"

	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions"
	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions/reviewer/prompts"
	"github.com/WinPooh32/reviewpls/internal/prompt"
	"github.com/openai/openai-go/v2"
)

type ChangesDescriberConfig struct {
	Model              string
	RetryMaxAttempts   int
	RetryDelay         time.Duration
	CommentaryLanguage string
}

type Reviewer struct {
	cfg     ChangesDescriberConfig
	cli     openai.Client
	prompts prompt.PromptSet
}

func New(cfg ChangesDescriberConfig, cli openai.Client, model string) (*Reviewer, error) {
	prompts, err := prompt.Load(prompts.Files, "*.tpl")
	if err != nil {
		return nil, fmt.Errorf("load prompts: %w", err)
	}

	if err := prompts.Require(
		"context",
		"instructions",
	); err != nil {
		return nil, fmt.Errorf("require prompts: %w", err)
	}

	return &Reviewer{
		cfg: cfg,
		cli: cli,
	}, nil
}

func (r *Reviewer) AnalyzeChangesSummary(ctx context.Context, summaries *functions.ChangesSummary) (*functions.ReviewComment, error) {
	panic("TODO: Implement")
}
