package handlers

import "context"

// GitService defines the interface for git operations needed by handlers
// This interface breaks the import cycle between git and git/handlers packages
type GitService interface {
	IsGitRepository(ctx context.Context) (bool, error)
	IsDetachedHead(ctx context.Context) (bool, error)
	GetCurrentBranch(ctx context.Context) (string, error)
	IsWorkingTreeClean(ctx context.Context) (bool, error)
	IsMainBehindOrigin(ctx context.Context) (bool, error)
	LocalBranchExists(ctx context.Context, branchName string) (bool, error)
	RemoteBranchExists(ctx context.Context, branchName string) (bool, error)
	CheckoutRemoteBranch(ctx context.Context, branchName string) error
	SwitchBranch(ctx context.Context, branchName string) error
	CreateBranch(ctx context.Context, branchName string) error
	ForceRecreateBranch(ctx context.Context, branchName string) error
	PushBranch(ctx context.Context, branchName string) error
}
