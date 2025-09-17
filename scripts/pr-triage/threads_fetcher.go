package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrNoEligibleReviewThreads = errors.New("no eligible review threads found")
)

// ErrNoEligibleThreads returns an error when no eligible review threads are found.
func ErrNoEligibleThreads(prNumber int) error {
	return fmt.Errorf("no eligible review threads found for PR %d: %w", prNumber, ErrNoEligibleReviewThreads)
}

// ThreadsFetcher handles fetching all review threads for a PR.
type ThreadsFetcher struct {
	gh GitHubClient
}

// NewThreadsFetcher creates a new ThreadsFetcher.
func NewThreadsFetcher(gh GitHubClient) *ThreadsFetcher {
	return &ThreadsFetcher{gh: gh}
}

// FetchAll retrieves all review threads for the given PR number.
func (t *ThreadsFetcher) FetchAll(ctx context.Context, prNumber int) ([]Thread, error) {
	owner, name, err := t.gh.GetRepoOwnerAndName(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository owner and name: %w", err)
	}

	all, err := t.fetchAllPages(ctx, owner, name, prNumber)
	if err != nil {
		return nil, err
	}

	eligible := filterEligibleThreads(all)
	if len(eligible) == 0 {
		return nil, ErrNoEligibleThreads(prNumber)
	}

	return eligible, nil
}

// fetchAllPages retrieves all pages of review threads.
func (t *ThreadsFetcher) fetchAllPages(ctx context.Context, owner, name string, prNumber int) ([]Thread, error) {
	var all []Thread

	after := ""

	for {
		page, err := t.fetchSinglePage(ctx, owner, name, prNumber, after)
		if err != nil {
			return nil, err
		}

		all = append(all, convertPageToThreads(page)...)

		pageInfo := page.Data.Repository.PullRequest.ReviewThreads.PageInfo
		if !pageInfo.HasNextPage || pageInfo.EndCursor == nil {
			break
		}

		after = *pageInfo.EndCursor
	}

	return all, nil
}

// fetchSinglePage retrieves a single page of review threads.
func (t *ThreadsFetcher) fetchSinglePage(
	ctx context.Context,
	owner, name string,
	prNumber int,
	after string,
) (threadsPageResponse, error) {
	query := buildThreadsPageQuery(owner, name)

	variables := map[string]string{
		"prNumber": strconv.Itoa(prNumber),
	}
	if after != "" {
		variables["after"] = after
	}

	out, err := t.gh.ExecuteGraphQL(ctx, query, variables)
	if err != nil {
		return threadsPageResponse{}, fmt.Errorf("graphql list threads: %w, out=%s", err, out)
	}

	var page threadsPageResponse

	err = json.Unmarshal([]byte(out), &page)
	if err != nil {
		return threadsPageResponse{}, fmt.Errorf("parse threads page: %w", err)
	}

	return page, nil
}

// convertPageToThreads converts a page response to Thread objects.
func convertPageToThreads(page threadsPageResponse) []Thread {
	nodes := page.Data.Repository.PullRequest.ReviewThreads.Nodes
	threads := make([]Thread, 0, len(nodes))

	for _, threadNode := range nodes {
		comments := convertNodesToComments(threadNode.Comments.Nodes)
		threads = append(threads, Thread{
			ID:         threadNode.ID,
			IsResolved: threadNode.IsResolved,
			Comments:   comments,
		})
	}

	return threads
}

// convertNodesToComments converts comment nodes to Comment objects.
func convertNodesToComments(nodes []commentNode) []Comment {
	comments := make([]Comment, 0, len(nodes))

	for _, commentNode := range nodes {
		file := ""
		line := 0

		if commentNode.Path != nil {
			file = *commentNode.Path
		}

		if commentNode.Line != nil {
			line = *commentNode.Line
		}

		comments = append(comments, Comment{
			File:     file,
			Line:     line,
			URL:      commentNode.URL,
			Body:     commentNode.Body,
			Outdated: commentNode.Outdated,
		})
	}

	return comments
}

// filterEligibleThreads returns only unresolved threads with comments.
func filterEligibleThreads(all []Thread) []Thread {
	var eligible []Thread

	for _, thread := range all {
		if !thread.IsResolved && len(thread.Comments) > 0 {
			eligible = append(eligible, thread)
		}
	}

	return eligible
}

// helpers for mapping GraphQL to local types

// commentNode represents a GraphQL comment node.
type commentNode struct {
	Path     *string `json:"path"`
	Line     *int    `json:"line"`
	Body     string  `json:"body"`
	Outdated bool    `json:"outdated"`
	URL      string  `json:"url"`
}

// pageInfo supports pagination cursors.
type pageInfo struct {
	HasNextPage bool    `json:"hasNextPage"`
	EndCursor   *string `json:"endCursor"`
}

// threadsPageResponse is a minimal GraphQL response shape for paginating threads.
type threadsPageResponse struct {
	Data struct {
		Repository struct {
			PullRequest struct {
				ReviewThreads struct {
					PageInfo pageInfo `json:"pageInfo"`
					Nodes    []struct {
						ID         string `json:"id"`
						IsResolved bool   `json:"isResolved"`
						Comments   struct {
							Nodes []commentNode `json:"nodes"`
						} `json:"comments"`
					} `json:"nodes"`
				} `json:"reviewThreads"`
			} `json:"pullRequest"`
		} `json:"repository"`
	} `json:"data"`
}

// buildThreadsPageQuery returns a query that supports pagination via $after.
func buildThreadsPageQuery(owner, name string) string {
	return fmt.Sprintf(`
query($prNumber: Int!, $after: String) {
  repository(owner: "%s", name: "%s") {
    pullRequest(number: $prNumber) {
      reviewThreads(first: 100, after: $after) {
        pageInfo { hasNextPage endCursor }
        nodes {
          id
          isResolved
          comments(first: 50) {
            nodes { path line body outdated url }
          }
        }
      }
    }
  }
}`, owner, name)
}
