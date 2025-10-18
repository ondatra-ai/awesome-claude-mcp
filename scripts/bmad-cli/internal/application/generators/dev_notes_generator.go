package generators

import (
	"context"
	"fmt"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/pkg/ai"
	"bmad-cli/internal/pkg/errors"
)

// DevNotesPromptData represents data needed for dev notes generation prompts.
type DevNotesPromptData struct {
	Story  *story.Story
	Tasks  []story.Task
	Docs   *docs.ArchitectureDocs
	TmpDir string // Path to run-specific tmp directory
}

// AIDevNotesGenerator generates story dev_notes using AI based on templates.
type AIDevNotesGenerator struct {
	aiClient ports.AIPort
	config   *config.ViperConfig
}

// NewDevNotesGenerator creates a new AIDevNotesGenerator instance.
func NewDevNotesGenerator(aiClient ports.AIPort, config *config.ViperConfig) *AIDevNotesGenerator {
	return &AIDevNotesGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateDevNotes generates story dev_notes using AI based on the story,
// tasks, and architecture documents.
func (g *AIDevNotesGenerator) GenerateDevNotes(
	ctx context.Context,
	storyDoc *story.StoryDocument,
	tmpDir string,
) (story.DevNotes, error) {
	generator := ai.NewAIGenerator[DevNotesPromptData, story.DevNotes](
		ctx,
		g.aiClient,
		g.config,
		storyDoc.Story.ID,
		"devnotes",
	)

	result, err := generator.
		WithTmpDir(tmpDir).
		WithData(func() (DevNotesPromptData, error) {
			return DevNotesPromptData{
				Story:  &storyDoc.Story,
				Tasks:  storyDoc.Tasks,
				Docs:   storyDoc.ArchitectureDocs,
				TmpDir: tmpDir,
			}, nil
		}).
		WithPrompt(func(data DevNotesPromptData) (string, string, error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.devnotes_system")
			systemLoader := template.NewTemplateLoader[DevNotesPromptData](systemTemplatePath)

			systemPrompt, err := systemLoader.LoadTemplate(DevNotesPromptData{})
			if err != nil {
				return "", "", errors.ErrLoadDevNotesPromptFailed(err)
			}

			// Load user prompt
			templatePath := g.config.GetString("templates.prompts.devnotes")
			userLoader := template.NewTemplateLoader[DevNotesPromptData](templatePath)

			userPrompt, err := userLoader.LoadTemplate(data)
			if err != nil {
				return "", "", errors.ErrLoadDevNotesUserPromptFailed(err)
			}

			return systemPrompt, userPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[story.DevNotes](
			g.config,
			storyDoc.Story.ID,
			"devnotes",
			"dev_notes",
			tmpDir,
		)).
		WithValidator(g.validateDevNotes).
		Generate(ctx)
	if err != nil {
		return story.DevNotes{}, fmt.Errorf("generate dev notes: %w", err)
	}

	return result, nil
}

// validateDevNotes validates that mandatory entities have required source and description fields.
func (g *AIDevNotesGenerator) validateDevNotes(devNotes story.DevNotes) error {
	mandatoryEntities := []string{"technology_stack", "architecture", "file_structure"}

	for _, entityName := range mandatoryEntities {
		entity, exists := devNotes[entityName]
		if !exists {
			return errors.ErrMandatoryEntityMissingError(entityName)
		}

		// Handle both map[string]interface{} and story.DevNotes (which is also map[string]interface{})
		var entityMap map[string]interface{}
		if em, ok := entity.(map[string]interface{}); ok {
			entityMap = em
		} else if dn, ok := entity.(story.DevNotes); ok {
			entityMap = dn
		} else {
			return errors.ErrEntityInvalidTypeError(entityName, entity)
		}

		// Check for mandatory source field
		if _, hasSource := entityMap["source"]; !hasSource {
			return errors.ErrEntityMissingFieldError(entityName, "source")
		}

		// Check for mandatory description field
		if _, hasDescription := entityMap["description"]; !hasDescription {
			return errors.ErrEntityMissingFieldError(entityName, "description")
		}
	}

	return nil
}
