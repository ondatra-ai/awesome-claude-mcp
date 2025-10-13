package commands

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"bmad-cli/internal/adapters/ai"
	storyModels "bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/story"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"

	"gopkg.in/yaml.v3"
)

type USImplementCommand struct {
	branchManager *git.BranchManager
	storyLoader   *story.StoryLoader
	claudeClient  *ai.ClaudeClient
	config        *config.ViperConfig
	runDir        *fs.RunDirectory
}

func NewUSImplementCommand(
	branchManager *git.BranchManager,
	storyLoader *story.StoryLoader,
	claudeClient *ai.ClaudeClient,
	cfg *config.ViperConfig,
	runDir *fs.RunDirectory,
) *USImplementCommand {
	return &USImplementCommand{
		branchManager: branchManager,
		storyLoader:   storyLoader,
		claudeClient:  claudeClient,
		config:        cfg,
		runDir:        runDir,
	}
}

func (c *USImplementCommand) Execute(ctx context.Context, storyNumber string, force bool) error {
	slog.Info("Starting user story implementation", "story_number", storyNumber, "force", force)

	// Get story slug from file
	_, err := c.storyLoader.GetStorySlug(storyNumber)
	if err != nil {
		return pkgerrors.ErrGetStorySlugFailed(err)
	}

	// TEMPORARILY COMMENTED OUT FOR TESTING
	// Ensure correct branch is checked out
	// if err := c.branchManager.EnsureBranch(ctx, storyNumber, storySlug, force); err != nil {
	// 	return fmt.Errorf("failed to ensure branch: %w", err)
	// }
	// slog.Info("Branch setup completed successfully")

	slog.Warn("Branch management temporarily disabled for testing")

	// Clone requirements.yml to run directory for safe testing
	outputFile := filepath.Join(c.runDir.GetTmpOutPath(), "requirements-merged.yml")

	err = c.cloneRequirements(outputFile)
	if err != nil {
		return pkgerrors.ErrCloneRequirementsFileFailed(err)
	}

	// Load story document
	storyDoc, err := c.storyLoader.Load(storyNumber)
	if err != nil {
		return pkgerrors.ErrLoadStoryFailed(err)
	}

	// Merge scenarios from story into requirements-merged.yml (test file)
	err = c.mergeScenarios(ctx, storyNumber, storyDoc, outputFile)
	if err != nil {
		return pkgerrors.ErrMergeScenariosFailed(err)
	}

	slog.Info(
		"✅ Scenario merge completed successfully",
		"scenarios", len(storyDoc.Scenarios.TestScenarios),
		"story", storyNumber,
	)

	// Replace original requirements.yml with merged version
	err = c.replaceRequirements(outputFile)
	if err != nil {
		return pkgerrors.ErrReplaceRequirementsFailed(err)
	}

	// Implement tests for pending scenarios
	err = c.implementTests(ctx, "docs/requirements.yml")
	if err != nil {
		return pkgerrors.ErrImplementTestsFailed(err)
	}

	slog.Info("✅ User story implementation completed successfully", "backup", "docs/requirements.yml.backup")

	return nil
}

func (c *USImplementCommand) cloneRequirements(outputFile string) error {
	slog.Info("Cloning requirements file for safe testing", "output", outputFile)

	// Ensure tmp directory exists
	outputDir := filepath.Dir(outputFile)

	err := os.MkdirAll(outputDir, fileModeDirectory)
	if err != nil {
		return pkgerrors.ErrCreateOutputDirectoryFailed(outputDir, err)
	}

	// Read original file
	data, err := os.ReadFile("docs/requirements.yml")
	if err != nil {
		return pkgerrors.ErrReadRequirementsFileFailed(err)
	}

	// Write to output file
	err = os.WriteFile(outputFile, data, fileModeReadWrite)
	if err != nil {
		return pkgerrors.ErrWriteOutputFileFailed(outputFile, err)
	}

	slog.Info("✓ Cloned requirements file", "source", "docs/requirements.yml", "destination", outputFile)

	return nil
}

func (c *USImplementCommand) replaceRequirements(mergedFile string) error {
	const (
		requirementsPath = "docs/requirements.yml"
		backupPath       = "docs/requirements.yml.backup"
	)

	slog.Info("Replacing requirements file", "source", mergedFile)

	// Create backup
	originalData, err := os.ReadFile(requirementsPath)
	if err != nil {
		return pkgerrors.ErrReadOriginalFileFailed(err)
	}

	err = os.WriteFile(backupPath, originalData, fileModeReadWrite)
	if err != nil {
		return pkgerrors.ErrCreateBackupFileFailed(err)
	}

	slog.Info("✓ Created backup", "path", backupPath)

	// Read merged file
	mergedData, err := os.ReadFile(mergedFile)
	if err != nil {
		return pkgerrors.ErrReadMergedFileFailed(err)
	}

	// Replace original
	err = os.WriteFile(requirementsPath, mergedData, fileModeReadWrite)
	if err != nil {
		_ = os.WriteFile(requirementsPath, originalData, fileModeReadWrite) // Restore

		return pkgerrors.ErrReplaceFileFailed(err)
	}

	slog.Info("✓ Replaced requirements with merged scenarios", "path", requirementsPath)

	return nil
}

