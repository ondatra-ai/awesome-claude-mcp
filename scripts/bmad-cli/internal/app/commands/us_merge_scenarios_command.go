package commands

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"bmad-cli/internal/app/generators/implement"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/story"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// File permission constants.
const (
	mergeCmdFileModeDirectory = 0o755
	mergeCmdFileModeReadWrite = 0o644
)

// USMergeScenariosCommand merges story scenarios into the requirements file.
type USMergeScenariosCommand struct {
	storyLoader       *story.StoryLoader
	mergeScenariosGen *implement.MergeScenariosGenerator
	runDir            *fs.RunDirectory
}

// NewUSMergeScenariosCommand creates a new USMergeScenariosCommand.
func NewUSMergeScenariosCommand(
	storyLoader *story.StoryLoader,
	mergeScenariosGen *implement.MergeScenariosGenerator,
	runDir *fs.RunDirectory,
) *USMergeScenariosCommand {
	return &USMergeScenariosCommand{
		storyLoader:       storyLoader,
		mergeScenariosGen: mergeScenariosGen,
		runDir:            runDir,
	}
}

// Execute merges scenarios from the story into the requirements file.
func (c *USMergeScenariosCommand) Execute(ctx context.Context, storyNumber string) error {
	slog.Info("Starting merge scenarios", "story_number", storyNumber)

	// Validate story exists
	_, err := c.storyLoader.GetStorySlug(storyNumber)
	if err != nil {
		return pkgerrors.ErrGetStorySlugFailed(err)
	}

	tmpDir := c.runDir.GetTmpOutPath()
	outputFile := filepath.Join(tmpDir, "requirements-merged.yaml")

	err = c.cloneRequirements(outputFile)
	if err != nil {
		return pkgerrors.ErrCloneRequirementsFileFailed(err)
	}

	storyDoc, err := c.storyLoader.Load(storyNumber)
	if err != nil {
		return pkgerrors.ErrLoadStoryFailed(err)
	}

	status, err := c.mergeScenariosGen.MergeScenarios(ctx, storyDoc, outputFile, tmpDir)
	if err != nil {
		return pkgerrors.ErrMergeScenariosFailed(err)
	}

	slog.Info("✅ Scenario merge completed",
		"scenarios", status.ItemsProcessed,
		"story", storyNumber,
	)

	err = c.replaceRequirements(outputFile)
	if err != nil {
		return pkgerrors.ErrReplaceRequirementsFailed(err)
	}

	return nil
}

func (c *USMergeScenariosCommand) cloneRequirements(outputFile string) error {
	slog.Info("Cloning requirements file for safe testing", "output", outputFile)

	outputDir := filepath.Dir(outputFile)

	err := os.MkdirAll(outputDir, mergeCmdFileModeDirectory)
	if err != nil {
		return pkgerrors.ErrCreateOutputDirectoryFailed(outputDir, err)
	}

	data, err := os.ReadFile("docs/requirements.yaml")
	if err != nil {
		return pkgerrors.ErrReadRequirementsFileFailed(err)
	}

	err = os.WriteFile(outputFile, data, mergeCmdFileModeReadWrite)
	if err != nil {
		return pkgerrors.ErrWriteOutputFileFailed(outputFile, err)
	}

	slog.Info("✓ Cloned requirements file", "destination", outputFile)

	return nil
}

func (c *USMergeScenariosCommand) replaceRequirements(mergedFile string) error {
	const (
		requirementsPath = "docs/requirements.yaml"
		backupPath       = "docs/requirements.yaml.backup"
	)

	slog.Info("Replacing requirements file", "source", mergedFile)

	originalData, err := os.ReadFile(requirementsPath)
	if err != nil {
		return pkgerrors.ErrReadOriginalFileFailed(err)
	}

	err = os.WriteFile(backupPath, originalData, mergeCmdFileModeReadWrite)
	if err != nil {
		return pkgerrors.ErrCreateBackupFileFailed(err)
	}

	slog.Info("✓ Created backup", "path", backupPath)

	mergedData, err := os.ReadFile(mergedFile)
	if err != nil {
		return pkgerrors.ErrReadMergedFileFailed(err)
	}

	err = os.WriteFile(requirementsPath, mergedData, mergeCmdFileModeReadWrite)
	if err != nil {
		_ = os.WriteFile(requirementsPath, originalData, mergeCmdFileModeReadWrite) // Restore

		return pkgerrors.ErrReplaceFileFailed(err)
	}

	slog.Info("✓ Replaced requirements with merged scenarios", "path", requirementsPath)

	return nil
}
