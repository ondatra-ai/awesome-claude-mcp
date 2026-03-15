package template

import (
	"fmt"
	"strings"
)

// TestGenerationData represents data needed for test generation template processing.
type TestGenerationData struct {
	ScenarioID          string      // e.g., "INT-011"
	Description         string      // Short description of test
	Level               string      // "integration" or "e2e"
	Service             string      // backend, frontend, mcp-service
	Priority            string      // P0, P1, P2, P3
	MergedSteps         MergedSteps // Given-When-Then steps
	RequirementsFile    string      // Path to requirements file to update
	ArchitectureContent string      // Current architecture.yaml content for validation
	ResultPath          string      // Path for FILE_START/FILE_END output markers
	TestFilePath        string      // Path where test file lives (e.g., "tests/e2e/E2E-021.spec.ts")
}

// TestArchUpdateData represents data needed for architecture update prompt templates.
type TestArchUpdateData struct {
	ArchitectureContent string // Current architecture.yaml content
	IssueType           string // e.g., "env_var", "helper_module"
	IssueName           string // e.g., "VIEWER_DOC_ID"
	TestFile            string // Where the issue was found
	ChosenOption        string // The user's chosen update option
	ResultPath          string // Path for FILE_START/FILE_END output markers
}

// MergedSteps represents the Given-When-Then structure from requirements.yaml.
type MergedSteps struct {
	Given []string
	When  []string
	Then  []string
}

// NewTestGenerationData creates a new TestGenerationData instance from requirements.yaml entry.
func NewTestGenerationData(
	scenarioID string,
	description string,
	level string,
	service string,
	priority string,
	given []string,
	when []string,
	then []string,
	requirementsFile string,
) *TestGenerationData {
	return &TestGenerationData{
		ScenarioID:  scenarioID,
		Description: description,
		Level:       level,
		Service:     service,
		Priority:    priority,
		MergedSteps: MergedSteps{
			Given: given,
			When:  when,
			Then:  then,
		},
		RequirementsFile: requirementsFile,
	}
}

// FormatSteps formats the Given-When-Then steps for display in template.
func (d *TestGenerationData) FormatSteps() string {
	var result strings.Builder

	if len(d.MergedSteps.Given) > 0 {
		result.WriteString("Given:\n")

		for _, step := range d.MergedSteps.Given {
			fmt.Fprintf(&result, "  - %s\n", step)
		}
	}

	if len(d.MergedSteps.When) > 0 {
		result.WriteString("When:\n")

		for _, step := range d.MergedSteps.When {
			fmt.Fprintf(&result, "  - %s\n", step)
		}
	}

	if len(d.MergedSteps.Then) > 0 {
		result.WriteString("Then:\n")

		for _, step := range d.MergedSteps.Then {
			fmt.Fprintf(&result, "  - %s\n", step)
		}
	}

	return result.String()
}
