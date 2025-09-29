package prompts

import (
	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/infrastructure/config"
)

type HeuristicPromptBuilder struct {
	engine *TemplateEngine
	config *config.ViperConfig
}

func NewHeuristicPromptBuilder(engine *TemplateEngine, config *config.ViperConfig) *HeuristicPromptBuilder {
	return &HeuristicPromptBuilder{engine: engine, config: config}
}

func (b *HeuristicPromptBuilder) Build(threadCtx models.ThreadContext) (string, error) {
	templatePath := b.config.GetString("templates.prompts.heuristic")
	checklistPath := b.config.GetString("templates.checklists.triage_heuristic")
	return b.engine.BuildFromTemplate(
		threadCtx,
		templatePath,
		checklistPath,
	)
}
