package handlers

import (
	"bmad-cli/internal/domain/ports"
	"context"
)

// RemoteBranchExistsHandler checks if the branch exists on remote and checks it out.
type RemoteBranchExistsHandler struct {
	BaseBranchHandler

	gitService ports.GitPort
}

// NewRemoteBranchExistsHandler creates a new handler.
func NewRemoteBranchExistsHandler(gitService ports.GitPort) *RemoteBranchExistsHandler {
	return &RemoteBranchExistsHandler{
		gitService: gitService,
	}
}

// Handle checks if branch exists on remote and checks it out.
func (h *RemoteBranchExistsHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	return h.handleBranchExistence(
		ctx,
		branchCtx,
		"remote",
		h.gitService.RemoteBranchExists,
		h.gitService.CheckoutRemoteBranch,
		ActionCheckout,
		"checkout",
	)
}