func (c *USImplementCommand) mergeScenarios(
	ctx context.Context,
	storyNumber string,
	storyDoc *storyModels.StoryDocument,
	outputFile string,
) error {
	slog.Info(
		"Starting scenario merge",
		"story_id", storyNumber,
		"scenario_count", len(storyDoc.Scenarios.TestScenarios),
		"output_file", outputFile,
	)

	// Create template loaders
	userPromptPath := c.config.GetString("templates.prompts.merge_scenarios")
	systemPromptPath := c.config.GetString("templates.prompts.merge_scenarios_system")
	userPromptLoader := template.NewTemplateLoader[*template.ScenarioMergeData](
		userPromptPath,
	)
	systemPromptLoader := template.NewTemplateLoader[*template.ScenarioMergeData](
		systemPromptPath,
	)

	// Process each scenario individually
	for i, scenario := range storyDoc.Scenarios.TestScenarios {
		startTime := time.Now()

		slog.Info(
			"Processing scenario",
			"index", i+1,
			"total", len(storyDoc.Scenarios.TestScenarios),
			"scenario_id", scenario.ID,
		)

		// Create merge data adapter
		mergeData := template.NewScenarioMergeData(storyNumber, scenario, outputFile)

		// Load templates using TemplateLoader
		userPrompt, err := userPromptLoader.LoadTemplate(mergeData)
		if err != nil {
			return pkgerrors.ErrLoadUserPromptForScenarioFailed(scenario.ID, err)
		}

		systemPrompt, err := systemPromptLoader.LoadTemplate(mergeData)
		if err != nil {
			return pkgerrors.ErrLoadSystemPromptForScenarioFailed(scenario.ID, err)
		}

		// Call Claude Code API to analyze and merge
		slog.Debug("Calling Claude Code for scenario merge", "scenario_id", scenario.ID)

		_, err = c.claudeClient.ExecutePromptWithSystem(
			ctx,
			systemPrompt,
			userPrompt,
			"sonnet",
			ai.ExecutionMode{
				AllowedTools: []string{"Read", "Edit"},
			},
		)
		if err != nil {
			return pkgerrors.ErrMergeScenarioFailed(scenario.ID, err)
		}

		duration := time.Since(startTime)
		slog.Info("✓ Scenario merged successfully", "scenario_id", scenario.ID, "duration", duration.Round(time.Second))
	}

	slog.Info("All scenarios merged successfully", "total_count", len(storyDoc.Scenarios.TestScenarios))

	return nil
}

func (c *USImplementCommand) implementTests(ctx context.Context, requirementsFile string) error {
	slog.Info("⚙️  Starting test implementation", "requirements_file", requirementsFile)

	// Parse requirements file to find pending scenarios
	pendingScenarios, err := c.parsePendingScenarios(requirementsFile)
	if err != nil {
		return pkgerrors.ErrParsePendingScenariosFailed(err)
	}

	if len(pendingScenarios) == 0 {
		slog.Info("✓ No pending scenarios to implement")

		return nil
	}

	slog.Info(
		"Found pending scenarios to implement",
		"count", len(pendingScenarios),
	)

	// Create template loaders
	userPromptLoader, systemPromptLoader := c.createTestTemplateLoaders()

	// Process each pending scenario
	implementedCount := c.processTestScenarios(ctx, pendingScenarios, userPromptLoader, systemPromptLoader)

	slog.Info(
		"✅ Test implementation completed",
		"implemented_count", implementedCount,
		"total_pending", len(pendingScenarios),
	)

	return nil
}

func (c *USImplementCommand) createTestTemplateLoaders() (
	*template.TemplateLoader[*template.TestImplementationData],
	*template.TemplateLoader[*template.TestImplementationData],
) {
	userPromptPath := c.config.GetString("templates.prompts.implement_tests")
	systemPromptPath := c.config.GetString("templates.prompts.implement_tests_system")
	userPromptLoader := template.NewTemplateLoader[*template.TestImplementationData](userPromptPath)
	systemPromptLoader := template.NewTemplateLoader[*template.TestImplementationData](systemPromptPath)

	return userPromptLoader, systemPromptLoader
}

