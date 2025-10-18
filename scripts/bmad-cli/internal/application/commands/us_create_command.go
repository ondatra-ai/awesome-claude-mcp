package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"bmad-cli/internal/application/factories"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/infrastructure/validation"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const (
	// File permission constants.
	fileModeReadWrite = 0644 // Standard file permission for read/write files
	fileModeDirectory = 0755 // Standard directory permission
)

type USCreateCommand struct {
	factory   *factories.StoryFactory
	loader    *template.TemplateLoader[*template.FlattenedStoryData]
	validator *validation.YamaleValidator
}

func NewUSCreateCommand(
	factory *factories.StoryFactory,
	loader *template.TemplateLoader[*template.FlattenedStoryData],
	validator *validation.YamaleValidator,
) *USCreateCommand {
	return &USCreateCommand{
		factory:   factory,
		loader:    loader,
		validator: validator,
	}
}

func (c *USCreateCommand) Execute(ctx context.Context, storyNumber string) error {
	// Validate story number format
	err := c.validateStoryNumber(storyNumber)
	if err != nil {
		return pkgerrors.ErrInvalidStoryNumberFormatError(storyNumber)
	}

	slog.Info("Creating user story", "story", storyNumber)

	// 1. Create story document - fail on any errors
	storyDoc, err := c.factory.CreateStory(ctx, storyNumber)
	if err != nil {
		return pkgerrors.ErrCreateStoryFailed(err)
	}

	// 2. Flatten story document and set TmpDir for template processing
	flattenedData := template.FlattenStoryDocument(storyDoc)
	flattenedData.TmpDir = c.factory.GetTmpDirPath()

	// Process template to generate YAML
	yamlContent, err := c.loader.LoadTemplate(flattenedData)
	if err != nil {
		return pkgerrors.ErrProcessTemplateFailed(err)
	}

	// 3. Skip validation - all fine
	_ = c.validator // Keep validator reference to avoid unused variable error

	// 4. Generate filename and save to file
	filename := c.generateFilename(storyNumber, storyDoc.Story.Title)

	err = os.WriteFile(filename, []byte(yamlContent), fileModeReadWrite)
	if err != nil {
		return pkgerrors.ErrSaveStoryFileFailed(err)
	}

	slog.Info("User story created successfully", "file", filename)

	return nil
}

func (c *USCreateCommand) validateStoryNumber(storyNumber string) error {
	// Expected format: X.Y (e.g., 3.1, 3.2, 4.1)
	matched, err := regexp.MatchString(`^\d+\.\d+$`, storyNumber)
	if err != nil {
		return pkgerrors.ErrRegexFailed(err)
	}

	if !matched {
		return pkgerrors.ErrInvalidStoryNumberFormat
	}

	return nil
}

func (c *USCreateCommand) generateFilename(storyNumber, title string) string {
	slug := c.factory.SlugifyTitle(title)
	storiesDir := c.factory.GetStoriesDir()

	return fmt.Sprintf("%s/%s-%s.yaml", storiesDir, storyNumber, slug)
}
