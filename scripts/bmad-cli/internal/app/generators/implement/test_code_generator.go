package implement

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/app/generators/validate"
	"bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/input"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"

	"gopkg.in/yaml.v3"
)

// TestCodeGenerator generates test code for pending scenarios using Claude.
type TestCodeGenerator struct {
	claudeClient       *ai.ClaudeClient
	config             *config.ViperConfig
	userInputCollector *input.UserInputCollector
}

// NewTestCodeGenerator creates a new TestCodeGenerator.
func NewTestCodeGenerator(
	claudeClient *ai.ClaudeClient,
	config *config.ViperConfig,
	userInputCollector *input.UserInputCollector,
) *TestCodeGenerator {
	return &TestCodeGenerator{
		claudeClient:       claudeClient,
		config:             config,
		userInputCollector: userInputCollector,
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

	// Load architecture content for validation
	archPath := g.config.GetString("documents.architecture_yaml")

	archContent, err := g.readArchitectureContent(archPath)
	if err != nil {
		slog.Warn("⚠️  Could not load architecture.yaml, skipping validation", "error", err)
	}

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

		if !g.implementSingleTest(ctx, scenario, userPromptLoader, systemPromptLoader, tmpDir) {
			continue
		}

		// Validate generated test against architecture.yaml
		if archContent != "" {
			changed := g.validateAndResolve(ctx, scenario, &archContent, archPath, tmpDir)
			if changed {
				slog.Info("🔄 Regenerating test with enriched architecture", "scenario_id", scenario.ScenarioID)
				g.implementSingleTest(ctx, scenario, userPromptLoader, systemPromptLoader, tmpDir)
			}
		}

		implementedCount++

		slog.Info("✓ Test implemented successfully",
			"scenario_id", scenario.ScenarioID,
			"duration", time.Since(startTime).Round(time.Second),
		)
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

// validateAndResolve validates a generated test and resolves issues with user input.
// Returns true if architecture was updated (caller should regenerate test).
func (g *TestCodeGenerator) validateAndResolve(
	ctx context.Context,
	scenario *template.TestGenerationData,
	archContent *string,
	archPath string,
	tmpDir string,
) bool {
	result := g.validateGeneratedTest(ctx, scenario, *archContent, tmpDir)
	if result == nil || !result.HasIssues() {
		slog.Info("✓ Test validation passed", "scenario_id", scenario.ScenarioID)

		return false
	}

	slog.Info("⚠️  Test validation found issues",
		"scenario_id", scenario.ScenarioID,
		"issue_count", len(result.Issues),
	)

	return g.resolveIssues(ctx, result.Issues, archContent, archPath, tmpDir)
}

// validateGeneratedTest validates a generated test against architecture.yaml.
func (g *TestCodeGenerator) validateGeneratedTest(
	ctx context.Context,
	scenario *template.TestGenerationData,
	archContent string,
	tmpDir string,
) *TestValidationOutput {
	validationUserPath := g.config.GetString("templates.prompts.validate_generated_test")
	validationSystemPath := g.config.GetString("templates.prompts.validate_generated_test_system")

	validationUserLoader := template.NewTemplateLoader[*template.TestGenerationData](validationUserPath)
	validationSystemLoader := template.NewTemplateLoader[*template.TestGenerationData](validationSystemPath)

	resultPath := filepath.Join(tmpDir, scenario.ScenarioID+"-validation-result.yaml")

	// Set validation context on scenario
	scenario.ArchitectureContent = archContent
	scenario.ResultPath = resultPath

	userPrompt, err := validationUserLoader.LoadTemplate(scenario)
	if err != nil {
		slog.Warn("⚠️  Failed to load validation user prompt", "error", err)

		return nil
	}

	systemPrompt, err := validationSystemLoader.LoadTemplate(scenario)
	if err != nil {
		slog.Warn("⚠️  Failed to load validation system prompt", "error", err)

		return nil
	}

	g.savePromptFile(tmpDir, scenario.ScenarioID+"-validation-user-prompt.txt", userPrompt)
	g.savePromptFile(tmpDir, scenario.ScenarioID+"-validation-system-prompt.txt", systemPrompt)

	slog.Info("🔍 Validating generated test against architecture", "scenario_id", scenario.ScenarioID)

	response, err := g.claudeClient.ExecutePromptWithSystem(
		ctx,
		systemPrompt,
		userPrompt,
		"sonnet",
		ai.ExecutionMode{AllowedTools: []string{"Read"}},
	)

	if response != "" {
		g.savePromptFile(tmpDir, scenario.ScenarioID+"-validation-response.txt", response)
	}

	if err != nil {
		slog.Warn("⚠️  Validation call failed", "scenario_id", scenario.ScenarioID, "error", err)

		return nil
	}

	// Parse validation result
	content := validate.ExtractFileContent(response, resultPath)
	if content == "" {
		slog.Warn("⚠️  No validation content found in response", "scenario_id", scenario.ScenarioID)

		return nil
	}

	content = validate.StripMarkdownCodeFences(content)

	var output TestValidationOutput

	err = yaml.Unmarshal([]byte(content), &output)
	if err != nil {
		slog.Warn("⚠️  Failed to parse validation output", "error", err)

		return nil
	}

	return &output
}

// resolveIssues presents validation issues to the user and applies architecture updates.
// Returns true if any architecture changes were made.
func (g *TestCodeGenerator) resolveIssues(
	ctx context.Context,
	issues []TestValidationIssue,
	archContent *string,
	archPath string,
	tmpDir string,
) bool {
	changed := false

	for _, issue := range issues {
		question := checklist.ClarifyQuestion{
			ID:       fmt.Sprintf("arch-%s-%s", issue.IssueType, issue.Name),
			Question: fmt.Sprintf("Test references undefined %s: %s", issue.IssueType, issue.Name),
			Context:  fmt.Sprintf("Found in %s. How should architecture.yaml be updated?", issue.TestFile),
			Options:  issue.ProposedUpdates,
		}

		answers := g.userInputCollector.AskQuestions([]checklist.ClarifyQuestion{question})

		chosenOption, ok := answers[question.ID]
		if !ok || chosenOption == "" {
			slog.Info("User skipped issue", "issue", issue.Name)

			continue
		}

		err := g.applyArchitectureUpdate(ctx, archContent, issue, chosenOption, archPath, tmpDir)
		if err != nil {
			slog.Warn("⚠️  Failed to apply architecture update",
				"issue", issue.Name,
				"error", err,
			)

			continue
		}

		changed = true
	}

	return changed
}

// applyArchitectureUpdate uses Claude to apply a chosen update to architecture.yaml.
func (g *TestCodeGenerator) applyArchitectureUpdate(
	ctx context.Context,
	archContent *string,
	issue TestValidationIssue,
	chosenOption string,
	archPath string,
	tmpDir string,
) error {
	updateUserPath := g.config.GetString("templates.prompts.apply_arch_update")
	updateSystemPath := g.config.GetString("templates.prompts.apply_arch_update_system")

	updateUserLoader := template.NewTemplateLoader[*template.TestArchUpdateData](updateUserPath)
	updateSystemLoader := template.NewTemplateLoader[*template.TestArchUpdateData](updateSystemPath)

	resultPath := filepath.Join(tmpDir, fmt.Sprintf("arch-update-%s-%s-result.yaml", issue.IssueType, issue.Name))

	promptData := &template.TestArchUpdateData{
		ArchitectureContent: *archContent,
		IssueType:           issue.IssueType,
		IssueName:           issue.Name,
		TestFile:            issue.TestFile,
		ChosenOption:        chosenOption,
		ResultPath:          resultPath,
	}

	userPrompt, err := updateUserLoader.LoadTemplate(promptData)
	if err != nil {
		return fmt.Errorf("load arch update user prompt: %w", err)
	}

	systemPrompt, err := updateSystemLoader.LoadTemplate(promptData)
	if err != nil {
		return fmt.Errorf("load arch update system prompt: %w", err)
	}

	g.savePromptFile(tmpDir, fmt.Sprintf("arch-update-%s-%s-user.txt", issue.IssueType, issue.Name), userPrompt)
	g.savePromptFile(tmpDir, fmt.Sprintf("arch-update-%s-%s-system.txt", issue.IssueType, issue.Name), systemPrompt)

	slog.Info("🔧 Applying architecture update", "issue_type", issue.IssueType, "name", issue.Name)

	response, err := g.claudeClient.ExecutePromptWithSystem(
		ctx,
		systemPrompt,
		userPrompt,
		"sonnet",
		ai.ExecutionMode{},
	)

	if response != "" {
		g.savePromptFile(tmpDir, fmt.Sprintf("arch-update-%s-%s-response.txt", issue.IssueType, issue.Name), response)
	}

	if err != nil {
		return fmt.Errorf("claude arch update call: %w", err)
	}

	updatedContent := validate.ExtractFileContent(response, resultPath)
	if updatedContent == "" {
		return pkgerrors.ErrArchUpdateNoContent
	}

	// Write updated architecture.yaml to disk
	err = os.WriteFile(archPath, []byte(updatedContent), fileModeReadWrite)
	if err != nil {
		return fmt.Errorf("write architecture.yaml: %w", err)
	}

	// Update in-memory content
	*archContent = updatedContent

	slog.Info("✓ Architecture updated", "path", archPath, "issue", issue.Name)

	return nil
}

// readArchitectureContent reads the architecture.yaml file.
func (g *TestCodeGenerator) readArchitectureContent(archPath string) (string, error) {
	data, err := os.ReadFile(archPath)
	if err != nil {
		return "", fmt.Errorf("read architecture.yaml: %w", err)
	}

	return string(data), nil
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
