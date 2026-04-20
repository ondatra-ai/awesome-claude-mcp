package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"bmad-cli/internal/app/generators/validate"
	checklistmodels "bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/input"
	"bmad-cli/internal/infrastructure/requirements"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/pkg/console"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const (
	testDirPermissions  = 0o755
	testFilePermissions = 0o644
)

// ExecuteScenarioChecklist runs the given checklist against every scenario
// in requirements.yaml. Used by `us generate_tests` and `us implement`.
// With --fix, enters the interactive fix loop per scenario.
func (c *USValidationCommand) ExecuteScenarioChecklist(
	ctx context.Context,
	requirementsFile string,
	checklistName string,
	fix bool,
	scenarioParser *requirements.ScenarioParser,
	evaluator *validate.ChecklistEvaluator,
	fixGenerator *validate.FixPromptGenerator,
	fixApplier *validate.FixApplier,
) error {
	console.Header(
		strings.ToUpper(checklistName)+" — SCENARIO VALIDATION",
		separatorWidth,
	)

	scenarios, err := scenarioParser.ParseScenarios(requirementsFile, false)
	if err != nil {
		return fmt.Errorf("failed to parse scenarios: %w", err)
	}

	if len(scenarios) == 0 {
		console.Println("No scenarios found to validate.")

		return nil
	}

	slog.Info("Found scenarios to validate",
		"checklist", checklistName,
		"count", len(scenarios),
	)

	prompts, err := c.checklistLoader.Load(checklistName)
	if err != nil {
		return fmt.Errorf("failed to load checklist %q: %w", checklistName, err)
	}

	// Empty-checklist short-circuit: when the checklist is a deliberate
	// placeholder (e.g. us-implement today), walking all scenarios to
	// run zero prompts is pointless. Report and return cleanly.
	if len(prompts) == 0 {
		console.Println(
			fmt.Sprintf("No validation prompts defined for %s. Nothing to do.", checklistName),
		)

		return nil
	}

	slog.Info("Loaded validation prompts",
		"checklist", checklistName,
		"count", len(prompts),
	)

	for i, scenario := range scenarios {
		console.Header(
			fmt.Sprintf("SCENARIO %d/%d: %s", i+1, len(scenarios), scenario.ScenarioID),
			separatorWidth,
		)

		err := c.validateScenario(
			ctx, scenario, prompts, fix,
			evaluator, fixGenerator, fixApplier,
		)
		if err != nil {
			slog.Error("Failed to validate scenario",
				"scenario_id", scenario.ScenarioID,
				"error", err,
			)
			console.Printf("Error validating %s: %v\n", scenario.ScenarioID, err)
		}
	}

	console.Header(
		strings.ToUpper(checklistName)+" — VALIDATION COMPLETE",
		separatorWidth,
	)

	return nil
}

func (c *USValidationCommand) validateScenario(
	ctx context.Context,
	scenario *template.TestGenerationData,
	prompts []checklistmodels.PromptWithContext,
	fix bool,
	testEvaluator *validate.ChecklistEvaluator,
	testFixGenerator *validate.FixPromptGenerator,
	testFixApplier *validate.FixApplier,
) error {
	c.loadTestContent(scenario)

	tmpDir := c.runDir.GetTmpOutPath()
	versionMgr := fs.NewContentVersionManager(c.runDir, scenario.ScenarioID, "test")

	err := versionMgr.SaveInitialVersion([]byte(scenario.ArchitectureContent))
	if err != nil {
		return fmt.Errorf("failed to save initial test version: %w", err)
	}

	for iteration := 1; ; iteration++ {
		shouldContinue, loopErr := c.runScenarioIteration(
			ctx, scenario, prompts, tmpDir, versionMgr, fix, iteration,
			testEvaluator, testFixGenerator, testFixApplier,
		)
		if loopErr != nil {
			return loopErr
		}

		if !shouldContinue {
			return nil
		}
	}
}

