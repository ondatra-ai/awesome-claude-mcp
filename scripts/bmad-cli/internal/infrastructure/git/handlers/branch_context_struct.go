package handlers

// BranchContext holds the state and configuration for branch operations.
type BranchContext struct {
	StoryNumber    string
	ExpectedBranch string
	CurrentBranch  string
	Force          bool
	Action         BranchAction
}
