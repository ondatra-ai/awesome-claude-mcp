package git

import "context"

// BranchHandler defines the interface for handlers in the Chain of Responsibility
type BranchHandler interface {
	SetNext(handler BranchHandler) BranchHandler
	Handle(ctx context.Context, branchCtx *BranchContext) error
}
