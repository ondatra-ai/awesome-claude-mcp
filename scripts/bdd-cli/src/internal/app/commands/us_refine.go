package commands

import (
	"context"
	"fmt"
	"log/slog"

	"bdd-cli/src/internal/app/generators/validate"
	"bdd-cli/src/internal/domain/models/story"
	"bdd-cli/src/internal/infrastructure/checklist"
	"bdd-cli/src/internal/infrastructure/fs"
	"bdd-cli/src/internal/infrastructure/input"
	storyinfra "bdd-cli/src/internal/infrastructure/story"
)

// RefineDeps bundles what `us refine` needs at the command boundary.
type RefineDeps struct {
	StoryLoader        *storyinfra.StoryLoader
	ChecklistLoader    *checklist.ChecklistLoader
	Evaluator          *validate.ChecklistEvaluator
	FixGenerator       *validate.FixPromptGenerator
	FixApplier         *validate.FixApplier
	UserInputCollector *input.UserInputCollector
	TableRenderer      *TableRenderer
	RunDir             *fs.RunDirectory
	StoriesDir         string
}

// RunRefine drives `us refine`. Loads a story from docs/stories/,
// walks the us-refine checklist, and on convergence updates the
// story file in place.
func RunRefine(ctx context.Context, deps RefineDeps, storyNumber string, fix bool) error {
	versionMgr := fs.NewStoryVersionManager(deps.RunDir, storyNumber)

	return Run(ctx, Spec[*story.Story]{
		Name:          "us refine",
		ChecklistName: "us-refine",
		StoryNumber:   storyNumber,
		Fix:           fix,

		LoadItems:  loadStoryFromFile(deps.StoryLoader, storyNumber, versionMgr),
		PostFix:    storyPostFix(versionMgr),
		Finalize:   storyFinalize(deps.StoriesDir, storyNumber, versionMgr, fix, false),
		GetSubject: storySubject,

		Evaluator:    deps.Evaluator,
		FixGenerator: deps.FixGenerator,
		FixApplier:   deps.FixApplier,

		ChecklistLoader: deps.ChecklistLoader,
		Renderer:        deps.TableRenderer,
		UI:              newFixLoopUI(deps.UserInputCollector),
		TmpDir:          deps.RunDir.GetTmpOutPath(),
	})
}

// loadStoryFromFile is the LoadItems factory for `us refine`. Loads
// the story from docs/stories/<id>-*.yaml and seeds the version
// manager.
func loadStoryFromFile(
	loader *storyinfra.StoryLoader,
	storyNumber string,
	versionMgr *fs.StoryVersionManager,
) func(ctx context.Context) ([]*story.Story, error) {
	return func(_ context.Context) ([]*story.Story, error) {
		doc, err := loader.Load(storyNumber)
		if err != nil {
			return nil, fmt.Errorf(
				"story file not found — run `bdd-cli us create %s` first: %w",
				storyNumber, err,
			)
		}

		loaded := &doc.Story
		slog.Info("Story loaded", "id", loaded.ID, "title", loaded.Title)

		err = versionMgr.SaveInitialVersion(loaded)
		if err != nil {
			return nil, fmt.Errorf("failed to save initial story version: %w", err)
		}

		return []*story.Story{loaded}, nil
	}
}
