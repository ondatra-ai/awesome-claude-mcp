package services

import (
	"context"
	"fmt"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// DevNotesPromptData represents data needed for dev notes generation prompts
type DevNotesPromptData struct {
	Story *story.Story
	Tasks []story.Task
	Docs  map[string]docs.ArchitectureDoc
}

// AIDevNotesGenerator generates story dev_notes using AI based on templates
type AIDevNotesGenerator struct {
	aiClient       AIClient
	templateLoader *template.PromptLoader[DevNotesPromptData]
}

// NewDevNotesGenerator creates a new AIDevNotesGenerator instance
func NewDevNotesGenerator(aiClient AIClient, templateLoader *template.PromptLoader[DevNotesPromptData]) *AIDevNotesGenerator {
	return &AIDevNotesGenerator{
		aiClient:       aiClient,
		templateLoader: templateLoader,
	}
}

// GenerateDevNotes generates story dev_notes using AI based on the story, tasks, and architecture documents
func (g *AIDevNotesGenerator) GenerateDevNotes(ctx context.Context, storyObj *story.Story, tasks []story.Task, architectureDocs map[string]docs.ArchitectureDoc) (story.DevNotes, error) {
	return NewAIGenerator[DevNotesPromptData, story.DevNotes](ctx, g.aiClient, storyObj.ID, "devnotes").
		WithData(func() (DevNotesPromptData, error) {
			return DevNotesPromptData{Story: storyObj, Tasks: tasks, Docs: architectureDocs}, nil
		}).
		WithPrompt(func(data DevNotesPromptData) (string, error) {
			return g.templateLoader.LoadPromptTemplate(data)
		}).
		WithResponseParser(CreateYAMLFileParser[story.DevNotes](storyObj.ID, "devnotes", "dev_notes")).
		WithValidator(g.validateDevNotes).
		Generate()
}

// NewDevNotesPromptLoader creates a new dev notes prompt loader
func NewDevNotesPromptLoader(templateFilePath string) *template.PromptLoader[DevNotesPromptData] {
	return template.NewPromptLoader[DevNotesPromptData](templateFilePath)
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
