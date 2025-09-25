package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/epic"
)

type StoryFactory struct {
	epicLoader *epic.EpicLoader
}

func NewStoryFactory(epicLoader *epic.EpicLoader) *StoryFactory {
	return &StoryFactory{
		epicLoader: epicLoader,
	}
}

func (f *StoryFactory) CreateStory(storyNumber string) *story.StoryDocument {
	// Try to load story from epic file first
	loadedStory, err := f.epicLoader.LoadStoryFromEpic(storyNumber)
	if err != nil {
		// Fallback to generating story if loading fails
		fmt.Printf("Warning: Could not load story from epic file: %v. Using default generation.\n", err)
		return f.createDefaultStory(storyNumber)
	}

	// Use loaded story data
	epic, storyNum := f.parseStoryNumber(storyNumber)
	component := f.getComponentFromTitle(loadedStory.Title)

	return &story.StoryDocument{
		Story: *loadedStory,
		Tasks: []story.Task{
			{
				Name:               fmt.Sprintf("Implement %s", component),
				AcceptanceCriteria: []string{"AC-1", "AC-2"},
				Subtasks: []string{
					"Create domain models",
					"Implement business logic",
					"Add error handling",
					"Write comprehensive tests",
				},
				Status: "pending",
			},
		},
		DevNotes: story.DevNotes{
			PreviousStoryInsights: "This is a new story without previous implementation insights",
			TechnologyStack: story.TechnologyStack{
				Language:       "Go",
				Framework:      "Standard library",
				MCPIntegration: f.getMCPIntegration(epic, storyNum),
				Logging:        "slog",
				Config:         "viper",
			},
			Architecture: story.Architecture{
				Component:        component,
				Responsibilities: f.getResponsibilities(epic, storyNum),
				Dependencies:     f.getDependencies(epic, storyNum),
				TechStack:        []string{"Go", "YAML", "HTTP", "JSON"},
			},
			FileStructure: story.FileStructure{
				Files: f.getFiles(epic, storyNum),
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
	}
}

func (f *StoryFactory) createDefaultStory(storyNumber string) *story.StoryDocument {
	epic, storyNum := f.parseStoryNumber(storyNumber)
	title := f.generateTitle(epic, storyNum)
	component := f.getComponentName(epic, storyNum)

	return &story.StoryDocument{
		Story: story.Story{
			ID:     storyNumber,
			Title:  title,
			AsA:    "developer",
			IWant:  fmt.Sprintf("to implement %s", strings.ToLower(title)),
			SoThat: "the system can provide the required functionality",
			Status: "PLANNED",
			AcceptanceCriteria: []story.AcceptanceCriterion{
				{ID: "AC-1", Description: fmt.Sprintf("%s is properly implemented", component)},
				{ID: "AC-2", Description: "Unit tests pass with >80% coverage"},
				{ID: "AC-3", Description: "Integration tests verify functionality"},
			},
		},
		Tasks: []story.Task{
			{
				Name:               fmt.Sprintf("Implement %s", component),
				AcceptanceCriteria: []string{"AC-1", "AC-2"},
				Subtasks: []string{
					"Create domain models",
					"Implement business logic",
					"Add error handling",
					"Write comprehensive tests",
				},
				Status: "pending",
			},
		},
		DevNotes: story.DevNotes{
			PreviousStoryInsights: "This is a new story without previous implementation insights",
			TechnologyStack: story.TechnologyStack{
				Language:       "Go",
				Framework:      "Standard library",
				MCPIntegration: f.getMCPIntegration(epic, storyNum),
				Logging:        "slog",
				Config:         "viper",
			},
			Architecture: story.Architecture{
				Component:        component,
				Responsibilities: f.getResponsibilities(epic, storyNum),
				Dependencies:     f.getDependencies(epic, storyNum),
				TechStack:        []string{"Go", "YAML", "HTTP", "JSON"},
			},
			FileStructure: story.FileStructure{
				Files: f.getFiles(epic, storyNum),
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
	}
}

func (f *StoryFactory) getComponentFromTitle(title string) string {
	// Extract component name from title
	switch title {
	case "MCP Server Implementation":
		return "MCP Server"
	case "Tool Registration":
		return "Tool Registry"
	case "Message Protocol Handler":
		return "Message Handler"
	case "MCP Error Handling":
		return "Error Handler"
	case "Connection Management":
		return "Connection Manager"
	default:
		// Default: use the title as component name
		return title
	}
}

func (f *StoryFactory) parseStoryNumber(storyNumber string) (int, int) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)$`)
	matches := re.FindStringSubmatch(storyNumber)
	if len(matches) != 3 {
		return 1, 1
	}

	epic := 1
	story := 1
	fmt.Sscanf(matches[1], "%d", &epic)
	fmt.Sscanf(matches[2], "%d", &story)

	return epic, story
}

func (f *StoryFactory) generateTitle(epic, storyNum int) string {
	switch epic {
	case 3:
		switch storyNum {
		case 1:
			return "MCP Server Implementation"
		case 2:
			return "Google Docs Integration"
		case 3:
			return "Document Processing Pipeline"
		default:
			return fmt.Sprintf("Epic 3 Story %d Implementation", storyNum)
		}
	case 4:
		switch storyNum {
		case 1:
			return "Authentication Service"
		case 2:
			return "User Management System"
		default:
			return fmt.Sprintf("Epic 4 Story %d Implementation", storyNum)
		}
	default:
		return fmt.Sprintf("Epic %d Story %d Implementation", epic, storyNum)
	}
}

func (f *StoryFactory) getComponentName(epic, storyNum int) string {
	switch epic {
	case 3:
		switch storyNum {
		case 1:
			return "MCP Server"
		case 2:
			return "Google Docs Adapter"
		case 3:
			return "Document Processor"
		default:
			return fmt.Sprintf("Epic3Story%d", storyNum)
		}
	case 4:
		switch storyNum {
		case 1:
			return "Auth Service"
		case 2:
			return "User Manager"
		default:
			return fmt.Sprintf("Epic4Story%d", storyNum)
		}
	default:
		return fmt.Sprintf("Epic%dStory%d", epic, storyNum)
	}
}

func (f *StoryFactory) getMCPIntegration(epic, storyNum int) string {
	if epic == 3 {
		return "Core MCP protocol implementation required"
	}
	return "MCP integration as needed"
}

func (f *StoryFactory) getResponsibilities(epic, storyNum int) []string {
	switch epic {
	case 3:
		switch storyNum {
		case 1:
			return []string{
				"Handle MCP protocol requests",
				"Manage server lifecycle",
				"Provide tool registration",
				"Handle client connections",
			}
		case 2:
			return []string{
				"Integrate with Google Docs API",
				"Handle document operations",
				"Manage authentication",
				"Process document changes",
			}
		default:
			return []string{"Implement core functionality", "Handle business logic"}
		}
	default:
		return []string{"Implement core functionality", "Handle business logic"}
	}
}

func (f *StoryFactory) getDependencies(epic, storyNum int) []string {
	base := []string{"context", "fmt", "log/slog"}

	switch epic {
	case 3:
		switch storyNum {
		case 1:
			return append(base, "net/http", "encoding/json", "gorilla/mux")
		case 2:
			return append(base, "google.golang.org/api/docs/v1", "oauth2")
		default:
			return base
		}
	default:
		return base
	}
}

func (f *StoryFactory) getFiles(epic, storyNum int) []string {
	base := "services/backend"

	switch epic {
	case 3:
		switch storyNum {
		case 1:
			return []string{
				fmt.Sprintf("%s/internal/mcp/server.go", base),
				fmt.Sprintf("%s/internal/mcp/handlers.go", base),
				fmt.Sprintf("%s/internal/mcp/protocol.go", base),
			}
		case 2:
			return []string{
				fmt.Sprintf("%s/internal/docs/client.go", base),
				fmt.Sprintf("%s/internal/docs/processor.go", base),
				fmt.Sprintf("%s/internal/docs/auth.go", base),
			}
		default:
			return []string{fmt.Sprintf("%s/internal/story/implementation.go", base)}
		}
	default:
		return []string{fmt.Sprintf("%s/internal/story/implementation.go", base)}
	}
}

func (f *StoryFactory) SlugifyTitle(title string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
