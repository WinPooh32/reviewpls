package changesdescriber

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions"
	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions/changes-describer/prompts"
	"github.com/WinPooh32/reviewpls/internal/chat"
	"github.com/WinPooh32/reviewpls/internal/prompt"
	"github.com/WinPooh32/reviewpls/internal/retry"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
)

type ChangesDescriberConfig struct {
	Model            string
	RetryMaxAttempts int
	RetryDelay       time.Duration
}

type ChangesDescriber struct {
	cfg     ChangesDescriberConfig
	cli     openai.Client
	prompts prompt.PromptSet
}

func New(cfg ChangesDescriberConfig, cli openai.Client) (*ChangesDescriber, error) {
	prompts, err := prompt.Load(prompts.Files, "*.tpl")
	if err != nil {
		return nil, fmt.Errorf("load prompts: %w", err)
	}

	if err := prompts.Require(
		"context",
		"describe_changes_0",
		"describe_changes_1",
	); err != nil {
		return nil, fmt.Errorf("require prompts: %w", err)
	}

	return &ChangesDescriber{
		cfg:     cfg,
		cli:     cli,
		prompts: prompts,
	}, nil
}

func (c *ChangesDescriber) DescribeFileChanges(ctx context.Context, file, filePatch string) (summ *functions.ChangesSummary, err error) {
	err = retry.Run(ctx,
		func() error {
			summ, err = c.describeFileChanges(ctx, file, filePatch)
			return err
		},
		retry.WithMaxAttempts(c.cfg.RetryMaxAttempts),
		retry.WithDelay(c.cfg.RetryDelay),
	)

	return summ, err
}

func (c *ChangesDescriber) describeFileChanges(ctx context.Context, file, filePatch string) (*functions.ChangesSummary, error) {
	contextMessage, err := c.prompts.Format("context", nil)
	if err != nil {
		return nil, errors.Join(err, retry.ErrFatal)
	}

	task0, err := c.prompts.Format("describe_changes_0", map[string]any{"Patch": filePatch})
	if err != nil {
		return nil, errors.Join(err, retry.ErrFatal)
	}

	task1, err := c.prompts.Format("describe_changes_1", map[string]any{"JSONSchema": responseChangesSchema})
	if err != nil {
		return nil, errors.Join(err, retry.ErrFatal)
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(contextMessage),
			openai.UserMessage(task0),
		},
		Model: c.cfg.Model,
	}

	choice, err := chat.NewCompletion(ctx, c.cli.Chat.Completions, params)
	if err != nil {
		return nil, fmt.Errorf("new chat completion: %w", err)
	}

	params.Messages = append(params.Messages,
		[]openai.ChatCompletionMessageParamUnion{
			choice.Message.ToParam(),
			openai.UserMessage(task1),
		}...,
	)

	params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
			Type: "json_object",
			JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:   "ChangesReview",
				Strict: openai.Bool(true),
				Schema: json.RawMessage(responseChangesSchema),
			},
		},
	}

	choice, err = chat.NewCompletion(ctx, c.cli.Chat.Completions, params)
	if err != nil {
		return nil, fmt.Errorf("new chat completion: %w", err)
	}

	var summ functions.ChangesSummary

	if err := json.Unmarshal([]byte(choice.Message.Content), &summ); err != nil {
		return nil, fmt.Errorf("parse model json response: %w", err)
	}

	summ.File = file

	return &summ, nil
}
