package git

import (
	"context"
	"fmt"
	"log/slog"
)

// DetachedHeadHandler prevents operations on detached HEAD state
type DetachedHeadHandler struct {
	BaseBranchHandler
	gitService *GitService
}

// NewDetachedHeadHandler creates a new detached HEAD handler
func NewDetachedHeadHandler(gitService *GitService) *DetachedHeadHandler {
	return &DetachedHeadHandler{
		gitService: gitService,
	}
}

// Handle checks if HEAD is detached
func (h *DetachedHeadHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Debug("Checking for detached HEAD state")

	isDetached, err := h.gitService.IsDetachedHead(ctx)
	if err != nil {
		slog.Error("Failed to check HEAD state", "error", err)
		return fmt.Errorf("failed to check HEAD state: %w", err)
	}

	if isDetached {
		slog.Error("HEAD is in detached state")
		return fmt.Errorf("HEAD is detached - please checkout a branch first")
	}

	// Get current branch for context
	currentBranch, err := h.gitService.GetCurrentBranch(ctx)
	if err != nil {
		slog.Error("Failed to get current branch", "error", err)
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	branchCtx.CurrentBranch = currentBranch
	slog.Debug("HEAD check passed", "current_branch", currentBranch)

	return h.callNext(ctx, branchCtx)
}
