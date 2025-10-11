package handlers

import (
	"context"
	"log/slog"
	"strings"

	pkgerrors "bmad-cli/internal/pkg/errors"
)

// UnrelatedBranchHandler fails if on a branch unrelated to the story.
type UnrelatedBranchHandler struct {
	BaseBranchHandler
}

// NewUnrelatedBranchHandler creates a new handler.
func NewUnrelatedBranchHandler() *UnrelatedBranchHandler {
	return &UnrelatedBranchHandler{}
}

// Handle checks if on an unrelated branch (not main, not the story branch).
func (h *UnrelatedBranchHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	// If on main, continue to creation logic
	if branchCtx.CurrentBranch == "main" {
		slog.Debug("On main branch, will check for existing branches")

		return h.callNext(ctx, branchCtx)
	}

	// If on a branch that starts with the story number, it's related
	if strings.HasPrefix(branchCtx.CurrentBranch, branchCtx.StoryNumber+"-") {
		slog.Debug("On related story branch", "branch", branchCtx.CurrentBranch)

		return h.callNext(ctx, branchCtx)
	}

	// Otherwise, it's an unrelated branch
	slog.Error("On unrelated branch",
		"current", branchCtx.CurrentBranch,
		"story", branchCtx.StoryNumber)

	return pkgerrors.ErrUnrelatedBranchError(branchCtx.CurrentBranch, branchCtx.StoryNumber)
}
