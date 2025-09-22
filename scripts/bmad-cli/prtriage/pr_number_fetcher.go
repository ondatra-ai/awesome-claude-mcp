package prtriage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNoPRFound = errors.New("no PR found")
)

// ErrNoPRFoundForBranch returns an error when no PR is found for the given branch.
func ErrNoPRFoundForBranch(branch string) error {
	return fmt.Errorf("no PR found for branch: %s: %w", branch, ErrNoPRFound)
}

// PRNumberFetcher handles fetching the current PR number.
type PRNumberFetcher struct {
	gh GitHubClient
}

// NewPRNumberFetcher creates a new PRNumberFetcher.
func NewPRNumberFetcher(gh GitHubClient) *PRNumberFetcher {
	return &PRNumberFetcher{gh: gh}
}

// Fetch retrieves the current PR number using primary and fallback strategies.
func (p *PRNumberFetcher) Fetch(ctx context.Context) (int, error) {
	// Primary strategy: try to get PR from current checkout
	prNumStr, err := p.gh.GetCurrentCheckoutPR(ctx)
	if err == nil {
		prNum, convErr := strconv.Atoi(strings.TrimSpace(prNumStr))
		if convErr == nil {
			return prNum, nil
		}
	}

	// Fallback strategy: detect current branch and list PRs for it
	branch, err := p.gh.GetCurrentBranch(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current branch: %w", err)
	}

	prs, err := p.gh.ListPRsForBranch(ctx, branch)
	if err != nil {
		return 0, fmt.Errorf("failed to list PRs for branch %s: %w", branch, err)
	}

	if len(prs) == 0 {
		return 0, ErrNoPRFoundForBranch(branch)
	}

	return prs[0].Number, nil
}
