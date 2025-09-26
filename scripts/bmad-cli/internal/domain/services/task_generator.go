package services

import (
	"context"
	"fmt"
	"strings"

	"bmad-cli/internal/domain/models/story"
	"gopkg.in/yaml.v3"
)

// AIClient defines the interface for AI communication
type AIClient interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}

// TemplateLoader defines the interface for loading templates
type TemplateLoader interface {
	LoadTaskPromptTemplate(story *story.Story, architectureDocs map[string]string) (string, error)
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
func (g *AITaskGenerator) GenerateTasks(ctx context.Context, story *story.Story, architectureDocs map[string]string) ([]story.Task, error) {
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

	// Parse the AI response
	tasks, err := g.parseTasksFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Validate that we have at least one task
	if len(tasks) == 0 {
		return nil, fmt.Errorf("AI generated no tasks")
	}

	return tasks, nil
}

// parseTasksFromResponse parses the AI response and extracts tasks
func (g *AITaskGenerator) parseTasksFromResponse(response string) ([]story.Task, error) {
	// Find the YAML block in the response
	yamlBlock := g.extractYAMLBlock(response)
	if yamlBlock == "" {
		return nil, fmt.Errorf("no YAML block found in AI response")
	}

	// Parse the YAML
	var taskData struct {
		Tasks []story.Task `yaml:"tasks"`
	}

	err := yaml.Unmarshal([]byte(yamlBlock), &taskData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return taskData.Tasks, nil
}

// extractYAMLBlock extracts the YAML block from the AI response
func (g *AITaskGenerator) extractYAMLBlock(response string) string {
	// Look for YAML code blocks (```yaml or ```yml)
	lines := strings.Split(response, "\n")
	var yamlLines []string
	inYamlBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Start of YAML block
		if strings.HasPrefix(trimmed, "```yaml") || strings.HasPrefix(trimmed, "```yml") {
			inYamlBlock = true
			continue
		}

		// End of code block
		if inYamlBlock && strings.HasPrefix(trimmed, "```") {
			break
		}

		// Collect YAML content
		if inYamlBlock {
			yamlLines = append(yamlLines, line)
		}
	}

	if len(yamlLines) == 0 {
		// Fallback: try to find tasks: section directly
		for i, line := range lines {
			if strings.TrimSpace(line) == "tasks:" {
				// Found tasks section, collect until empty line or end
				for j := i; j < len(lines); j++ {
					yamlLines = append(yamlLines, lines[j])
					if j+1 < len(lines) && strings.TrimSpace(lines[j+1]) == "" {
						break
					}
				}
				break
			}
		}
	}

	return strings.Join(yamlLines, "\n")
}
