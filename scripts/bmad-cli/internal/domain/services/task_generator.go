package services

import (
	"context"
	"fmt"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// TaskPromptData represents data needed for task generation prompts
type TaskPromptData struct {
	Story *story.Story
	Docs  map[string]docs.ArchitectureDoc
}

// AITaskGenerator generates story tasks using AI based on templates
type AITaskGenerator struct {
	aiClient       AIClient
	templateLoader *template.PromptLoader[TaskPromptData]
}

// NewTaskGenerator creates a new AITaskGenerator instance
func NewTaskGenerator(aiClient AIClient, templateLoader *template.PromptLoader[TaskPromptData]) *AITaskGenerator {
	return &AITaskGenerator{
		aiClient:       aiClient,
		templateLoader: templateLoader,
	}
}

// GenerateTasks generates story tasks using AI based on the story and architecture documents
func (g *AITaskGenerator) GenerateTasks(ctx context.Context, storyObj *story.Story, architectureDocs map[string]docs.ArchitectureDoc) ([]story.Task, error) {
	return NewAIGenerator[TaskPromptData, []story.Task](ctx, g.aiClient, storyObj.ID, "tasks").
		WithData(func() (TaskPromptData, error) {
			return TaskPromptData{Story: storyObj, Docs: architectureDocs}, nil
		}).
		WithPrompt(func(data TaskPromptData) (string, error) {
			return g.templateLoader.LoadPromptTemplate(data)
		}).
		WithResponseParser(CreateYAMLFileParser[[]story.Task](storyObj.ID, "tasks", "tasks")).
		WithValidator(func(tasks []story.Task) error {
			if len(tasks) == 0 {
				return fmt.Errorf("AI generated no tasks")
			}
			return nil
		}).
		Generate()
}

// NewTaskPromptLoader creates a new task prompt loader
func NewTaskPromptLoader(templateFilePath string) *template.PromptLoader[TaskPromptData] {
	return template.NewPromptLoader[TaskPromptData](templateFilePath)
}
