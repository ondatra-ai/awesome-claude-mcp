package git

// BranchAction represents the action to be taken on the branch
type BranchAction int

const (
	ActionNone BranchAction = iota
	ActionCreate
	ActionSwitch
	ActionCheckout
	ActionForceRecreate
)

// BranchContext holds the state and configuration for branch operations
type BranchContext struct {
	StoryNumber    string
	ExpectedBranch string
	CurrentBranch  string
	Force          bool
	Action         BranchAction
}

// NewBranchContext creates a new branch context with the given story number and force flag
func NewBranchContext(storyNumber string, force bool) *BranchContext {
	return &BranchContext{
		StoryNumber: storyNumber,
		Force:       force,
		Action:      ActionNone,
	}
}
