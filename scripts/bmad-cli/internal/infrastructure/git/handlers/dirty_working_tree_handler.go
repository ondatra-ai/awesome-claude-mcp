package handlers

import (
	"context"
	"log/slog"

	"bmad-cli/internal/domain/ports"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// DirtyWorkingTreeHandler checks for uncommitted changes.
type DirtyWorkingTreeHandler struct {
	BaseBranchHandler

	gitService ports.GitPort
}

// NewDirtyWorkingTreeHandler creates a new dirty working tree handler.
func NewDirtyWorkingTreeHandler(gitService ports.GitPort) *DirtyWorkingTreeHandler {
	return &DirtyWorkingTreeHandler{
		gitService: gitService,
	}
}

// Handle checks if the working tree has uncommitted changes.
func (h *DirtyWorkingTreeHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Debug("Checking working tree status")

	isClean, err := h.gitService.IsWorkingTreeClean(ctx)
	if err != nil {
		slog.Error("Failed to check working tree status", "error", err)

		return pkgerrors.ErrCheckWorkingTreeStatusFailed(err)
	}

	if !isClean {
		slog.Error("Working tree has uncommitted changes")

		return pkgerrors.ErrWorkingTreeDirty
	}

	slog.Debug("Working tree is clean")

	return h.callNext(ctx, branchCtx)
}
