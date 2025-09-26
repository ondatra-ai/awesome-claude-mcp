package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"bmad-cli/internal/domain/models/story"
	"gopkg.in/yaml.v3"
)

// DevNotesTemplateLoader defines the interface for loading dev notes templates
type DevNotesTemplateLoader interface {
	LoadDevNotesPromptTemplate(story *story.Story, tasks []story.Task, architectureDocs map[string]string) (string, error)
}

// AIDevNotesGenerator generates story dev_notes using AI based on templates
type AIDevNotesGenerator struct {
	aiClient       AIClient
	templateLoader DevNotesTemplateLoader
}

// NewDevNotesGenerator creates a new AIDevNotesGenerator instance
func NewDevNotesGenerator(aiClient AIClient, templateLoader DevNotesTemplateLoader) *AIDevNotesGenerator {
	return &AIDevNotesGenerator{
		aiClient:       aiClient,
		templateLoader: templateLoader,
	}
}

// GenerateDevNotes generates story dev_notes using AI based on the story, tasks, and architecture documents
func (g *AIDevNotesGenerator) GenerateDevNotes(ctx context.Context, story *story.Story, tasks []story.Task, architectureDocs map[string]string) (*story.DevNotes, error) {
	// Load and prepare the prompt template
	prompt, err := g.templateLoader.LoadDevNotesPromptTemplate(story, tasks, architectureDocs)
	if err != nil {
		return nil, fmt.Errorf("failed to load dev notes prompt template: %w", err)
	}

	// Generate dev_notes using AI
	response, err := g.aiClient.GenerateContent(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dev_notes with AI: %w", err)
	}

	// Create tmp directory if it doesn't exist
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		return nil, fmt.Errorf("failed to create tmp directory: %w", err)
	}

	// Write full AI response to file for debugging
	responseFile := fmt.Sprintf("./tmp/%s-devnotes-full-response.txt", story.ID)
	if err := os.WriteFile(responseFile, []byte(response), 0644); err != nil {
		return nil, fmt.Errorf("failed to write response file: %w", err)
	}
	fmt.Printf("💾 Full AI dev_notes response saved to: %s\n", responseFile)

	// Parse the AI response
	devNotes, err := g.parseDevNotesFromResponse(response)
	if err != nil {
		// If parsing fails, let's save the extracted YAML block for debugging
		yamlBlock := g.extractYAMLBlock(response)
		yamlFile := fmt.Sprintf("./tmp/%s-devnotes-extracted-yaml.yml", story.ID)
		if yamlBlock != "" {
			if err := os.WriteFile(yamlFile, []byte(yamlBlock), 0644); err == nil {
				fmt.Printf("💾 Extracted dev_notes YAML block saved to: %s\n", yamlFile)
			}
		}
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// Save successfully parsed dev_notes to YAML file
	devNotesMap := map[string]interface{}{"dev_notes": devNotes}
	if devNotesYAML, yamlErr := yaml.Marshal(devNotesMap); yamlErr == nil {
		devNotesFile := fmt.Sprintf("./tmp/%s-devnotes.yml", story.ID)
		if writeErr := os.WriteFile(devNotesFile, devNotesYAML, 0644); writeErr == nil {
			fmt.Printf("✅ Parsed dev_notes saved to: %s\n", devNotesFile)
		}
	}

	return devNotes, nil
}

// parseDevNotesFromResponse parses the AI response and extracts dev_notes
func (g *AIDevNotesGenerator) parseDevNotesFromResponse(response string) (*story.DevNotes, error) {
	// Find the YAML block in the response
	yamlBlock := g.extractYAMLBlock(response)
	if yamlBlock == "" {
		return nil, fmt.Errorf("no YAML block found in AI response")
	}

	// Parse the YAML
	var devNotesData struct {
		DevNotes story.DevNotes `yaml:"dev_notes"`
	}

	err := yaml.Unmarshal([]byte(yamlBlock), &devNotesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &devNotesData.DevNotes, nil
}

// extractYAMLBlock extracts the YAML block from the AI response
func (g *AIDevNotesGenerator) extractYAMLBlock(response string) string {
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
		// Fallback: try to find dev_notes: section directly
		for i, line := range lines {
			if strings.TrimSpace(line) == "dev_notes:" {
				// Found dev_notes section, collect until empty line or end
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
