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
func (l *TaskPromptLoader) convertStoryToYAML(storyObj *story.Story) (string, error) {
	// Create a wrapper struct with "story" key to match template expectation
	storyWrapper := struct {
		Story *story.Story `yaml:"story"`
	}{
		Story: storyObj,
	}

	yamlBytes, err := yaml.Marshal(storyWrapper)
	if err != nil {
		return "", fmt.Errorf("failed to marshal story to YAML: %w", err)
	}

	return string(yamlBytes), nil
}

// injectTemplateData replaces placeholders in the template with actual data
func (l *TaskPromptLoader) injectTemplateData(template, storyYAML string, architectureDocs map[string]string) string {
	result := template

	// Replace the story YAML placeholder
	result = strings.ReplaceAll(result, "{{.StoryYAML}}", storyYAML)

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
