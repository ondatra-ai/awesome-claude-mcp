package services

import (
	"context"
	"fmt"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
)

// Type aliases to work around Go generics type resolution issues
type StoryType = story.Story
type TaskType = story.Task

// TaskPromptData represents data needed for task generation prompts
type TaskPromptData struct {
	Story *StoryType
	Docs  map[string]docs.ArchitectureDoc
}

// TemplateLoader defines the interface for loading templates
type TemplateLoader interface {
	LoadTaskPromptTemplate(story *story.Story, architectureDocs map[string]docs.ArchitectureDoc) (string, error)
}

// AITaskGenerator generates story tasks using AI based on templates
type AITaskGenerator struct {
	aiClient       AIClient
	templateLoader TemplateLoader
}

// NewTaskGenerator creates a new AITaskGenerator instance
func NewTaskGenerator(aiClient AIClient, templateLoader TemplateLoader) *AITaskGenerator {
	return &AITaskGenerator{
		aiClient:       aiClient,
		templateLoader: templateLoader,
	}
}

// GenerateTasks generates story tasks using AI based on the story and architecture documents
func (g *AITaskGenerator) GenerateTasks(ctx context.Context, story *story.Story, architectureDocs map[string]docs.ArchitectureDoc) ([]TaskType, error) {
	return NewAIGenerator[TaskPromptData, []TaskType](ctx, g.aiClient, story.ID, "tasks").
		WithData(func() (TaskPromptData, error) {
			return TaskPromptData{Story: story, Docs: architectureDocs}, nil
		}).
		WithPrompt(func(data TaskPromptData) (string, error) {
			return g.templateLoader.LoadTaskPromptTemplate(data.Story, data.Docs)
		}).
		WithResponseParser(CreateYAMLFileParser[[]TaskType](story.ID, "tasks", "tasks")).
		WithValidator(func(tasks []TaskType) error {
			if len(tasks) == 0 {
				return fmt.Errorf("AI generated no tasks")
			}
			return nil
		}).
		Generate()
}
