package services

import (
	"bmad-cli/internal/common/ai"
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
	return ai.NewAIGenerator[TaskPromptData, []story.Task](ctx, g.aiClient, g.config, storyDoc.Story.ID, "tasks").
		WithData(func() (TaskPromptData, error) {
			return TaskPromptData{
				Story: &storyDoc.Story,
				Docs:  storyDoc.ArchitectureDocs,
			}, nil
		}).
		WithPrompt(func(data TaskPromptData) (systemPrompt string, userPrompt string, err error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.tasks_system")
			systemLoader := template.NewTemplateLoader[TaskPromptData](systemTemplatePath)
			systemPrompt, err = systemLoader.LoadTemplate(TaskPromptData{})
			if err != nil {
				return "", "", fmt.Errorf("failed to load tasks system prompt: %w", err)
			}

			// Load user prompt
			templatePath := g.config.GetString("templates.prompts.tasks")
			userLoader := template.NewTemplateLoader[TaskPromptData](templatePath)
			userPrompt, err = userLoader.LoadTemplate(data)
			if err != nil {
				return "", "", fmt.Errorf("failed to load tasks user prompt: %w", err)
			}

			return systemPrompt, userPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[[]story.Task](g.config, storyDoc.Story.ID, "tasks", "tasks")).
		WithValidator(func(tasks []story.Task) error {
			if len(tasks) == 0 {
				return fmt.Errorf("AI generated no tasks")
			}
			return nil
		}).
		Generate()
}
