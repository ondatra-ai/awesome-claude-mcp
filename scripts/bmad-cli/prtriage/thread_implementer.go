package prtriage

import (
	"context"
	"fmt"
	"strings"
)

// ThreadImplementer handles implementing code changes for PR review threads.
type ThreadImplementer struct {
	client AIClient
}

// NewThreadImplementer creates a new ThreadImplementer.
func NewThreadImplementer(client AIClient) *ThreadImplementer {
	return &ThreadImplementer{client: client}
}

// Implement applies code changes based on the thread context and returns a summary.
func (ti *ThreadImplementer) Implement(ctx context.Context, threadContext ThreadContext) (string, error) {
	// Build the implementation prompt
	prompt, err := buildImplementCodePrompt(threadContext)
	if err != nil {
		return "", fmt.Errorf("failed to build implementation prompt: %w", err)
	}

	debugLogWithSeparator(ti.client.Name()+" implementation prompt", prompt)

	// Execute the prompt using the AI client strategy
	rawOutput, err := ti.client.ExecutePrompt(ctx, prompt, ApplyMode)
	if err != nil {
		return "", fmt.Errorf("AI client implementation failed: %w", err)
	}

	debugLogWithSeparator(ti.client.Name()+" implementation output", rawOutput)

	// Extract the summary line from the output
	summary := strings.TrimSpace(firstLine(rawOutput))
	if summary == "" {
		summary = "Applied changes as requested"
	}

	return summary, nil
}
