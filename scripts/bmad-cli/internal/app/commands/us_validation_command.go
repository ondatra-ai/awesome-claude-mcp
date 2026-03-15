package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/app/generators/validate"
	checklistmodels "bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/checklist"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/input"
	storyinfra "bmad-cli/internal/infrastructure/story"
	"bmad-cli/internal/pkg/console"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const (
	maxClarificationIterations = 5
	maxRefinementIterations    = 3
	separatorWidth             = 80
	storyFilePermissions       = 0o644
	storyDirPermissions        = 0o755
)

// StageConfig defines the configuration for a stage-specific validation command.
type StageConfig struct {
	StageID       string // Stage to validate (e.g., "story_creation", "refinement", "ready_gate")
	RequiredStage string // Stage the story must be at to proceed ("" = no check, for us create)
	NextStage     string // Stage to set when all checks pass
	LoadFromEpic  bool   // true: load from epic (create), false: load from story file (refine/ready)
	StageName     string // Human-readable name (e.g., "Story Creation")
	CommandName   string // CLI command name for error messages (e.g., "us create")
}

// USValidationCommand validates user stories against stage-specific validation checklists.
type USValidationCommand struct {
	epicLoader         *epic.EpicLoader
	storyLoader        *storyinfra.StoryLoader
	checklistLoader    *checklist.ChecklistLoader
	checklistEvaluator *validate.ChecklistEvaluator
	fixPromptGenerator *validate.FixPromptGenerator
	fixApplier         *validate.FixApplier
	userInputCollector *input.UserInputCollector
	tableRenderer      *TableRenderer
	runDir             *fs.RunDirectory
	storiesDir         string
}

// NewUSValidationCommand creates a new stage-aware validation command.
func NewUSValidationCommand(
	epicLoader *epic.EpicLoader,
	storyLoader *storyinfra.StoryLoader,
	checklistLoader *checklist.ChecklistLoader,
	evaluator *validate.ChecklistEvaluator,
	fixPromptGen *validate.FixPromptGenerator,
	fixApplier *validate.FixApplier,
	inputCollector *input.UserInputCollector,
	renderer *TableRenderer,
	runDir *fs.RunDirectory,
	storiesDir string,
) *USValidationCommand {
	return &USValidationCommand{
		epicLoader:         epicLoader,
		storyLoader:        storyLoader,
		checklistLoader:    checklistLoader,
		checklistEvaluator: evaluator,
		fixPromptGenerator: fixPromptGen,
		fixApplier:         fixApplier,
		userInputCollector: inputCollector,
		tableRenderer:      renderer,
		runDir:             runDir,
		storiesDir:         storiesDir,
	}
}

// validationContext holds the context for a validation run.
type validationContext struct {
	versionMgr  *fs.StoryVersionManager
	prompts     []checklistmodels.PromptWithContext
	tmpDir      string
	iteration   int
	stageConfig StageConfig
	storyNumber string
}

// Execute runs stage-specific validation for the specified story.
func (c *USValidationCommand) Execute(
	ctx context.Context,
	storyNumber string,
	fix bool,
	config StageConfig,
) error {
	valCtx, err := c.initializeValidation(storyNumber, config)
	if err != nil {
		return err
	}

	return c.runValidationLoop(ctx, valCtx, fix)
}

// AdvanceStage checks the required stage and advances to the next stage without running validation prompts.
// Used for stages with no automated checks (e.g., architecture).
func (c *USValidationCommand) AdvanceStage(storyNumber string, config StageConfig) error {
	err := c.validateStoryNumber(storyNumber)
	if err != nil {
		return fmt.Errorf("invalid story number: %w", err)
	}

	storyData, err := c.loadStoryFromFile(storyNumber)
	if err != nil {
		return err
	}

	err = c.checkRequiredStage(storyData, config)
	if err != nil {
		return err
	}

	console.Header(fmt.Sprintf("%s — Story %s", strings.ToUpper(config.StageName), storyNumber), separatorWidth)
	console.Println("No automated checks defined for this stage.")

	storyData.Stage = config.NextStage

	storyPath, err := c.updateStoryFile(storyNumber, storyData)
	if err != nil {
		return fmt.Errorf("failed to update story file: %w", err)
	}

	console.Printf("Stage advanced to %q. Story saved to: %s\n", config.NextStage, storyPath)

	return nil
}

