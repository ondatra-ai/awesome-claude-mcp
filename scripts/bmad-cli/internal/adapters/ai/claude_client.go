package ai

import (
	"context"
	"fmt"

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

func (c *ClaudeClient) ExecutePrompt(ctx context.Context, prompt string, mode ExecutionMode) (string, error) {
	fmt.Printf("🔄 Calling claude with prompt length: %d\n", len(prompt))

	// Build options based on execution mode with strict file system restrictions
	var opts []claudecode.Option
	switch mode {
	case PlanMode:
		opts = []claudecode.Option{
			claudecode.WithPermissionMode(claudecode.PermissionModePlan),
			// Only allow reading from specific patterns and writing to tmp directory
			claudecode.WithAllowedTools(
				"Read(**.yaml)",           // Allow reading YAML files
				"Read(**.yml)",            // Allow reading YML files
				"Read(**.md)",             // Allow reading markdown files
				"Write(./tmp/**)",         // Only allow writing to tmp directory
				"Glob(**)",                // Allow file discovery
				"Grep(**)",                // Allow searching content
			),
			// Explicitly disallow potentially dangerous tools
			claudecode.WithDisallowedTools(
				"Bash",                    // No shell commands
				"Edit",                    // No file editing outside tmp
				"MultiEdit",               // No multi-file editing
				"Write(**.go)",            // No Go file modifications
				"Write(**.yaml)",          // No YAML file modifications outside tmp
				"Write(**.yml)",           // No YML file modifications outside tmp
				"Write(**.json)",          // No JSON file modifications outside tmp
			),
		}
	case ApplyMode:
		opts = []claudecode.Option{
			claudecode.WithPermissionMode(claudecode.PermissionModeAcceptEdits),
			// Even in apply mode, restrict to tmp directory for file creation
			claudecode.WithAllowedTools(
				"Read(**)",                // Allow reading all files
				"Write(./tmp/**)",         // Only allow writing to tmp directory
				"Glob(**)",                // Allow file discovery
				"Grep(**)",                // Allow searching content
			),
			claudecode.WithDisallowedTools(
				"Bash",                    // No shell commands
				"Edit(**.go)",             // No Go file editing
				"MultiEdit(**.go)",        // No Go multi-file editing
			),
		}
	default:
		return "", fmt.Errorf("unsupported execution mode: %v", mode)
	}

	// Use Query for one-shot execution
	iterator, err := claudecode.Query(ctx, prompt, opts...)
	if err != nil {
		return "", fmt.Errorf("claude execution failed: %w", err)
	}
	defer iterator.Close()

	// Collect all response text
	var result string
	for {
		msg, err := iterator.Next(ctx)
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
