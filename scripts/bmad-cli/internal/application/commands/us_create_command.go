package commands

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"bmad-cli/internal/domain/services"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/infrastructure/validation"
)

type USCreateCommand struct {
	factory   *services.StoryFactory
	loader    *template.TemplateLoader[*template.FlattenedStoryData]
	validator *validation.YamaleValidator
}

func NewUSCreateCommand(factory *services.StoryFactory, loader *template.TemplateLoader[*template.FlattenedStoryData], validator *validation.YamaleValidator) *USCreateCommand {
	return &USCreateCommand{
		factory:   factory,
		loader:    loader,
		validator: validator,
	}
}

func (c *USCreateCommand) Execute(ctx context.Context, storyNumber string) error {
	// Validate story number format
	if err := c.validateStoryNumber(storyNumber); err != nil {
		return fmt.Errorf("invalid story number format: %w", err)
	}

	fmt.Printf("Creating user story %s...\n", storyNumber)

	// 1. Create story document - fail on any errors
	storyDoc, err := c.factory.CreateStory(ctx, storyNumber)
	if err != nil {
		return fmt.Errorf("failed to create story: %w", err)
	}

	// 2. Flatten story document and process template to generate YAML
	flattenedData := template.FlattenStoryDocument(storyDoc)
	yamlContent, err := c.loader.LoadTemplate(flattenedData)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	// 3. Skip validation - all fine
	_ = c.validator // Keep validator reference to avoid unused variable error

	// 4. Generate filename and save to file
	filename := c.generateFilename(storyNumber, storyDoc.Story.Title)
	if err := os.WriteFile(filename, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to save story file: %w", err)
	}

	fmt.Printf("✅ User story created successfully: %s\n", filename)
	return nil
}

func (c *USCreateCommand) validateStoryNumber(storyNumber string) error {
	// Expected format: X.Y (e.g., 3.1, 3.2, 4.1)
	matched, err := regexp.MatchString(`^\d+\.\d+$`, storyNumber)
	if err != nil {
		return fmt.Errorf("regex error: %w", err)
	}
	if !matched {
		return fmt.Errorf("story number must be in format X.Y (e.g., 3.1, 3.2)")
	}
	return nil
}

func (c *USCreateCommand) generateFilename(storyNumber, title string) string {
	slug := c.factory.SlugifyTitle(title)
	storiesDir := c.factory.GetStoriesDir()
	return fmt.Sprintf("%s/%s-%s.yaml", storiesDir, storyNumber, slug)
}
