package github

import (
	"context"

	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/infrastructure/shell"
)

type GitHubService struct {
	prFetcher      *PRNumberFetcher
	threadsFetcher *ThreadsFetcher
	threadResolver *ThreadResolver
}

func NewGitHubService(shell shell.Executor) *GitHubService {
	client := NewGitHubCLIClient(shell)

	return &GitHubService{
		prFetcher:      NewPRNumberFetcher(client),
		threadsFetcher: NewThreadsFetcher(client),
		threadResolver: NewThreadResolver(client),
	}
}

func (s *GitHubService) GetPRNumber(ctx context.Context) (int, error) {
	return s.prFetcher.Fetch(ctx)
}

func (s *GitHubService) FetchThreads(ctx context.Context, prNumber int) ([]models.Thread, error) {
	return s.threadsFetcher.FetchAll(ctx, prNumber)
}

func (s *GitHubService) ResolveThread(ctx context.Context, threadID, message string) error {
	return s.threadResolver.Resolve(ctx, threadID, message)
}
