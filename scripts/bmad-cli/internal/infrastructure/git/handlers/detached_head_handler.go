package handlers

import (
	"bmad-cli/internal/domain/ports"
	"context"
	"errors"
	"log/slog"

	pkgerrors "bmad-cli/internal/pkg/errors"
)

// DetachedHeadHandler prevents operations on detached HEAD state.
type DetachedHeadHandler struct {
	BaseBranchHandler

	gitService ports.GitPort
}

// NewDetachedHeadHandler creates a new detached HEAD handler.
func NewDetachedHeadHandler(gitService ports.GitPort) *DetachedHeadHandler {
	return &DetachedHeadHandler{
		gitService: gitService,
	}
}

// Handle checks if HEAD is detached.
func (h *DetachedHeadHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Debug("Checking for detached HEAD state")

	isDetached, err := h.gitService.IsDetachedHead(ctx)
	if err != nil {
		slog.Error("Failed to check HEAD state", "error", err)

		return pkgerrors.ErrCheckHEADStateFailed(err)
	}

	if isDetached {
		slog.Error("HEAD is in detached state")

		return errors.New("HEAD is detached - please checkout a branch first")
	}

	// Get current branch for context
	currentBranch, err := h.gitService.GetCurrentBranch(ctx)
	if err != nil {
		slog.Error("Failed to get current branch", "error", err)

		return pkgerrors.ErrGetCurrentBranchFailed(err)
	}

	branchCtx.CurrentBranch = currentBranch
	slog.Debug("HEAD check passed", "current_branch", currentBranch)

	return h.callNext(ctx, branchCtx)
}
