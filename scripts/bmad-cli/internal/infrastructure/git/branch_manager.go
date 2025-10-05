package git

import (
	"context"
	"fmt"
	"log/slog"
)

// BranchManager manages the branch creation/switching logic using Chain of Responsibility
type BranchManager struct {
	chain BranchHandler
}

// HandlerFactory creates a handler instance
type HandlerFactory func(*GitService) BranchHandler

// handlerFactories defines the handler chain order
var handlerFactories = []HandlerFactory{
	func(gs *GitService) BranchHandler { return NewGitRepoCheckHandler(gs) },
	func(gs *GitService) BranchHandler { return NewDirtyWorkingTreeHandler(gs) },
	func(gs *GitService) BranchHandler { return NewDetachedHeadHandler(gs) },
	func(gs *GitService) BranchHandler { return NewForceRecreateHandler(gs) },
	func(gs *GitService) BranchHandler { return NewAlreadyOnStoryBranchHandler() },
	func(gs *GitService) BranchHandler { return NewUnrelatedBranchHandler() },
	func(gs *GitService) BranchHandler { return NewMainBehindOriginHandler(gs) },
	func(gs *GitService) BranchHandler { return NewLocalBranchExistsHandler(gs) },
	func(gs *GitService) BranchHandler { return NewRemoteBranchExistsHandler(gs) },
	func(gs *GitService) BranchHandler { return NewCreateBranchHandler(gs) },
}

// NewBranchManager creates a new branch manager with the complete handler chain
func NewBranchManager(gitService *GitService) *BranchManager {
	// Create handlers from factories
	var handlers []BranchHandler
	for _, factory := range handlerFactories {
		handlers = append(handlers, factory(gitService))
	}

	// Chain handlers automatically
	for i := 0; i < len(handlers)-1; i++ {
		handlers[i].SetNext(handlers[i+1])
	}

	return &BranchManager{
		chain: handlers[0],
	}
}

// EnsureBranch ensures the correct branch is checked out for the story
func (m *BranchManager) EnsureBranch(ctx context.Context, storyNumber, storySlug string, force bool) error {
	expectedBranch := fmt.Sprintf("%s-%s", storyNumber, storySlug)

	slog.Info("Ensuring branch for story",
		"story_number", storyNumber,
		"story_slug", storySlug,
		"expected_branch", expectedBranch,
		"force", force)

	branchCtx := NewBranchContext(storyNumber, force)
	branchCtx.ExpectedBranch = expectedBranch

	if err := m.chain.Handle(ctx, branchCtx); err != nil {
		slog.Error("Failed to ensure branch", "error", err)
		return err
	}

	slog.Info("Branch operation completed successfully",
		"action", branchCtx.Action,
		"branch", expectedBranch)

	return nil
}
