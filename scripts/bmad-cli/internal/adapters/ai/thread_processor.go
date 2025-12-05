package ai

import (
	"context"
	"log/slog"
	"strings"

	"bmad-cli/internal/app/prompt_builders"
	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/pkg/errors"
)

type ThreadProcessor struct {
	client                *ClaudeClient
	heuristicBuilder      *prompt_builders.HeuristicPromptBuilder
	implementationBuilder *prompt_builders.ImplementationPromptBuilder
	yamlParser            *prompt_builders.YAMLParser
	modeFactory           *ModeFactory
}

func NewThreadProcessor(
	client *ClaudeClient,
	heuristicBuilder *prompt_builders.HeuristicPromptBuilder,
	implementationBuilder *prompt_builders.ImplementationPromptBuilder,
	yamlParser *prompt_builders.YAMLParser,
	modeFactory *ModeFactory,
) *ThreadProcessor {
	return &ThreadProcessor{
		client:                client,
		heuristicBuilder:      heuristicBuilder,
		implementationBuilder: implementationBuilder,
		yamlParser:            yamlParser,
		modeFactory:           modeFactory,
	}
}

func (tp *ThreadProcessor) AnalyzeThread(
	ctx context.Context,
	threadContext models.ThreadContext,
) (models.HeuristicAnalysisResult, error) {
	prompt, err := tp.heuristicBuilder.Build(threadContext)
	if err != nil {
		return models.HeuristicAnalysisResult{}, errors.ErrBuildHeuristicPromptFailed(err)
	}

	rawOutput, err := tp.client.ExecutePromptWithSystem(ctx, "", prompt, "sonnet", tp.modeFactory.GetThinkMode())
	if err != nil {
		return models.HeuristicAnalysisResult{}, errors.ErrAIClientExecutionFailed(err)
	}

	if strings.TrimSpace(rawOutput) == "" {
		return models.HeuristicAnalysisResult{}, errors.ErrEmptyClientOutput(tp.client.Name())
	}

	result, err := tp.yamlParser.ParseHeuristicResult(rawOutput)
	if err != nil {
		return models.HeuristicAnalysisResult{}, errors.ErrParseAIOutputFailed(tp.client.Name(), err)
	}

	return result, nil
}

func (tp *ThreadProcessor) ImplementChanges(ctx context.Context, threadContext models.ThreadContext) (string, error) {
	prompt, err := tp.implementationBuilder.Build(threadContext)
	if err != nil {
		return "", errors.ErrBuildImplementationPromptFailed(err)
	}

	slog.Debug("Implementation prompt", "client", tp.client.Name(), "prompt", prompt)

	rawOutput, err := tp.client.ExecutePromptWithSystem(ctx, "", prompt, "sonnet", tp.modeFactory.GetThinkMode())
	if err != nil {
		return "", errors.ErrAIClientImplementationFailed(err)
	}

	slog.Debug("Implementation output", "client", tp.client.Name(), "output", rawOutput)

	// Extract first line as summary
	lines := strings.Split(rawOutput, "\n")

	summary := ""
	if len(lines) > 0 {
		summary = strings.TrimSpace(lines[0])
	}

	if summary == "" {
		summary = "Applied changes as requested"
	}

	return summary, nil
}
