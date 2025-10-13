package models

// Thread represents a GitHub PR review thread.
// Uses camelCase to match GitHub's API contract.
type Thread struct {
	ID         string    `json:"id"`
	IsResolved bool      `json:"isResolved"`
	Comments   []Comment `json:"comments"`
}
