package git

import (
	"context"
	"fmt"
	"log/slog"

	"bmad-cli/internal/infrastructure/git/handlers"
)

// BranchManager manages the branch creation/switching logic using Chain of Responsibility.
type BranchManager struct {
	chain handlers.BranchHandler
}

// HandlerFactory creates a handler instance.
type HandlerFactory func(*GitService) handlers.BranchHandler

// getHandlerFactories returns the handler chain order.
func getHandlerFactories() []HandlerFactory {
	return []HandlerFactory{
		func(gs *GitService) handlers.BranchHandler { return handlers.NewGitRepoCheckHandler(gs) },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewDirtyWorkingTreeHandler(gs) },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewDetachedHeadHandler(gs) },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewForceRecreateHandler(gs) },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewAlreadyOnStoryBranchHandler() },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewUnrelatedBranchHandler() },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewMainBehindOriginHandler(gs) },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewLocalBranchExistsHandler(gs) },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewRemoteBranchExistsHandler(gs) },
		func(gs *GitService) handlers.BranchHandler { return handlers.NewCreateBranchHandler(gs) },
	}
}

// NewBranchManager creates a new branch manager with the complete handler chain.
func NewBranchManager(gitService *GitService) *BranchManager {
	// Create handlers from factories
	factories := getHandlerFactories()
	handlerList := make([]handlers.BranchHandler, 0, len(factories))

	for _, factory := range factories {
		handlerList = append(handlerList, factory(gitService))
	}

	// Chain handlers automatically
	for i := range len(handlerList) - 1 {
		handlerList[i].SetNext(handlerList[i+1])
	}

	return &BranchManager{
		chain: handlerList[0],
	}
}

// EnsureBranch ensures the correct branch is checked out for the story.
func (m *BranchManager) EnsureBranch(ctx context.Context, storyNumber, storySlug string, force bool) error {
	expectedBranch := fmt.Sprintf("%s-%s", storyNumber, storySlug)

	slog.Info("Ensuring branch for story",
		"story_number", storyNumber,
		"story_slug", storySlug,
		"expected_branch", expectedBranch,
		"force", force)

	branchCtx := handlers.NewBranchContext(storyNumber, force)
	branchCtx.ExpectedBranch = expectedBranch

	err := m.chain.Handle(ctx, branchCtx)
	if err != nil {
		slog.Error("Failed to ensure branch", "error", err)

		return fmt.Errorf("handle branch chain for %s: %w", expectedBranch, err)
	}

	slog.Info("Branch operation completed successfully",
		"action", branchCtx.Action,
		"branch", expectedBranch)

	return nil
}
