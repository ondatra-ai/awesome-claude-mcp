package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	claudecode "bmad-cli/claudecode"
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
	slog.Info("Calling Claude", "prompt_length", len(prompt))

	// Set timeout for large prompts - 10 minutes for complex multi-file operations
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
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
		slog.Debug("Claude tools configured", "allowed_tools", mode.AllowedTools)
		opts = append(opts, claudecode.WithAllowedTools(mode.AllowedTools...))
	}

	if len(mode.DisallowedTools) > 0 {
		slog.Debug("Claude tools configured", "disallowed_tools", mode.DisallowedTools)
		opts = append(opts, claudecode.WithDisallowedTools(mode.DisallowedTools...))
	}

	// Use WithClient pattern with streaming to prevent buffer overflow
	var resultStr string
	err := claudecode.WithClient(timeoutCtx, func(client claudecode.Client) error {
		slog.Info("Connected to Claude client")

		// Send query using the client
		slog.Info("Sending query to Claude", "length", len(prompt))
		if err := client.Query(timeoutCtx, prompt); err != nil {
			slog.Error("Query failed", "error", err)
			return fmt.Errorf("failed to send query: %w", err)
		}
		slog.Debug("Query sent successfully")

		// Use local builder inside the function scope
		var result strings.Builder

		// Stream messages using exact pattern from SDK docs
		slog.Debug("Starting message stream")
		msgChan := client.ReceiveMessages(timeoutCtx)
		messageCount := 0
		for {
			select {
			case message := <-msgChan:
				messageCount++
				slog.Debug("Message received", "count", messageCount, "type", fmt.Sprintf("%T", message))
				if message == nil {
					slog.Error("Received nil message from Claude stream")
					return fmt.Errorf("Claude stream returned nil message")
				}

				switch msg := message.(type) {
				case *claudecode.AssistantMessage:
					slog.Debug("AssistantMessage received", "content_blocks", len(msg.Content))
					slog.Debug("AssistantMessage content", "msg", fmt.Sprintf("%+v", msg))
					for i, block := range msg.Content {
						slog.Debug("Processing content block", "index", i, "type", fmt.Sprintf("%T", block))
						if textBlock, ok := block.(*claudecode.TextBlock); ok {
							slog.Debug("TextBlock received")
							slog.Debug("TextBlock content", "text", textBlock.Text)
							result.WriteString(textBlock.Text)
						} else if toolUseBlock, ok := block.(*claudecode.ToolUseBlock); ok {
							slog.Debug("ToolUseBlock received")
							if toolBytes, err := json.MarshalIndent(toolUseBlock, "      ", "  "); err == nil {
								slog.Debug("ToolUseBlock details", "content", string(toolBytes))
							} else {
								slog.Debug("ToolUseBlock details (raw)", "content", fmt.Sprintf("%+v", toolUseBlock))
							}
						} else {
							slog.Debug("Unknown block type", "type", fmt.Sprintf("%T", block))
							if blockBytes, err := json.MarshalIndent(block, "      ", "  "); err == nil {
								slog.Debug("Unknown block content", "content", string(blockBytes))
							} else {
								slog.Debug("Unknown block content (raw)", "content", fmt.Sprintf("%+v", block))
							}
						}
					}
				case *claudecode.UserMessage:
					slog.Debug("UserMessage received")
					slog.Debug("UserMessage content", "msg", fmt.Sprintf("%+v", msg))
				case *claudecode.SystemMessage:
					slog.Debug("SystemMessage received")
					slog.Debug("SystemMessage content", "msg", fmt.Sprintf("%+v", msg))
				case *claudecode.ResultMessage:
					slog.Debug("ResultMessage received", "is_error", msg.IsError, "result", msg.Result)
					if msg.IsError {
						return fmt.Errorf("Claude returned error: %s", msg.Result)
					}
					resultStr = result.String()
					slog.Debug("ResultMessage success", "captured_chars", len(resultStr))
					return nil
				default:
					slog.Debug("Unhandled message type", "type", fmt.Sprintf("%T", message))
					slog.Debug("Unhandled message content", "msg", fmt.Sprintf("%+v", message))
				}
			case <-timeoutCtx.Done():
				slog.Warn("Timeout reached", "error", timeoutCtx.Err())
				return timeoutCtx.Err()
			}
		}
	}, opts...)

	if err != nil {
		slog.Error("WithClient error", "error", err)
		// Check for buffer overflow errors and provide helpful context
		errStr := err.Error()
		if strings.Contains(errStr, "token too long") || strings.Contains(errStr, "bufio.Scanner") {
			return "", fmt.Errorf("Claude response too large for buffer (using streaming approach): %w", err)
		}
		return "", fmt.Errorf("claude execution failed: %w", err)
	}

	slog.Info("Claude returned result", "length", len(resultStr))
	return resultStr, nil
}

