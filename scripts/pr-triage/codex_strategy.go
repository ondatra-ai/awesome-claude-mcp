package main

import (
	"context"
	"time"
)

const (
	defaultCodexTimeout  = 30 * time.Second
	extendedCodexTimeout = 2 * time.Minute
	outputTruncateLength = 4096
)

// codexStrategy implements AIClient interface using Codex CLI.
type codexStrategy struct{}

// NewCodexStrategy creates a new Codex AI client strategy.
func NewCodexStrategy() AIClient {
	return &codexStrategy{}
}

// Name returns the client identifier.
func (c *codexStrategy) Name() string {
	return "Codex"
}

// ExecutePrompt executes a prompt using Codex CLI.
func (c *codexStrategy) ExecutePrompt(ctx context.Context, prompt string, mode ExecutionMode) (string, error) {
	// Ensure bounded execution if caller didn't set a deadline
	if _, ok := ctx.Deadline(); !ok {
		timeout := defaultCodexTimeout
		if mode == ApplyMode {
			timeout = extendedCodexTimeout
		}

		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	logDebugf("BEGIN_CODEX_RUN\n")
	logDebugf("mode: %s\n", mode)

	args := []string{"exec"}

	// Configure flags based on execution mode
	if mode == ApplyMode {
		args = append(args, "--full-auto")
	}

	// Execute using the shell utility with stdin
	out, err := runShellWithStdin(ctx, "codex", prompt, args...)
	if err != nil {
		logDebugf("exit: %d\n", exitCode(err))
	} else {
		logDebugf("exit: 0\n")
	}

	logDebugf("stdout:\n%s\n", truncate(out, outputTruncateLength))
	logDebugf("END_CODEX_RUN\n")

	if err != nil {
		return "", err
	}

	return out, nil
}
