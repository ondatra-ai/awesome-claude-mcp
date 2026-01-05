package commands

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"bmad-cli/internal/app/generators/validate"
	checklistmodels "bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/checklist"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// USChecklistCommand validates user stories against the validation checklist.
type USChecklistCommand struct {
	epicLoader         *epic.EpicLoader
	checklistLoader    *checklist.ChecklistLoader
	checklistEvaluator *validate.ChecklistEvaluator
	fixPromptGenerator *validate.FixPromptGenerator
	tableRenderer      *TableRenderer
	runDir             *fs.RunDirectory
}

// NewUSChecklistCommand creates a new checklist validation command.
func NewUSChecklistCommand(
	epicLoader *epic.EpicLoader,
	checklistLoader *checklist.ChecklistLoader,
	evaluator *validate.ChecklistEvaluator,
	fixPromptGen *validate.FixPromptGenerator,
	renderer *TableRenderer,
	runDir *fs.RunDirectory,
) *USChecklistCommand {
	return &USChecklistCommand{
		epicLoader:         epicLoader,
		checklistLoader:    checklistLoader,
		checklistEvaluator: evaluator,
		fixPromptGenerator: fixPromptGen,
		tableRenderer:      renderer,
		runDir:             runDir,
	}
}

// Execute runs the checklist validation for the specified story.
func (c *USChecklistCommand) Execute(ctx context.Context, storyNumber string) error {
	// Validate story number format
	err := c.validateStoryNumber(storyNumber)
	if err != nil {
		return fmt.Errorf("invalid story number: %w", err)
	}

	slog.Info("Starting checklist validation", "story", storyNumber)

	// 1. Load story from epic file
	storyData, err := c.epicLoader.LoadStoryFromEpic(storyNumber)
	if err != nil {
		return fmt.Errorf("failed to load story: %w", pkgerrors.ErrLoadStoryFromEpicFailed(err))
	}

	slog.Info("Story loaded", "id", storyData.ID, "title", storyData.Title)

	// 2. Load and parse checklist
	checklistData, err := c.checklistLoader.Load()
	if err != nil {
		return fmt.Errorf("failed to load checklist: %w", err)
	}

	// 3. Extract all prompts (excluding skipped ones)
	prompts := c.checklistLoader.ExtractAllPrompts(checklistData)
	slog.Info("Extracted prompts", "count", len(prompts))

	// 4. Get run-specific tmp directory for prompt files
	tmpDir := c.runDir.GetTmpOutPath()

	// 5. Evaluate all prompts using AI
	report, err := c.checklistEvaluator.Evaluate(ctx, storyData, prompts, tmpDir)
	if err != nil {
		return fmt.Errorf("failed to evaluate checklist: %w", err)
	}

	// 5. Render results as table
	c.tableRenderer.RenderReport(report)

	// 6. Generate fix prompts for each failed validation
	c.generateFixPrompts(ctx, storyData, report, tmpDir)

	return nil
}

// generateFixPrompts generates a fix prompt for each failed validation.
func (c *USChecklistCommand) generateFixPrompts(
	ctx context.Context,
	storyData *story.Story,
	report *checklistmodels.ChecklistReport,
	tmpDir string,
) {
	// Collect failed checks with fix prompts
	failedChecks := c.collectFailedChecks(report)
	if len(failedChecks) == 0 {
		return
	}

	slog.Info("Generating fix prompts", "failedChecks", len(failedChecks))

	// Generate a fix prompt for each failed check using its original prompt index
	for _, failedCheck := range failedChecks {
		params := validate.GenerateParams{
			StoryData:   storyData,
			FailedCheck: failedCheck,
			TmpDir:      tmpDir,
		}

		_, err := c.fixPromptGenerator.Generate(ctx, params)
		if err != nil {
			slog.Error("Failed to generate fix prompt",
				"promptIndex", failedCheck.PromptIndex,
				"section", failedCheck.SectionPath,
				"error", err,
			)

			continue
		}
	}
}

// collectFailedChecks returns validation results that failed and have fix prompts.
func (c *USChecklistCommand) collectFailedChecks(
	report *checklistmodels.ChecklistReport,
) []checklistmodels.ValidationResult {
	var failedChecks []checklistmodels.ValidationResult

	for _, result := range report.Results {
		if result.Status == checklistmodels.StatusFail && result.FixPrompt != "" {
			failedChecks = append(failedChecks, result)
		}
	}

	return failedChecks
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
