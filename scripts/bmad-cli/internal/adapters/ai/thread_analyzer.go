package ai

import (
	"context"
	"fmt"
	"strings"

	"bmad-cli/internal/application/prompts"
	"bmad-cli/internal/common/errors"
	"bmad-cli/internal/domain/models"
)

type ThreadAnalyzer struct {
	client        *ClaudeClient
	promptBuilder *prompts.HeuristicPromptBuilder
	parser        *prompts.YAMLParser
	modeFactory   *ModeFactory
}

func NewThreadAnalyzer(
	client *ClaudeClient,
	promptBuilder *prompts.HeuristicPromptBuilder,
	parser *prompts.YAMLParser,
	modeFactory *ModeFactory,
) *ThreadAnalyzer {
	return &ThreadAnalyzer{
		client:        client,
		promptBuilder: promptBuilder,
		parser:        parser,
		modeFactory:   modeFactory,
	}
}

func (ta *ThreadAnalyzer) Analyze(ctx context.Context, threadContext models.ThreadContext) (models.HeuristicAnalysisResult, error) {
	prompt, err := ta.promptBuilder.Build(threadContext)
	if err != nil {
		return models.HeuristicAnalysisResult{}, fmt.Errorf("failed to build heuristic prompt: %w", err)
	}

	rawOutput, err := ta.client.ExecutePromptWithSystem(ctx, "", prompt, "sonnet", ta.modeFactory.GetThinkMode())
	if err != nil {
		return models.HeuristicAnalysisResult{}, fmt.Errorf("AI client execution failed: %w", err)
	}

	if strings.TrimSpace(rawOutput) == "" {
		return models.HeuristicAnalysisResult{}, errors.ErrEmptyClientOutput(ta.client.Name())
	}

	result, err := ta.parser.ParseHeuristicResult(rawOutput)
	if err != nil {
		return models.HeuristicAnalysisResult{}, fmt.Errorf("failed to parse %s output: %w", ta.client.Name(), err)
	}

	return result, nil
}
