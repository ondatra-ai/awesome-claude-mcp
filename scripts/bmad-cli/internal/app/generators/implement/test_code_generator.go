package implement

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"

	"gopkg.in/yaml.v3"
)

// TestCodeGenerator generates test code for pending scenarios using Claude.
type TestCodeGenerator struct {
	claudeClient *ai.ClaudeClient
	config       *config.ViperConfig
}

// NewTestCodeGenerator creates a new TestCodeGenerator.
func NewTestCodeGenerator(
	claudeClient *ai.ClaudeClient,
	config *config.ViperConfig,
) *TestCodeGenerator {
	return &TestCodeGenerator{
		claudeClient: claudeClient,
		config:       config,
	}
}

// GenerateTests generates test code for all pending scenarios in the requirements file.
func (g *TestCodeGenerator) GenerateTests(
	ctx context.Context,
	requirementsFile string,
	tmpDir string,
) (GenerationStatus, error) {
	slog.Info("⚙️  Starting test generation", "requirements_file", requirementsFile)

	// Parse requirements file to find pending scenarios
	pendingScenarios, err := g.parsePendingScenarios(requirementsFile)
	if err != nil {
		return NewFailureStatus("parse pending scenarios failed"),
			pkgerrors.ErrParsePendingScenariosFailed(err)
	}

	if len(pendingScenarios) == 0 {
		slog.Info("✓ No pending scenarios to implement")

		return NewSuccessStatus(0, nil, "No pending scenarios"), nil
	}

	slog.Info("Found pending scenarios to implement", "count", len(pendingScenarios))

	// Create template loaders
	userPromptPath := g.config.GetString("templates.prompts.generate_tests")
	systemPromptPath := g.config.GetString("templates.prompts.generate_tests_system")
	userPromptLoader := template.NewTemplateLoader[*template.TestGenerationData](userPromptPath)
	systemPromptLoader := template.NewTemplateLoader[*template.TestGenerationData](systemPromptPath)

	// Process each pending scenario
	implementedCount := 0

	var filesModified []string

	for i, scenario := range pendingScenarios {
		startTime := time.Now()

		slog.Info("Implementing test scenario",
			"progress", i+1,
			"total", len(pendingScenarios),
			"scenario_id", scenario.ScenarioID,
		)

		if g.implementSingleTest(ctx, scenario, userPromptLoader, systemPromptLoader, tmpDir) {
			implementedCount++

			slog.Info("✓ Test implemented successfully",
				"scenario_id", scenario.ScenarioID,
				"duration", time.Since(startTime).Round(time.Second),
			)
		}
	}

	slog.Info("✅ Test generation completed",
		"implemented_count", implementedCount,
		"total_pending", len(pendingScenarios),
	)

	return NewSuccessStatus(
		implementedCount,
		filesModified,
		"Test generation completed",
	), nil
}

// parsePendingScenarios reads requirements file and extracts scenarios with status: "pending".
func (g *TestCodeGenerator) parsePendingScenarios(
	requirementsFile string,
) ([]*template.TestGenerationData, error) {
	slog.Debug("Parsing requirements file", "file", requirementsFile)

	data, err := os.ReadFile(requirementsFile)
	if err != nil {
		return nil, pkgerrors.ErrReadRequirementsFailed(err)
	}

	// Parse YAML structure
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
		return nil, pkgerrors.ErrUnmarshalRequirementsFailed(err)
	}

	// Filter pending scenarios
	pendingScenarios := make([]*template.TestGenerationData, 0, len(requirements.Scenarios))

	for scenarioID, scenario := range requirements.Scenarios {
		if scenario.ImplementationStatus.Status != "pending" {
			slog.Debug("Skipping non-pending scenario",
				"scenario_id", scenarioID,
				"status", scenario.ImplementationStatus.Status,
			)

			continue
		}

		givenSteps := g.convertStepsToStrings(scenario.MergedSteps.Given)
		whenSteps := g.convertStepsToStrings(scenario.MergedSteps.When)
		thenSteps := g.convertStepsToStrings(scenario.MergedSteps.Then)

		testData := template.NewTestGenerationData(
			scenarioID,
			scenario.Description,
			scenario.Level,
			scenario.Service,
			scenario.Priority,
			givenSteps,
			whenSteps,
			thenSteps,
			requirementsFile,
		)

		pendingScenarios = append(pendingScenarios, testData)

		slog.Debug("Found pending scenario", "scenario_id", scenarioID)
	}

	slog.Info("Parsed requirements file",
		"total_scenarios", len(requirements.Scenarios),
		"pending_count", len(pendingScenarios),
	)

	return pendingScenarios, nil
}

// convertStepsToStrings converts []interface{} to []string.
func (g *TestCodeGenerator) convertStepsToStrings(steps []interface{}) []string {
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

// implementSingleTest implements a single test scenario using Claude.
func (g *TestCodeGenerator) implementSingleTest(
	ctx context.Context,
	scenario *template.TestGenerationData,
	userLoader *template.TemplateLoader[*template.TestGenerationData],
	systemLoader *template.TemplateLoader[*template.TestGenerationData],
	tmpDir string,
) bool {
	userPrompt, err := userLoader.LoadTemplate(scenario)
	if err != nil {
		slog.Warn("⚠️  Skipping scenario: failed to load user prompt",
			"scenario_id", scenario.ScenarioID,
			"error", err,
		)

		return false
	}

	// Save user prompt
	g.savePromptFile(tmpDir, scenario.ScenarioID+"-test-generation-user-prompt.txt", userPrompt)

	systemPrompt, err := systemLoader.LoadTemplate(scenario)
	if err != nil {
		slog.Warn("⚠️  Skipping scenario: failed to load system prompt",
			"scenario_id", scenario.ScenarioID,
			"error", err,
		)

		return false
	}

	// Save system prompt
	g.savePromptFile(tmpDir, scenario.ScenarioID+"-test-generation-system-prompt.txt", systemPrompt)

	slog.Info("🤖 Calling Claude for test generation", "scenario_id", scenario.ScenarioID)

	response, err := g.claudeClient.ExecutePromptWithSystem(
		ctx,
		systemPrompt,
		userPrompt,
		"sonnet",
		ai.ExecutionMode{AllowedTools: []string{"Read", "Write", "Edit"}},
	)

	// Save response
	if response != "" {
		g.savePromptFile(tmpDir, scenario.ScenarioID+"-test-generation-response.txt", response)
	}

	if err != nil {
		slog.Warn("⚠️  Failed to implement test scenario",
			"scenario_id", scenario.ScenarioID,
			"error", err,
		)

		return false
	}

	return true
}

// savePromptFile saves content to a file in the tmp directory.
func (g *TestCodeGenerator) savePromptFile(tmpDir, filename, content string) {
	filePath := filepath.Join(tmpDir, filename)

	err := os.WriteFile(filePath, []byte(content), fileModeReadWrite)
	if err != nil {
		slog.Warn("Failed to save file", "file", filePath, "error", err)
	} else {
		slog.Info("💾 File saved", "file", filePath)
	}
}
