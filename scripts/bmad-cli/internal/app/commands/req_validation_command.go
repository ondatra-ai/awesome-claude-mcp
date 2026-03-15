package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"bmad-cli/internal/app/generators/validate"
	checklistmodels "bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/infrastructure/checklist"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/input"
	"bmad-cli/internal/infrastructure/requirements"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/pkg/console"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const (
	reqMaxClarificationIterations = 5
	reqMaxRefinementIterations    = 3
	reqSeparatorWidth             = 80
	testFilePermissions           = 0o644
	testDirPermissions            = 0o755
)

// ReqValidationCommand validates generated tests against BDD scenarios using checklists.
type ReqValidationCommand struct {
	checklistLoader    *checklist.ChecklistLoader
	checklistEvaluator *validate.ChecklistEvaluator
	fixPromptGenerator *validate.FixPromptGenerator
	fixApplier         *validate.FixApplier
	userInputCollector *input.UserInputCollector
	tableRenderer      *TableRenderer
	scenarioParser     *requirements.ScenarioParser
	runDir             *fs.RunDirectory
}

// NewReqValidationCommand creates a new test validation command.
func NewReqValidationCommand(
	checklistLoader *checklist.ChecklistLoader,
	evaluator *validate.ChecklistEvaluator,
	fixPromptGen *validate.FixPromptGenerator,
	fixApplier *validate.FixApplier,
	inputCollector *input.UserInputCollector,
	renderer *TableRenderer,
	scenarioParser *requirements.ScenarioParser,
	runDir *fs.RunDirectory,
) *ReqValidationCommand {
	return &ReqValidationCommand{
		checklistLoader:    checklistLoader,
		checklistEvaluator: evaluator,
		fixPromptGenerator: fixPromptGen,
		fixApplier:         fixApplier,
		userInputCollector: inputCollector,
		tableRenderer:      renderer,
		scenarioParser:     scenarioParser,
		runDir:             runDir,
	}
}

// Execute runs test validation for scenarios in the requirements file.
func (c *ReqValidationCommand) Execute(
	ctx context.Context,
	requirementsFile string,
	fix bool,
	all bool,
) error {
	console.Header("TEST VALIDATION", reqSeparatorWidth)

	// Parse scenarios
	scenarios, err := c.scenarioParser.ParseScenarios(requirementsFile, !all)
	if err != nil {
		return fmt.Errorf("failed to parse scenarios: %w", err)
	}

	if len(scenarios) == 0 {
		console.Println("No scenarios found to validate.")

		return nil
	}

	slog.Info("Found scenarios to validate", "count", len(scenarios))

	// Load checklist
	checklistData, err := c.checklistLoader.Load()
	if err != nil {
		return fmt.Errorf("failed to load test checklist: %w", err)
	}

	prompts := c.checklistLoader.ExtractPromptsForStage(checklistData, "test_validation")
	if len(prompts) == 0 {
		console.Println("No validation prompts found for test_validation stage.")

		return pkgerrors.ErrNoPromptsForStageFailed("test_validation")
	}

	slog.Info("Extracted test validation prompts", "count", len(prompts))

	// Validate each scenario
	for i, scenario := range scenarios {
		console.Header(fmt.Sprintf("SCENARIO %d/%d: %s", i+1, len(scenarios), scenario.ScenarioID), reqSeparatorWidth)

		err = c.validateScenario(ctx, scenario, prompts, fix)
		if err != nil {
			slog.Error("Failed to validate scenario",
				"scenario_id", scenario.ScenarioID,
				"error", err,
			)
			console.Printf("Error validating %s: %v\n", scenario.ScenarioID, err)
		}
	}

	console.Header("TEST VALIDATION COMPLETE", reqSeparatorWidth)

	return nil
}

// validateScenario validates a single test scenario.
func (c *ReqValidationCommand) validateScenario(
	ctx context.Context,
	scenario *template.TestGenerationData,
	prompts []checklistmodels.PromptWithContext,
	fix bool,
) error {
	// Load existing test file content if available
	c.loadTestContent(scenario)

	tmpDir := c.runDir.GetTmpOutPath()
	versionMgr := fs.NewContentVersionManager(c.runDir, scenario.ScenarioID, "test")

	// Save initial version
	err := versionMgr.SaveInitialVersion([]byte(scenario.ArchitectureContent))
	if err != nil {
		return fmt.Errorf("failed to save initial test version: %w", err)
	}

	// Run validation loop
	for iteration := 1; ; iteration++ {
		shouldContinue, loopErr := c.runScenarioIteration(
			ctx, scenario, prompts, tmpDir, versionMgr, fix, iteration)
		if loopErr != nil {
			return loopErr
		}

		if !shouldContinue {
			return nil
		}
	}
}

