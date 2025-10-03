package story

// TestScenario represents a single BDD test scenario with Given-When-Then format
type TestScenario struct {
	ID                   string   `yaml:"id" json:"id"`
	AcceptanceCriteria   []string `yaml:"acceptance_criteria" json:"acceptance_criteria"`
	Given                string   `yaml:"given" json:"given"`
	When                 string   `yaml:"when" json:"when"`
	Then                 string   `yaml:"then" json:"then"`
	Level                string   `yaml:"level" json:"level"`
	Priority             string   `yaml:"priority" json:"priority"`
	MitigatesRisks       []string `yaml:"mitigates_risks,omitempty" json:"mitigates_risks,omitempty"`
}

// Scenarios contains all test scenarios for a story
type Scenarios struct {
	TestScenarios []TestScenario `yaml:"test_scenarios" json:"test_scenarios"`
}
