package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"bmad-cli/internal/common/errors"
)

type PRNumberFetcher struct {
	client *GitHubCLIClient
}

func NewPRNumberFetcher(client *GitHubCLIClient) *PRNumberFetcher {
	return &PRNumberFetcher{client: client}
}

func (p *PRNumberFetcher) Fetch(ctx context.Context) (int, error) {
	prNumStr, err := p.client.GetCurrentCheckoutPR(ctx)
	if err == nil {
		prNum, convErr := strconv.Atoi(strings.TrimSpace(prNumStr))
		if convErr == nil {
			return prNum, nil
		}
	}

	branch, err := p.client.GetCurrentBranch(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get current branch: %w", err)
	}

	prs, err := p.client.ListPRsForBranch(ctx, branch)
	if err != nil {
		return 0, fmt.Errorf("failed to list PRs for branch %s: %w", branch, err)
	}

	if len(prs) == 0 {
		return 0, errors.ErrNoPRFoundForBranch(branch)
	}

	return prs[0].Number, nil
}
