package handlers

import (
	"bmad-cli/internal/domain/ports"
	"context"
	"errors"
	"log/slog"

	pkgerrors "bmad-cli/internal/pkg/errors"
)

// GitRepoCheckHandler validates that the current directory is a git repository.
type GitRepoCheckHandler struct {
	BaseBranchHandler

	gitService ports.GitPort
}

// NewGitRepoCheckHandler creates a new git repository check handler.
func NewGitRepoCheckHandler(gitService ports.GitPort) *GitRepoCheckHandler {
	return &GitRepoCheckHandler{
		gitService: gitService,
	}
}

// Handle validates that the current directory is a git repository.
func (h *GitRepoCheckHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Debug("Checking if current directory is a git repository")

	isRepo, err := h.gitService.IsGitRepository(ctx)
	if err != nil {
		slog.Error("Failed to check git repository", "error", err)

		return pkgerrors.ErrCheckGitRepositoryFailed(err)
	}

	if !isRepo {
		slog.Error("Current directory is not a git repository")

		return errors.New("current directory is not a git repository")
	}

	slog.Debug("Git repository check passed")

	return h.callNext(ctx, branchCtx)
}
