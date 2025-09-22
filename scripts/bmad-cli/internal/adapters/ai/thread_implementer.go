package ai

import (
	"context"
	"fmt"
	"strings"

	"bmad-cli/internal/application/prompts"
	"bmad-cli/internal/common/utils"
	"bmad-cli/internal/domain/models"
)

type ThreadImplementer struct {
	client        AIClient
	promptBuilder *prompts.ImplementationPromptBuilder
}

func NewThreadImplementer(
	client AIClient,
	promptBuilder *prompts.ImplementationPromptBuilder,
) *ThreadImplementer {
	return &ThreadImplementer{
		client:        client,
		promptBuilder: promptBuilder,
	}
}

func (ti *ThreadImplementer) Implement(ctx context.Context, threadContext models.ThreadContext) (string, error) {
	prompt, err := ti.promptBuilder.Build(threadContext)
	if err != nil {
		return "", fmt.Errorf("failed to build implementation prompt: %w", err)
	}

	utils.DebugLogWithSeparator(ti.client.Name()+" implementation prompt", prompt)

	rawOutput, err := ti.client.ExecutePrompt(ctx, prompt, ApplyMode)
	if err != nil {
		return "", fmt.Errorf("AI client implementation failed: %w", err)
	}

	utils.DebugLogWithSeparator(ti.client.Name()+" implementation output", rawOutput)

	summary := strings.TrimSpace(utils.FirstLine(rawOutput))
	if summary == "" {
		summary = "Applied changes as requested"
	}

	return summary, nil
}
