package ports

import (
	"context"

	"bmad-cli/internal/domain/models"
)

type GitHubService interface {
	GetPRNumber(ctx context.Context) (int, error)
	FetchThreads(ctx context.Context, prNumber int) ([]models.Thread, error)
	ResolveThread(ctx context.Context, threadID, message string) error
}
