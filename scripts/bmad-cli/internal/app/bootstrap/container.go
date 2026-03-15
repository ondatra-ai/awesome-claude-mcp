package bootstrap

import (
	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/adapters/github"
	"bmad-cli/internal/app/commands"
	"bmad-cli/internal/app/factories"
	"bmad-cli/internal/app/generators/implement"
	"bmad-cli/internal/app/generators/validate"
	"bmad-cli/internal/infrastructure/checklist"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/input"
	"bmad-cli/internal/infrastructure/shell"
	"bmad-cli/internal/infrastructure/story"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

type Container struct {
	Config              *config.ViperConfig
	PRTriageCmd         *commands.PRTriageCommand
	USImplementCmd      *commands.USImplementCommand
	USMergeScenariosCmd *commands.USMergeScenariosCommand
	USValidationCmd     *commands.USValidationCommand
	ReqGenerateTestsCmd *commands.ReqGenerateTestsCommand
	RunDir              *fs.RunDirectory
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

	epicLoader := epic.NewEpicLoader(cfg)

	// Setup AI client - required for operation
	claudeClient, err := ai.NewClaudeClient()
	if err != nil {
		return nil, pkgerrors.ErrCreateAIClientFailed(err)
	}

	// Setup PR triage command - required for operation
	prTriageCmd := createPRTriageCommand(githubService, claudeClient, cfg)

	// Setup user story commands
	gitService := git.NewGitService(shellExec)
	branchManager := git.NewBranchManager(gitService)
	storyLoader := story.NewStoryLoader(cfg)

	mergeScenariosGen := implement.NewMergeScenariosGenerator(claudeClient, cfg)
	usValidateCmd := commands.NewUSValidateCommand(storyLoader)
	usMergeScenariosCmd := commands.NewUSMergeScenariosCommand(
		storyLoader, mergeScenariosGen, runDir,
	)

	implementFactory := factories.NewImplementFactory(
		branchManager,
		storyLoader,
		claudeClient,
		cfg,
		runDir,
		shellExec,
		usValidateCmd,
		usMergeScenariosCmd,
	)
	usImplementCmd := commands.NewUSImplementCommand(implementFactory)

	// Setup requirements commands
	testCodeGen := implement.NewTestCodeGenerator(claudeClient, cfg)
	reqGenerateTestsCmd := commands.NewReqGenerateTestsCommand(testCodeGen, runDir)

	// Setup user story validation command (replaces checklist command)
	checklistLoader := checklist.NewChecklistLoader(cfg)
	checklistEvaluator := validate.NewChecklistEvaluator(claudeClient, cfg)
	fixPromptGenerator := validate.NewFixPromptGenerator(claudeClient, cfg)
	fixApplier := validate.NewFixApplier(claudeClient, cfg)
	userInputCollector := input.NewUserInputCollector()
	tableRenderer := commands.NewTableRenderer()
	storiesDir := cfg.GetString("paths.stories_dir")
	usValidationCmd := commands.NewUSValidationCommand(
		epicLoader,
		storyLoader,
		checklistLoader,
		checklistEvaluator,
		fixPromptGenerator,
		fixApplier,
		userInputCollector,
		tableRenderer,
		runDir,
		storiesDir,
	)

	return &Container{
		Config:              cfg,
		PRTriageCmd:         prTriageCmd,
		USImplementCmd:      usImplementCmd,
		USMergeScenariosCmd: usMergeScenariosCmd,
		USValidationCmd:     usValidationCmd,
		ReqGenerateTestsCmd: reqGenerateTestsCmd,
		RunDir:              runDir,
	}, nil
}
