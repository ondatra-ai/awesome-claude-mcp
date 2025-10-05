package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/adapters/github"
	"bmad-cli/internal/application/commands"
	"bmad-cli/internal/application/factories"
	"bmad-cli/internal/application/prompt_builders"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/shell"
	"bmad-cli/internal/infrastructure/story"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/infrastructure/validation"
)

type Container struct {
	Config         *config.ViperConfig
	PRTriageCmd    *commands.PRTriageCommand
	USCreateCmd    *commands.USCreateCommand
	USImplementCmd *commands.USImplementCommand
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
	epicLoader := epic.NewEpicLoader(cfg)

	// Setup architecture document loader
	architectureLoader := docs.NewArchitectureLoader(cfg)

	// Setup AI task generation - required for operation
	claudeClient, err := ai.NewClaudeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %w", err)
	}

	// Setup user story creation command - required for operation
	usCreateCmd := createUSCreateCommand(epicLoader, claudeClient, cfg, architectureLoader)

	// Setup PR triage command - required for operation
	prTriageCmd := createPRTriageCommand(githubService, claudeClient, cfg)

	// Setup user story implement command
	gitService := git.NewGitService(shellExec)
	branchManager := git.NewBranchManager(gitService)
	storyLoader := story.NewStoryLoader(cfg)
	usImplementCmd := commands.NewUSImplementCommand(branchManager, storyLoader)

	return &Container{
		Config:         cfg,
		PRTriageCmd:    prTriageCmd,
		USCreateCmd:    usCreateCmd,
		USImplementCmd: usImplementCmd,
	}, nil
}

func createUSCreateCommand(epicLoader *epic.EpicLoader, claudeClient *ai.ClaudeClient, cfg *config.ViperConfig, architectureLoader *docs.ArchitectureLoader) *commands.USCreateCommand {
	storyFactory := factories.NewStoryFactory(epicLoader, claudeClient, cfg, architectureLoader)
	storyTemplateLoader := template.NewTemplateLoader[*template.FlattenedStoryData](cfg.GetString("templates.story.template"))
	yamaleValidator := validation.NewYamaleValidator(cfg.GetString("templates.story.schema"))
	return commands.NewUSCreateCommand(storyFactory, storyTemplateLoader, yamaleValidator)
}

func createPRTriageCommand(githubService *github.GitHubService, claudeClient *ai.ClaudeClient, cfg *config.ViperConfig) *commands.PRTriageCommand {
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

func configureLogging() {
	log.SetFlags(0)

	// Ensure tmp directory exists
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		fmt.Printf("Warning: failed to create tmp directory: %v\n", err)
	}

	// Open log file for JSON output (all levels)
	logFile, err := os.OpenFile("./tmp/bmad-cli.log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Warning: failed to open log file: %v\n", err)
		// Fallback to console only
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				if a.Key == slog.LevelKey {
					level := a.Value.String()
					switch level {
					case "INFO":
						return slog.String("", "ℹ️")
					case "WARN":
						return slog.String("", "⚠️")
					case "ERROR":
						return slog.String("", "❌")
					case "DEBUG":
						return slog.String("", "🐛")
					default:
						return slog.String("", level)
					}
				}
				return a
			},
		}
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
		return
	}

	// Create JSON handler for file (all levels including debug)
	fileOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // Log everything to file
	}
	fileHandler := slog.NewJSONHandler(logFile, fileOpts)

	// Create text handler for console (info and above only)
	consoleOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Only info, warn, error to console
	}
	consoleHandler := slog.NewTextHandler(os.Stdout, consoleOpts)

	// Create multi-handler that writes to both file and console
	multiHandler := &MultiHandler{
		fileHandler:    fileHandler,
		consoleHandler: consoleHandler,
	}

	slog.SetDefault(slog.New(multiHandler))
}

// MultiHandler writes to both file and console with different levels
type MultiHandler struct {
	fileHandler    slog.Handler
	consoleHandler slog.Handler
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.fileHandler.Enabled(ctx, level) || h.consoleHandler.Enabled(ctx, level)
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	// Always write to file
	if err := h.fileHandler.Handle(ctx, record); err != nil {
		// Don't fail if file write fails, continue to console
	}

	// Write to console only for info and above
	if record.Level >= slog.LevelInfo {
		if err := h.consoleHandler.Handle(ctx, record); err != nil {
			return err
		}
	}

	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MultiHandler{
		fileHandler:    h.fileHandler.WithAttrs(attrs),
		consoleHandler: h.consoleHandler.WithAttrs(attrs),
	}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	return &MultiHandler{
		fileHandler:    h.fileHandler.WithGroup(name),
		consoleHandler: h.consoleHandler.WithGroup(name),
	}
}