func (c *USValidationCommand) runScenarioIteration(
	ctx context.Context,
	scenario *template.TestGenerationData,
	prompts []checklistmodels.PromptWithContext,
	tmpDir string,
	versionMgr *fs.ContentVersionManager,
	fix bool,
	iteration int,
	testEvaluator *validate.ChecklistEvaluator,
	testFixGenerator *validate.FixPromptGenerator,
	testFixApplier *validate.FixApplier,
) (bool, error) {
	latestContent, err := versionMgr.LoadLatest()
	if err != nil {
		return false, fmt.Errorf("failed to load test version: %w", err)
	}

	scenario.ArchitectureContent = string(latestContent)

	report, err := c.evaluateTestScenario(ctx, scenario, prompts, tmpDir, fix, testEvaluator)
	if err != nil {
		return false, err
	}

	c.tableRenderer.RenderReport(report, fix)

	if report.AllPassed() {
		console.Header("ALL CHECKS PASSED!", separatorWidth)
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

	return c.runTestFixLoop(
		ctx, scenario, *failedCheck, tmpDir, versionMgr, iteration,
		testFixGenerator, testFixApplier,
	)
}

func (c *USValidationCommand) evaluateTestScenario(
	ctx context.Context,
	scenario *template.TestGenerationData,
	prompts []checklistmodels.PromptWithContext,
	tmpDir string,
	fix bool,
	testEvaluator *validate.ChecklistEvaluator,
) (*checklistmodels.ChecklistReport, error) {
	var (
		report *checklistmodels.ChecklistReport
		err    error
	)

	if fix {
		report, err = testEvaluator.EvaluateUntilFailure(
			ctx, scenario, scenario.ScenarioID, scenario.Description, prompts, tmpDir)
	} else {
		report, err = testEvaluator.Evaluate(
			ctx, scenario, scenario.ScenarioID, scenario.Description, prompts, tmpDir)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to evaluate checklist: %w", err)
	}

	return report, nil
}

func (c *USValidationCommand) runTestFixLoop(
	ctx context.Context,
	scenario *template.TestGenerationData,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	versionMgr *fs.ContentVersionManager,
	iteration int,
	testFixGenerator *validate.FixPromptGenerator,
	testFixApplier *validate.FixApplier,
) (bool, error) {
	userAnswers := make(map[string]string)
	refinementCount := 0

	fixPrompt, answers, err := c.generateTestFixPrompt(
		ctx, scenario, failedCheck, tmpDir, userAnswers, testFixGenerator,
	)
	if err != nil {
		return false, fmt.Errorf("failed to generate fix prompt: %w", err)
	}

	userAnswers = answers

	if fixPrompt == "" {
		return false, nil
	}

	for {
		c.displayFixPrompt(fixPrompt)

		action := c.userInputCollector.AskApplyRefineOrExit()

		switch action {
		case input.ActionApply:
			return c.applyTestFix(
				ctx, scenario, fixPrompt, tmpDir, versionMgr, iteration, testFixApplier,
			)

		case input.ActionRefine:
			if refinementCount >= maxRefinementIterations {
				console.Printf(
					"\nMax refinement attempts (%d) reached.\n",
					maxRefinementIterations,
				)

				continue
			}

			refinementCount++

			newPrompt, updatedAnswers, refineErr := c.refineTestFixPrompt(
				ctx, scenario, failedCheck, tmpDir, userAnswers, refinementCount,
				testFixGenerator,
			)
			if refineErr != nil {
				return false, pkgerrors.ErrFixPromptRefinementFailed(refineErr)
			}

			if newPrompt == "" {
				console.Println("\nNo feedback provided. Keeping current fix prompt.")

				continue
			}

			fixPrompt = newPrompt
			userAnswers = updatedAnswers

			console.Printf(
				"\n(Refinement %d of %d)\n",
				refinementCount, maxRefinementIterations,
			)

		case input.ActionExit:
			console.Printf(
				"\nExiting. Latest version saved at: %s\n",
				versionMgr.GetLatestPath(),
			)

			return false, nil
		}
	}
}

func (c *USValidationCommand) generateTestFixPrompt(
	ctx context.Context,
	scenario *template.TestGenerationData,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	initialAnswers map[string]string,
	testFixGenerator *validate.FixPromptGenerator,
) (string, map[string]string, error) {
	userAnswers := make(map[string]string)

	for id, answer := range initialAnswers {
		userAnswers[id] = answer
	}

	for iteration := 1; iteration <= maxClarificationIterations; iteration++ {
		params := validate.GenerateParams{
			Subject:     scenario,
			SubjectID:   scenario.ScenarioID,
			FailedCheck: failedCheck,
			TmpDir:      tmpDir,
			UserAnswers: userAnswers,
			Iteration:   iteration,
		}

		result, err := testFixGenerator.Generate(ctx, params)
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

func (c *USValidationCommand) refineTestFixPrompt(
	ctx context.Context,
	scenario *template.TestGenerationData,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	existingAnswers map[string]string,
	refinementIteration int,
	testFixGenerator *validate.FixPromptGenerator,
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
		Iteration:   refinementIteration + maxClarificationIterations,
	}

	result, err := testFixGenerator.Generate(ctx, params)
	if err != nil {
		return "", existingAnswers, pkgerrors.ErrFixPromptGenerationFailed(err)
	}

	if !result.HasFixPrompt() {
		return "", existingAnswers, nil
	}

	return result.FixPrompt, existingAnswers, nil
}

func (c *USValidationCommand) applyTestFix(
	ctx context.Context,
	scenario *template.TestGenerationData,
	fixPrompt string,
	tmpDir string,
	versionMgr *fs.ContentVersionManager,
	iteration int,
	testFixApplier *validate.FixApplier,
) (bool, error) {
	content, err := testFixApplier.Apply(
		ctx, scenario, scenario.ScenarioID, fixPrompt, tmpDir, iteration,
	)
	if err != nil {
		return false, fmt.Errorf("failed to apply fix: %w", err)
	}

	_, err = versionMgr.SaveNextVersion([]byte(content))
	if err != nil {
		return false, pkgerrors.ErrSaveStoryVersionFailed(err)
	}

	if scenario.TestFilePath != "" {
		writeErr := writeTestFile(scenario.TestFilePath, content)
		if writeErr != nil {
			return false, fmt.Errorf("failed to write test file: %w", writeErr)
		}

		console.Printf("\nTest file written to: %s\n", scenario.TestFilePath)
	}

	scenario.ArchitectureContent = content

	console.Printf("Fix applied. Saved as version %d.\n", versionMgr.GetCurrentVersion())
	console.Println("Re-running validation...")

	return true, nil
}

func (c *USValidationCommand) loadTestContent(scenario *template.TestGenerationData) {
	if scenario.ArchitectureContent != "" {
		return
	}

	if scenario.TestFilePath == "" {
		return
	}

	data, err := os.ReadFile(scenario.TestFilePath)
	if err != nil {
		return
	}

	scenario.ArchitectureContent = string(data)
}

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
