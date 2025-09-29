package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/epic"
)

// TaskGenerator interface for generating tasks
type TaskGenerator interface {
	GenerateTasks(ctx context.Context, storyDoc *story.StoryDocument) ([]story.Task, error)
}

// DevNotesGenerator interface for generating dev notes
type DevNotesGenerator interface {
	GenerateDevNotes(ctx context.Context, storyDoc *story.StoryDocument) (story.DevNotes, error)
}

// QAResultsGenerator interface for generating QA results
type QAResultsGenerator interface {
	GenerateQAResults(ctx context.Context, storyDoc *story.StoryDocument) (story.QAResults, error)
}

// TestingRequirementsGenerator interface for generating testing requirements
type TestingRequirementsGenerator interface {
	GenerateTesting(ctx context.Context, storyDoc *story.StoryDocument) (story.Testing, error)
}

type StoryFactory struct {
	epicLoader         *epic.EpicLoader
	aiClient           AIClient
	config             *config.ViperConfig
	architectureLoader *docs.ArchitectureLoader
}

func NewStoryFactory(epicLoader *epic.EpicLoader, aiClient AIClient, config *config.ViperConfig, architectureLoader *docs.ArchitectureLoader) *StoryFactory {
	return &StoryFactory{
		epicLoader:         epicLoader,
		aiClient:           aiClient,
		config:             config,
		architectureLoader: architectureLoader,
	}
}

func (f *StoryFactory) CreateStory(ctx context.Context, storyNumber string) (*story.StoryDocument, error) {
	// Load story from epic file - fail if not found
	loadedStory, err := f.epicLoader.LoadStoryFromEpic(storyNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to load story from epic file: %w", err)
	}

	// Load architecture documents once for all generators
	architectureDocs, err := f.architectureLoader.LoadAllArchitectureDocsStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to load architecture documents: %w", err)
	}

	// Create initial story document with all required data
	storyDoc := &story.StoryDocument{
		Story: *loadedStory,
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
	taskGenerator := NewTaskGenerator(f.aiClient, f.config)
	devNotesGenerator := NewDevNotesGenerator(f.aiClient, f.config)
	testingGenerator := NewTestingGenerator(f.aiClient, f.config)
	qaResultsGenerator := NewQAAssessmentGenerator(f.aiClient, f.config)

	// Generate tasks using AI - fail on any error
	tasks, err := taskGenerator.GenerateTasks(ctx, storyDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tasks: %w", err)
	}
	storyDoc.Tasks = tasks

	// Generate dev_notes using AI - fail on any error
	devNotes, err := devNotesGenerator.GenerateDevNotes(ctx, storyDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dev_notes: %w", err)
	}
	storyDoc.DevNotes = devNotes

	// Generate testing requirements using AI - fail on any error
	testing, err := testingGenerator.GenerateTesting(ctx, storyDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to generate testing requirements: %w", err)
	}
	storyDoc.Testing = testing

	// Generate QA results using AI - fail on any error
	qaResults, err := qaResultsGenerator.GenerateQAResults(ctx, storyDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QA results: %w", err)
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

// GetTmpDir returns the configured temporary directory path
func (f *StoryFactory) GetTmpDir() string {
	return f.config.GetString("paths.tmp_dir")
}

// GetStoriesDir returns the configured stories output directory path
func (f *StoryFactory) GetStoriesDir() string {
	return f.config.GetString("paths.stories_dir")
}
