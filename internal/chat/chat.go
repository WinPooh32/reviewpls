package chat

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/WinPooh32/reviewpls/internal/slogutil"
	"github.com/openai/openai-go/v2"
)

func NewCompletion(ctx context.Context, completions openai.ChatCompletionService, params openai.ChatCompletionNewParams) (*openai.ChatCompletionChoice, error) {
	compl, err := completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("new completion: %w", err)
	}

	if len(compl.Choices) != 1 {
		return nil, fmt.Errorf("unexpected completion choices count")
	}

	choice := &compl.Choices[0]

	if err := checkStopReason(choice); err != nil {
		return nil, err
	}

	logCompletion(ctx, compl)

	return choice, nil
}

func logCompletion(ctx context.Context, completion *openai.ChatCompletion) {
	usage := completion.Usage
	choice := completion.Choices[0]

	slogutil.Ctx(ctx).Debug("new completion",
		slog.Group("usage",
			slog.Int64("total_tokens", usage.TotalTokens),
			slog.Int64("prompt_tokens", usage.PromptTokens),
			slog.Int64("completion_tokens", usage.CompletionTokens),
		),
		slog.String("choice", choice.Message.Content),
	)
}

func checkStopReason(completion *openai.ChatCompletionChoice) error {
	if completion.FinishReason != "stop" {
		return fmt.Errorf("unexpected finish reason: %s", completion.FinishReason)
	}

	return nil
}
