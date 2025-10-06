package generators

import (
	"bmad-cli/internal/domain/ports"
	"context"
	"fmt"

	"bmad-cli/internal/pkg/ai"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// DevNotesPromptData represents data needed for dev notes generation prompts
type DevNotesPromptData struct {
	Story *story.Story
	Tasks []story.Task
	Docs  *docs.ArchitectureDocs
}

// AIDevNotesGenerator generates story dev_notes using AI based on templates
type AIDevNotesGenerator struct {
	aiClient ports.AIPort
	config   *config.ViperConfig
}

// NewDevNotesGenerator creates a new AIDevNotesGenerator instance
func NewDevNotesGenerator(aiClient ports.AIPort, config *config.ViperConfig) *AIDevNotesGenerator {
	return &AIDevNotesGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateDevNotes generates story dev_notes using AI based on the story, tasks, and architecture documents
func (g *AIDevNotesGenerator) GenerateDevNotes(ctx context.Context, storyDoc *story.StoryDocument) (story.DevNotes, error) {
	return ai.NewAIGenerator[DevNotesPromptData, story.DevNotes](ctx, g.aiClient, g.config, storyDoc.Story.ID, "devnotes").
		WithData(func() (DevNotesPromptData, error) {
			return DevNotesPromptData{
				Story: &storyDoc.Story,
				Tasks: storyDoc.Tasks,
				Docs:  storyDoc.ArchitectureDocs,
			}, nil
		}).
		WithPrompt(func(data DevNotesPromptData) (systemPrompt string, userPrompt string, err error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.devnotes_system")
			systemLoader := template.NewTemplateLoader[DevNotesPromptData](systemTemplatePath)
			systemPrompt, err = systemLoader.LoadTemplate(DevNotesPromptData{})
			if err != nil {
				return "", "", fmt.Errorf("failed to load devnotes system prompt: %w", err)
			}

			// Load user prompt
			templatePath := g.config.GetString("templates.prompts.devnotes")
			userLoader := template.NewTemplateLoader[DevNotesPromptData](templatePath)
			userPrompt, err = userLoader.LoadTemplate(data)
			if err != nil {
				return "", "", fmt.Errorf("failed to load devnotes user prompt: %w", err)
			}

			return systemPrompt, userPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[story.DevNotes](g.config, storyDoc.Story.ID, "devnotes", "dev_notes")).
		WithValidator(g.validateDevNotes).
		Generate()
}

// validateDevNotes validates that mandatory entities have required source and description fields
func (g *AIDevNotesGenerator) validateDevNotes(devNotes story.DevNotes) error {
	mandatoryEntities := []string{"technology_stack", "architecture", "file_structure"}

	for _, entityName := range mandatoryEntities {
		entity, exists := devNotes[entityName]
		if !exists {
			return fmt.Errorf("mandatory entity '%s' is missing", entityName)
		}

		// Handle both map[string]interface{} and story.DevNotes (which is also map[string]interface{})
		var entityMap map[string]interface{}
		if em, ok := entity.(map[string]interface{}); ok {
			entityMap = em
		} else if dn, ok := entity.(story.DevNotes); ok {
			entityMap = dn
		} else {
			return fmt.Errorf("entity '%s' must be a map, got %T", entityName, entity)
		}

		// Check for mandatory source field
		if _, hasSource := entityMap["source"]; !hasSource {
			return fmt.Errorf("entity '%s' is missing mandatory 'source' field", entityName)
		}

		// Check for mandatory description field
		if _, hasDescription := entityMap["description"]; !hasDescription {
			return fmt.Errorf("entity '%s' is missing mandatory 'description' field", entityName)
		}
	}

	return nil
}
