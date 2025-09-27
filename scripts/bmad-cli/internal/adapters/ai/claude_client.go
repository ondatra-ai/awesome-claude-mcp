package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/lancekrogers/claude-code-go/pkg/claude"
	"github.com/lancekrogers/claude-code-go/pkg/claude/dangerous"
)

type ClaudeClient struct {
	client          *claude.ClaudeClient
	dangerousClient *dangerous.DangerousClient
}

func NewClaudeClient() (*ClaudeClient, error) {
	// Set required environment variable for dangerous client
	os.Setenv("CLAUDE_ENABLE_DANGEROUS", "i-accept-all-risks")

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
	fmt.Printf("🔄 Calling claude with prompt length: %d\n", len(prompt))

	var opts *claude.RunOptions
	switch mode {
	case PlanMode:
		opts = &claude.RunOptions{
			Format:         claude.TextOutput,
			PermissionTool: "plan",
		}
	case ApplyMode:
		opts = &claude.RunOptions{
			Format: claude.TextOutput,
		}
	default:
		return "", fmt.Errorf("unsupported execution mode: %v", mode)
	}

	result, err := c.dangerousClient.BYPASS_ALL_PERMISSIONS(prompt, opts)
	if err != nil {
		return "", fmt.Errorf("claude execution failed: %w", err)
	}

	fmt.Printf("✅ Claude returned result length: %d\n", len(result.Result))
	return result.Result, nil
}