// runScenarioIteration runs one iteration of test validation.
func (c *ReqValidationCommand) runScenarioIteration(
	ctx context.Context,
	scenario *template.TestGenerationData,
	prompts []checklistmodels.PromptWithContext,
	tmpDir string,
	versionMgr *fs.ContentVersionManager,
	fix bool,
	iteration int,
) (bool, error) {
	// Reload latest content
	latestContent, err := versionMgr.LoadLatest()
	if err != nil {
		return false, fmt.Errorf("failed to load test version: %w", err)
	}

	scenario.ArchitectureContent = string(latestContent)

	var report *checklistmodels.ChecklistReport

	if fix {
		report, err = c.checklistEvaluator.EvaluateUntilFailure(
			ctx, scenario, scenario.ScenarioID, scenario.Description, 0, prompts, tmpDir)
	} else {
		report, err = c.checklistEvaluator.Evaluate(
			ctx, scenario, scenario.ScenarioID, scenario.Description, 0, prompts, tmpDir)
	}

	if err != nil {
		return false, fmt.Errorf("failed to evaluate checklist: %w", err)
	}

	c.tableRenderer.RenderReport(report, fix)

	if report.AllPassed() {
		console.Header("ALL CHECKS PASSED!", reqSeparatorWidth)
		console.Printf("Latest version: %s\n", versionMgr.GetLatestPath())

		return false, nil
	}

	failedCheck := c.getFirstFailedCheck(report)
	if failedCheck == nil {
		slog.Warn("No failed check found despite not all passed")

		return false, nil
	}

	c.displayFailureInfo(failedCheck)

	if !fix {
		console.BlankLine()
		console.Printf("Validation failed. Use --fix flag to enter interactive fix mode.\n")

		return false, nil
	}

	return c.runTestFixLoop(ctx, scenario, *failedCheck, tmpDir, versionMgr, iteration)
}

// runTestFixLoop handles the fix prompt generation and refinement loop for tests.
func (c *ReqValidationCommand) runTestFixLoop(
	ctx context.Context,
	scenario *template.TestGenerationData,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	versionMgr *fs.ContentVersionManager,
	iteration int,
) (bool, error) {
	userAnswers := make(map[string]string)
	refinementCount := 0

	fixPrompt, answers, err := c.generateTestFixPrompt(
		ctx, scenario, failedCheck, tmpDir, userAnswers)
	if err != nil {
		return false, fmt.Errorf("failed to generate fix prompt: %w", err)
	}

	userAnswers = answers

	if fixPrompt == "" {
		slog.Warn("No fix prompt generated")

		return false, nil
	}

	for {
		c.displayFixPrompt(fixPrompt)

		action := c.userInputCollector.AskApplyRefineOrExit()

		switch action {
		case input.ActionApply:
			return c.applyTestFix(ctx, scenario, fixPrompt, tmpDir, versionMgr, iteration)

		case input.ActionRefine:
			if refinementCount >= reqMaxRefinementIterations {
				console.Printf("\nMax refinement attempts (%d) reached.\n", reqMaxRefinementIterations)

				continue
			}

			refinementCount++

			newPrompt, updatedAnswers, refineErr := c.refineTestFixPrompt(
				ctx, scenario, failedCheck, tmpDir, userAnswers, refinementCount)
			if refineErr != nil {
				return false, pkgerrors.ErrFixPromptRefinementFailed(refineErr)
			}

			if newPrompt == "" {
				console.Println("\nNo feedback provided. Keeping current fix prompt.")

				continue
			}

			fixPrompt = newPrompt
			userAnswers = updatedAnswers

			console.Printf("\n(Refinement %d of %d)\n", refinementCount, reqMaxRefinementIterations)

		case input.ActionExit:
			console.Printf("\nExiting. Latest version saved at: %s\n", versionMgr.GetLatestPath())

			return false, nil
		}
	}
}

// generateTestFixPrompt generates fix prompt with clarification loop.
func (c *ReqValidationCommand) generateTestFixPrompt(
	ctx context.Context,
	scenario *template.TestGenerationData,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	initialAnswers map[string]string,
) (string, map[string]string, error) {
	userAnswers := make(map[string]string)

	for id, answer := range initialAnswers {
		userAnswers[id] = answer
	}

	for iteration := 1; iteration <= reqMaxClarificationIterations; iteration++ {
		params := validate.GenerateParams{
			Subject:     scenario,
			SubjectID:   scenario.ScenarioID,
			FailedCheck: failedCheck,
			TmpDir:      tmpDir,
			UserAnswers: userAnswers,
			Iteration:   iteration,
		}

		result, err := c.fixPromptGenerator.Generate(ctx, params)
		if err != nil {
			return "", userAnswers, pkgerrors.ErrFixPromptGenerationFailed(err)
		}

		if result.HasFixPrompt() {
			return result.FixPrompt, userAnswers, nil
		}

		if !result.HasQuestions() {
			return "", userAnswers, nil
		}

		answers := c.userInputCollector.AskQuestions(result.Questions)

		for id, answer := range answers {
			userAnswers[id] = answer
		}
	}

	return "", userAnswers, nil
}

