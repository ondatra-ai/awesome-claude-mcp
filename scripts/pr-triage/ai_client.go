package main

import "context"

// ExecutionMode represents the mode for AI prompt execution.
type ExecutionMode string

const (
	PlanMode  ExecutionMode = "plan"  // Analyze and plan without making changes
	ApplyMode ExecutionMode = "apply" // Execute changes to files
)

// AIClient defines the strategy interface for different AI providers
// Each implementation focuses purely on executing prompts and returning raw output.
type AIClient interface {
	// ExecutePrompt executes a prompt with the given mode and returns raw output
	ExecutePrompt(ctx context.Context, prompt string, mode ExecutionMode) (string, error)

	// Name returns the client identifier for logging/debugging
	Name() string
}
