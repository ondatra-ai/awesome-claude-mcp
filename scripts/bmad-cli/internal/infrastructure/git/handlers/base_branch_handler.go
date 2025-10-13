package handlers

import (
	"context"
	"log/slog"
)

// BaseBranchHandler provides the base implementation for all handlers.
type BaseBranchHandler struct {
	next BranchHandler
}

// SetNext sets the next handler in the chain.
func (h *BaseBranchHandler) SetNext(handler BranchHandler) BranchHandler {
	h.next = handler

	return handler
}

// callNext calls the next handler in the chain if it exists.
func (h *BaseBranchHandler) callNext(ctx context.Context, branchCtx *BranchContext) error {
	if h.next == nil {
		return nil
	}

	return h.next.Handle(ctx, branchCtx)
}

// handleBranchExistence provides common logic for checking branch existence and taking action.
func (h *BaseBranchHandler) handleBranchExistence(
	ctx context.Context,
	branchCtx *BranchContext,
	locationName string,
	checkExists func(context.Context, string) (bool, error),
	takeAction func(context.Context, string) error,
	action BranchAction,
	actionVerb string,
) error {
	slog.Debug("Checking if "+locationName+" branch exists", "branch", branchCtx.ExpectedBranch)

	exists, err := checkExists(ctx, branchCtx.ExpectedBranch)
	if err != nil {
		slog.Error("Failed to check "+locationName+" branch", "error", err)

		return err
	}

	if !exists {
		slog.Debug(locationName + " branch does not exist, continuing chain")

		return h.callNext(ctx, branchCtx)
	}

	slog.Info(locationName+" branch exists, "+actionVerb+" it", "branch", branchCtx.ExpectedBranch)

	err = takeAction(ctx, branchCtx.ExpectedBranch)
	if err != nil {
		slog.Error("Failed to "+actionVerb+" "+locationName+" branch", "error", err)

		return err
	}

	branchCtx.Action = action
	slog.Info("Successfully "+actionVerb+"ed "+locationName+" branch", "branch", branchCtx.ExpectedBranch)

	return nil
}
