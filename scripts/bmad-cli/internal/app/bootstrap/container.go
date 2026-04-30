package bootstrap

import (
	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/app/commands"
	"bmad-cli/internal/app/generators/validate"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/checklist"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/input"
	"bmad-cli/internal/infrastructure/requirements"
	"bmad-cli/internal/infrastructure/story"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// scenarioTriple bundles the (evaluator, fix-prompt-generator,
// fix-applier) trio that every scenario-walking command depends on.
// Used internally by NewContainer to keep the bootstrap function under
// the configured complexity / length budgets.
type scenarioTriple struct {
	evaluator    *validate.ChecklistEvaluator
	fixGenerator *validate.FixPromptGenerator
	fixApplier   *validate.FixApplier
}

// scenarioTripleConfigKeys names the bmad-cli.yaml config paths for one
// scenario-walking command's evaluator / fix-generator / fix-applier
// templates.
type scenarioTripleConfigKeys struct {
	checklistSystem    string
	checklist          string
	fixGeneratorSystem string
	fixGenerator       string
	fixApplierSystem   string
	fixApplier         string
}

func newScenarioTriple(
	aiClient ports.AIPort,
	cfg *config.ViperConfig,
	keys scenarioTripleConfigKeys,
) scenarioTriple {
	return scenarioTriple{
		evaluator: validate.NewChecklistEvaluatorWithPaths(
			aiClient, cfg,
			cfg.GetString(keys.checklistSystem),
			cfg.GetString(keys.checklist),
		),
		fixGenerator: validate.NewFixPromptGeneratorWithPaths(
			aiClient, cfg,
			cfg.GetString(keys.fixGeneratorSystem),
			cfg.GetString(keys.fixGenerator),
		),
		fixApplier: validate.NewFixApplierWithPaths(
			aiClient, cfg,
			cfg.GetString(keys.fixApplierSystem),
			cfg.GetString(keys.fixApplier),
		),
	}
}

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
	// Apply-flavored equivalents driving `us apply`. The parser reads a
	// refined story file (acceptance_criteria[].steps shape) and emits
	// one ScenarioApplyData per AC; the evaluator / fix-prompt / fix-
	// applier triple uses the templates.prompts.apply_* templates and
	// targets the scratch copy of docs/requirements.yaml.
	StoryScenarioParser     *story.StoryScenarioParser
	ApplyEvaluator          *validate.ChecklistEvaluator
	ApplyFixPromptGenerator *validate.FixPromptGenerator
	ApplyFixApplier         *validate.FixApplier
	RunDir                  *fs.RunDirectory
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
	scenarioTrip := newScenarioTriple(claudeClient, cfg, scenarioTripleConfigKeys{
		checklistSystem:    "templates.prompts.test_checklist_system",
		checklist:          "templates.prompts.test_checklist",
		fixGeneratorSystem: "templates.prompts.test_fix_generator_system",
		fixGenerator:       "templates.prompts.test_fix_generator",
		fixApplierSystem:   "templates.prompts.test_fix_applier_system",
		fixApplier:         "templates.prompts.test_fix_applier",
	})

	// Apply-flavored evaluator / fix-prompt / fix-applier set used by
	// `us apply`. Templates live under templates.prompts.apply_* and
	// are written for the merge-into-requirements.yaml subject.
	applyTrip := newScenarioTriple(claudeClient, cfg, scenarioTripleConfigKeys{
		checklistSystem:    "templates.prompts.apply_checklist_system",
		checklist:          "templates.prompts.apply_checklist",
		fixGeneratorSystem: "templates.prompts.apply_fix_generator_system",
		fixGenerator:       "templates.prompts.apply_fix_generator",
		fixApplierSystem:   "templates.prompts.apply_fix_applier_system",
		fixApplier:         "templates.prompts.apply_fix_applier",
	})

	scenarioParser := requirements.NewScenarioParser()
	storyScenarioParser := story.NewStoryScenarioParser(cfg)

	return &Container{
		Config:                     cfg,
		USValidationCmd:            usValidationCmd,
		ScenarioParser:             scenarioParser,
		ScenarioEvaluator:          scenarioTrip.evaluator,
		ScenarioFixPromptGenerator: scenarioTrip.fixGenerator,
		ScenarioFixApplier:         scenarioTrip.fixApplier,
		StoryScenarioParser:        storyScenarioParser,
		ApplyEvaluator:             applyTrip.evaluator,
		ApplyFixPromptGenerator:    applyTrip.fixGenerator,
		ApplyFixApplier:            applyTrip.fixApplier,
		RunDir:                     runDir,
	}, nil
}
