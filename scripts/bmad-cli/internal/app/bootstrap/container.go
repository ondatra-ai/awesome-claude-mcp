package bootstrap

import (
	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/app/commands"
	"bmad-cli/internal/app/generators/validate"
	"bmad-cli/internal/infrastructure/checklist"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/input"
	"bmad-cli/internal/infrastructure/requirements"
	"bmad-cli/internal/infrastructure/story"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// Container wires together the components needed by the CLI.
type Container struct {
	Config                 *config.ViperConfig
	USValidationCmd        *commands.USValidationCommand
	ScenarioParser         *requirements.ScenarioParser
	TestChecklistEvaluator *validate.ChecklistEvaluator
	TestFixPromptGenerator *validate.FixPromptGenerator
	TestFixApplier         *validate.FixApplier
	RunDir                 *fs.RunDirectory
}

// NewContainer builds the Container.
func NewContainer() (*Container, error) {
	cfg, err := config.NewViperConfig()
	if err != nil {
		return nil, pkgerrors.ErrInitializeConfigFailed(err)
	}

	configureLogging()

	runDir, err := fs.NewRunDirectory(cfg.GetString("paths.tmp_dir"))
	if err != nil {
		return nil, pkgerrors.ErrCreateRunDirectoryFailed(err)
	}

	epicLoader := epic.NewEpicLoader(cfg)

	claudeClient, err := ai.NewClaudeClient()
	if err != nil {
		return nil, pkgerrors.ErrCreateAIClientFailed(err)
	}

	storyLoader := story.NewStoryLoader(cfg)
	userInputCollector := input.NewUserInputCollector()

	checklistLoader := checklist.NewChecklistLoader(cfg)

	evaluator := validate.NewChecklistEvaluator(claudeClient, cfg)
	fixPromptGenerator := validate.NewFixPromptGenerator(claudeClient, cfg)
	fixApplier := validate.NewFixApplier(claudeClient, cfg)

	tableRenderer := commands.NewTableRenderer()
	storiesDir := cfg.GetString("paths.stories_dir")

	usValidationCmd := commands.NewUSValidationCommand(
		epicLoader,
		storyLoader,
		checklistLoader,
		evaluator,
		fixPromptGenerator,
		fixApplier,
		userInputCollector,
		tableRenderer,
		runDir,
		storiesDir,
	)

	// Separate prompt-template set for test validation (`us generate_tests`).
	testEvaluator := validate.NewChecklistEvaluatorWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_checklist_system"),
		cfg.GetString("templates.prompts.test_checklist"),
	)

	testFixPromptGenerator := validate.NewFixPromptGeneratorWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_fix_generator_system"),
		cfg.GetString("templates.prompts.test_fix_generator"),
	)

	testFixApplier := validate.NewFixApplierWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_fix_applier_system"),
		cfg.GetString("templates.prompts.test_fix_applier"),
	)

	scenarioParser := requirements.NewScenarioParser()

	return &Container{
		Config:                 cfg,
		USValidationCmd:        usValidationCmd,
		ScenarioParser:         scenarioParser,
		TestChecklistEvaluator: testEvaluator,
		TestFixPromptGenerator: testFixPromptGenerator,
		TestFixApplier:         testFixApplier,
		RunDir:                 runDir,
	}, nil
}
