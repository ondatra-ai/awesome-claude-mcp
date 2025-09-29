package ai

import (
	"context"

	"bmad-cli/internal/application/prompts"
	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/infrastructure/config"
)

type AIService struct {
	analyzer    *ThreadAnalyzer
	implementer *ThreadImplementer
}

func NewAIService(config *config.ViperConfig) (*AIService, error) {
	client, err := NewClaudeClient()
	if err != nil {
		return nil, err
	}

	templateEngine := prompts.NewTemplateEngine()
	yamlParser := prompts.NewYAMLParser()

	heuristicBuilder := prompts.NewHeuristicPromptBuilder(templateEngine, config)
	implementationBuilder := prompts.NewImplementationPromptBuilder(templateEngine, config)

	analyzer := NewThreadAnalyzer(client, heuristicBuilder, yamlParser)
	implementer := NewThreadImplementer(client, implementationBuilder)

	return &AIService{
		analyzer:    analyzer,
		implementer: implementer,
	}, nil
}

func (s *AIService) AnalyzeThread(ctx context.Context, threadContext models.ThreadContext) (models.HeuristicAnalysisResult, error) {
	return s.analyzer.Analyze(ctx, threadContext)
}

func (s *AIService) ImplementChanges(ctx context.Context, threadContext models.ThreadContext) (string, error) {
	return s.implementer.Implement(ctx, threadContext)
}
