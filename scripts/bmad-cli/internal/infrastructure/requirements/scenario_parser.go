package requirements

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/infrastructure/template"
)

// ScenarioParser parses requirements YAML files to extract test scenarios.
type ScenarioParser struct{}

// NewScenarioParser creates a new ScenarioParser.
func NewScenarioParser() *ScenarioParser {
	return &ScenarioParser{}
}

// ParseScenarios reads a requirements file and extracts scenarios.
// If pendingOnly is true, only scenarios with status "pending" are returned.
func (p *ScenarioParser) ParseScenarios(
	filePath string,
	pendingOnly bool,
) ([]*template.TestGenerationData, error) {
	slog.Debug("Parsing requirements file", "file", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read requirements file: %w", err)
	}

	var requirements struct {
		Scenarios map[string]struct {
			Description          string `yaml:"description"`
			Service              string `yaml:"service"`
			Level                string `yaml:"level"`
			Priority             string `yaml:"priority"`
			ImplementationStatus struct {
				Status   string `yaml:"status"`
				FilePath string `yaml:"file_path"`
			} `yaml:"implementation_status"`
			MergedSteps struct {
				Given []interface{} `yaml:"given"`
				When  []interface{} `yaml:"when"`
				Then  []interface{} `yaml:"then"`
			} `yaml:"merged_steps"`
		} `yaml:"scenarios"`
	}

	err = yaml.Unmarshal(data, &requirements)
	if err != nil {
		return nil, fmt.Errorf("failed to parse requirements YAML: %w", err)
	}

	scenarios := make([]*template.TestGenerationData, 0, len(requirements.Scenarios))

	for scenarioID, scenario := range requirements.Scenarios {
		if pendingOnly && scenario.ImplementationStatus.Status != "pending" {
			slog.Debug("Skipping non-pending scenario",
				"scenario_id", scenarioID,
				"status", scenario.ImplementationStatus.Status,
			)

			continue
		}

		givenSteps := convertStepsToStrings(scenario.MergedSteps.Given)
		whenSteps := convertStepsToStrings(scenario.MergedSteps.When)
		thenSteps := convertStepsToStrings(scenario.MergedSteps.Then)

		testData := template.NewTestGenerationData(
			scenarioID,
			scenario.Description,
			scenario.Level,
			scenario.Service,
			scenario.Priority,
			givenSteps,
			whenSteps,
			thenSteps,
			filePath,
		)

		scenarios = append(scenarios, testData)

		slog.Debug("Found scenario", "scenario_id", scenarioID)
	}

	slog.Info("Parsed requirements file",
		"total_scenarios", len(requirements.Scenarios),
		"extracted_count", len(scenarios),
	)

	return scenarios, nil
}

// convertStepsToStrings converts []interface{} to []string.
func convertStepsToStrings(steps []interface{}) []string {
	result := make([]string, 0, len(steps))

	for _, step := range steps {
		switch v := step.(type) {
		case string:
			result = append(result, v)
		case map[string]interface{}:
			for keyword, value := range v {
				if strValue, ok := value.(string); ok {
					result = append(result, keyword+" "+strValue)
				}
			}
		}
	}

	return result
}