// refineTestFixPrompt collects user feedback and regenerates the fix prompt.
func (c *ReqValidationCommand) refineTestFixPrompt(
	ctx context.Context,
	scenario *template.TestGenerationData,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	existingAnswers map[string]string,
	refinementIteration int,
) (string, map[string]string, error) {
	feedback := c.userInputCollector.AskRefinementFeedback()
	if feedback == "" {
		return "", existingAnswers, nil
	}

	existingAnswers["_user_refinement"] = feedback

	params := validate.GenerateParams{
		Subject:     scenario,
		SubjectID:   scenario.ScenarioID,
		FailedCheck: failedCheck,
		TmpDir:      tmpDir,
		UserAnswers: existingAnswers,
		Iteration:   refinementIteration + reqMaxClarificationIterations,
	}

	result, err := c.fixPromptGenerator.Generate(ctx, params)
	if err != nil {
		return "", existingAnswers, pkgerrors.ErrFixPromptGenerationFailed(err)
	}

	if !result.HasFixPrompt() {
		return "", existingAnswers, nil
	}

	return result.FixPrompt, existingAnswers, nil
}

// applyTestFix applies the fix and saves a new version.
func (c *ReqValidationCommand) applyTestFix(
	ctx context.Context,
	scenario *template.TestGenerationData,
	fixPrompt string,
	tmpDir string,
	versionMgr *fs.ContentVersionManager,
	iteration int,
) (bool, error) {
	content, err := c.fixApplier.Apply(
		ctx, scenario, scenario.ScenarioID, fixPrompt, tmpDir, iteration)
	if err != nil {
		return false, fmt.Errorf("failed to apply fix: %w", err)
	}

	_, err = versionMgr.SaveNextVersion([]byte(content))
	if err != nil {
		return false, pkgerrors.ErrSaveStoryVersionFailed(err)
	}

	// Write to working directory
	if scenario.TestFilePath != "" {
		writeErr := writeTestFile(scenario.TestFilePath, content)
		if writeErr != nil {
			return false, fmt.Errorf("failed to write test file: %w", writeErr)
		}

		console.Printf("\nTest file written to: %s\n", scenario.TestFilePath)
	}

	// Update scenario content for next iteration
	scenario.ArchitectureContent = content

	console.Printf("Fix applied. Saved as version %d.\n", versionMgr.GetCurrentVersion())
	console.Println("Re-running validation...")

	return true, nil
}

// loadTestContent reads the existing test file content into the scenario.
func (c *ReqValidationCommand) loadTestContent(scenario *template.TestGenerationData) {
	if scenario.ArchitectureContent != "" {
		return
	}

	if scenario.TestFilePath == "" {
		slog.Debug("No test file path for scenario", "scenario_id", scenario.ScenarioID)

		return
	}

	data, err := os.ReadFile(scenario.TestFilePath)
	if err != nil {
		slog.Debug("No existing test file on disk",
			"scenario_id", scenario.ScenarioID,
			"path", scenario.TestFilePath,
			"error", err,
		)

		return
	}

	scenario.ArchitectureContent = string(data)
	slog.Info("Loaded existing test content",
		"scenario_id", scenario.ScenarioID,
		"path", scenario.TestFilePath,
	)
}

// writeTestFile writes content to the test file path, creating directories as needed.
func writeTestFile(filePath string, content string) error {
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, testDirPermissions)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	err = os.WriteFile(filePath, []byte(content), testFilePermissions)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

func (c *ReqValidationCommand) getFirstFailedCheck(
	report *checklistmodels.ChecklistReport,
) *checklistmodels.ValidationResult {
	for _, result := range report.Results {
		if result.Status == checklistmodels.StatusFail {
			return &result
		}
	}

	return nil
}

func (c *ReqValidationCommand) displayFailureInfo(failedCheck *checklistmodels.ValidationResult) {
	console.BlankLine()
	console.Separator("=", reqSeparatorWidth)
	console.Printf("CHECK FAILED: %s\n", failedCheck.SectionPath)
	console.Separator("=", reqSeparatorWidth)
	console.Printf("Question: %s\n", failedCheck.Question)
	console.Printf("Expected: %s\n", failedCheck.ExpectedAnswer)
	console.Printf("Actual: %s\n", failedCheck.ActualAnswer)

	if failedCheck.Rationale != "" {
		console.Printf("Rationale: %s\n", failedCheck.Rationale)
	}
}

func (c *ReqValidationCommand) displayFixPrompt(fixPrompt string) {
	console.BlankLine()
	console.Header("FIX PROMPT GENERATED", reqSeparatorWidth)
	console.Println(fixPrompt)
	console.Separator("=", reqSeparatorWidth)
}