// initializeValidation sets up the validation context.
func (c *USValidationCommand) initializeValidation(
	storyNumber string,
	config StageConfig,
) (*validationContext, error) {
	err := c.validateStoryNumber(storyNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid story number: %w", err)
	}

	slog.Info("Starting validation",
		"story", storyNumber,
		"stage", config.StageID,
		"stageName", config.StageName,
	)

	console.Header(fmt.Sprintf("%s VALIDATION — Story %s", strings.ToUpper(config.StageName), storyNumber), separatorWidth)

	// Load story based on config
	var originalStory *story.Story

	if config.LoadFromEpic {
		originalStory, err = c.loadFromEpic(storyNumber)
	} else {
		originalStory, err = c.loadStoryFromFile(storyNumber)
	}

	if err != nil {
		return nil, err
	}

	slog.Info("Story loaded", "id", originalStory.ID, "title", originalStory.Title)

	err = c.checkRequiredStage(originalStory, config)
	if err != nil {
		return nil, err
	}

	tmpDir := c.runDir.GetTmpOutPath()
	versionMgr := fs.NewStoryVersionManager(c.runDir, storyNumber)

	err = versionMgr.SaveInitialVersion(originalStory)
	if err != nil {
		return nil, fmt.Errorf("failed to save initial story version: %w", err)
	}

	// Load checklist and extract stage-specific prompts
	checklistData, err := c.checklistLoader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load checklist: %w", err)
	}

	prompts := c.checklistLoader.ExtractPromptsForStage(checklistData, config.StageID)
	slog.Info("Extracted prompts for stage", "stage", config.StageID, "count", len(prompts))

	if len(prompts) == 0 {
		console.Println("No validation prompts found for this stage.")

		return nil, pkgerrors.ErrNoPromptsForStageFailed(config.StageID)
	}

	return &validationContext{
		versionMgr:  versionMgr,
		prompts:     prompts,
		tmpDir:      tmpDir,
		iteration:   0,
		stageConfig: config,
		storyNumber: storyNumber,
	}, nil
}

// loadFromEpic loads a story from its epic.
func (c *USValidationCommand) loadFromEpic(storyNumber string) (*story.Story, error) {
	console.Header("LOADING STORY FROM EPIC", separatorWidth)

	originalStory, err := c.epicLoader.LoadStoryFromEpic(storyNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to load story: %w", pkgerrors.ErrLoadStoryFromEpicFailed(err))
	}

	c.displayStory(originalStory, "STORY FROM EPIC")

	return originalStory, nil
}

// loadStoryFromFile loads a story from the docs/stories/ directory.
func (c *USValidationCommand) loadStoryFromFile(storyNumber string) (*story.Story, error) {
	doc, err := c.storyLoader.Load(storyNumber)
	if err != nil {
		return nil, fmt.Errorf("story file not found — run `bmad-cli us create %s` first: %w", storyNumber, err)
	}

	return &doc.Story, nil
}

// checkRequiredStage verifies the story is at the required stage to proceed.
func (c *USValidationCommand) checkRequiredStage(storyData *story.Story, config StageConfig) error {
	if config.RequiredStage == "" {
		return nil
	}

	if storyData.Stage != config.RequiredStage {
		return pkgerrors.ErrStageMismatchError(storyData.ID, storyData.Stage, config.CommandName, config.RequiredStage)
	}

	return nil
}

