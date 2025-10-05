package services

import (
	"context"

	"bmad-cli/internal/adapters/ai"
)

// AIClient defines the interface for AI communication
// This interface belongs in domain/services as it represents a domain service contract
type AIClient interface {
	ExecutePromptWithSystem(ctx context.Context, systemPrompt string, userPrompt string, model string, mode ai.ExecutionMode) (string, error)
}
