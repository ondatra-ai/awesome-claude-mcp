package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bmad-cli/internal/domain/models/story"
	"gopkg.in/yaml.v3"
)

// TaskPromptLoader loads and processes the task generation prompt template
type TaskPromptLoader struct {
	templatePath string
}

// NewTaskPromptLoader creates a new TaskPromptLoader instance
func NewTaskPromptLoader(templatePath string) *TaskPromptLoader {
	return &TaskPromptLoader{
		templatePath: templatePath,
	}
}

// LoadTaskPromptTemplate loads the task prompt template and injects story and architecture data
func (l *TaskPromptLoader) LoadTaskPromptTemplate(story *story.Story, architectureDocs map[string]string) (string, error) {
	// Load the template file
	templateContent, err := l.loadTemplateFile()
	if err != nil {
		return "", fmt.Errorf("failed to load template file: %w", err)
	}

	// Convert story to YAML for injection
	storyYAML, err := l.convertStoryToYAML(story)
	if err != nil {
		return "", fmt.Errorf("failed to convert story to YAML: %w", err)
	}

	// Replace placeholders in the template
	prompt := l.injectTemplateData(templateContent, storyYAML, architectureDocs)

	return prompt, nil
}

// loadTemplateFile loads the template file from disk
func (l *TaskPromptLoader) loadTemplateFile() (string, error) {
	templateFile := filepath.Join(l.templatePath, "us-create.tasks.prompt.tpl")
	content, err := os.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", templateFile, err)
	}
	return string(content), nil
}

// convertStoryToYAML converts the story struct to YAML format
func (l *TaskPromptLoader) convertStoryToYAML(story *story.Story) (string, error) {
	// Create a simple story structure for the template
	storyData := map[string]interface{}{
		"story": map[string]interface{}{
			"id":                 story.ID,
			"title":              story.Title,
			"status":             story.Status,
			"as_a":               story.AsA,
			"i_want":             story.IWant,
			"so_that":            story.SoThat,
			"acceptance_criteria": story.AcceptanceCriteria,
		},
	}

	yamlBytes, err := yaml.Marshal(storyData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal story to YAML: %w", err)
	}

	return string(yamlBytes), nil
}

// injectTemplateData replaces placeholders in the template with actual data
func (l *TaskPromptLoader) injectTemplateData(template, storyYAML string, architectureDocs map[string]string) string {
	result := template

	// Replace the story YAML block (find and replace the existing story block in template)
	result = l.replaceStoryBlock(result, storyYAML)

	// Replace architecture document placeholders
	for key, content := range architectureDocs {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		result = strings.ReplaceAll(result, placeholder, content)
	}

	// Handle any remaining empty placeholders
	placeholders := []string{
		"{{.Architecture}}",
		"{{.FrontendArchitecture}}",
		"{{.CodingStandards}}",
		"{{.SourceTree}}",
	}

	for _, placeholder := range placeholders {
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("# %s\n(Document not available)", strings.TrimSuffix(strings.TrimPrefix(placeholder, "{{."), "}}")))
		}
	}

	return result
}

// replaceStoryBlock replaces the existing story block in the template with the new story YAML
func (l *TaskPromptLoader) replaceStoryBlock(template, newStoryYAML string) string {
	// Find the story block between ```yaml and ``` that contains the story
	lines := strings.Split(template, "\n")
	var result []string
	inStoryBlock := false
	storyBlockFound := false

	for i, line := range lines {
		if strings.Contains(line, "```yaml") && !storyBlockFound {
			// Check if this yaml block contains a story by looking ahead
			if l.isStoryBlock(lines, i) {
				inStoryBlock = true
				storyBlockFound = true
				result = append(result, line)
				result = append(result, strings.Split(newStoryYAML, "\n")...)
				continue
			}
		}

		if inStoryBlock && strings.TrimSpace(line) == "```" {
			inStoryBlock = false
			result = append(result, line)
			continue
		}

		if !inStoryBlock {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// isStoryBlock checks if a YAML block starting at the given index contains a story
func (l *TaskPromptLoader) isStoryBlock(lines []string, startIndex int) bool {
	// Look ahead in the YAML block for story-related content
	for i := startIndex + 1; i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```"); i++ {
		if strings.Contains(lines[i], "story:") || strings.Contains(lines[i], "id:") || strings.Contains(lines[i], "title:") {
			return true
		}
	}
	return false
}
