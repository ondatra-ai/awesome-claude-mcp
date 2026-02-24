package bootstrap

import (
	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/adapters/github"
	"bmad-cli/internal/app/commands"
	"bmad-cli/internal/app/prompt_builders"
	"bmad-cli/internal/infrastructure/config"
)

func createPRTriageCommand(
	githubService *github.GitHubService,
	claudeClient *ai.ClaudeClient,
	cfg *config.ViperConfig,
) *commands.PRTriageCommand {
	// Create prompt dependencies
	templateEngine := prompt_builders.NewTemplateEngine()
	yamlParser := prompt_builders.NewYAMLParser()
	heuristicBuilder := prompt_builders.NewHeuristicPromptBuilder(templateEngine, cfg)
	implementationBuilder := prompt_builders.NewImplementationPromptBuilder(templateEngine, cfg)
	modeFactory := ai.NewModeFactory(cfg)

	// Create thread processor with all AI-related dependencies
	threadProcessor := ai.NewThreadProcessor(
		claudeClient,
		heuristicBuilder,
		implementationBuilder,
		yamlParser,
		modeFactory,
	)

	return commands.NewPRTriageCommand(
		githubService,
		threadProcessor,
		cfg,
	)
}
