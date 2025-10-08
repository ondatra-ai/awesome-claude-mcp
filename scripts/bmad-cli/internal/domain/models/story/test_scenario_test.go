package story

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestStepStatementUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected StepStatement
	}{
		{
			name:     "plain string",
			yaml:     `"Server is ready"`,
			expected: StepStatement{Type: "", Statement: "Server is ready"},
		},
		{
			name:     "and modifier",
			yaml:     `and: "Authentication configured"`,
			expected: StepStatement{Type: ModifierTypeAnd, Statement: "Authentication configured"},
		},
		{
			name:     "but modifier",
			yaml:     `but: "No requests made"`,
			expected: StepStatement{Type: ModifierTypeBut, Statement: "No requests made"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stmt StepStatement
			err := yaml.Unmarshal([]byte(tt.yaml), &stmt)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if stmt.Type != tt.expected.Type {
				t.Errorf("Type = %v, want %v", stmt.Type, tt.expected.Type)
			}
			if stmt.Statement != tt.expected.Statement {
				t.Errorf("Statement = %v, want %v", stmt.Statement, tt.expected.Statement)
			}
		})
	}
}

func TestScenarioStepUnmarshal(t *testing.T) {
	yamlData := `
given:
  - "Server is ready to accept connections"
  - and: "Authentication is configured"
when:
  - "Client attempts to connect"
then:
  - "Server accepts connection"
  - and: "Welcome message sent"
`

	var step ScenarioStep
	err := yaml.Unmarshal([]byte(yamlData), &step)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// Check Given
	if len(step.Given) != 2 {
		t.Errorf("Given length = %d, want 2", len(step.Given))
	}
	if step.Given[0].Type != "" {
		t.Errorf("Given[0].Type = %v, want empty", step.Given[0].Type)
	}
	if step.Given[0].Statement != "Server is ready to accept connections" {
		t.Errorf("Given[0].Statement = %v", step.Given[0].Statement)
	}
	if step.Given[1].Type != ModifierTypeAnd {
		t.Errorf("Given[1].Type = %v, want and", step.Given[1].Type)
	}

	// Check When
	if len(step.When) != 1 {
		t.Errorf("When length = %d, want 1", len(step.When))
	}

	// Check Then
	if len(step.Then) != 2 {
		t.Errorf("Then length = %d, want 2", len(step.Then))
	}
}
