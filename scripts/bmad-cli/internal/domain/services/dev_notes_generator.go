package services

import (
	"context"
	"fmt"

	"bmad-cli/internal/common/utils"
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

// DevNotesTemplateLoader defines the interface for loading dev notes templates
type DevNotesTemplateLoader interface {
	LoadPromptTemplate(data DevNotesPromptData) (string, error)
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

// NewDevNotesPromptLoader creates a new dev notes prompt loader with inline template builder
func NewDevNotesPromptLoader(templateFilePath string) DevNotesTemplateLoader {
	return template.NewPromptLoader(templateFilePath, func(data DevNotesPromptData) (map[string]interface{}, error) {
		storyYAML, err := utils.MarshalToYAML(data.Story)
		if err != nil {
			return nil, fmt.Errorf("failed to convert story to YAML: %w", err)
		}

		tasksYAML, err := utils.MarshalWithWrapper(data.Tasks, "tasks")
		if err != nil {
			return nil, fmt.Errorf("failed to convert tasks to YAML: %w", err)
		}

		return map[string]interface{}{
			"StoryYAML":                storyYAML,
			"StoryID":                  data.Story.ID,
			"TasksYAML":                tasksYAML,
			"Architecture":             data.Docs["Architecture"].Content,
			"FrontendArchitecture":     data.Docs["FrontendArchitecture"].Content,
			"CodingStandards":          data.Docs["CodingStandards"].Content,
			"SourceTree":               data.Docs["SourceTree"].Content,
			"TechStack":                data.Docs["TechStack"].Content,
			"ArchitecturePath":         data.Docs["Architecture"].FilePath,
			"FrontendArchitecturePath": data.Docs["FrontendArchitecture"].FilePath,
			"CodingStandardsPath":      data.Docs["CodingStandards"].FilePath,
			"SourceTreePath":           data.Docs["SourceTree"].FilePath,
			"TechStackPath":            data.Docs["TechStack"].FilePath,
		}, nil
	})
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