func (c *ClaudeClient) ExecutePromptWithSystem(ctx context.Context, systemPrompt string, userPrompt string, model string, mode ExecutionMode) (string, error) {
	slog.Info("Calling Claude with system prompt", "system_length", len(systemPrompt), "user_length", len(userPrompt))

	// Set timeout for large prompts - 10 minutes for complex multi-file operations
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
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
		slog.Debug("Claude tools configured", "allowed_tools", mode.AllowedTools)
		opts = append(opts, claudecode.WithAllowedTools(mode.AllowedTools...))
	}

	if len(mode.DisallowedTools) > 0 {
		slog.Debug("Claude tools configured", "disallowed_tools", mode.DisallowedTools)
		opts = append(opts, claudecode.WithDisallowedTools(mode.DisallowedTools...))
	}

	// Combine system and user prompts into a single message for now
	// TODO: When claudecode supports system messages, use them separately
	combinedPrompt := fmt.Sprintf("System: %s\n\nUser: %s", systemPrompt, userPrompt)

	// Use WithClient pattern with streaming to prevent buffer overflow
	var resultStr string
	err := claudecode.WithClient(timeoutCtx, func(client claudecode.Client) error {
		slog.Info("Connected to Claude client")

		// Send query using the client
		slog.Info("Sending combined prompt to Claude", "length", len(combinedPrompt))
		if err := client.Query(timeoutCtx, combinedPrompt); err != nil {
			slog.Error("Query failed", "error", err)
			return fmt.Errorf("failed to send query: %w", err)
		}
		slog.Debug("Query sent successfully")

		// Use local builder inside the function scope
		var result strings.Builder

		// Stream messages using exact pattern from SDK docs
		slog.Debug("Starting message stream")
		msgChan := client.ReceiveMessages(timeoutCtx)
		messageCount := 0
		for {
			select {
			case message := <-msgChan:
				messageCount++
				slog.Debug("Message received", "count", messageCount, "type", fmt.Sprintf("%T", message))
				if message == nil {
					slog.Error("Received nil message from Claude stream")
					return fmt.Errorf("Claude stream returned nil message")
				}

				switch msg := message.(type) {
				case *claudecode.AssistantMessage:
					slog.Debug("AssistantMessage received", "content_blocks", len(msg.Content))
					slog.Debug("AssistantMessage content", "msg", fmt.Sprintf("%+v", msg))
					for i, block := range msg.Content {
						slog.Debug("Processing content block", "index", i, "type", fmt.Sprintf("%T", block))
						if textBlock, ok := block.(*claudecode.TextBlock); ok {
							slog.Debug("TextBlock received")
							slog.Debug("TextBlock content", "text", textBlock.Text)
							result.WriteString(textBlock.Text)
						} else if toolUseBlock, ok := block.(*claudecode.ToolUseBlock); ok {
							slog.Debug("ToolUseBlock received")
							if toolBytes, err := json.MarshalIndent(toolUseBlock, "      ", "  "); err == nil {
								slog.Debug("ToolUseBlock details", "content", string(toolBytes))
							} else {
								slog.Debug("ToolUseBlock details (raw)", "content", fmt.Sprintf("%+v", toolUseBlock))
							}
						} else {
							slog.Debug("Unknown block type", "type", fmt.Sprintf("%T", block))
							if blockBytes, err := json.MarshalIndent(block, "      ", "  "); err == nil {
								slog.Debug("Unknown block content", "content", string(blockBytes))
							} else {
								slog.Debug("Unknown block content (raw)", "content", fmt.Sprintf("%+v", block))
							}
						}
					}
				case *claudecode.UserMessage:
					slog.Debug("UserMessage received")
					slog.Debug("UserMessage content", "msg", fmt.Sprintf("%+v", msg))
				case *claudecode.SystemMessage:
					slog.Debug("SystemMessage received")
					slog.Debug("SystemMessage content", "msg", fmt.Sprintf("%+v", msg))
				case *claudecode.ResultMessage:
					slog.Debug("ResultMessage received", "is_error", msg.IsError, "result", msg.Result)
					if msg.IsError {
						return fmt.Errorf("Claude returned error: %s", msg.Result)
					}
					resultStr = result.String()
					slog.Debug("ResultMessage success", "captured_chars", len(resultStr))
					return nil
				default:
					slog.Debug("Unhandled message type", "type", fmt.Sprintf("%T", message))
					slog.Debug("Unhandled message content", "msg", fmt.Sprintf("%+v", message))
				}
			case <-timeoutCtx.Done():
				slog.Warn("Timeout reached", "error", timeoutCtx.Err())
				return timeoutCtx.Err()
			}
		}
	}, opts...)

	if err != nil {
		slog.Error("WithClient error", "error", err)
		// Check for buffer overflow errors and provide helpful context
		errStr := err.Error()
		if strings.Contains(errStr, "token too long") || strings.Contains(errStr, "bufio.Scanner") {
			return "", fmt.Errorf("Claude response too large for buffer (using streaming approach): %w", err)
		}
		return "", fmt.Errorf("claude execution failed: %w", err)
	}

	slog.Info("Claude returned result", "length", len(resultStr))
	return resultStr, nil
}
