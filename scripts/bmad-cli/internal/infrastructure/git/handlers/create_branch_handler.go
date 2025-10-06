package handlers

import (
	"context"
	"log/slog"

	"bmad-cli/internal/domain/ports"
)

// CreateBranchHandler creates a new branch as the final handler in the chain
type CreateBranchHandler struct {
	BaseBranchHandler
	gitService ports.GitPort
}

// NewCreateBranchHandler creates a new handler
func NewCreateBranchHandler(gitService ports.GitPort) *CreateBranchHandler {
	return &CreateBranchHandler{
		gitService: gitService,
	}
}

// Handle creates a new branch
func (h *CreateBranchHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	slog.Info("Creating new branch", "branch", branchCtx.ExpectedBranch)

	if err := h.gitService.CreateBranch(ctx, branchCtx.ExpectedBranch); err != nil {
		slog.Error("Failed to create branch", "error", err)
		return err
	}

	branchCtx.Action = ActionCreate
	slog.Info("Successfully created new branch", "branch", branchCtx.ExpectedBranch)
	return nil
}
