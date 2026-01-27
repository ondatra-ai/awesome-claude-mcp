package commands

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/app/generators/validate"
	checklistmodels "bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/checklist"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/input"
	"bmad-cli/internal/pkg/console"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const (
	maxClarificationIterations = 5
	maxRefinementIterations    = 3
	separatorWidth             = 80
)

// USChecklistCommand validates user stories against the validation checklist.
type USChecklistCommand struct {
	epicLoader         *epic.EpicLoader
	checklistLoader    *checklist.ChecklistLoader
	checklistEvaluator *validate.ChecklistEvaluator
	fixPromptGenerator *validate.FixPromptGenerator
	fixApplier         *validate.FixApplier
	userInputCollector *input.UserInputCollector
	tableRenderer      *TableRenderer
	runDir             *fs.RunDirectory
}

// NewUSChecklistCommand creates a new checklist validation command.
func NewUSChecklistCommand(
	epicLoader *epic.EpicLoader,
	checklistLoader *checklist.ChecklistLoader,
	evaluator *validate.ChecklistEvaluator,
	fixPromptGen *validate.FixPromptGenerator,
	fixApplier *validate.FixApplier,
	inputCollector *input.UserInputCollector,
	renderer *TableRenderer,
	runDir *fs.RunDirectory,
) *USChecklistCommand {
	return &USChecklistCommand{
		epicLoader:         epicLoader,
		checklistLoader:    checklistLoader,
		checklistEvaluator: evaluator,
		fixPromptGenerator: fixPromptGen,
		fixApplier:         fixApplier,
		userInputCollector: inputCollector,
		tableRenderer:      renderer,
		runDir:             runDir,
	}
}

// validationContext holds the context for a validation run.
type validationContext struct {
	versionMgr *fs.StoryVersionManager
	prompts    []checklistmodels.PromptWithContext
	tmpDir     string
	iteration  int
}

// Execute runs the iterative checklist validation for the specified story.
// If fix is true, the command enters interactive fix mode when validation fails.
func (c *USChecklistCommand) Execute(ctx context.Context, storyNumber string, fix bool) error {
	valCtx, err := c.initializeValidation(storyNumber)
	if err != nil {
		return err
	}

	return c.runValidationLoop(ctx, valCtx, fix)
}

// initializeValidation sets up the validation context.
func (c *USChecklistCommand) initializeValidation(storyNumber string) (*validationContext, error) {
	err := c.validateStoryNumber(storyNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid story number: %w", err)
	}

	slog.Info("Starting checklist validation", "story", storyNumber)

	originalStory, err := c.epicLoader.LoadStoryFromEpic(storyNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to load story: %w", pkgerrors.ErrLoadStoryFromEpicFailed(err))
	}

	slog.Info("Story loaded", "id", originalStory.ID, "title", originalStory.Title)

	tmpDir := c.runDir.GetTmpOutPath()
	versionMgr := fs.NewStoryVersionManager(c.runDir, storyNumber)

	err = versionMgr.SaveInitialVersion(originalStory)
	if err != nil {
		return nil, fmt.Errorf("failed to save initial story version: %w", err)
	}

	checklistData, err := c.checklistLoader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load checklist: %w", err)
	}

	prompts := c.checklistLoader.ExtractAllPrompts(checklistData)
	slog.Info("Extracted prompts", "count", len(prompts))

	return &validationContext{
		versionMgr: versionMgr,
		prompts:    prompts,
		tmpDir:     tmpDir,
		iteration:  0,
	}, nil
}

// runValidationLoop executes the main validation loop.
//

func (c *USChecklistCommand) runValidationLoop(
	ctx context.Context,
	valCtx *validationContext,
	fix bool,
) error {
	for {
		valCtx.iteration++

		shouldContinue, err := c.runSingleIteration(ctx, valCtx, fix)
		if err != nil {
			return err
		}

		if !shouldContinue {
			return nil
		}
	}
}

// runSingleIteration runs one iteration of validation and returns whether to continue.
func (c *USChecklistCommand) runSingleIteration(
	ctx context.Context,
	valCtx *validationContext,
	fix bool,
) (bool, error) {
	currentStory, err := valCtx.versionMgr.LoadLatest()
	if err != nil {
		return false, fmt.Errorf("failed to load story version: %w", err)
	}

	var report *checklistmodels.ChecklistReport

	if fix {
		// In fix mode: stop at first failure for iterative fixing
		report, err = c.checklistEvaluator.EvaluateUntilFailure(ctx, currentStory, valCtx.prompts, valCtx.tmpDir)
	} else {
		// In report mode: evaluate ALL items for complete report
		report, err = c.checklistEvaluator.Evaluate(ctx, currentStory, valCtx.prompts, valCtx.tmpDir)
	}

	if err != nil {
		return false, fmt.Errorf("failed to evaluate checklist: %w", err)
	}

	c.tableRenderer.RenderReport(report, fix)

	if report.AllPassed() {
		c.handleAllPassed(valCtx.versionMgr, fix)

		return false, nil
	}

	failedCheck := c.getFirstFailedCheck(report)
	if failedCheck == nil {
		slog.Warn("No failed check found despite not all passed")

		return false, nil
	}

	c.displayFailureInfo(failedCheck)

	// Only enter fix loop if --fix flag is set
	if !fix {
		console.BlankLine()
		console.Println("Validation failed. Use --fix flag to enter interactive fix mode.")

		return false, nil
	}

	// Generate initial fix prompt and enter the fix prompt loop
	return c.runFixPromptLoop(ctx, valCtx, currentStory, *failedCheck)
}

