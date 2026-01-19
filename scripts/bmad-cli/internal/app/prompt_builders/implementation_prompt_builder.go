package prompt_builders

import (
	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/infrastructure/config"
)

type ImplementationPromptBuilder struct {
	engine *TemplateEngine
	config *config.ViperConfig
}

func NewImplementationPromptBuilder(engine *TemplateEngine, config *config.ViperConfig) *ImplementationPromptBuilder {
	return &ImplementationPromptBuilder{engine: engine, config: config}
}

func (b *ImplementationPromptBuilder) Build(threadCtx models.ThreadContext) (string, error) {
	templatePath := b.config.GetString("templates.prompts.apply")

	return b.engine.BuildFromTemplate(
		threadCtx,
		templatePath,
		"",
	)
}
