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
	"bmad-cli/internal/infrastructure/logging"
	"bmad-cli/internal/infrastructure/shell"
)

type Container struct {
	Config      config.ConfigProvider
	Logger      logging.Logger
	Shell       shell.Executor
	GitHub      ports.GitHubService
	AI          ports.AIService
	PRTriageCmd *commands.PRTriageCommand
}

func NewContainer() (*Container, error) {
	cfg := config.NewViperConfig()

	configureLogging()
	logger := logging.NewSlogLogger()

	shellExec := shell.NewCommandRunner()

	githubService := github.NewGitHubService(shellExec)

	aiService, err := ai.NewAIService(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI service: %w", err)
	}

	orchestrator := services.NewPRTriageOrchestrator(githubService, aiService, logger)
	prTriageCmd := commands.NewPRTriageCommand(orchestrator)

	return &Container{
		Config:      cfg,
		Logger:      logger,
		Shell:       shellExec,
		GitHub:      githubService,
		AI:          aiService,
		PRTriageCmd: prTriageCmd,
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
