package main

import (
	"context"
	"fmt"

	"github.com/lancekrogers/claude-code-go/pkg/claude"
	"github.com/lancekrogers/claude-code-go/pkg/claude/dangerous"
)

// claudeStrategy implements AIClient interface using Claude Code Go SDK.
type claudeStrategy struct {
	client          *claude.ClaudeClient
	dangerousClient *dangerous.DangerousClient
}

// NewClaudeStrategy creates a new Claude AI client strategy.
func NewClaudeStrategy() (AIClient, error) {
	client := claude.NewClient("claude")

	// For ApplyMode, also initialize dangerous client
	dangerousClient, err := dangerous.NewDangerousClient("claude")
	if err != nil {
		return nil, fmt.Errorf("failed to create dangerous client: %w", err)
	}

	return &claudeStrategy{
		client:          client,
		dangerousClient: dangerousClient,
	}, nil
}

// Name returns the client identifier.
func (c *claudeStrategy) Name() string {
	return "Claude"
}

// ExecutePrompt executes a prompt using Claude Code Go SDK.
func (c *claudeStrategy) ExecutePrompt(ctx context.Context, prompt string, mode ExecutionMode) (string, error) {
	switch mode {
	case PlanMode:
		opts := &claude.RunOptions{
			Format:         claude.TextOutput,
			PermissionTool: "plan",
		}
		result, err := c.client.RunPrompt(prompt, opts)
		if err != nil {
			return "", fmt.Errorf("claude execution failed: %w", err)
		}
		return result.Result, nil

	case ApplyMode:
		result, err := c.dangerousClient.BYPASS_ALL_PERMISSIONS(prompt, &claude.RunOptions{
			Format: claude.TextOutput,
		})
		if err != nil {
			return "", fmt.Errorf("claude apply execution failed: %w", err)
		}
		return result.Result, nil

	default:
		return "", fmt.Errorf("unsupported execution mode: %v", mode)
	}
}
