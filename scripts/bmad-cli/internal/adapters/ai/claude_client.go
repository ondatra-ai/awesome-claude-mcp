package ai

import (
	"context"
	"fmt"

	"github.com/lancekrogers/claude-code-go/pkg/claude"
	"github.com/lancekrogers/claude-code-go/pkg/claude/dangerous"
)

type ClaudeClient struct {
	client          *claude.ClaudeClient
	dangerousClient *dangerous.DangerousClient
}

func NewClaudeClient() (*ClaudeClient, error) {
	// Try to create the main claude client
	client := claude.NewClient("claude")

	// Try to create dangerous client with error handling
	dangerousClient, err := dangerous.NewDangerousClient("claude")
	if err != nil {
		return nil, fmt.Errorf("failed to create dangerous client: %w", err)
	}

	return &ClaudeClient{
		client:          client,
		dangerousClient: dangerousClient,
	}, nil
}

func (c *ClaudeClient) Name() string {
	return "Claude"
}

func (c *ClaudeClient) ExecutePrompt(ctx context.Context, prompt string, mode ExecutionMode) (string, error) {
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

// GenerateContent generates content using Claude for general purposes
func (c *ClaudeClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// Use dangerous client with bypass permissions for task generation
	result, err := c.dangerousClient.BYPASS_ALL_PERMISSIONS(prompt, &claude.RunOptions{
		Format: claude.TextOutput,
	})
	if err != nil {
		return "", fmt.Errorf("claude content generation failed: %w", err)
	}

	return result.Result, nil
}
