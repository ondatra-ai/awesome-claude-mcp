package prtriage

// Thread represents a PR review thread with its comments.
type Thread struct {
	ID         string    `json:"id"`
	IsResolved bool      `json:"isResolved"`
	Comments   []Comment `json:"comments"`
}

// Comment represents a single comment in a review thread.
type Comment struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	URL      string `json:"url"`
	Body     string `json:"body"`
	Outdated bool   `json:"outdated"`
}

// HeuristicAnalysisResult captures the heuristic assessment outcome from AI.
type HeuristicAnalysisResult struct {
	Score           int
	Summary         string
	ProposedActions []string
	Items           map[string]bool
	Alternatives    []map[string]string
}

// ThreadContext provides inputs for AI analysis/implementation.
type ThreadContext struct {
	PRNumber int
	Thread   Thread
	Comment  Comment
}
