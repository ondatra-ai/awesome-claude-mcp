package git

import (
	"context"
	"fmt"
	"log/slog"
)

// GitRepoCheckHandler validates that the current directory is a git repository
type GitRepoCheckHandler struct {
	BaseBranchHandler
	gitService *GitService
}

// NewGitRepoCheckHandler creates a new git repository check handler
func NewGitRepoCheckHandler(gitService *GitService) *GitRepoCheckHandler {
	return &GitRepoCheckHandler{
		gitService: gitService,
	}
}

// Handle validates that the current directory is a git repository
func (h *GitRepoCheckHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Debug("Checking if current directory is a git repository")

	isRepo, err := h.gitService.IsGitRepository(ctx)
	if err != nil {
		slog.Error("Failed to check git repository", "error", err)
		return fmt.Errorf("failed to check git repository: %w", err)
	}

	if !isRepo {
		slog.Error("Current directory is not a git repository")
		return fmt.Errorf("current directory is not a git repository")
	}

	slog.Debug("Git repository check passed")
	return h.callNext(ctx, branchCtx)
}
