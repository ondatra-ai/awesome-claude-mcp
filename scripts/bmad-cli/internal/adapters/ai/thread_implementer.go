package ai

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"bmad-cli/internal/application/prompts"
	"bmad-cli/internal/domain/models"
)

type ThreadImplementer struct {
	client        *ClaudeClient
	promptBuilder *prompts.ImplementationPromptBuilder
	modeFactory   *ModeFactory
}

func NewThreadImplementer(
	client *ClaudeClient,
	promptBuilder *prompts.ImplementationPromptBuilder,
	modeFactory *ModeFactory,
) *ThreadImplementer {
	return &ThreadImplementer{
		client:        client,
		promptBuilder: promptBuilder,
		modeFactory:   modeFactory,
	}
}

func (ti *ThreadImplementer) Implement(ctx context.Context, threadContext models.ThreadContext) (string, error) {
	prompt, err := ti.promptBuilder.Build(threadContext)
	if err != nil {
		return "", fmt.Errorf("failed to build implementation prompt: %w", err)
	}

	slog.Debug("Implementation prompt", "client", ti.client.Name(), "prompt", prompt)

	rawOutput, err := ti.client.ExecutePrompt(ctx, prompt, "sonnet", ti.modeFactory.GetThinkMode())
	if err != nil {
		return "", fmt.Errorf("AI client implementation failed: %w", err)
	}

	slog.Debug("Implementation output", "client", ti.client.Name(), "output", rawOutput)

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
