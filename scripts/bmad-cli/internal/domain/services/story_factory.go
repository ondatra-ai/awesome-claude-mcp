package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/epic"
)

// TaskGenerator interface for generating tasks
type TaskGenerator interface {
	GenerateTasks(ctx context.Context, story *story.Story, architectureDocs map[string]docs.ArchitectureDoc) ([]story.Task, error)
}

// DevNotesGenerator interface for generating dev notes
type DevNotesGenerator interface {
	GenerateDevNotes(ctx context.Context, story *story.Story, tasks []story.Task, architectureDocs map[string]docs.ArchitectureDoc) (story.DevNotes, error)
}

// ArchitectureLoader interface for loading architecture documents
type ArchitectureLoader interface {
	LoadAllArchitectureDocs() (map[string]docs.ArchitectureDoc, error)
}

type StoryFactory struct {
	epicLoader         *epic.EpicLoader
	taskGenerator      TaskGenerator
	devNotesGenerator  DevNotesGenerator
	architectureLoader ArchitectureLoader
}

func NewStoryFactory(epicLoader *epic.EpicLoader, taskGenerator TaskGenerator, devNotesGenerator DevNotesGenerator, architectureLoader ArchitectureLoader) *StoryFactory {
	return &StoryFactory{
		epicLoader:         epicLoader,
		taskGenerator:      taskGenerator,
		devNotesGenerator:  devNotesGenerator,
		architectureLoader: architectureLoader,
	}
}

func (f *StoryFactory) CreateStory(ctx context.Context, storyNumber string) (*story.StoryDocument, error) {
	// Load story from epic file - fail if not found
	loadedStory, err := f.epicLoader.LoadStoryFromEpic(storyNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to load story from epic file: %w", err)
	}

	// Generate tasks using AI - fail on any error
	tasks, err := f.generateTasks(ctx, loadedStory)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tasks: %w", err)
	}

	// Generate dev_notes using AI - fail on any error
	devNotes, err := f.generateDevNotes(ctx, loadedStory, tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dev_notes: %w", err)
	}

	return &story.StoryDocument{
		Story:    *loadedStory,
		Tasks:    tasks,
		DevNotes: devNotes,
		Testing: story.Testing{
			TestLocation: "services/backend/tests",
			Frameworks:   []string{"testing", "testify"},
			Requirements: []string{
				"Unit tests for all public methods",
				"Integration tests for external dependencies",
				"End-to-end tests for complete workflows",
			},
			Coverage: map[string]string{
				"business_logic": "80%",
				"overall":        "75%",
			},
		},
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
	}, nil
}

func (f *StoryFactory) SlugifyTitle(title string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

// generateTasks generates tasks using AI - fails on any error
func (f *StoryFactory) generateTasks(ctx context.Context, loadedStory *story.Story) ([]story.Task, error) {

	// Load architecture documents - fail immediately if any are missing
	architectureDocs, err := f.architectureLoader.LoadAllArchitectureDocs()
	if err != nil {
		return nil, fmt.Errorf("failed to load architecture documents: %w", err)
	}

	// Generate tasks using AI - fail if AI generation fails
	tasks, err := f.taskGenerator.GenerateTasks(ctx, loadedStory, architectureDocs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tasks using AI: %w", err)
	}

	fmt.Printf("✅ Generated %d tasks using AI\n", len(tasks))
	return tasks, nil
}

// generateDevNotes generates dev_notes using AI - fails on any error
func (f *StoryFactory) generateDevNotes(ctx context.Context, loadedStory *story.Story, tasks []story.Task) (story.DevNotes, error) {

	// Load architecture documents - fail immediately if any are missing
	architectureDocs, err := f.architectureLoader.LoadAllArchitectureDocs()
	if err != nil {
		return nil, fmt.Errorf("failed to load architecture documents: %w", err)
	}

	// Generate dev_notes using AI - fail if AI generation fails
	devNotes, err := f.devNotesGenerator.GenerateDevNotes(ctx, loadedStory, tasks, architectureDocs)
	if err != nil {
		return story.DevNotes{}, fmt.Errorf("failed to generate dev_notes using AI: %w", err)
	}

	fmt.Printf("✅ Generated dev_notes using AI\n")
	return devNotes, nil
}
