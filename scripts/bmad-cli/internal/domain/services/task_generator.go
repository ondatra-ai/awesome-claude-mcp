package services

import (
	"context"
	"fmt"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// TaskPromptData represents data needed for task generation prompts
type TaskPromptData struct {
	Story *story.Story
	Docs  *docs.ArchitectureDocs
}

// AITaskGenerator generates story tasks using AI based on templates
type AITaskGenerator struct {
	aiClient AIClient
	config   *config.ViperConfig
}

// NewTaskGenerator creates a new AITaskGenerator instance
func NewTaskGenerator(aiClient AIClient, config *config.ViperConfig) *AITaskGenerator {
	return &AITaskGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateTasks generates story tasks using AI based on the story and architecture documents
func (g *AITaskGenerator) GenerateTasks(ctx context.Context, storyDoc *story.StoryDocument) ([]story.Task, error) {
	return NewAIGenerator[TaskPromptData, []story.Task](ctx, g.aiClient, storyDoc.Story.ID, "tasks").
		WithData(func() (TaskPromptData, error) {
			return TaskPromptData{
				Story: &storyDoc.Story,
				Docs:  storyDoc.ArchitectureDocs,
			}, nil
		}).
		WithPrompt(func(data TaskPromptData) (string, error) {
			templatePath := g.config.GetString("templates.prompts.tasks")
			loader := template.NewTemplateLoader[TaskPromptData](templatePath)
			return loader.LoadTemplate(data)
		}).
		WithResponseParser(CreateYAMLFileParser[[]story.Task](storyDoc.Story.ID, "tasks", "tasks")).
		WithValidator(func(tasks []story.Task) error {
			if len(tasks) == 0 {
				return fmt.Errorf("AI generated no tasks")
			}
			return nil
		}).
		Generate()
}
