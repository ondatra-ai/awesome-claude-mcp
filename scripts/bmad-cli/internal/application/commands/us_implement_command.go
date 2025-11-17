package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"bmad-cli/internal/adapters/ai"
	storyModels "bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/shell"
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
	shellExec     *shell.CommandRunner
}

func NewUSImplementCommand(
	branchManager *git.BranchManager,
	storyLoader *story.StoryLoader,
	claudeClient *ai.ClaudeClient,
	cfg *config.ViperConfig,
	runDir *fs.RunDirectory,
	shellExec *shell.CommandRunner,
) *USImplementCommand {
	return &USImplementCommand{
		branchManager: branchManager,
		storyLoader:   storyLoader,
		claudeClient:  claudeClient,
		config:        cfg,
		runDir:        runDir,
		shellExec:     shellExec,
	}
}

func (c *USImplementCommand) Execute(ctx context.Context, storyNumber string, force bool, stepsStr string) error {
	// Parse steps
	steps, err := ParseSteps(stepsStr)
	if err != nil {
		return pkgerrors.ErrInvalidSteps(err)
	}

	slog.Info("Starting user story implementation",
		"story_number", storyNumber,
		"force", force,
		"steps", steps.String(),
	)

	// Execute all steps
	err = c.executeSteps(ctx, storyNumber, steps)
	if err != nil {
		return err
	}

	slog.Info("✅ User story implementation completed successfully", "backup", "docs/requirements.yml.backup")

	return nil
}

func (c *USImplementCommand) executeSteps(ctx context.Context, storyNumber string, steps *ExecutionSteps) error {
	err := c.executePreparationSteps(storyNumber, steps)
	if err != nil {
		return err
	}

	return c.executeImplementationSteps(ctx, storyNumber, steps)
}

