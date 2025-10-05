package github

import (
	"context"

	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/infrastructure/shell"
)

type GitHubPort struct {
	prFetcher      *PRNumberFetcher
	threadsFetcher *ThreadsFetcher
	threadResolver *ThreadResolver
}

func NewGitHubPort(shell *shell.CommandRunner) *GitHubPort {
	client := NewGitHubCLIClient(shell)

	return &GitHubPort{
		prFetcher:      NewPRNumberFetcher(client),
		threadsFetcher: NewThreadsFetcher(client),
		threadResolver: NewThreadResolver(client),
	}
}

func (s *GitHubPort) GetPRNumber(ctx context.Context) (int, error) {
	return s.prFetcher.Fetch(ctx)
}

func (s *GitHubPort) FetchThreads(ctx context.Context, prNumber int) ([]models.Thread, error) {
	return s.threadsFetcher.FetchAll(ctx, prNumber)
}

func (s *GitHubPort) ResolveThread(ctx context.Context, threadID, message string) error {
	return s.threadResolver.Resolve(ctx, threadID, message)
}
