package changesdescriber

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions"
	"github.com/WinPooh32/reviewpls/cmd/reviewpls/functions/changes-describer/prompts"
	"github.com/WinPooh32/reviewpls/internal/prompt"
	"github.com/WinPooh32/reviewpls/internal/retry"
	sc "github.com/WinPooh32/reviewpls/internal/schema"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
)

type ChangesDescriberOpenAIConfig struct {
	Model            string
	RetryMaxAttempts int
	RetryDelay       time.Duration
}

type ChangesDescriberOpenAI struct {
	cfg     ChangesDescriberOpenAIConfig
	cli     openai.Client
	prompts map[string]prompt.Prompt
}

func NewChangesDescriberOpenAI(cfg ChangesDescriberOpenAIConfig, cli openai.Client) (*ChangesDescriberOpenAI, error) {
	prompts, err := prompt.Load(prompts.Files, "*.tpl")
	if err != nil {
		return nil, fmt.Errorf("load prompts: %w", err)
	}

	if _, ok := prompts["context"]; !ok {
		return nil, fmt.Errorf("'context' prompt is not loaded")
	}

	if _, ok := prompts["describe_changes_0"]; !ok {
		return nil, fmt.Errorf("'describe_changes_0' prompt is not loaded")
	}

	if _, ok := prompts["describe_changes_1"]; !ok {
		return nil, fmt.Errorf("'describe_changes_1' prompt is not loaded")
	}

	return &ChangesDescriberOpenAI{
		cfg:     cfg,
		cli:     cli,
		prompts: prompts,
	}, nil
}

func (c *ChangesDescriberOpenAI) DescribeFileChanges(ctx context.Context, file, filePatch string) (summ functions.ChangesSummary, err error) {
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

func (c *ChangesDescriberOpenAI) describeFileChanges(ctx context.Context, file, filePatch string) (summ functions.ChangesSummary, err error) {
	JSONSchemaRoot := sc.Object().
		AdditionalProperties(false).
		Required(
			"Summary",
			"Comments",
		).
		Property("Summary",
			sc.String().
				Description("Brief description of the changes."),
		).
		Property("Comments",
			sc.Array().
				MinItems(0).
				Items(
					sc.Object().
						AdditionalProperties(false).
						Required(
							"Line",
							"Text",
						).
						Property("Line",
							sc.Integer().
								Description("Start line number of the commented code."),
						).
						Property("Text",
							sc.String().
								Description("Commentary content text."),
						),
				),
		)

	messageSchema, err := json.Marshal(&JSONSchemaRoot)
	if err != nil {
		return summ, fmt.Errorf("encode schema json: %w", err)
	}

	pc, ok := c.prompts["context"]
	if !ok {
		return summ, errors.Join(errors.New("can't find 'context' prompt"), retry.ErrFatal)
	}

	pdc0, ok := c.prompts["describe_changes_0"]
	if !ok {
		return summ, errors.Join(errors.New("can't find 'describe_changes_0' prompt"), retry.ErrFatal)
	}

	pdc1, ok := c.prompts["describe_changes_1"]
	if !ok {
		return summ, errors.Join(errors.New("can't find 'describe_changes_1' prompt"), retry.ErrFatal)
	}

	contextMessage, err := pc.Execute(nil)
	if err != nil {
		return summ, errors.Join(fmt.Errorf("execute 'context' prompt template: %w", err), retry.ErrFatal)
	}

	task0, err := pdc0.Execute(map[string]any{"Patch": filePatch})
	if err != nil {
		return summ, errors.Join(fmt.Errorf("execute `describe_changes_0` prompt template: %w", err), retry.ErrFatal)
	}

	task1, err := pdc1.Execute(map[string]any{
		"JSONSchema": messageSchema,
	})
	if err != nil {
		return summ, errors.Join(fmt.Errorf("execute `describe_changes_1` prompt template: %w", err), retry.ErrFatal)
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(contextMessage),
			openai.UserMessage(task0),
		},
		Model: c.cfg.Model,
	}

	chatCompletion0, err := c.cli.Chat.Completions.New(ctx, params)
	if err != nil {
		return summ, fmt.Errorf("new completion: %w", err)
	}

	completion := chatCompletion0.Choices[0]

	if err := checkCompletion(completion); err != nil {
		return summ, err
	}

	params.Messages = append(params.Messages,
		[]openai.ChatCompletionMessageParamUnion{
			completion.Message.ToParam(),
			openai.UserMessage(task1),
		}...,
	)

	params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
			Type: "json_object",
			JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:   "ChangesReview",
				Strict: openai.Bool(true),
				Schema: json.RawMessage(messageSchema),
			},
		},
	}

	chatCompletion1, err := c.cli.Chat.Completions.New(ctx, params)
	if err != nil {
		return summ, fmt.Errorf("new completion: %w", err)
	}

	completion = chatCompletion1.Choices[0]

	if err := checkCompletion(completion); err != nil {
		return summ, err
	}

	if err := json.Unmarshal([]byte(completion.Message.Content), &summ); err != nil {
		return summ, fmt.Errorf("parse model json response: %w", err)
	}

	summ.File = file

	return summ, nil
}

func checkCompletion(completion openai.ChatCompletionChoice) error {
	if completion.FinishReason != "stop" {
		return fmt.Errorf("unexpected finish reason: %s", completion.FinishReason)
	}

	return nil
}
