package handlers

// BranchAction represents the action to be taken on the branch.
type BranchAction int

const (
	ActionNone BranchAction = iota
	ActionCreate
	ActionSwitch
	ActionCheckout
	ActionForceRecreate
)
