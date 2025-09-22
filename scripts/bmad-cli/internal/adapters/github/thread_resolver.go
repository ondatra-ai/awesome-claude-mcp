package github

import (
	"context"
	"fmt"
)

type ThreadResolver struct {
	client  *GitHubCLIClient
	builder *GraphQLBuilder
}

func NewThreadResolver(client *GitHubCLIClient) *ThreadResolver {
	return &ThreadResolver{
		client:  client,
		builder: NewGraphQLBuilder(),
	}
}

func (t *ThreadResolver) Resolve(ctx context.Context, threadID string, message string) error {
	return t.resolveReply(ctx, threadID, message, true)
}

func (t *ThreadResolver) resolveReply(ctx context.Context, threadID, body string, resolve bool) error {
	if resolve {
		query := t.builder.BuildResolveThreadMutation()
		variables := map[string]string{"tid": threadID}

		out, err := t.client.ExecuteGraphQL(ctx, query, variables)
		if err != nil {
			return fmt.Errorf("resolve thread: %w, out=%s", err, out)
		}

		return nil
	}

	query := t.builder.BuildReplyThreadMutation()
	variables := map[string]string{"tid": threadID, "body": body}

	out, err := t.client.ExecuteGraphQL(ctx, query, variables)
	if err != nil {
		return fmt.Errorf("reply thread: %w, out=%s", err, out)
	}

	return nil
}
