package ai

import (
	"context"
	"fmt"
	"time"

	claudecode "github.com/severity1/claude-code-sdk-go"
)

type ClaudeClient struct {
	// No persistent client needed with severity1 SDK
}

func NewClaudeClient() (*ClaudeClient, error) {
	// No initialization needed with severity1 SDK - clients are created per-request
	return &ClaudeClient{}, nil
}

func (c *ClaudeClient) Name() string {
	return "Claude"
}

func (c *ClaudeClient) ExecutePrompt(ctx context.Context, prompt string, model string, mode ExecutionMode) (string, error) {
	fmt.Printf("🔄 Calling claude with prompt length: %d\n", len(prompt))

	// Set timeout for large prompts - 5 minutes should be sufficient for most cases
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Build options based on execution mode with strict file system restrictions
	var opts []claudecode.Option

	// Set model based on parameter
	if model == "opus" {
		opts = append(opts, claudecode.WithModel("claude-3-opus-20240229")) // Use Opus 4.1 when available
	} else {
		opts = append(opts, claudecode.WithModel("claude-3-5-sonnet-20241022"))
	}

	// Apply mode permissions directly from struct fields
	opts = append(opts, claudecode.WithPermissionMode(claudecode.PermissionModeAcceptEdits))

	if len(mode.AllowedTools) > 0 {
		fmt.Printf("🔧 Allowed tools: %v\n", mode.AllowedTools)
		opts = append(opts, claudecode.WithAllowedTools(mode.AllowedTools...))
	}

	if len(mode.DisallowedTools) > 0 {
		fmt.Printf("🚫 Disallowed tools: %v\n", mode.DisallowedTools)
		opts = append(opts, claudecode.WithDisallowedTools(mode.DisallowedTools...))
	}

	// Use Query for one-shot execution with timeout context
	iterator, err := claudecode.Query(timeoutCtx, prompt, opts...)
	if err != nil {
		return "", fmt.Errorf("claude execution failed: %w", err)
	}
	defer iterator.Close()

	// Collect all response text
	var result string
	for {
		msg, err := iterator.Next(timeoutCtx)
		if err != nil {
			if err == claudecode.ErrNoMoreMessages {
				break
			}
			return "", fmt.Errorf("failed to read response: %w", err)
		}

		if msg == nil {
			break
		}

		// Handle different message types
		switch message := msg.(type) {
		case *claudecode.AssistantMessage:
			for _, block := range message.Content {
				if textBlock, ok := block.(*claudecode.TextBlock); ok {
					result += textBlock.Text
				}
			}
		case *claudecode.ResultMessage:
			if message.IsError {
				return "", fmt.Errorf("Claude returned error: %v", message.Result)
			}
		}
	}

	fmt.Printf("✅ Claude returned result length: %d\n", len(result))
	return result, nil
}
