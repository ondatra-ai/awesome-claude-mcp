package prompts

import "bmad-cli/internal/domain/models"

type ImplementationPromptBuilder struct {
	engine *TemplateEngine
}

func NewImplementationPromptBuilder(engine *TemplateEngine) *ImplementationPromptBuilder {
	return &ImplementationPromptBuilder{engine: engine}
}

func (b *ImplementationPromptBuilder) Build(threadCtx models.ThreadContext) (string, error) {
	return b.engine.BuildFromTemplate(
		threadCtx,
		"scripts/pr-triage/apply.prompt.tpl",
		"",
	)
}
