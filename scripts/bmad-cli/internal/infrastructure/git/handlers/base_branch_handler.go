package handlers

import "context"

// BaseBranchHandler provides the base implementation for all handlers
type BaseBranchHandler struct {
	next BranchHandler
}

// SetNext sets the next handler in the chain
func (h *BaseBranchHandler) SetNext(handler BranchHandler) BranchHandler {
	h.next = handler
	return handler
}

// callNext calls the next handler in the chain if it exists
func (h *BaseBranchHandler) callNext(ctx context.Context, branchCtx *BranchContext) error {
	if h.next == nil {
		return nil
	}
	return h.next.Handle(ctx, branchCtx)
}
