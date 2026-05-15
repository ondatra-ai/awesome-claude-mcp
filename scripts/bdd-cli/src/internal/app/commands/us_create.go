package commands

import (
	"context"
	"fmt"
	"log/slog"

	"bdd-cli/src/internal/app/generators/validate"
	checklistmodels "bdd-cli/src/internal/domain/models/checklist"
	"bdd-cli/src/internal/domain/models/story"
	"bdd-cli/src/internal/infrastructure/checklist"
	"bdd-cli/src/internal/infrastructure/epic"
	"bdd-cli/src/internal/infrastructure/fs"
	"bdd-cli/src/internal/infrastructure/input"
	"bdd-cli/src/internal/pkg/console"
	pkgerrors "bdd-cli/src/internal/pkg/errors"
)

// RenderedPrompt is the Q value the engine passes around. Pairs a
// PromptWithContext with its 1-based index so per-cell file naming
// stays stable across walks.
type RenderedPrompt struct {
	Prompt checklistmodels.PromptWithContext
	Index  int
}

// renderPrompt is the engine's GenerateQFn for all three `us`
// commands. Cheap on purpose — the heavy template rendering happens
// inside the evaluator's Claude call, which has both item and q.
func renderPrompt(idx int, p checklistmodels.PromptWithContext) *RenderedPrompt {
	return &RenderedPrompt{Prompt: p, Index: idx}
}

// CreateDeps bundles what `us create` needs at the command boundary.
type CreateDeps struct {
	EpicLoader         *epic.EpicLoader
	ChecklistLoader    *checklist.ChecklistLoader
	Evaluator          *validate.ChecklistEvaluator
	FixGenerator       *validate.FixPromptGenerator
	FixApplier         *validate.FixApplier
	UserInputCollector *input.UserInputCollector
	TableRenderer      *TableRenderer
	RunDir             *fs.RunDirectory
	StoriesDir         string
}

// RunCreate drives `us create`. Loads a story from its epic, walks
// the us-create checklist, and on convergence writes the new story
// file under StoriesDir.
func RunCreate(ctx context.Context, deps CreateDeps, storyNumber string, fix bool) error {
	versionMgr := fs.NewStoryVersionManager(deps.RunDir, storyNumber)

	return Run(ctx, Spec[*story.Story]{
		Name:          "us create",
		ChecklistName: "us-create",
		StoryNumber:   storyNumber,
		Fix:           fix,

		LoadItems:  loadStoryFromEpic(deps.EpicLoader, storyNumber, versionMgr),
		PostFix:    storyPostFix(versionMgr),
		Finalize:   storyFinalize(deps.StoriesDir, storyNumber, versionMgr, fix, true),
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

// loadStoryFromEpic is the LoadItems factory for `us create`. Loads
// the story from its epic, displays it, and seeds the version
// manager with the initial snapshot.
func loadStoryFromEpic(
	loader *epic.EpicLoader,
	storyNumber string,
	versionMgr *fs.StoryVersionManager,
) func(ctx context.Context) ([]*story.Story, error) {
	return func(_ context.Context) ([]*story.Story, error) {
		console.Header("LOADING STORY FROM EPIC", separatorWidth)

		loaded, err := loader.LoadStoryFromEpic(storyNumber)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to load story: %w",
				pkgerrors.ErrLoadStoryFromEpicFailed(err),
			)
		}

		displayStory(loaded, "STORY FROM EPIC")

		slog.Info("Story loaded", "id", loaded.ID, "title", loaded.Title)

		err = versionMgr.SaveInitialVersion(loaded)
		if err != nil {
			return nil, fmt.Errorf("failed to save initial story version: %w", err)
		}

		return []*story.Story{loaded}, nil
	}
}
