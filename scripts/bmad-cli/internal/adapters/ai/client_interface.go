package ai

import "context"

type AIClient interface {
	ExecutePrompt(ctx context.Context, prompt string, mode ExecutionMode) (string, error)
	Name() string
}
