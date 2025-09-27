package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
	"gopkg.in/yaml.v3"
)

// AIClient defines the interface for AI communication
type AIClient interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

// TemplateLoader defines the interface for loading templates
type TemplateLoader interface {
	LoadTaskPromptTemplate(story *story.Story, architectureDocs map[string]docs.ArchitectureDoc) (string, error)
}

// AITaskGenerator generates story tasks using AI based on templates
type AITaskGenerator struct {
	aiClient       AIClient
	templateLoader TemplateLoader
}

// NewTaskGenerator creates a new AITaskGenerator instance
func NewTaskGenerator(aiClient AIClient, templateLoader TemplateLoader) *AITaskGenerator {
	return &AITaskGenerator{
		aiClient:       aiClient,
		templateLoader: templateLoader,
	}
}

// GenerateTasks generates story tasks using AI based on the story and architecture documents
func (g *AITaskGenerator) GenerateTasks(ctx context.Context, story *story.Story, architectureDocs map[string]docs.ArchitectureDoc) ([]story.Task, error) {
	// Load and prepare the prompt template
	prompt, err := g.templateLoader.LoadTaskPromptTemplate(story, architectureDocs)
	if err != nil {
		return nil, fmt.Errorf("failed to load task prompt template: %w", err)
	}

	// Generate tasks using AI
	response, err := g.aiClient.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tasks with AI: %w", err)
	}

	// Create tmp directory if it doesn't exist
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		return nil, fmt.Errorf("failed to create tmp directory: %w", err)
	}

	// Write full AI response to file for debugging
	responseFile := fmt.Sprintf("./tmp/%s-full-response.txt", story.ID)
	if err := os.WriteFile(responseFile, []byte(response), 0644); err != nil {
		return nil, fmt.Errorf("failed to write response file: %w", err)
	}
	fmt.Printf("💾 Full AI response saved to: %s\n", responseFile)

	// Try to read from file first (new approach)
	tasks, err := g.readTasksFromFile(story.ID)
	if err != nil {
		// Fall back to parsing the AI response using marker-based extraction
		fmt.Printf("🔄 File not found, falling back to marker-based extraction...\n")
		tasks, err = g.parseTasksFromResponse(response)
		if err != nil {
			// If parsing fails, let's save the extracted content for debugging
			markerContent := g.extractMarkerBasedContent(response)
			contentFile := fmt.Sprintf("./tmp/%s-extracted-content.yml", story.ID)
			if markerContent != "" {
				if err := os.WriteFile(contentFile, []byte(markerContent), 0644); err == nil {
					fmt.Printf("💾 Extracted marker content saved to: %s\n", contentFile)
				}
			}
			return nil, fmt.Errorf("failed to parse AI response: %w", err)
		}
	} else {
		fmt.Printf("✅ Tasks read from file: ./tmp/%s-tasks.yaml\n", story.ID)
	}

	// Save successfully parsed tasks to YAML file
	taskMap := map[string]interface{}{"tasks": tasks}
	if tasksYAML, yamlErr := yaml.Marshal(taskMap); yamlErr == nil {
		tasksFile := fmt.Sprintf("./tmp/%s-tasks.yml", story.ID)
		if writeErr := os.WriteFile(tasksFile, tasksYAML, 0644); writeErr == nil {
			fmt.Printf("✅ Parsed tasks saved to: %s\n", tasksFile)
		}
	}

	// Validate that we have at least one task
	if len(tasks) == 0 {
		return nil, fmt.Errorf("AI generated no tasks")
	}

	return tasks, nil
}

// parseTasksFromResponse parses the AI response and extracts tasks using marker-based extraction
func (g *AITaskGenerator) parseTasksFromResponse(response string) ([]story.Task, error) {
	// Find the content between FILE_START and FILE_END markers
	markerContent := g.extractMarkerBasedContent(response)
	if markerContent == "" {
		return nil, fmt.Errorf("no marker-based content found in AI response")
	}

	// Parse the YAML
	var taskData struct {
		Tasks []story.Task `yaml:"tasks"`
	}

	err := yaml.Unmarshal([]byte(markerContent), &taskData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return taskData.Tasks, nil
}

// extractMarkerBasedContent extracts content between FILE_START and FILE_END markers
func (g *AITaskGenerator) extractMarkerBasedContent(response string) string {
	lines := strings.Split(response, "\n")
	var contentLines []string
	inFileContent := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start of file content
		if strings.HasPrefix(trimmed, "=== FILE_START:") {
			inFileContent = true
			continue
		}

		// End of file content
		if inFileContent && strings.HasPrefix(trimmed, "=== FILE_END:") {
			break
		}

		// Collect file content
		if inFileContent {
			contentLines = append(contentLines, line)
		}
	}

	return strings.Join(contentLines, "\n")
}

// readTasksFromFile reads tasks from file created by Claude
func (g *AITaskGenerator) readTasksFromFile(storyID string) ([]story.Task, error) {
	filePath := fmt.Sprintf("./tmp/%s-tasks.yaml", storyID)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("tasks file not found: %s", filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	// Parse the YAML
	var taskData struct {
		Tasks []story.Task `yaml:"tasks"`
	}

	err = yaml.Unmarshal(content, &taskData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks YAML: %w", err)
	}

	return taskData.Tasks, nil
}
