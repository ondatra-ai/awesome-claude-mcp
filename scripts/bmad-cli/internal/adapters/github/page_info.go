package github

// pageInfo represents GitHub GraphQL API pagination info.
// Uses camelCase to match GitHub's API contract.
type pageInfo struct {
	HasNextPage bool    `json:"hasNextPage"`
	EndCursor   *string `json:"endCursor"`
}
