package prompts

import "bmad-cli/internal/domain/models"

type HeuristicPromptBuilder struct {
	engine *TemplateEngine
}

func NewHeuristicPromptBuilder(engine *TemplateEngine) *HeuristicPromptBuilder {
	return &HeuristicPromptBuilder{engine: engine}
}

func (b *HeuristicPromptBuilder) Build(threadCtx models.ThreadContext) (string, error) {
	return b.engine.BuildFromTemplate(
		threadCtx,
		"scripts/pr-triage/heuristic.prompt.tpl",
		".bmad-core/checklists/triage-heuristic-checklist.md",
	)
}
