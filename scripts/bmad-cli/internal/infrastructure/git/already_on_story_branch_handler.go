package git

import (
	"context"
	"log/slog"
)

// AlreadyOnStoryBranchHandler checks if already on the expected story branch
type AlreadyOnStoryBranchHandler struct {
	BaseBranchHandler
}

// NewAlreadyOnStoryBranchHandler creates a new handler
func NewAlreadyOnStoryBranchHandler() *AlreadyOnStoryBranchHandler {
	return &AlreadyOnStoryBranchHandler{}
}

// Handle checks if already on the correct branch
func (h *AlreadyOnStoryBranchHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	if branchCtx.CurrentBranch == branchCtx.ExpectedBranch {
		slog.Info("Already on the expected story branch", "branch", branchCtx.ExpectedBranch)
		branchCtx.Action = ActionNone
		return nil
	}

	slog.Debug("Not on expected branch, continuing chain",
		"current", branchCtx.CurrentBranch,
		"expected", branchCtx.ExpectedBranch)
	return h.callNext(ctx, branchCtx)
}