// runValidationLoop executes the main validation loop.
func (c *USValidationCommand) runValidationLoop(
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
func (c *USValidationCommand) runSingleIteration(
	ctx context.Context,
	valCtx *validationContext,
	fix bool,
) (bool, error) {
	currentStory, err := valCtx.versionMgr.LoadLatest()
	if err != nil {
		return false, fmt.Errorf("failed to load story version: %w", err)
	}

	acCount := len(currentStory.AcceptanceCriteria)

	var report *checklistmodels.ChecklistReport

	if fix {
		report, err = c.checklistEvaluator.EvaluateUntilFailure(
			ctx, currentStory, currentStory.ID, currentStory.Title, acCount, valCtx.prompts, valCtx.tmpDir)
	} else {
		report, err = c.checklistEvaluator.Evaluate(
			ctx, currentStory, currentStory.ID, currentStory.Title, acCount, valCtx.prompts, valCtx.tmpDir)
	}

	if err != nil {
		return false, fmt.Errorf("failed to evaluate checklist: %w", err)
	}

	c.tableRenderer.RenderReport(report, fix)

	if report.AllPassed() {
		c.handleAllPassed(valCtx)

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

	return c.runFixPromptLoop(ctx, valCtx, currentStory, *failedCheck)
}

// runFixPromptLoop handles the fix prompt generation and refinement loop.
func (c *USValidationCommand) runFixPromptLoop(
	ctx context.Context,
	valCtx *validationContext,
	currentStory *story.Story,
	failedCheck checklistmodels.ValidationResult,
) (bool, error) {
	userAnswers := make(map[string]string)
	refinementCount := 0

	fixPrompt, answers, err := c.generateFixPromptWithAnswers(ctx, currentStory, failedCheck, valCtx.tmpDir, userAnswers)
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
func (c *USValidationCommand) refineFixPrompt(
	ctx context.Context,
	currentStory *story.Story,
	failedCheck checklistmodels.ValidationResult,
	tmpDir string,
	existingAnswers map[string]string,
	refinementIteration int,
) (string, map[string]string, error) {
	feedback := c.userInputCollector.AskRefinementFeedback()
	if feedback == "" {
		return "", existingAnswers, nil
	}

	slog.Info("User provided refinement feedback",
		"promptIndex", failedCheck.PromptIndex,
		"refinementIteration", refinementIteration,
		"feedbackLength", len(feedback),
	)

	existingAnswers["_user_refinement"] = feedback

	params := validate.GenerateParams{
		Subject:     currentStory,
		SubjectID:   currentStory.ID,
		FailedCheck: failedCheck,
		TmpDir:      tmpDir,
		UserAnswers: existingAnswers,
		Iteration:   refinementIteration + maxClarificationIterations,
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
func (c *USValidationCommand) generateFixPromptWithAnswers(
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
			Subject:     storyData,
			SubjectID:   storyData.ID,
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
func (c *USValidationCommand) applyFix(
	ctx context.Context,
	valCtx *validationContext,
	currentStory *story.Story,
	fixPrompt string,
) (bool, error) {
	content, err := c.fixApplier.Apply(ctx, currentStory, currentStory.ID, fixPrompt, valCtx.tmpDir, valCtx.iteration)
	if err != nil {
		return false, fmt.Errorf("failed to apply fix: %w", err)
	}

	// Parse the returned content back into acceptance criteria
	var updatedACs []story.AcceptanceCriterion

	err = yaml.Unmarshal([]byte(content), &updatedACs)
	if err != nil {
		return false, fmt.Errorf("failed to parse updated acceptance criteria: %w", err)
	}

	updatedStory := *currentStory
	updatedStory.AcceptanceCriteria = updatedACs

	_, err = valCtx.versionMgr.SaveNextVersion(&updatedStory)
	if err != nil {
		return false, pkgerrors.ErrSaveStoryVersionFailed(err)
	}

	console.Printf("\nFix applied. Saved as version %d.\n", valCtx.versionMgr.GetCurrentVersion())
	console.Println("Re-running validation...")

	return true, nil
}

// handleAllPassed handles the case when all checks have passed.
func (c *USValidationCommand) handleAllPassed(valCtx *validationContext) {
	console.Header("ALL CHECKS PASSED!", separatorWidth)
	console.Printf("Latest version: %s\n", valCtx.versionMgr.GetLatestPath())

	c.displayFinalStory(valCtx.versionMgr)

	// Write story to docs/stories/
	latestStory, err := valCtx.versionMgr.LoadLatest()
	if err != nil {
		slog.Warn("Could not load latest story for writing", "error", err)

		return
	}

	latestStory.Stage = valCtx.stageConfig.NextStage

	var storyPath string

	if valCtx.stageConfig.LoadFromEpic {
		storyPath, err = c.writeNewStoryFile(latestStory)
	} else {
		storyPath, err = c.updateStoryFile(valCtx.storyNumber, latestStory)
	}

	if err != nil {
		slog.Warn("Could not write story file", "error", err)
		console.Printf("Warning: Could not write story file: %v\n", err)

		return
	}

	console.Printf("Story saved to: %s\n", storyPath)
}

// writeNewStoryFile creates a new story file in docs/stories/.
func (c *USValidationCommand) writeNewStoryFile(storyData *story.Story) (string, error) {
	slug := slugify(storyData.Title)
	filename := fmt.Sprintf("%s-%s.yaml", storyData.ID, slug)
	filePath := filepath.Join(c.storiesDir, filename)

	// Ensure directory exists
	err := os.MkdirAll(c.storiesDir, storyDirPermissions)
	if err != nil {
		return "", pkgerrors.ErrWriteStoryFileFailed(err)
	}

	// Write only the story portion (no tasks/dev_notes/testing scaffolding)
	wrapper := struct {
		Story story.Story `yaml:"story"`
	}{Story: *storyData}

	data, err := yaml.Marshal(wrapper)
	if err != nil {
		return "", pkgerrors.ErrWriteStoryFileFailed(err)
	}

	err = os.WriteFile(filePath, data, storyFilePermissions)
	if err != nil {
		return "", pkgerrors.ErrWriteStoryFileFailed(err)
	}

	slog.Info("Story file created", "path", filePath)

	return filePath, nil
}

// updateStoryFile updates an existing story file in docs/stories/ with story-only format.
func (c *USValidationCommand) updateStoryFile(storyNumber string, updatedStory *story.Story) (string, error) {
	pattern := filepath.Join(c.storiesDir, storyNumber+"-*.yaml")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", pkgerrors.ErrWriteStoryFileFailed(err)
	}

	if len(matches) == 0 {
		return c.writeNewStoryFile(updatedStory)
	}

	filePath := matches[0]

	wrapper := struct {
		Story story.Story `yaml:"story"`
	}{Story: *updatedStory}

	data, err := yaml.Marshal(wrapper)
	if err != nil {
		return "", pkgerrors.ErrWriteStoryFileFailed(err)
	}

	err = os.WriteFile(filePath, data, storyFilePermissions)
	if err != nil {
		return "", pkgerrors.ErrWriteStoryFileFailed(err)
	}

	slog.Info("Story file updated", "path", filePath)

	return filePath, nil
}

// displayStory shows a story's content in the terminal.
func (c *USValidationCommand) displayStory(storyData *story.Story, header string) {
	console.BlankLine()
	console.Header(header, separatorWidth)

	yamlBytes, err := yaml.Marshal(storyData)
	if err != nil {
		slog.Warn("Could not marshal story to YAML", "error", err)

		return
	}

	console.Println(string(yamlBytes))
	console.Separator("=", separatorWidth)
}

// displayFinalStory loads and displays the final story version.
func (c *USValidationCommand) displayFinalStory(versionMgr *fs.StoryVersionManager) {
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
func (c *USValidationCommand) getFirstFailedCheck(
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
func (c *USValidationCommand) displayFailureInfo(failedCheck *checklistmodels.ValidationResult) {
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
func (c *USValidationCommand) displayFixPrompt(fixPrompt string) {
	console.BlankLine()
	console.Header("FIX PROMPT GENERATED", separatorWidth)
	console.Println(fixPrompt)
	console.Separator("=", separatorWidth)
}

// validateStoryNumber validates the story number format (X.Y).
func (c *USValidationCommand) validateStoryNumber(storyNumber string) error {
	matched, err := regexp.MatchString(`^\d+\.\d+$`, storyNumber)
	if err != nil {
		return fmt.Errorf("regex failed: %w", err)
	}

	if !matched {
		return pkgerrors.ErrInvalidStoryNumberFormat
	}

	return nil
}

// slugify converts a title string into a URL-friendly slug.
func slugify(title string) string {
	lower := strings.ToLower(title)

	// Replace non-alphanumeric characters with hyphens
	var builder strings.Builder

	for _, r := range lower {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune('-')
		}
	}

	slug := builder.String()

	// Collapse multiple hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}