// runFixPromptLoop handles the fix prompt generation and refinement loop.
//

func (c *USChecklistCommand) runFixPromptLoop(
	ctx context.Context,
	valCtx *validationContext,
	currentStory *story.Story,
	failedCheck checklistmodels.ValidationResult,
) (bool, error) {
	userAnswers := make(map[string]string)
	refinementCount := 0

	// Generate initial fix prompt
	fixPrompt, answers, err := c.generateFixPromptWithAnswers(ctx, currentStory, failedCheck, valCtx.tmpDir, userAnswers)
	if err != nil {
		return false, fmt.Errorf("failed to generate fix prompt: %w", err)
	}

	userAnswers = answers

	if fixPrompt == "" {
		slog.Warn("No fix prompt generated")

		return false, nil
	}

	// Fix prompt loop: display, ask action, handle refinement
	for {
		c.displayFixPrompt(fixPrompt)

		action := c.userInputCollector.AskApplyRefineOrExit()

		switch action {
		case input.ActionApply:
			return c.applyFix(ctx, valCtx, currentStory, fixPrompt)

		case input.ActionRefine:
			if refinementCount >= maxRefinementIterations {
				console.Printf("\nMax refinement attempts (%d) reached. Please apply or exit.\n", maxRefinementIterations)

				continue
			}

			refinementCount++

			newPrompt, updatedAnswers, refineErr := c.refineFixPrompt(
				ctx, currentStory, failedCheck, valCtx.tmpDir, userAnswers, refinementCount)
			if refineErr != nil {
				return false, pkgerrors.ErrFixPromptRefinementFailed(refineErr)
			}

			if newPrompt == "" {
				console.Println("\nNo feedback provided. Keeping current fix prompt.")

				continue
			}

			fixPrompt = newPrompt
			userAnswers = updatedAnswers

			console.Printf("\n(Refinement %d of %d)\n", refinementCount, maxRefinementIterations)

		case input.ActionExit:
			console.Printf("\nExiting. Latest version saved at: %s\n", valCtx.versionMgr.GetLatestPath())

			return false, nil
		}
	}
}

// refineFixPrompt collects user feedback and regenerates the fix prompt.
//

func (c *USChecklistCommand) refineFixPrompt(
	ctx context.Context,
	currentStory *story.Story,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	existingAnswers map[string]string,
	refinementIteration int,
) (string, map[string]string, error) {
	// Get user feedback
	feedback := c.userInputCollector.AskRefinementFeedback()
	if feedback == "" {
		return "", existingAnswers, nil
	}

	slog.Info("User provided refinement feedback",
		"promptIndex", failedCheck.PromptIndex,
		"refinementIteration", refinementIteration,
		"feedbackLength", len(feedback),
	)

	// Add feedback to answers map with special key
	existingAnswers["_user_refinement"] = feedback

	// Regenerate fix prompt with feedback
	params := validate.GenerateParams{
		StoryData:   currentStory,
		FailedCheck: failedCheck,
		TmpDir:      tmpDir,
		UserAnswers: existingAnswers,
		Iteration:   refinementIteration + maxClarificationIterations, // Offset for unique file names
	}

	result, err := c.fixPromptGenerator.Generate(ctx, params)
	if err != nil {
		return "", existingAnswers, pkgerrors.ErrFixPromptGenerationFailed(err)
	}

	if !result.HasFixPrompt() {
		slog.Warn("Refinement did not produce a fix prompt",
			"promptIndex", failedCheck.PromptIndex,
		)

		return "", existingAnswers, nil
	}

	slog.Info("Fix prompt refined successfully",
		"promptIndex", failedCheck.PromptIndex,
		"refinementIteration", refinementIteration,
	)

	return result.FixPrompt, existingAnswers, nil
}

