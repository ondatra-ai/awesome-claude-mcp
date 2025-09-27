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

// DevNotesTemplateLoader defines the interface for loading dev notes templates
type DevNotesTemplateLoader interface {
	LoadDevNotesPromptTemplate(story *story.Story, tasks []story.Task, architectureDocs map[string]docs.ArchitectureDoc) (string, error)
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
func (g *AIDevNotesGenerator) GenerateDevNotes(ctx context.Context, story *story.Story, tasks []story.Task, architectureDocs map[string]docs.ArchitectureDoc) (*story.DevNotes, error) {
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

	// Try to read from file first (new approach)
	devNotes, err := g.readDevNotesFromFile(story.ID)
	if err != nil {
		// Fall back to parsing the AI response using marker-based extraction
		fmt.Printf("🔄 File not found, falling back to marker-based extraction...\n")
		devNotes, err = g.parseDevNotesFromResponse(response)
		if err != nil {
			// If parsing fails, let's save the extracted content for debugging
			markerContent := g.extractMarkerBasedContent(response)
			contentFile := fmt.Sprintf("./tmp/%s-devnotes-extracted-content.yml", story.ID)
			if markerContent != "" {
				if err := os.WriteFile(contentFile, []byte(markerContent), 0644); err == nil {
					fmt.Printf("💾 Extracted marker content saved to: %s\n", contentFile)
				}
			}
			return nil, fmt.Errorf("failed to parse AI response: %w", err)
		}
	} else {
		fmt.Printf("✅ Dev notes read from file: ./tmp/%s-devnotes.yaml\n", story.ID)
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

// parseDevNotesFromResponse parses the AI response and extracts dev_notes using marker-based extraction
func (g *AIDevNotesGenerator) parseDevNotesFromResponse(response string) (*story.DevNotes, error) {
	// Find the content between FILE_START and FILE_END markers
	markerContent := g.extractMarkerBasedContent(response)
	if markerContent == "" {
		return nil, fmt.Errorf("no marker-based content found in AI response")
	}

	// Parse the YAML
	var devNotesData struct {
		DevNotes story.DevNotes `yaml:"dev_notes"`
	}

	err := yaml.Unmarshal([]byte(markerContent), &devNotesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate mandatory fields for main entities
	err = g.validateDevNotes(&devNotesData.DevNotes)
	if err != nil {
		return nil, fmt.Errorf("dev_notes validation failed: %w", err)
	}

	return &devNotesData.DevNotes, nil
}

// extractMarkerBasedContent extracts content between FILE_START and FILE_END markers
func (g *AIDevNotesGenerator) extractMarkerBasedContent(response string) string {
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

// validateDevNotes validates that mandatory entities have required source and description fields
func (g *AIDevNotesGenerator) validateDevNotes(devNotes *story.DevNotes) error {
	mandatoryEntities := []string{"technology_stack", "architecture", "file_structure"}

	for _, entityName := range mandatoryEntities {
		entity, exists := (*devNotes)[entityName]
		if !exists {
			return fmt.Errorf("mandatory entity '%s' is missing", entityName)
		}

		// Handle both map[string]interface{} and story.DevNotes (which is also map[string]interface{})
		var entityMap map[string]interface{}
		if em, ok := entity.(map[string]interface{}); ok {
			entityMap = em
		} else if dn, ok := entity.(story.DevNotes); ok {
			entityMap = dn
		} else {
			return fmt.Errorf("entity '%s' must be a map, got %T", entityName, entity)
		}

		// Check for mandatory source field
		if _, hasSource := entityMap["source"]; !hasSource {
			return fmt.Errorf("entity '%s' is missing mandatory 'source' field", entityName)
		}

		// Check for mandatory description field
		if _, hasDescription := entityMap["description"]; !hasDescription {
			return fmt.Errorf("entity '%s' is missing mandatory 'description' field", entityName)
		}
	}

	return nil
}

// readDevNotesFromFile reads dev_notes from file created by Claude
func (g *AIDevNotesGenerator) readDevNotesFromFile(storyID string) (*story.DevNotes, error) {
	filePath := fmt.Sprintf("./tmp/%s-devnotes.yaml", storyID)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dev_notes file not found: %s", filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dev_notes file: %w", err)
	}

	// Parse the YAML
	var devNotesData struct {
		DevNotes story.DevNotes `yaml:"dev_notes"`
	}

	err = yaml.Unmarshal(content, &devNotesData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dev_notes YAML: %w", err)
	}

	// Validate mandatory fields for main entities
	err = g.validateDevNotes(&devNotesData.DevNotes)
	if err != nil {
		return nil, fmt.Errorf("dev_notes validation failed: %w", err)
	}

	return &devNotesData.DevNotes, nil
}
