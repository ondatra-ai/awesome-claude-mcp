package ports

import (
	"context"

	"bmad-cli/internal/adapters/ai"
)

// AIPort defines the interface for AI communication
// This port interface represents the contract for AI operations in the domain
// layer.
type AIPort interface {
	ExecutePromptWithSystem(
		ctx context.Context,
		systemPrompt string,
		userPrompt string,
		model string,
		mode ai.ExecutionMode,
	) (string, error)
}
