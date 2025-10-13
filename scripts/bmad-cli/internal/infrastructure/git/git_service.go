package git

import (
	"context"
	"log/slog"
	"strings"

	"bmad-cli/internal/infrastructure/shell"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// GitService provides git operations.
type GitService struct {
	shellExec *shell.CommandRunner
}

// NewGitService creates a new git service.
func NewGitService(shellExec *shell.CommandRunner) *GitService {
	return &GitService{
		shellExec: shellExec,
	}
}

// IsGitRepository checks if the current directory is a git repository.
func (s *GitService) IsGitRepository(ctx context.Context) (bool, error) {
	slog.Debug("Checking if directory is a git repository")

	_, err := s.shellExec.Run(ctx, "git", "rev-parse", "--git-dir")
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetCurrentBranch returns the name of the current branch.
func (s *GitService) GetCurrentBranch(ctx context.Context) (string, error) {
	slog.Debug("Getting current git branch")

	output, err := s.shellExec.Run(ctx, "git", "branch", "--show-current")
	if err != nil {
		return "", pkgerrors.ErrGetCurrentBranchFailed(err)
	}

	branch := strings.TrimSpace(output)
	slog.Debug("Current branch", "branch", branch)

	return branch, nil
}

// IsWorkingTreeClean checks if the working tree has no uncommitted changes.
func (s *GitService) IsWorkingTreeClean(ctx context.Context) (bool, error) {
	slog.Debug("Checking if working tree is clean")

	output, err := s.shellExec.Run(ctx, "git", "status", "--porcelain")
	if err != nil {
		return false, pkgerrors.ErrCheckWorkingTreeStatusFailed(err)
	}

	clean := strings.TrimSpace(output) == ""
	slog.Debug("Working tree status", "clean", clean)

	return clean, nil
}

// IsDetachedHead checks if HEAD is detached.
func (s *GitService) IsDetachedHead(ctx context.Context) (bool, error) {
	slog.Debug("Checking if HEAD is detached")

	branch, err := s.GetCurrentBranch(ctx)
	if err != nil {
		return false, err
	}

	detached := branch == ""
	slog.Debug("HEAD detachment status", "detached", detached)

	return detached, nil
}

// LocalBranchExists checks if a branch exists locally.
func (s *GitService) LocalBranchExists(ctx context.Context, branch string) (bool, error) {
	slog.Debug("Checking if local branch exists", "branch", branch)

	_, err := s.shellExec.Run(ctx, "git", "rev-parse", "--verify", branch)
	if err != nil {
		return false, err
	}

	slog.Debug("Local branch exists", "branch", branch)

	return true, nil
}

// RemoteBranchExists checks if a branch exists on the remote.
func (s *GitService) RemoteBranchExists(ctx context.Context, branch string) (bool, error) {
	slog.Debug("Checking if remote branch exists", "branch", branch)

	output, err := s.shellExec.Run(ctx, "git", "ls-remote", "--heads", "origin", branch)
	if err != nil {
		return false, err
	}

	exists := strings.TrimSpace(output) != ""
	slog.Debug("Remote branch check result", "branch", branch, "exists", exists)

	return exists, nil
}

// IsMainBehindOrigin checks if main branch is behind origin/main.
func (s *GitService) IsMainBehindOrigin(ctx context.Context) (bool, error) {
	slog.Debug("Checking if main is behind origin/main")

	// Fetch latest from origin
	_, err := s.shellExec.Run(ctx, "git", "fetch", "origin", "main")
	if err != nil {
		return false, pkgerrors.ErrFetchOriginMainFailed(err)
	}

	// Check if main is behind origin/main
	output, err := s.shellExec.Run(ctx, "git", "rev-list", "--count", "main..origin/main")
	if err != nil {
		return false, pkgerrors.ErrCompareMainWithOriginFailed(err)
	}

	behind := strings.TrimSpace(output) != "0"
	slog.Debug("Main branch status", "behind", behind)

	return behind, nil
}

// SwitchBranch switches to an existing branch.
func (s *GitService) SwitchBranch(ctx context.Context, branch string) error {
	slog.Info("Switching to branch", "branch", branch)

	_, err := s.shellExec.Run(ctx, "git", "switch", branch)
	if err != nil {
		return pkgerrors.ErrSwitchBranchFailed(branch, err)
	}

	slog.Info("Successfully switched to branch", "branch", branch)

	return nil
}

// CheckoutRemoteBranch checks out a branch from remote.
func (s *GitService) CheckoutRemoteBranch(ctx context.Context, branch string) error {
	slog.Info("Checking out remote branch", "branch", branch)

	_, err := s.shellExec.Run(ctx, "git", "checkout", "-b", branch, "origin/"+branch)
	if err != nil {
		return pkgerrors.ErrCheckoutRemoteBranchFailed(branch, err)
	}

	slog.Info("Successfully checked out remote branch", "branch", branch)

	return nil
}

// CreateBranch creates a new branch from the current HEAD and pushes it to remote.
func (s *GitService) CreateBranch(ctx context.Context, branch string) error {
	slog.Info("Creating new branch", "branch", branch)

	_, err := s.shellExec.Run(ctx, "git", "switch", "-c", branch)
	if err != nil {
		return pkgerrors.ErrCreateBranchFailed(branch, err)
	}

	slog.Info("Successfully created branch", "branch", branch)

	// Push to remote immediately
	err = s.PushBranch(ctx, branch)
	if err != nil {
		return err
	}

	return nil
}

// PushBranch pushes a branch to remote with tracking.
func (s *GitService) PushBranch(ctx context.Context, branch string) error {
	slog.Info("Pushing branch to remote", "branch", branch)

	_, err := s.shellExec.Run(ctx, "git", "push", "-u", "origin", branch)
	if err != nil {
		return pkgerrors.ErrPushBranchFailed(branch, err)
	}

	slog.Info("Successfully pushed branch to remote", "branch", branch)

	return nil
}

// ForceRecreateBranch deletes and recreates a branch.
func (s *GitService) ForceRecreateBranch(ctx context.Context, branch string) error {
	slog.Info("Force recreating branch", "branch", branch)

	// Switch to main first
	err := s.SwitchBranch(ctx, "main")
	if err != nil {
		return pkgerrors.ErrSwitchToMainFailed(err)
	}

	// Delete local branch if exists
	localExists, err := s.LocalBranchExists(ctx, branch)
	if err != nil {
		return err
	}

	if localExists {
		slog.Debug("Deleting local branch", "branch", branch)

		_, err := s.shellExec.Run(ctx, "git", "branch", "-D", branch)
		if err != nil {
			return pkgerrors.ErrDeleteLocalBranchFailed(branch, err)
		}
	}

	// Delete remote branch if exists
	remoteExists, err := s.RemoteBranchExists(ctx, branch)
	if err != nil {
		return err
	}

	if remoteExists {
		slog.Debug("Deleting remote branch", "branch", branch)

		_, err := s.shellExec.Run(ctx, "git", "push", "origin", "--delete", branch)
		if err != nil {
			return pkgerrors.ErrDeleteRemoteBranchFailed(branch, err)
		}
	}

	// Create new branch
	err = s.CreateBranch(ctx, branch)
	if err != nil {
		return err
	}

	slog.Info("Successfully force recreated branch", "branch", branch)

	return nil
}
