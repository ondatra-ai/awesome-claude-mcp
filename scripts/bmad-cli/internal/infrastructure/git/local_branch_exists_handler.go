package git

import (
	"context"
	"log/slog"
)

// LocalBranchExistsHandler checks if the branch exists locally and switches to it
type LocalBranchExistsHandler struct {
	BaseBranchHandler
	gitService *GitService
}

// NewLocalBranchExistsHandler creates a new handler
func NewLocalBranchExistsHandler(gitService *GitService) *LocalBranchExistsHandler {
	return &LocalBranchExistsHandler{
		gitService: gitService,
	}
}

// Handle checks if branch exists locally and switches to it
func (h *LocalBranchExistsHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Debug("Checking if local branch exists", "branch", branchCtx.ExpectedBranch)

	exists, err := h.gitService.LocalBranchExists(ctx, branchCtx.ExpectedBranch)
	if err != nil {
		slog.Error("Failed to check local branch", "error", err)
		return err
	}

	if !exists {
		slog.Debug("Local branch does not exist, continuing chain")
		return h.callNext(ctx, branchCtx)
	}

	slog.Info("Local branch exists, switching to it", "branch", branchCtx.ExpectedBranch)

	if err := h.gitService.SwitchBranch(ctx, branchCtx.ExpectedBranch); err != nil {
		slog.Error("Failed to switch to local branch", "error", err)
		return err
	}

	branchCtx.Action = ActionSwitch
	slog.Info("Successfully switched to local branch", "branch", branchCtx.ExpectedBranch)
	return nil
}
