package factories

import (
	"context"
	"regexp"
	"strings"
	"time"

	"bmad-cli/internal/app/generators/create"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/epic"
	"bmad-cli/internal/infrastructure/fs"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// TaskGenerator interface for generating tasks.
type TaskGenerator interface {
	GenerateTasks(ctx context.Context, storyDoc *story.StoryDocument) ([]story.Task, error)
}

// DevNotesGenerator interface for generating dev notes.
type DevNotesGenerator interface {
	GenerateDevNotes(ctx context.Context, storyDoc *story.StoryDocument) (story.DevNotes, error)
}

// QAResultsGenerator interface for generating QA results.
type QAResultsGenerator interface {
	GenerateQAResults(ctx context.Context, storyDoc *story.StoryDocument) (story.QAResults, error)
}

// TestingRequirementsGenerator interface for generating testing requirements.
type TestingRequirementsGenerator interface {
	GenerateTesting(ctx context.Context, storyDoc *story.StoryDocument) (story.Testing, error)
}

type StoryFactory struct {
	epicLoader         *epic.EpicLoader
	aiClient           ports.AIPort
	config             *config.ViperConfig
	architectureLoader *docs.ArchitectureLoader
	runDirectory       *fs.RunDirectory
}

func NewStoryFactory(
	epicLoader *epic.EpicLoader,
	aiClient ports.AIPort,
	config *config.ViperConfig,
	architectureLoader *docs.ArchitectureLoader,
	runDirectory *fs.RunDirectory,
) *StoryFactory {
	return &StoryFactory{
		epicLoader:         epicLoader,
		aiClient:           aiClient,
		config:             config,
		architectureLoader: architectureLoader,
		runDirectory:       runDirectory,
	}
}

func (f *StoryFactory) CreateStory(ctx context.Context, storyNumber string) (*story.StoryDocument, error) {
	// Load story from epic file - fail if not found
	loadedStory, err := f.epicLoader.LoadStoryFromEpic(storyNumber)
	if err != nil {
		return nil, pkgerrors.ErrLoadStoryFromEpicFailed(err)
	}

	// Load architecture documents once for all generators
	architectureDocs, err := f.architectureLoader.LoadAllArchitectureDocsStruct()
	if err != nil {
		return nil, pkgerrors.ErrLoadArchitectureDocsFailed(err)
	}

	// Create initial story document with all required data
	storyDoc := &story.StoryDocument{
		Story:            *loadedStory,
		ArchitectureDocs: architectureDocs,
		ChangeLog: []story.ChangeLogEntry{
			{
				Date:        time.Now().Format("2006-01-02"),
				Version:     "1.0.0",
				Description: "Initial story creation",
				Author:      "bmad-cli",
			},
		},
		DevAgentRecord: story.DevAgentRecord{
			AgentModelUsed:     nil,
			DebugLogReferences: []string{},
			CompletionNotes:    []string{},
			FileList:           []string{},
		},
	}

	// Create generators
	taskGenerator := create.NewTaskGenerator(f.aiClient, f.config)
	devNotesGenerator := create.NewDevNotesGenerator(f.aiClient, f.config)
	testingGenerator := create.NewAITestingGenerator(f.aiClient, f.config)
	scenariosGenerator := create.NewAIScenariosGenerator(f.aiClient, f.config)
	qaResultsGenerator := create.NewAIQAAssessmentGenerator(f.aiClient, f.config)

	// Get run directory path for passing to generators
	runDirPath := f.runDirectory.GetTmpOutPath()

	// Generate tasks using AI - fail on any error
	tasks, err := taskGenerator.GenerateTasks(ctx, storyDoc, runDirPath)
	if err != nil {
		return nil, pkgerrors.ErrGenerateTasksFailed(err)
	}

	storyDoc.Tasks = tasks

	// Generate dev_notes using AI - fail on any error
	devNotes, err := devNotesGenerator.GenerateDevNotes(ctx, storyDoc, runDirPath)
	if err != nil {
		return nil, pkgerrors.ErrGenerateDevNotesFailed(err)
	}

	storyDoc.DevNotes = devNotes

	// Generate testing requirements using AI - fail on any error
	testing, err := testingGenerator.GenerateTesting(ctx, storyDoc, runDirPath)
	if err != nil {
		return nil, pkgerrors.ErrGenerateTestingReqsFailed(err)
	}

	storyDoc.Testing = testing

	// Generate test scenarios using AI - fail on any error
	scenarios, err := scenariosGenerator.GenerateScenarios(ctx, storyDoc, runDirPath)
	if err != nil {
		return nil, pkgerrors.ErrGenerateTestScenariosFailed(err)
	}

	storyDoc.Scenarios = scenarios

	// Generate QA results using AI - fail on any error
	qaResults, err := qaResultsGenerator.GenerateQAResults(ctx, storyDoc, runDirPath)
	if err != nil {
		return nil, pkgerrors.ErrGenerateQAFailed(err)
	}

	storyDoc.QAResults = &qaResults

	return storyDoc, nil
}

func (f *StoryFactory) SlugifyTitle(title string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}

// GetTmpDirPath returns the run-specific directory path for this execution.
func (f *StoryFactory) GetTmpDirPath() string {
	if f.runDirectory != nil {
		return f.runDirectory.GetTmpOutPath()
	}
	// Fallback to configured tmp_dir if no run directory created yet
	return f.config.GetString("paths.tmp_dir")
}

// GetStoriesDir returns the configured stories output directory path.
func (f *StoryFactory) GetStoriesDir() string {
	return f.config.GetString("paths.stories_dir")
}
