package handlers

import (
	"bmad-cli/internal/domain/ports"
	"context"
)

// LocalBranchExistsHandler checks if the branch exists locally and switches to it.
type LocalBranchExistsHandler struct {
	BaseBranchHandler

	gitService ports.GitPort
}

// NewLocalBranchExistsHandler creates a new handler.
func NewLocalBranchExistsHandler(gitService ports.GitPort) *LocalBranchExistsHandler {
	return &LocalBranchExistsHandler{
		gitService: gitService,
	}
}

// Handle checks if branch exists locally and switches to it.
func (h *LocalBranchExistsHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	return h.handleBranchExistence(
		ctx,
		branchCtx,
		"local",
		h.gitService.LocalBranchExists,
		h.gitService.SwitchBranch,
		ActionSwitch,
		"switch",
	)
}
