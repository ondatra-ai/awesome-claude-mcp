package template

import (
	"fmt"
	"strings"
)

// TestImplementationData represents data needed for test implementation template processing.
type TestImplementationData struct {
	ScenarioID       string      // e.g., "INT-011"
	Description      string      // Short description of test
	Level            string      // "integration" or "e2e"
	Category         string      // backend, frontend, performance
	Priority         string      // P0, P1, P2, P3
	MergedSteps      MergedSteps // Given-When-Then steps
	RequirementsFile string      // Path to requirements file to update
}

// MergedSteps represents the Given-When-Then structure from requirements.yml.
type MergedSteps struct {
	Given []string
	When  []string
	Then  []string
}

// NewTestImplementationData creates a new TestImplementationData instance from requirements.yml entry.
func NewTestImplementationData(
	scenarioID string,
	description string,
	level string,
	category string,
	priority string,
	given []string,
	when []string,
	then []string,
	requirementsFile string,
) *TestImplementationData {
	return &TestImplementationData{
		ScenarioID:  scenarioID,
		Description: description,
		Level:       level,
		Category:    category,
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
func (d *TestImplementationData) FormatSteps() string {
	var result strings.Builder

	if len(d.MergedSteps.Given) > 0 {
		result.WriteString("Given:\n")

		for _, step := range d.MergedSteps.Given {
			result.WriteString(fmt.Sprintf("  - %s\n", step))
		}
	}

	if len(d.MergedSteps.When) > 0 {
		result.WriteString("When:\n")

		for _, step := range d.MergedSteps.When {
			result.WriteString(fmt.Sprintf("  - %s\n", step))
		}
	}

	if len(d.MergedSteps.Then) > 0 {
		result.WriteString("Then:\n")

		for _, step := range d.MergedSteps.Then {
			result.WriteString(fmt.Sprintf("  - %s\n", step))
		}
	}

	return result.String()
}
