package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/epic"
)

// TaskGenerator interface for generating tasks
type TaskGenerator interface {
	GenerateTasks(ctx context.Context, story *story.Story, architectureDocs map[string]string) ([]story.Task, error)
}

// ArchitectureLoader interface for loading architecture documents
type ArchitectureLoader interface {
	LoadAllArchitectureDocs() (map[string]string, error)
}

type StoryFactory struct {
	epicLoader          *epic.EpicLoader
	taskGenerator       TaskGenerator
	architectureLoader  ArchitectureLoader
}

func NewStoryFactory(epicLoader *epic.EpicLoader, taskGenerator TaskGenerator, architectureLoader ArchitectureLoader) *StoryFactory {
	return &StoryFactory{
		epicLoader:         epicLoader,
		taskGenerator:      taskGenerator,
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

	return &story.StoryDocument{
		Story: *loadedStory,
		Tasks: tasks,
		DevNotes: story.DevNotes{
			PreviousStoryInsights: "This is a new story without previous implementation insights",
			TechnologyStack: story.TechnologyStack{
				Language:       "Go",
				Framework:      "Standard library",
				MCPIntegration: "MCP integration as needed",
				Logging:        "slog",
				Config:         "viper",
			},
			Architecture: story.Architecture{
				Component:        loadedStory.Title,
				Responsibilities: []string{"Implement core functionality", "Handle business logic"},
				Dependencies:     []string{"context", "fmt", "log/slog"},
				TechStack:        []string{"Go", "YAML", "HTTP", "JSON"},
			},
			FileStructure: story.FileStructure{
				Files: []string{"services/backend/internal/story/implementation.go"},
			},
			Configuration: story.Configuration{
				EnvironmentVariables: map[string]string{
					"LOG_LEVEL":     "info",
					"PORT":          "8080",
					"TEMPLATE_PATH": "templates/",
				},
			},
			PerformanceRequirements: story.PerformanceRequirements{
				ConnectionEstablishment: "< 100ms",
				MessageProcessing:       "< 50ms",
				ConcurrentConnections:   "100",
				MemoryUsage:            "< 100MB",
			},
		},
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
			AgentModelUsed:      nil,
			DebugLogReferences:  []string{},
			CompletionNotes:     []string{},
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
