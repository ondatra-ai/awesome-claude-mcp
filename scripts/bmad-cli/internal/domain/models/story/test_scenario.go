package story

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// ModifierType represents the type of statement modifier
type ModifierType string

const (
	ModifierTypeAnd ModifierType = "and"
	ModifierTypeBut ModifierType = "but"
)

// StepStatement represents a single statement in a Gherkin step
// Can be either a main statement (plain string) or a modifier (and/but)
type StepStatement struct {
	Type      ModifierType // Empty for main statement, "and" or "but" for modifiers
	Statement string
}

// UnmarshalYAML implements custom YAML unmarshaling
// Handles both plain strings and objects with and/but keys
func (s *StepStatement) UnmarshalYAML(node *yaml.Node) error {
	// Handle plain string (main statement)
	if node.Kind == yaml.ScalarNode {
		s.Type = ""
		s.Statement = node.Value
		return nil
	}

	// Handle object with and/but key
	if node.Kind == yaml.MappingNode {
		if len(node.Content) != 2 {
			return fmt.Errorf("modifier must have exactly one key")
		}

		key := node.Content[0].Value
		value := node.Content[1].Value

		switch key {
		case "and":
			s.Type = ModifierTypeAnd
			s.Statement = value
		case "but":
			s.Type = ModifierTypeBut
			s.Statement = value
		default:
			return fmt.Errorf("invalid modifier type: %s (must be 'and' or 'but')", key)
		}

		return nil
	}

	return fmt.Errorf("invalid step statement format")
}

// MarshalYAML implements custom YAML marshaling
// Outputs plain string for main statement, object for modifiers
func (s StepStatement) MarshalYAML() (interface{}, error) {
	// Main statement: output as plain string
	if s.Type == "" {
		return s.Statement, nil
	}

	// Modifier: output as object with and/but key
	return map[string]string{
		string(s.Type): s.Statement,
	}, nil
}

// ScenarioStep represents a single step in a Gherkin scenario
// Each Given/When/Then is an array of statements
type ScenarioStep struct {
	Given []StepStatement `yaml:"given,omitempty" json:"given,omitempty"`
	When  []StepStatement `yaml:"when,omitempty" json:"when,omitempty"`
	Then  []StepStatement `yaml:"then,omitempty" json:"then,omitempty"`
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
