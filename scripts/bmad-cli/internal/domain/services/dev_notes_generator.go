package services

import (
	"context"
	"fmt"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
	"gopkg.in/yaml.v3"
)

// DevNotesPromptData represents data needed for dev notes generation prompts
type DevNotesPromptData struct {
	Story                *story.Story
	Tasks                []story.Task
	Docs                 map[string]docs.ArchitectureDoc
	StoryYAML            string
	Architecture         string
	FrontendArchitecture string
	CodingStandards      string
	SourceTree           string
	TechStack            string
}

// AIDevNotesGenerator generates story dev_notes using AI based on templates
type AIDevNotesGenerator struct {
	aiClient AIClient
}

// NewDevNotesGenerator creates a new AIDevNotesGenerator instance
func NewDevNotesGenerator(aiClient AIClient) *AIDevNotesGenerator {
	return &AIDevNotesGenerator{
		aiClient: aiClient,
	}
}

// GenerateDevNotes generates story dev_notes using AI based on the story, tasks, and architecture documents
func (g *AIDevNotesGenerator) GenerateDevNotes(ctx context.Context, storyObj *story.Story, tasks []story.Task, architectureDocs map[string]docs.ArchitectureDoc) (story.DevNotes, error) {
	return NewAIGenerator[DevNotesPromptData, story.DevNotes](ctx, g.aiClient, storyObj.ID, "devnotes").
		WithData(func() (DevNotesPromptData, error) {
			// Marshal story to YAML
			storyYAML, err := yaml.Marshal(storyObj)
			if err != nil {
				return DevNotesPromptData{}, fmt.Errorf("failed to marshal story to YAML: %w", err)
			}

			// Extract document content
			promptData := DevNotesPromptData{
				Story:     storyObj,
				Tasks:     tasks,
				Docs:      architectureDocs,
				StoryYAML: string(storyYAML),
			}

			// Populate architecture documents if available
			if doc, exists := architectureDocs["Architecture"]; exists {
				promptData.Architecture = doc.Content
			}
			if doc, exists := architectureDocs["FrontendArchitecture"]; exists {
				promptData.FrontendArchitecture = doc.Content
			}
			if doc, exists := architectureDocs["CodingStandards"]; exists {
				promptData.CodingStandards = doc.Content
			}
			if doc, exists := architectureDocs["SourceTree"]; exists {
				promptData.SourceTree = doc.Content
			}
			if doc, exists := architectureDocs["TechStack"]; exists {
				promptData.TechStack = doc.Content
			}

			return promptData, nil
		}).
		WithPrompt(func(data DevNotesPromptData) (string, error) {
			loader := template.NewTemplateLoader[DevNotesPromptData]("templates/us-create.devnotes.prompt.tpl")
			return loader.LoadTemplate(data)
		}).
		WithResponseParser(CreateYAMLFileParser[story.DevNotes](storyObj.ID, "devnotes", "dev_notes")).
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
