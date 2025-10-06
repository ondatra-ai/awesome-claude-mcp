package app

import (
	"fmt"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/adapters/github"
	"bmad-cli/internal/application/commands"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/shell"
	"bmad-cli/internal/infrastructure/story"
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
