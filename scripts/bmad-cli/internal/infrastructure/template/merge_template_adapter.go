package template

import (
	"fmt"
	"strings"

	storyModels "bmad-cli/internal/domain/models/story"
)

// ScenarioMergeData represents data needed for scenario merge template processing
type ScenarioMergeData struct {
	StoryNumber        string
	ScenarioID         string
	Level              string
	Priority           string
	AcceptanceCriteria string // Pre-formatted for template
	Steps              string // Pre-formatted Gherkin steps
	RequirementsFile   string
}

// NewScenarioMergeData creates a new ScenarioMergeData instance from a test scenario
func NewScenarioMergeData(storyNumber string, scenario storyModels.TestScenario, outputFile string) *ScenarioMergeData {
	return &ScenarioMergeData{
		StoryNumber:        storyNumber,
		ScenarioID:         scenario.ID,
		Level:              scenario.Level,
		Priority:           scenario.Priority,
		AcceptanceCriteria: formatAcceptanceCriteria(scenario.AcceptanceCriteria),
		Steps:              formatScenarioSteps(scenario),
		RequirementsFile:   outputFile,
	}
}

// formatScenarioSteps formats scenario steps into Gherkin-style text
func formatScenarioSteps(scenario storyModels.TestScenario) string {
	var result strings.Builder

	for _, step := range scenario.Steps {
		if len(step.Given) > 0 {
			result.WriteString("  Given:\n")
			for _, g := range step.Given {
				result.WriteString(fmt.Sprintf("    - %s\n", g))
			}
		}
		if len(step.When) > 0 {
			result.WriteString("  When:\n")
			for _, w := range step.When {
				result.WriteString(fmt.Sprintf("    - %s\n", w))
			}
		}
		if len(step.Then) > 0 {
			result.WriteString("  Then:\n")
			for _, t := range step.Then {
				result.WriteString(fmt.Sprintf("    - %s\n", t))
			}
		}
	}

	return result.String()
}

// formatAcceptanceCriteria formats acceptance criteria array for template
func formatAcceptanceCriteria(criteria []string) string {
	if len(criteria) == 0 {
		return "[]"
	}
	return fmt.Sprintf(`["%s"]`, strings.Join(criteria, `", "`))
}
