package app

import (
	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/adapters/github"
	"bmad-cli/internal/application/commands"
	"bmad-cli/internal/application/factories"
	"bmad-cli/internal/application/prompt_builders"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/infrastructure/validation"
)

func createUSCreateCommand(
	epicLoader *epic.EpicLoader,
	claudeClient *ai.ClaudeClient,
	cfg *config.ViperConfig,
	architectureLoader *docs.ArchitectureLoader,
	runDir *fs.RunDirectory,
) *commands.USCreateCommand {
	storyFactory := factories.NewStoryFactory(
		epicLoader, claudeClient, cfg, architectureLoader, runDir,
	)
	storyTemplateLoader := template.NewTemplateLoader[*template.FlattenedStoryData](
		cfg.GetString("templates.story.template"),
	)
	yamaleValidator := validation.NewYamaleValidator(cfg.GetString("templates.story.schema"))

	return commands.NewUSCreateCommand(storyFactory, storyTemplateLoader, yamaleValidator)
}

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
