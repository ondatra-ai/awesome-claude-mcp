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

// NewBranchManager creates a new branch manager with the complete handler chain
func NewBranchManager(gitService *GitService) *BranchManager {
	// Create all handlers
	gitRepoCheck := NewGitRepoCheckHandler(gitService)
	dirtyTreeCheck := NewDirtyWorkingTreeHandler(gitService)
	detachedHeadCheck := NewDetachedHeadHandler(gitService)
	forceRecreate := NewForceRecreateHandler(gitService)
	alreadyOnBranch := NewAlreadyOnStoryBranchHandler()
	unrelatedBranch := NewUnrelatedBranchHandler()
	mainBehindOrigin := NewMainBehindOriginHandler(gitService)
	localBranchExists := NewLocalBranchExistsHandler(gitService)
	remoteBranchExists := NewRemoteBranchExistsHandler(gitService)
	createBranch := NewCreateBranchHandler(gitService)

	// Build the chain
	gitRepoCheck.
		SetNext(dirtyTreeCheck).
		SetNext(detachedHeadCheck).
		SetNext(forceRecreate).
		SetNext(alreadyOnBranch).
		SetNext(unrelatedBranch).
		SetNext(mainBehindOrigin).
		SetNext(localBranchExists).
		SetNext(remoteBranchExists).
		SetNext(createBranch)

	return &BranchManager{
		chain: gitRepoCheck,
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
