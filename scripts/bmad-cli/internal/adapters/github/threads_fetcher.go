package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/pkg/errors"
)

type ThreadsFetcher struct {
	client  *GitHubCLIClient
	builder *GraphQLBuilder
}

func NewThreadsFetcher(client *GitHubCLIClient) *ThreadsFetcher {
	return &ThreadsFetcher{
		client:  client,
		builder: NewGraphQLBuilder(),
	}
}

func (t *ThreadsFetcher) FetchAll(ctx context.Context, prNumber int) ([]models.Thread, error) {
	owner, name, err := t.client.GetRepoOwnerAndName(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository owner and name: %w", err)
	}

	all, err := t.fetchAllPages(ctx, owner, name, prNumber)
	if err != nil {
		return nil, err
	}

	eligible := t.filterEligibleThreads(all)
	if len(eligible) == 0 {
		return nil, errors.ErrNoEligibleThreads(prNumber)
	}

	return eligible, nil
}

func (t *ThreadsFetcher) fetchAllPages(ctx context.Context, owner, name string, prNumber int) ([]models.Thread, error) {
	var all []models.Thread

	after := ""

	for {
		page, err := t.fetchSinglePage(ctx, owner, name, prNumber, after)
		if err != nil {
			return nil, err
		}

		all = append(all, t.convertPageToThreads(page)...)

		pageInfo := page.Data.Repository.PullRequest.ReviewThreads.PageInfo
		if !pageInfo.HasNextPage || pageInfo.EndCursor == nil {
			break
		}

		after = *pageInfo.EndCursor
	}

	return all, nil
}

func (t *ThreadsFetcher) fetchSinglePage(
	ctx context.Context,
	owner, name string,
	prNumber int,
	after string,
) (threadsPageResponse, error) {
	query := t.builder.BuildThreadsPageQuery(owner, name)

	variables := map[string]string{
		"prNumber": strconv.Itoa(prNumber),
	}
	if after != "" {
		variables["after"] = after
	}

	out, err := t.client.ExecuteGraphQL(ctx, query, variables)
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

func (t *ThreadsFetcher) convertPageToThreads(page threadsPageResponse) []models.Thread {
	nodes := page.Data.Repository.PullRequest.ReviewThreads.Nodes
	threads := make([]models.Thread, 0, len(nodes))

	for _, threadNode := range nodes {
		comments := t.convertNodesToComments(threadNode.Comments.Nodes)
		threads = append(threads, models.Thread{
			ID:         threadNode.ID,
			IsResolved: threadNode.IsResolved,
			Comments:   comments,
		})
	}

	return threads
}

func (t *ThreadsFetcher) convertNodesToComments(nodes []commentNode) []models.Comment {
	comments := make([]models.Comment, 0, len(nodes))

	for _, commentNode := range nodes {
		file := ""
		line := 0

		if commentNode.Path != nil {
			file = *commentNode.Path
		}

		if commentNode.Line != nil {
			line = *commentNode.Line
		}

		comments = append(comments, models.Comment{
			File:     file,
			Line:     line,
			URL:      commentNode.URL,
			Body:     commentNode.Body,
			Outdated: commentNode.Outdated,
		})
	}

	return comments
}

func (t *ThreadsFetcher) filterEligibleThreads(all []models.Thread) []models.Thread {
	var eligible []models.Thread

	for _, thread := range all {
		if !thread.IsResolved && len(thread.Comments) > 0 {
			eligible = append(eligible, thread)
		}
	}

	return eligible
}
