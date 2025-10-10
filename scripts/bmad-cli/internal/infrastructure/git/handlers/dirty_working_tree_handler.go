package handlers

import (
	"bmad-cli/internal/domain/ports"
	"context"
	"fmt"
	"log/slog"
)

// DirtyWorkingTreeHandler checks for uncommitted changes
type DirtyWorkingTreeHandler struct {
	BaseBranchHandler
	gitService ports.GitPort
}

// NewDirtyWorkingTreeHandler creates a new dirty working tree handler
func NewDirtyWorkingTreeHandler(gitService ports.GitPort) *DirtyWorkingTreeHandler {
	return &DirtyWorkingTreeHandler{
		gitService: gitService,
	}
}

// Handle checks if the working tree has uncommitted changes
func (h *DirtyWorkingTreeHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Debug("Checking working tree status")

	isClean, err := h.gitService.IsWorkingTreeClean(ctx)
	if err != nil {
		slog.Error("Failed to check working tree status", "error", err)
		return fmt.Errorf("failed to check working tree status: %w", err)
	}

	if !isClean {
		slog.Error("Working tree has uncommitted changes")
		return fmt.Errorf("working tree has uncommitted changes - please commit or stash them first")
	}

	slog.Debug("Working tree is clean")
	return h.callNext(ctx, branchCtx)
}
