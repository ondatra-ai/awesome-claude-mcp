package git

import (
	"context"
	"log/slog"
)

// ForceRecreateHandler handles the --force flag to recreate branches
type ForceRecreateHandler struct {
	BaseBranchHandler
	gitService *GitService
}

// NewForceRecreateHandler creates a new force recreate handler
func NewForceRecreateHandler(gitService *GitService) *ForceRecreateHandler {
	return &ForceRecreateHandler{
		gitService: gitService,
	}
}

// Handle recreates the branch if --force flag is set
func (h *ForceRecreateHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	if !branchCtx.Force {
		slog.Debug("Force flag not set, continuing normal flow")
		return h.callNext(ctx, branchCtx)
	}

	slog.Info("Force flag detected, recreating branch", "branch", branchCtx.ExpectedBranch)

	if err := h.gitService.ForceRecreateBranch(ctx, branchCtx.ExpectedBranch); err != nil {
		slog.Error("Failed to force recreate branch", "error", err)
		return err
	}

	branchCtx.Action = ActionForceRecreate
	slog.Info("Branch successfully recreated", "branch", branchCtx.ExpectedBranch)

	// No need to continue chain after force recreate
	return nil
}
