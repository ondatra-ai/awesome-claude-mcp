package services

import (
	"context"
	"fmt"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
)

// Type aliases to work around Go generics type resolution issues
type DevNotesType = story.DevNotes

// DevNotesPromptData represents data needed for dev notes generation prompts
type DevNotesPromptData struct {
	Story *story.Story
	Tasks []story.Task
	Docs  map[string]docs.ArchitectureDoc
}

// DevNotesTemplateLoader defines the interface for loading dev notes templates
type DevNotesTemplateLoader interface {
	LoadDevNotesPromptTemplate(story *story.Story, tasks []story.Task, architectureDocs map[string]docs.ArchitectureDoc) (string, error)
}

// AIDevNotesGenerator generates story dev_notes using AI based on templates
type AIDevNotesGenerator struct {
	aiClient       AIClient
	templateLoader DevNotesTemplateLoader
}

// NewDevNotesGenerator creates a new AIDevNotesGenerator instance
func NewDevNotesGenerator(aiClient AIClient, templateLoader DevNotesTemplateLoader) *AIDevNotesGenerator {
	return &AIDevNotesGenerator{
		aiClient:       aiClient,
		templateLoader: templateLoader,
	}
}

// GenerateDevNotes generates story dev_notes using AI based on the story, tasks, and architecture documents
func (g *AIDevNotesGenerator) GenerateDevNotes(ctx context.Context, story *story.Story, tasks []story.Task, architectureDocs map[string]docs.ArchitectureDoc) (*DevNotesType, error) {
	return NewAIGenerator[DevNotesPromptData, *DevNotesType](ctx, g.aiClient, story.ID, "devnotes").
		WithData(func() (DevNotesPromptData, error) {
			return DevNotesPromptData{Story: story, Tasks: tasks, Docs: architectureDocs}, nil
		}).
		WithPrompt(func(data DevNotesPromptData) (string, error) {
			return g.templateLoader.LoadDevNotesPromptTemplate(data.Story, data.Tasks, data.Docs)
		}).
		WithResponseParser(CreateYAMLFileParser[*DevNotesType](story.ID, "devnotes", "dev_notes")).
		WithValidator(g.validateDevNotes).
		Generate()
}


// validateDevNotes validates that mandatory entities have required source and description fields
func (g *AIDevNotesGenerator) validateDevNotes(devNotes *DevNotesType) error {
	mandatoryEntities := []string{"technology_stack", "architecture", "file_structure"}

	for _, entityName := range mandatoryEntities {
		entity, exists := (*devNotes)[entityName]
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
