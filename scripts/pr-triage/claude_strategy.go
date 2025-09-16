package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// claudeStrategy implements AIClient interface using Claude Code CLI.
type claudeStrategy struct{}

// NewClaudeStrategy creates a new Claude AI client strategy.
func NewClaudeStrategy() AIClient {
	return &claudeStrategy{}
}

// Name returns the client identifier.
func (c *claudeStrategy) Name() string {
	return "Claude"
}

// ExecutePrompt executes a prompt using Claude Code CLI.
func (c *claudeStrategy) ExecutePrompt(ctx context.Context, prompt string, mode ExecutionMode) (string, error) {
	args := []string{"claude", "--print"}

	// Configure permissions based on execution mode
	switch mode {
	case ApplyMode:
		// For apply mode, bypass all permissions to enable automated changes
		args = append(args, "--dangerously-skip-permissions")
	case PlanMode:
		// For plan mode, use plan permission mode to prevent changes
		args = append(args, "--permission-mode", "plan")
	}

	// Execute Claude with prompt via stdin
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(prompt)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("claude command execution failed: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}
