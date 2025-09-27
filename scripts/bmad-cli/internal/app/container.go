package app

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/adapters/github"
	"bmad-cli/internal/application/commands"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/domain/services"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/shell"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/infrastructure/validation"
)

type Container struct {
	Config      *config.ViperConfig
	PRTriageCmd *commands.PRTriageCommand
	USCreateCmd *commands.USCreateCommand
}

func NewContainer() (*Container, error) {
	cfg, err := config.NewViperConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	configureLogging()

	shellExec := shell.NewCommandRunner()

	githubService := github.NewGitHubService(shellExec)

	// Setup user story creation dependencies
	epicLoader := epic.NewEpicLoader()

	// Setup architecture document loader
	architectureLoader := docs.NewArchitectureLoader(cfg)


	// Setup AI task generation - required for operation
	claudeClient, err := ai.NewClaudeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %w", err)
	}
	taskGenerator := services.NewTaskGenerator(claudeClient, cfg)
	devNotesGenerator := services.NewDevNotesGenerator(claudeClient, cfg)
	qaAssessmentGenerator := services.NewQAAssessmentGenerator(claudeClient)
	testingGenerator := services.NewTestingGenerator(claudeClient)

	storyFactory := services.NewStoryFactory(epicLoader, taskGenerator, devNotesGenerator, qaAssessmentGenerator, testingGenerator, architectureLoader)

	storyTemplateLoader := template.NewTemplateLoader[*template.FlattenedStoryData]("templates/story.yaml.tpl")
	yamaleValidator := validation.NewYamaleValidator("templates/story-schema.yaml")
	usCreateCmd := commands.NewUSCreateCommand(storyFactory, storyTemplateLoader, yamaleValidator)

	// AI service and PR triage are optional - only create if needed
	var aiService ports.AIService
	var prTriageCmd *commands.PRTriageCommand

	// Try to create AI service, but don't fail if it's not available
	if aiSvc, err := ai.NewAIService(cfg); err == nil {
		aiService = aiSvc
		orchestrator := services.NewPRTriageOrchestrator(githubService, aiService)
		prTriageCmd = commands.NewPRTriageCommand(orchestrator)
	}

	return &Container{
		Config:      cfg,
		PRTriageCmd: prTriageCmd,
		USCreateCmd: usCreateCmd,
	}, nil
}

func configureLogging() {
	log.SetFlags(0)
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			return a
		},
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
}
