package prtriage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	expectedParts = 2
)

var (
	ErrUnexpectedOutput = errors.New("unexpected output format")
)

// Error wrapping function for better error handling.
func ErrUnexpectedRepoOutput(output string) error {
	return fmt.Errorf("unexpected repo view output: %s: %w", output, ErrUnexpectedOutput)
}

type GitHubClient interface {
	// Raw GitHub CLI operations
	GetCurrentCheckoutPR(ctx context.Context) (string, error)
	GetCurrentBranch(ctx context.Context) (string, error)
	ListPRsForBranch(ctx context.Context, branch string) ([]PullRequest, error)
	GetRepoOwnerAndName(ctx context.Context) (owner, name string, err error)

	// Raw GraphQL operations
	ExecuteGraphQL(ctx context.Context, query string, variables map[string]string) (string, error)
}

type ghCLIClient struct{}

func NewGitHubCLIClient() GitHubClient { return &ghCLIClient{} }

// GetCurrentCheckoutPR gets the PR number for the current checkout.
func (c *ghCLIClient) GetCurrentCheckoutPR(ctx context.Context) (string, error) {
	return runShell(ctx, "gh", "pr", "view", "--json", "number", "-q", ".number")
}

// GetCurrentBranch gets the current git branch name.
func (c *ghCLIClient) GetCurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git branch: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// ListPRsForBranch lists PRs for the given branch.
func (c *ghCLIClient) ListPRsForBranch(ctx context.Context, branch string) ([]PullRequest, error) {
	out, err := runShell(ctx, "gh", "pr", "list", "--head", branch,
		"--json", "number,title,state,url", "--limit", "1")
	if err != nil {
		return nil, fmt.Errorf("gh pr list failed: %w, out=%s", err, out)
	}

	var prs []PullRequest

	err = json.Unmarshal([]byte(out), &prs)
	if err != nil {
		return nil, fmt.Errorf("parse pr list: %w", err)
	}

	return prs, nil
}

// GetRepoOwnerAndName gets the repository owner and name.
func (c *ghCLIClient) GetRepoOwnerAndName(ctx context.Context) (string, string, error) {
	out, err := runShell(ctx, "gh", "repo", "view", "--json", "owner,name", "-q", ".owner.login + \" \" + .name")
	if err != nil {
		return "", "", fmt.Errorf("repo view: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(out), " ")
	if len(parts) != expectedParts {
		return "", "", ErrUnexpectedRepoOutput(out)
	}

	return parts[0], parts[1], nil
}

// ExecuteGraphQL executes a GraphQL query with variables using GitHub CLI.
func (c *ghCLIClient) ExecuteGraphQL(ctx context.Context, query string, variables map[string]string) (string, error) {
	args := []string{"gh", "api", "graphql", "-f", "query=" + query}
	for key, value := range variables {
		args = append(args, "-F", key+"="+value)
	}

	return runShell(ctx, args[0], args[1:]...)
}

// PullRequest represents a GitHub pull request.
type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
}