// generateFixPromptWithAnswers generates fix prompt and returns accumulated user answers.
func (c *USChecklistCommand) generateFixPromptWithAnswers(
	ctx context.Context,
	storyData *story.Story,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	initialAnswers map[string]string,
) (string, map[string]string, error) {
	userAnswers := make(map[string]string)

	for id, answer := range initialAnswers {
		userAnswers[id] = answer
	}

	for iteration := 1; iteration <= maxClarificationIterations; iteration++ {
		params := validate.GenerateParams{
			StoryData:   storyData,
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
			slog.Info("Fix prompt generated successfully",
				"promptIndex", failedCheck.PromptIndex,
				"iterations", iteration,
			)

			return result.FixPrompt, userAnswers, nil
		}

		if !result.HasQuestions() {
			slog.Warn("No fix prompt or questions returned",
				"promptIndex", failedCheck.PromptIndex,
			)

			return "", userAnswers, nil
		}

		slog.Info("Requesting user clarification",
			"promptIndex", failedCheck.PromptIndex,
			"questionCount", len(result.Questions),
			"iteration", iteration,
		)

		answers := c.userInputCollector.AskQuestions(result.Questions)

		for id, answer := range answers {
			userAnswers[id] = answer
		}
	}

	slog.Warn("Max clarification iterations reached",
		"promptIndex", failedCheck.PromptIndex,
		"maxIterations", maxClarificationIterations,
	)

	return "", userAnswers, nil
}

// applyFix applies the fix prompt and saves a new version.
//

func (c *USChecklistCommand) applyFix(
	ctx context.Context,
	valCtx *validationContext,
	currentStory *story.Story,
	fixPrompt string,
) (bool, error) {
	updatedStory, err := c.fixApplier.Apply(ctx, currentStory, fixPrompt, valCtx.tmpDir, valCtx.iteration)
	if err != nil {
		return false, fmt.Errorf("failed to apply fix: %w", err)
	}

	_, err = valCtx.versionMgr.SaveNextVersion(updatedStory)
	if err != nil {
		return false, pkgerrors.ErrSaveStoryVersionFailed(err)
	}

	console.Printf("\nFix applied. Saved as version %d.\n", valCtx.versionMgr.GetCurrentVersion())
	console.Println("Re-running validation...")

	return true, nil
}

// handleAllPassed handles the case when all checks have passed.
func (c *USChecklistCommand) handleAllPassed(versionMgr *fs.StoryVersionManager, fix bool) {
	console.Header("ALL CHECKS PASSED!", separatorWidth)
	console.Printf("Latest version: %s\n", versionMgr.GetLatestPath())

	// Display the final story content
	c.displayFinalStory(versionMgr)

	// Only ask about copying if we were in fix mode (changes might have been made)
	if fix {
		if c.userInputCollector.AskCopyToOriginal() {
			// Copy to original not yet available - show manual instructions
			console.Println("\nCopy to original not yet available.")
			console.Printf("Please manually copy from: %s\n", versionMgr.GetLatestPath())
		}
	}
}

// displayFinalStory loads and displays the final story version.
func (c *USChecklistCommand) displayFinalStory(versionMgr *fs.StoryVersionManager) {
	storyData, err := versionMgr.LoadLatest()
	if err != nil {
		slog.Warn("Could not load final story for display", "error", err)

		return
	}

	console.BlankLine()
	console.Header("FINAL STORY VERSION", separatorWidth)

	yamlBytes, err := yaml.Marshal(storyData)
	if err != nil {
		slog.Warn("Could not marshal story to YAML", "error", err)

		return
	}

	console.Println(string(yamlBytes))
	console.Separator("=", separatorWidth)
}

// getFirstFailedCheck returns the first failed check from the report.
func (c *USChecklistCommand) getFirstFailedCheck(
	report *checklistmodels.ChecklistReport,
) *checklistmodels.ValidationResult {
	for _, result := range report.Results {
		if result.Status == checklistmodels.StatusFail {
			return &result
		}
	}

	return nil
}

// displayFailureInfo displays information about the failed check.
func (c *USChecklistCommand) displayFailureInfo(failedCheck *checklistmodels.ValidationResult) {
	console.BlankLine()
	console.Separator("=", separatorWidth)
	console.Printf("CHECK FAILED: %s\n", failedCheck.SectionPath)
	console.Separator("=", separatorWidth)
	console.Printf("Question: %s\n", failedCheck.Question)
	console.Printf("Expected: %s\n", failedCheck.ExpectedAnswer)
	console.Printf("Actual: %s\n", failedCheck.ActualAnswer)

	if failedCheck.Rationale != "" {
		console.Printf("Rationale: %s\n", failedCheck.Rationale)
	}
}

// displayFixPrompt displays the generated fix prompt.
func (c *USChecklistCommand) displayFixPrompt(fixPrompt string) {
	console.BlankLine()
	console.Header("FIX PROMPT GENERATED", separatorWidth)
	console.Println(fixPrompt)
	console.Separator("=", separatorWidth)
}

// validateStoryNumber validates the story number format (X.Y).
func (c *USChecklistCommand) validateStoryNumber(storyNumber string) error {
	matched, err := regexp.MatchString(`^\d+\.\d+$`, storyNumber)
	if err != nil {
		return fmt.Errorf("regex failed: %w", err)
	}

	if !matched {
		return pkgerrors.ErrInvalidStoryNumberFormat
	}

	return nil
}
