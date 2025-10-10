package github

import "fmt"

type GraphQLBuilder struct{}

func NewGraphQLBuilder() *GraphQLBuilder {
	return &GraphQLBuilder{}
}

func (b *GraphQLBuilder) BuildThreadsPageQuery(owner, name string) string {
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

func (b *GraphQLBuilder) BuildResolveThreadMutation() string {
	return "mutation($tid:ID!){ resolveReviewThread(input:{threadId:$tid}){ thread { isResolved } } }"
}

func (b *GraphQLBuilder) BuildReplyThreadMutation() string {
	return "mutation($tid:ID!, $body:String!){ " +
		"addPullRequestReviewThreadReply(input:{pullRequestReviewThreadId:$tid, body:$body})" +
		"{ comment { id } } }"
}
