package main

import (
	"context"
	"fmt"
)

// ThreadResolver handles resolving PR review threads.
type ThreadResolver struct {
	gh GitHubClient
}

// NewThreadResolver creates a new ThreadResolver.
func NewThreadResolver(gh GitHubClient) *ThreadResolver {
	return &ThreadResolver{gh: gh}
}

// Resolve resolves a thread with the given message.
func (t *ThreadResolver) Resolve(ctx context.Context, threadID string, message string) error {
	return t.resolveReply(ctx, threadID, message, true)
}

// resolveReply handles both resolving and replying to threads.
func (t *ThreadResolver) resolveReply(ctx context.Context, threadID, body string, resolve bool) error {
	if resolve {
		// Resolve the review thread via GraphQL
		query := "mutation($tid:ID!){ resolveReviewThread(input:{threadId:$tid}){ thread { isResolved } } }"
		variables := map[string]string{"tid": threadID}

		out, err := t.gh.ExecuteGraphQL(ctx, query, variables)
		if err != nil {
			return fmt.Errorf("resolve thread: %w, out=%s", err, out)
		}

		return nil
	}

	query := "mutation($tid:ID!, $body:String!){ " +
		"addPullRequestReviewThreadReply(input:{pullRequestReviewThreadId:$tid, body:$body})" +
		"{ comment { id } } }"
	variables := map[string]string{"tid": threadID, "body": body}

	out, err := t.gh.ExecuteGraphQL(ctx, query, variables)
	if err != nil {
		return fmt.Errorf("reply thread: %w, out=%s", err, out)
	}

	return nil
}
