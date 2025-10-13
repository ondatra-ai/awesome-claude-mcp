package generators

import (
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/pkg/ai"
	"context"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// TaskPromptData represents data needed for task generation prompts.
type TaskPromptData struct {
	Story  *story.Story
	Docs   *docs.ArchitectureDocs
	TmpDir string // Path to run-specific tmp directory
}

// AITaskGenerator generates story tasks using AI based on templates.
type AITaskGenerator struct {
	aiClient ports.AIPort
	config   *config.ViperConfig
}

// NewTaskGenerator creates a new AITaskGenerator instance.
func NewTaskGenerator(aiClient ports.AIPort, config *config.ViperConfig) *AITaskGenerator {
	return &AITaskGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateTasks generates story tasks using AI based on the story and
// architecture documents.
func (g *AITaskGenerator) GenerateTasks(
	ctx context.Context,
	storyDoc *story.StoryDocument,
	tmpDir string,
) ([]story.Task, error) {
	generator := ai.NewAIGenerator[TaskPromptData, []story.Task](
		ctx,
		g.aiClient,
		g.config,
		storyDoc.Story.ID,
		"tasks",
	)

	return generator.
		WithTmpDir(tmpDir).
		WithData(func() (TaskPromptData, error) {
			return TaskPromptData{
				Story:  &storyDoc.Story,
				Docs:   storyDoc.ArchitectureDocs,
				TmpDir: tmpDir,
			}, nil
		}).
		WithPrompt(func(data TaskPromptData) (string, string, error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.tasks_system")
			systemLoader := template.NewTemplateLoader[TaskPromptData](systemTemplatePath)

			sysPrompt, err := systemLoader.LoadTemplate(TaskPromptData{})
			if err != nil {
				return "", "", pkgerrors.ErrLoadTasksSystemPromptFailed(err)
			}

			// Load user prompt
			templatePath := g.config.GetString("templates.prompts.tasks")
			userLoader := template.NewTemplateLoader[TaskPromptData](templatePath)

			usrPrompt, err := userLoader.LoadTemplate(data)
			if err != nil {
				return "", "", pkgerrors.ErrLoadTasksUserPromptFailed(err)
			}

			return sysPrompt, usrPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[[]story.Task](
			g.config,
			storyDoc.Story.ID,
			"tasks",
			"tasks",
			tmpDir,
		)).
		WithValidator(func(tasks []story.Task) error {
			if len(tasks) == 0 {
				return pkgerrors.ErrAIGeneratedNoTasks
			}

			return nil
		}).
		Generate(ctx)
}
