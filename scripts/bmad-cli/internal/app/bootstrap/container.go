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
	Config          *config.ViperConfig
	USValidationCmd *commands.USValidationCommand
	ScenarioParser  *requirements.ScenarioParser
	// Shared evaluator / fix-prompt / fix-applier triple used by every
	// scenario-walking command (`us generate_tests`, `us implement`).
	// The underlying prompt templates live at
	// templates.prompts.test_checklist* / test_fix_* in bmad-cli.yaml;
	// they're generic enough to drive any checklist that iterates over
	// TestGenerationData scenarios.
	ScenarioEvaluator          *validate.ChecklistEvaluator
	ScenarioFixPromptGenerator *validate.FixPromptGenerator
	ScenarioFixApplier         *validate.FixApplier
	RunDir                     *fs.RunDirectory
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

	// Scenario-validation evaluator / fix-prompt / fix-applier set, shared
	// by `us generate_tests` and `us implement`. The underlying templates
	// are configured under templates.prompts.test_checklist* / test_fix_*.
	scenarioEvaluator := validate.NewChecklistEvaluatorWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_checklist_system"),
		cfg.GetString("templates.prompts.test_checklist"),
	)

	scenarioFixPromptGenerator := validate.NewFixPromptGeneratorWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_fix_generator_system"),
		cfg.GetString("templates.prompts.test_fix_generator"),
	)

	scenarioFixApplier := validate.NewFixApplierWithPaths(
		claudeClient, cfg,
		cfg.GetString("templates.prompts.test_fix_applier_system"),
		cfg.GetString("templates.prompts.test_fix_applier"),
	)

	scenarioParser := requirements.NewScenarioParser()

	return &Container{
		Config:                     cfg,
		USValidationCmd:            usValidationCmd,
		ScenarioParser:             scenarioParser,
		ScenarioEvaluator:          scenarioEvaluator,
		ScenarioFixPromptGenerator: scenarioFixPromptGenerator,
		ScenarioFixApplier:         scenarioFixApplier,
		RunDir:                     runDir,
	}, nil
}
