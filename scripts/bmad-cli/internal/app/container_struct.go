package app

import (
	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/adapters/github"
	"bmad-cli/internal/application/commands"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/shell"
	"bmad-cli/internal/infrastructure/story"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

type Container struct {
	Config         *config.ViperConfig
	PRTriageCmd    *commands.PRTriageCommand
	USCreateCmd    *commands.USCreateCommand
	USImplementCmd *commands.USImplementCommand
	RunDir         *fs.RunDirectory
}

func NewContainer() (*Container, error) {
	cfg, err := config.NewViperConfig()
	if err != nil {
		return nil, pkgerrors.ErrInitializeConfigFailed(err)
	}

	configureLogging()

	// Create run directory once for entire CLI execution
	runDir, err := fs.NewRunDirectory(cfg.GetString("paths.tmp_dir"))
	if err != nil {
		return nil, pkgerrors.ErrCreateRunDirectoryFailed(err)
	}

	shellExec := shell.NewCommandRunner()

	githubService := github.NewGitHubService(shellExec)

	// Setup user story creation dependencies
	epicLoader := epic.NewEpicLoader(cfg)

	// Setup architecture document loader
	architectureLoader := docs.NewArchitectureLoader(cfg)

	// Setup AI task generation - required for operation
	claudeClient, err := ai.NewClaudeClient()
	if err != nil {
		return nil, pkgerrors.ErrCreateAIClientFailed(err)
	}

	// Setup user story creation command - required for operation
	usCreateCmd := createUSCreateCommand(epicLoader, claudeClient, cfg, architectureLoader, runDir)

	// Setup PR triage command - required for operation
	prTriageCmd := createPRTriageCommand(githubService, claudeClient, cfg)

	// Setup user story implement command
	gitService := git.NewGitService(shellExec)
	branchManager := git.NewBranchManager(gitService)
	storyLoader := story.NewStoryLoader(cfg)
	usImplementCmd := commands.NewUSImplementCommand(branchManager, storyLoader, claudeClient, cfg, runDir, shellExec)

	return &Container{
		Config:         cfg,
		PRTriageCmd:    prTriageCmd,
		USCreateCmd:    usCreateCmd,
		USImplementCmd: usImplementCmd,
		RunDir:         runDir,
	}, nil
}
