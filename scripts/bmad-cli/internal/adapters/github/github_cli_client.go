package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/infrastructure/shell"
	"bmad-cli/internal/pkg/errors"
)

const expectedParts = 2

type GitHubCLIClient struct {
	shell *shell.CommandRunner
}

func NewGitHubCLIClient(shell *shell.CommandRunner) *GitHubCLIClient {
	return &GitHubCLIClient{shell: shell}
}

func (c *GitHubCLIClient) GetCurrentCheckoutPR(ctx context.Context) (string, error) {
	out, err := c.shell.Run(ctx, "gh", "pr", "view", "--json", "number", "-q", ".number")
	if err != nil {
		return "", fmt.Errorf("get current PR number: %w", err)
	}

	return out, nil
}

func (c *GitHubCLIClient) GetCurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git branch failed: %w", errors.ErrGitBranchFailed(err))
	}

	return strings.TrimSpace(string(out)), nil
}

func (c *GitHubCLIClient) ListPRsForBranch(ctx context.Context, branch string) ([]models.PullRequest, error) {
	out, err := c.shell.Run(ctx, "gh", "pr", "list", "--head", branch,
		"--json", "number,title,state,url", "--limit", "1")
	if err != nil {
		return nil, fmt.Errorf("gh pr list failed: %w", errors.ErrGHPRListFailed(err, out))
	}

	var prs []models.PullRequest

	err = json.Unmarshal([]byte(out), &prs)
	if err != nil {
		return nil, fmt.Errorf("parse PR list failed: %w", errors.ErrParsePRListFailed(err))
	}

	return prs, nil
}

func (c *GitHubCLIClient) GetRepoOwnerAndName(ctx context.Context) (string, string, error) {
	out, err := c.shell.Run(ctx, "gh", "repo", "view", "--json", "owner,name", "-q", ".owner.login + \" \" + .name")
	if err != nil {
		return "", "", fmt.Errorf("repo view failed: %w", errors.ErrRepoViewFailed(err))
	}

	parts := strings.Split(strings.TrimSpace(out), " ")
	if len(parts) != expectedParts {
		return "", "", fmt.Errorf("unexpected repo output: %w", errors.ErrUnexpectedRepoOutputWithDetails(out))
	}

	return parts[0], parts[1], nil
}

func (c *GitHubCLIClient) ExecuteGraphQL(
	ctx context.Context,
	query string,
	variables map[string]string,
) (string, error) {
	args := []string{"api", "graphql", "-f", "query=" + query}
	for key, value := range variables {
		args = append(args, "-F", key+"="+value)
	}

	out, err := c.shell.Run(ctx, "gh", args...)
	if err != nil {
		return "", fmt.Errorf("execute GraphQL query: %w", err)
	}

	return out, nil
}
