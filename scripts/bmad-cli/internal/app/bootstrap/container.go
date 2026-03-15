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
	"bmad-cli/internal/infrastructure/requirements"
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
	ReqValidationCmd    *commands.ReqValidationCommand
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

	// Setup user input collector (shared across commands)
	userInputCollector := input.NewUserInputCollector()

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

	usValidationCmd := createUSValidationCommand(
		epicLoader, storyLoader, claudeClient, cfg, userInputCollector, runDir,
	)

	reqValidationCmd := createReqValidationCommand(
		claudeClient, cfg, userInputCollector, runDir,
	)

	return &Container{
		Config:              cfg,
		PRTriageCmd:         prTriageCmd,
		USImplementCmd:      usImplementCmd,
		USMergeScenariosCmd: usMergeScenariosCmd,
		USValidationCmd:     usValidationCmd,
		ReqGenerateTestsCmd: reqGenerateTestsCmd,
		ReqValidationCmd:    reqValidationCmd,
		RunDir:              runDir,
	}, nil
}

func createUSValidationCommand(
	epicLoader *epic.EpicLoader,
	storyLoader *story.StoryLoader,
	claudeClient *ai.ClaudeClient,
	cfg *config.ViperConfig,
	userInputCollector *input.UserInputCollector,
	runDir *fs.RunDirectory,
) *commands.USValidationCommand {
	checklistLoader := checklist.NewChecklistLoader(cfg)
	checklistEvaluator := validate.NewChecklistEvaluator(claudeClient, cfg)
	fixPromptGenerator := validate.NewFixPromptGenerator(claudeClient, cfg)
	fixApplier := validate.NewFixApplier(claudeClient, cfg)
	tableRenderer := commands.NewTableRenderer()
	storiesDir := cfg.GetString("paths.stories_dir")

	return commands.NewUSValidationCommand(
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
}

func createReqValidationCommand(
	claudeClient *ai.ClaudeClient,
	cfg *config.ViperConfig,
	userInputCollector *input.UserInputCollector,
	runDir *fs.RunDirectory,
) *commands.ReqValidationCommand {
	testChecklistPath := cfg.GetString("paths.test_checklist")
	testChecklistLoader := checklist.NewChecklistLoaderWithPath(testChecklistPath)

	testEvaluator := validate.NewChecklistEvaluatorWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_checklist_system"),
		cfg.GetString("templates.prompts.test_checklist"),
	)

	testFixGenerator := validate.NewFixPromptGeneratorWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_fix_generator_system"),
		cfg.GetString("templates.prompts.test_fix_generator"),
	)

	testFixApplier := validate.NewFixApplierWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_fix_applier_system"),
		cfg.GetString("templates.prompts.test_fix_applier"),
	)

	tableRenderer := commands.NewTableRenderer()
	scenarioParser := requirements.NewScenarioParser()

	return commands.NewReqValidationCommand(
		testChecklistLoader,
		testEvaluator,
		testFixGenerator,
		testFixApplier,
		userInputCollector,
		tableRenderer,
		scenarioParser,
		runDir,
	)
}
