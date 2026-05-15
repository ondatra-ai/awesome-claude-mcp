package bootstrap

import (
	"bdd-cli/src/adapters/ai"
	"bdd-cli/src/internal/app/commands"
	"bdd-cli/src/internal/app/generators/validate"
	"bdd-cli/src/internal/domain/ports"
	"bdd-cli/src/internal/infrastructure/checklist"
	"bdd-cli/src/internal/infrastructure/config"
	"bdd-cli/src/internal/infrastructure/epic"
	"bdd-cli/src/internal/infrastructure/fs"
	"bdd-cli/src/internal/infrastructure/input"
	"bdd-cli/src/internal/infrastructure/story"
	pkgerrors "bdd-cli/src/internal/pkg/errors"
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

// scenarioTripleConfigKeys names the bdd-cli.yaml config paths for one
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

// Container wires together the components needed by the CLI. Per-command
// deps (CreateDeps, RefineDeps, ApplyDeps) are projected from this on
// the cobra side so each command depends on only what it actually uses.
type Container struct {
	Config              *config.ViperConfig
	RunDir              *fs.RunDirectory
	StoriesDir          string
	EpicLoader          *epic.EpicLoader
	StoryLoader         *story.StoryLoader
	StoryScenarioParser *story.StoryScenarioParser
	ChecklistLoader     *checklist.ChecklistLoader
	UserInputCollector  *input.UserInputCollector
	TableRenderer       *commands.TableRenderer
	// Standard generator triple drives `us create` and `us refine`.
	Evaluator    *validate.ChecklistEvaluator
	FixGenerator *validate.FixPromptGenerator
	FixApplier   *validate.FixApplier
	// Apply-flavored triple drives `us apply`. Templates live under
	// templates.prompts.apply_* and the fix-applier is configured with
	// EditMode so Claude can mutate the scratch registry directly.
	ApplyEvaluator          *validate.ChecklistEvaluator
	ApplyFixPromptGenerator *validate.FixPromptGenerator
	ApplyFixApplier         *validate.FixApplier
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

	claudeClient, err := ai.NewClaudeClient()
	if err != nil {
		return nil, pkgerrors.ErrCreateAIClientFailed(err)
	}

	applyTrip := newScenarioTriple(claudeClient, cfg, scenarioTripleConfigKeys{
		checklistSystem:    "templates.prompts.apply_checklist_system",
		checklist:          "templates.prompts.apply_checklist",
		fixGeneratorSystem: "templates.prompts.apply_fix_generator_system",
		fixGenerator:       "templates.prompts.apply_fix_generator",
		fixApplierSystem:   "templates.prompts.apply_fix_applier_system",
		fixApplier:         "templates.prompts.apply_fix_applier",
	})
	applyTrip.fixApplier.UseEditMode()

	return &Container{
		Config:                  cfg,
		RunDir:                  runDir,
		StoriesDir:              cfg.GetString("paths.stories_dir"),
		EpicLoader:              epic.NewEpicLoader(cfg),
		StoryLoader:             story.NewStoryLoader(cfg),
		StoryScenarioParser:     story.NewStoryScenarioParser(cfg),
		ChecklistLoader:         checklist.NewChecklistLoader(cfg),
		UserInputCollector:      input.NewUserInputCollector(),
		TableRenderer:           commands.NewTableRenderer(),
		Evaluator:               validate.NewChecklistEvaluator(claudeClient, cfg),
		FixGenerator:            validate.NewFixPromptGenerator(claudeClient, cfg),
		FixApplier:              validate.NewFixApplier(claudeClient, cfg),
		ApplyEvaluator:          applyTrip.evaluator,
		ApplyFixPromptGenerator: applyTrip.fixGenerator,
		ApplyFixApplier:         applyTrip.fixApplier,
	}, nil
}
