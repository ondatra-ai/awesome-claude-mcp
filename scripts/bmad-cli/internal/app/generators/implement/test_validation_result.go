package implement

// TestValidationIssue represents a single issue found during test validation.
type TestValidationIssue struct {
	IssueType       string   `yaml:"issue_type"`       // "env_var", "helper_module", "import", "fixture"
	Name            string   `yaml:"name"`             // e.g., "VIEWER_DOC_ID"
	TestFile        string   `yaml:"test_file"`        // Where found
	ProposedUpdates []string `yaml:"proposed_updates"` // Options for architecture.yaml updates
}

// TestValidationOutput represents the result of validating a generated test.
type TestValidationOutput struct {
	Status string                `yaml:"status"` // "pass" or "issues_found"
	Issues []TestValidationIssue `yaml:"issues"`
}

// HasIssues returns true if validation found issues.
func (o *TestValidationOutput) HasIssues() bool {
	return o.Status == "issues_found" && len(o.Issues) > 0
}