func (c *USImplementCommand) executePreparationSteps(storyNumber string, steps *ExecutionSteps) error {
	// Step 1: Validate story
	if steps.ValidateStory {
		err := c.executeValidateStory(storyNumber)
		if err != nil {
			return err
		}
	}

	// Step 2: Create branch
	if steps.CreateBranch {
		err := c.executeCreateBranch(storyNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *USImplementCommand) executeImplementationSteps(
	ctx context.Context,
	storyNumber string,
	steps *ExecutionSteps,
) error {
	// Step 3: Merge scenarios
	if steps.MergeScenarios {
		_, err := c.executeMergeScenarios(ctx, storyNumber)
		if err != nil {
			return err
		}
	}

	// Step 4: Generate tests
	if steps.GenerateTests {
		err := c.executeGenerateTests(ctx)
		if err != nil {
			return err
		}
	}

	// Step 5: Implement feature
	if steps.ImplementFeature {
		err := c.executeImplementFeature(ctx, storyNumber)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *USImplementCommand) executeValidateStory(storyNumber string) error {
	slog.Info("Step 1: Validating story file")

	storySlug, err := c.storyLoader.GetStorySlug(storyNumber)
	if err != nil {
		return pkgerrors.ErrGetStorySlugFailed(err)
	}

	slog.Info("✓ Story validated successfully", "slug", storySlug)

	return nil
}

func (c *USImplementCommand) executeCreateBranch(storyNumber string) error {
	slog.Info("Step 2: Creating story branch")

	storySlug, err := c.storyLoader.GetStorySlug(storyNumber)
	if err != nil {
		return pkgerrors.ErrGetStorySlugFailed(err)
	}

	// TEMPORARILY COMMENTED OUT FOR TESTING
	// if err := c.branchManager.EnsureBranch(ctx, storyNumber, storySlug, force); err != nil {
	// 	return fmt.Errorf("failed to ensure branch: %w", err)
	// }
	_ = storySlug

	slog.Warn("⚠️  Branch management temporarily disabled for testing")

	return nil
}

func (c *USImplementCommand) executeMergeScenarios(
	ctx context.Context,
	storyNumber string,
) (*storyModels.StoryDocument, error) {
	slog.Info("Step 3: Merging scenarios into requirements")

	// Clone requirements.yml to run directory for safe testing
	outputFile := filepath.Join(c.runDir.GetTmpOutPath(), "requirements-merged.yml")

	err := c.cloneRequirements(outputFile)
	if err != nil {
		return nil, pkgerrors.ErrCloneRequirementsFileFailed(err)
	}

	// Load story document
	storyDoc, err := c.storyLoader.Load(storyNumber)
	if err != nil {
		return nil, pkgerrors.ErrLoadStoryFailed(err)
	}

	// Merge scenarios from story into requirements-merged.yml
	err = c.mergeScenarios(ctx, storyNumber, storyDoc, outputFile)
	if err != nil {
		return nil, pkgerrors.ErrMergeScenariosFailed(err)
	}

	slog.Info(
		"✅ Scenario merge completed successfully",
		"scenarios", len(storyDoc.Scenarios.TestScenarios),
		"story", storyNumber,
	)

	// Replace original requirements.yml with merged version
	err = c.replaceRequirements(outputFile)
	if err != nil {
		return nil, pkgerrors.ErrReplaceRequirementsFailed(err)
	}

	return storyDoc, nil
}

func (c *USImplementCommand) executeGenerateTests(ctx context.Context) error {
	slog.Info("Step 4: Generating test code")

	err := c.generateTests(ctx, "docs/requirements.yml")
	if err != nil {
		return pkgerrors.ErrGenerateTestsFailed(err)
	}

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

func (c *USImplementCommand) generateTests(ctx context.Context, requirementsFile string) error {
	slog.Info("⚙️  Starting test generation", "requirements_file", requirementsFile)

	// Step 1: Validate baseline tests (must pass)
	err := c.validateBaselineTests(ctx)
	if err != nil {
		return err
	}

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
		"✅ Test generation completed",
		"implemented_count", implementedCount,
		"total_pending", len(pendingScenarios),
	)

	// Step 2: Validate generated tests (must fail - TDD red phase)
	err = c.validateGeneratedTests(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *USImplementCommand) createTestTemplateLoaders() (
	*template.TemplateLoader[*template.TestGenerationData],
	*template.TemplateLoader[*template.TestGenerationData],
) {
	userPromptPath := c.config.GetString("templates.prompts.generate_tests")
	systemPromptPath := c.config.GetString("templates.prompts.generate_tests_system")
	userPromptLoader := template.NewTemplateLoader[*template.TestGenerationData](userPromptPath)
	systemPromptLoader := template.NewTemplateLoader[*template.TestGenerationData](systemPromptPath)

	return userPromptLoader, systemPromptLoader
}

// runTests executes the test command and saves output to tmp directory.
func (c *USImplementCommand) runTests(ctx context.Context, phase string) (string, error) {
	slog.Info("🧪 Running tests", "phase", phase)

	// Get test command from config
	testCommand := c.config.GetString("testing.command")
	if testCommand == "" {
		testCommand = "make test-e2e" // Default fallback
	}

	slog.Debug("Executing test command", "command", testCommand)

	// Execute test command
	output, err := c.shellExec.Run(ctx, "sh", "-c", testCommand)
	if err != nil {
		return output, pkgerrors.ErrRunTestsFailed(phase, err)
	}

	// Save output to tmp directory
	outputFile := filepath.Join(c.runDir.GetTmpOutPath(), "test-output-"+phase+".txt")

	const filePermission = 0o644

	writeErr := os.WriteFile(outputFile, []byte(output), filePermission)
	if writeErr != nil {
		slog.Warn("Failed to write test output", "file", outputFile, "error", writeErr)
	} else {
		slog.Debug("Test output saved", "file", outputFile)
	}

	return output, nil
}

// validateBaselineTests runs tests before test generation and ensures they pass.
func (c *USImplementCommand) validateBaselineTests(ctx context.Context) error {
	slog.Info("📋 Validating baseline tests (must pass)")

	output, err := c.runTests(ctx, "before")
	if err != nil {
		// Tests failed - this is an error for baseline
		slog.Error(
			"❌ Baseline tests failed",
			"error", err,
			"output_file", filepath.Join(c.runDir.GetTmpOutPath(), "test-output-before.txt"),
		)

		return pkgerrors.ErrBaselineTestsFailedError(output)
	}

	slog.Info("✅ Baseline tests passed - ready for test generation")

	return nil
}

// validateGeneratedTests runs tests after test generation and ensures they fail (TDD red phase).
func (c *USImplementCommand) validateGeneratedTests(ctx context.Context) error {
	slog.Info("🔴 Validating generated tests (must fail - TDD red phase)")

	output, err := c.runTests(ctx, "after")

	// Check if tests passed (err == nil means success, which is bad for TDD red phase)
	if err == nil {
		// Tests passed - this is an error (they should be failing)
		slog.Error(
			"❌ Generated tests are passing but should fail (TDD red phase)",
			"output_file", filepath.Join(c.runDir.GetTmpOutPath(), "test-output-after.txt"),
		)

		return pkgerrors.ErrGeneratedTestsPassError(output)
	}

	// Tests failed - this is expected (TDD red phase)
	slog.Info(
		"✅ Generated tests are failing as expected (TDD red phase)",
		"output_file", filepath.Join(c.runDir.GetTmpOutPath(), "test-output-after.txt"),
	)

	return nil
}

func (c *USImplementCommand) processTestScenarios(
	ctx context.Context,
	scenarios []*template.TestGenerationData,
	userLoader *template.TemplateLoader[*template.TestGenerationData],
	systemLoader *template.TemplateLoader[*template.TestGenerationData],
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
	scenario *template.TestGenerationData,
	userLoader *template.TemplateLoader[*template.TestGenerationData],
	systemLoader *template.TemplateLoader[*template.TestGenerationData],
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

	// Save user prompt to tmp directory for debugging
	userPromptFile := filepath.Join(c.runDir.GetTmpOutPath(),
		scenario.ScenarioID+"-test-generation-user-prompt.txt")

	writeErr := os.WriteFile(userPromptFile, []byte(userPrompt), fileModeReadWrite)
	if writeErr != nil {
		slog.Warn("Failed to save user prompt", "file", userPromptFile, "error", writeErr)
	} else {
		slog.Info("💾 User prompt saved", "file", userPromptFile, "scenario_id", scenario.ScenarioID)
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

	// Save system prompt to tmp directory for debugging
	systemPromptFile := filepath.Join(c.runDir.GetTmpOutPath(),
		scenario.ScenarioID+"-test-generation-system-prompt.txt")

	writeErr = os.WriteFile(systemPromptFile, []byte(systemPrompt), fileModeReadWrite)
	if writeErr != nil {
		slog.Warn("Failed to save system prompt", "file", systemPromptFile, "error", writeErr)
	} else {
		slog.Info("💾 System prompt saved", "file", systemPromptFile, "scenario_id", scenario.ScenarioID)
	}

	slog.Info("🤖 Calling Claude for test generation", "scenario_id", scenario.ScenarioID)

	response, err := c.claudeClient.ExecutePromptWithSystem(
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

	// Save Claude response to tmp directory for debugging
	responseFile := filepath.Join(c.runDir.GetTmpOutPath(),
		scenario.ScenarioID+"-test-generation-response.txt")

	writeErr = os.WriteFile(responseFile, []byte(response), fileModeReadWrite)
	if writeErr != nil {
		slog.Warn("Failed to save Claude response", "file", responseFile, "error", writeErr)
	} else {
		slog.Info("💾 Claude response saved", "file", responseFile, "scenario_id", scenario.ScenarioID)
	}

	return true
}

// parsePendingScenarios reads requirements file and extracts scenarios with
// status: "pending".
func (c *USImplementCommand) parsePendingScenarios(
	requirementsFile string,
) ([]*template.TestGenerationData, error) {
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

	// Filter pending scenarios and convert to TestGenerationData
	pendingScenarios := make([]*template.TestGenerationData, 0, len(requirements.Scenarios))

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

		// Create TestGenerationData
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

// ImplementFeatureData holds the data for the implement feature prompt.
type ImplementFeatureData struct {
	StoryID     string
	StoryTitle  string
	AsA         string
	IWant       string
	SoThat      string
	TestCommand string
	TestOutput  string
	Attempt     int
	MaxAttempts int
}

// savePromptFile saves a prompt to tmp directory for debugging.
func (c *USImplementCommand) savePromptFile(content, filename string) {
	filePath := filepath.Join(c.runDir.GetTmpOutPath(), filename)

	writeErr := os.WriteFile(filePath, []byte(content), fileModeReadWrite)
	if writeErr != nil {
		slog.Warn("Failed to save prompt file", "file", filePath, "error", writeErr)
	} else {
		slog.Info("💾 Prompt saved", "file", filePath)
	}
}

func (c *USImplementCommand) executeImplementFeature(ctx context.Context, storyNumber string) error {
	slog.Info("Step 5: Implementing feature")

	// Load story to get basic context
	storyDoc, err := c.storyLoader.Load(storyNumber)
	if err != nil {
		return pkgerrors.ErrLoadStoryFailed(err)
	}

	// Read test command from config
	testCommand := c.config.GetString("testing.command")
	if testCommand == "" {
		testCommand = "make test-e2e" // Default fallback
	}

	// Load prompt templates once
	userPromptPath := c.config.GetString("templates.prompts.implement_feature")
	systemPromptPath := c.config.GetString("templates.prompts.implement_feature_system")

	userPromptLoader := template.NewTemplateLoader[*ImplementFeatureData](userPromptPath)
	systemPromptLoader := template.NewTemplateLoader[*ImplementFeatureData](systemPromptPath)

	// Iteration loop: try up to 5 times to make tests pass
	const maxAttempts = 5

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		slog.Info("🔄 Implementation attempt", "attempt", attempt, "max", maxAttempts)

		// Run tests to get current state
		testOutput, testErr := c.runTests(ctx, "implement-feature")

		// Check if tests are passing
		if testErr == nil {
			slog.Info("✅ All tests passing - feature implementation complete!", "attempts_used", attempt)

			return nil
		}

		// Tests are failing - need to implement/fix code
		slog.Info("❌ Tests failing - calling Claude to fix", "attempt", attempt)

		// Create prompt data with test output
		promptData := &ImplementFeatureData{
			StoryID:     storyNumber,
			StoryTitle:  storyDoc.Story.Title,
			AsA:         storyDoc.Story.AsA,
			IWant:       storyDoc.Story.IWant,
			SoThat:      storyDoc.Story.SoThat,
			TestCommand: testCommand,
			TestOutput:  testOutput,
			Attempt:     attempt,
			MaxAttempts: maxAttempts,
		}

		userPrompt, err := userPromptLoader.LoadTemplate(promptData)
		if err != nil {
			return pkgerrors.ErrLoadPromptsFailed(err)
		}

		c.savePromptFile(userPrompt, fmt.Sprintf("%s-implement-feature-attempt-%d-user-prompt.txt", storyNumber, attempt))

		systemPrompt, err := systemPromptLoader.LoadTemplate(promptData)
		if err != nil {
			return pkgerrors.ErrLoadPromptsFailed(err)
		}

		c.savePromptFile(systemPrompt, fmt.Sprintf("%s-implement-feature-attempt-%d-system-prompt.txt", storyNumber, attempt))

		slog.Info("🤖 Calling Claude to implement feature", "attempt", attempt)

		response, err := c.claudeClient.ExecutePromptWithSystem(
			ctx,
			systemPrompt,
			userPrompt,
			"sonnet",
			ai.ExecutionMode{AllowedTools: []string{"Read", "Write", "Edit", "Bash"}},
		)
		if err != nil {
			return pkgerrors.ErrImplementFeaturesFailed(err)
		}

		c.savePromptFile(response, fmt.Sprintf("%s-implement-feature-attempt-%d-response.txt", storyNumber, attempt))

		slog.Info("✓ Claude finished attempt", "attempt", attempt)
	}

	// If we get here, we've exhausted all attempts and tests are still failing
	slog.Error("❌ Failed to make tests pass after maximum attempts", "max_attempts", maxAttempts)

	return pkgerrors.ErrImplementFeaturesMaxAttemptsExceeded(maxAttempts)
}
