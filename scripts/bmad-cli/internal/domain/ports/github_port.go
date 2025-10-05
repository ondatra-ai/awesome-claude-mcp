package ports

import (
	"context"

	"bmad-cli/internal/domain/models"
)

// GitHubPort defines the port interface for GitHub operations
// Following Hexagonal Architecture, ports define contracts for external adapters
type GitHubPort interface {
	GetPRNumber(ctx context.Context) (int, error)
	FetchThreads(ctx context.Context, prNumber int) ([]models.Thread, error)
	ResolveThread(ctx context.Context, threadID, message string) error
}