func (c *USImplementCommand) processTestScenarios(
	ctx context.Context,
	scenarios []*template.TestImplementationData,
	userLoader *template.TemplateLoader[*template.TestImplementationData],
	systemLoader *template.TemplateLoader[*template.TestImplementationData],
) int {
	implementedCount := 0

	for i, scenario := range scenarios {
		startTime := time.Now()

		slog.Info(
			"Implementing test scenario",
			"progress", i+1,
			"total", len(scenarios),
			"scenario_id", scenario.ScenarioID,
		)

		if c.implementSingleTest(ctx, scenario, userLoader, systemLoader) {
			implementedCount++

			slog.Info(
				"✓ Test implemented successfully",
				"scenario_id", scenario.ScenarioID,
				"duration", time.Since(startTime).Round(time.Second),
			)
		}
	}

	return implementedCount
}

func (c *USImplementCommand) implementSingleTest(
	ctx context.Context,
	scenario *template.TestImplementationData,
	userLoader *template.TemplateLoader[*template.TestImplementationData],
	systemLoader *template.TemplateLoader[*template.TestImplementationData],
) bool {
	userPrompt, err := userLoader.LoadTemplate(scenario)
	if err != nil {
		slog.Warn(
			"⚠️  Skipping scenario: failed to load user prompt",
			"scenario_id", scenario.ScenarioID,
			"error", err,
		)

		return false
	}

	systemPrompt, err := systemLoader.LoadTemplate(scenario)
	if err != nil {
		slog.Warn(
			"⚠️  Skipping scenario: failed to load system prompt",
			"scenario_id", scenario.ScenarioID,
			"error", err,
		)

		return false
	}

	slog.Debug("Calling Claude Code for test implementation", "scenario_id", scenario.ScenarioID)

	_, err = c.claudeClient.ExecutePromptWithSystem(
		ctx,
		systemPrompt,
		userPrompt,
		"sonnet",
		ai.ExecutionMode{AllowedTools: []string{"Read", "Write", "Edit"}},
	)
	if err != nil {
		slog.Warn(
			"⚠️  Failed to implement test scenario",
			"scenario_id", scenario.ScenarioID,
			"error", err,
		)

		return false
	}

	return true
}

// parsePendingScenarios reads requirements file and extracts scenarios with
// status: "pending".
func (c *USImplementCommand) parsePendingScenarios(
	requirementsFile string,
) ([]*template.TestImplementationData, error) {
	slog.Debug("Parsing requirements file", "file", requirementsFile)

	// Read requirements file
	data, err := os.ReadFile(requirementsFile)
	if err != nil {
		return nil, pkgerrors.ErrReadRequirementsFailed(err)
	}

	// Parse YAML structure
	var requirements struct {
		Scenarios map[string]struct {
			Description          string `yaml:"description"`
			Category             string `yaml:"category"`
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

	// Filter pending scenarios and convert to TestImplementationData
	pendingScenarios := make([]*template.TestImplementationData, 0, len(requirements.Scenarios))

	for scenarioID, scenario := range requirements.Scenarios {
		// Only process scenarios with status "pending"
		if scenario.ImplementationStatus.Status != "pending" {
			slog.Debug(
				"Skipping non-pending scenario",
				"scenario_id", scenarioID,
				"status", scenario.ImplementationStatus.Status,
			)

			continue
		}

		// Convert interface{} arrays to string arrays
		givenSteps := convertStepsToStrings(scenario.MergedSteps.Given)
		whenSteps := convertStepsToStrings(scenario.MergedSteps.When)
		thenSteps := convertStepsToStrings(scenario.MergedSteps.Then)

		// Create TestImplementationData
		testData := template.NewTestImplementationData(
			scenarioID,
			scenario.Description,
			scenario.Level,
			scenario.Category,
			scenario.Priority,
			givenSteps,
			whenSteps,
			thenSteps,
			requirementsFile,
		)

		pendingScenarios = append(pendingScenarios, testData)

		slog.Debug("Found pending scenario", "scenario_id", scenarioID)
	}

	slog.Info(
		"Parsed requirements file",
		"total_scenarios", len(requirements.Scenarios),
		"pending_count", len(pendingScenarios),
	)

	return pendingScenarios, nil
}

// convertStepsToStrings converts []interface{} to []string, handling both string and map formats
// Example: "step" -> "step", {and: "step"} -> "And step".
func convertStepsToStrings(steps []interface{}) []string {
	result := make([]string, 0, len(steps))
	for _, step := range steps {
		switch v := step.(type) {
		case string:
			result = append(result, v)
		case map[string]interface{}:
			// Handle Gherkin keywords like {and: "step"}, {but: "step"}
			for keyword, value := range v {
				if strValue, ok := value.(string); ok {
					// Capitalize keyword and prepend
					result = append(result, keyword+" "+strValue)
				}
			}
		}
	}

	return result
}
