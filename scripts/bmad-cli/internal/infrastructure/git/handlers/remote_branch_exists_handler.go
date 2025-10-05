package handlers

import (
	"bmad-cli/internal/infrastructure/git"
	"context"
	"log/slog"
)

// RemoteBranchExistsHandler checks if the branch exists on remote and checks it out
type RemoteBranchExistsHandler struct {
	BaseBranchHandler
	gitService *git.GitService
}

// NewRemoteBranchExistsHandler creates a new handler
func NewRemoteBranchExistsHandler(gitService *git.GitService) *RemoteBranchExistsHandler {
	return &RemoteBranchExistsHandler{
		gitService: gitService,
	}
}

// Handle checks if branch exists on remote and checks it out
func (h *RemoteBranchExistsHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Debug("Checking if remote branch exists", "branch", branchCtx.ExpectedBranch)

	exists, err := h.gitService.RemoteBranchExists(ctx, branchCtx.ExpectedBranch)
	if err != nil {
		slog.Error("Failed to check remote branch", "error", err)
		return err
	}

	if !exists {
		slog.Debug("Remote branch does not exist, continuing chain")
		return h.callNext(ctx, branchCtx)
	}

	slog.Info("Remote branch exists, checking it out", "branch", branchCtx.ExpectedBranch)

	if err := h.gitService.CheckoutRemoteBranch(ctx, branchCtx.ExpectedBranch); err != nil {
		slog.Error("Failed to checkout remote branch", "error", err)
		return err
	}

	branchCtx.Action = ActionCheckout
	slog.Info("Successfully checked out remote branch", "branch", branchCtx.ExpectedBranch)
	return nil
}
