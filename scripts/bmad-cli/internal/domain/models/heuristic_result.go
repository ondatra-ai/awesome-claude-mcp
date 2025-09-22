package models

type HeuristicAnalysisResult struct {
	Score           int
	Summary         string
	ProposedActions []string
	Items           map[string]bool
	Alternatives    []map[string]string
}
