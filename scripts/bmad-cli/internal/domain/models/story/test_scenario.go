package story

// ScenarioStep represents a single step in a Gherkin scenario
// Each step should have exactly one keyword set (Given, When, Then, And, or But)
type ScenarioStep struct {
	Given string `yaml:"given,omitempty" json:"given,omitempty"`
	When  string `yaml:"when,omitempty" json:"when,omitempty"`
	Then  string `yaml:"then,omitempty" json:"then,omitempty"`
	And   string `yaml:"and,omitempty" json:"and,omitempty"`
	But   string `yaml:"but,omitempty" json:"but,omitempty"`
}

// TestScenario represents a single BDD test scenario with Gherkin format
// Supports standard scenarios and scenario outlines for data-driven testing
type TestScenario struct {
	ID                 string                   `yaml:"id" json:"id"`
	AcceptanceCriteria []string                 `yaml:"acceptance_criteria" json:"acceptance_criteria"`
	Steps              []ScenarioStep           `yaml:"steps" json:"steps"`
	ScenarioOutline    bool                     `yaml:"scenario_outline,omitempty" json:"scenario_outline,omitempty"`
	Examples           []map[string]interface{} `yaml:"examples,omitempty" json:"examples,omitempty"`
	Level              string                   `yaml:"level" json:"level"`
	Priority           string                   `yaml:"priority" json:"priority"`
	MitigatesRisks     []string                 `yaml:"mitigates_risks,omitempty" json:"mitigates_risks,omitempty"`
}

// Scenarios contains all test scenarios for a story
type Scenarios struct {
	TestScenarios []TestScenario `yaml:"test_scenarios" json:"test_scenarios"`
}
