package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"bmad-cli/internal/domain/models/story"
	"gopkg.in/yaml.v3"
)

// TaskPromptLoader loads and processes the task generation prompt template
type TaskPromptLoader struct {
	templateFilePath string
}

// NewTaskPromptLoader creates a new TaskPromptLoader instance
func NewTaskPromptLoader(templateFilePath string) *TaskPromptLoader {
	return &TaskPromptLoader{
		templateFilePath: templateFilePath,
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

	// Use proper template system to inject data
	prompt, err := l.executeTemplate(templateContent, storyYAML, architectureDocs)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return prompt, nil
}

// loadTemplateFile loads the template file from disk
func (l *TaskPromptLoader) loadTemplateFile() (string, error) {
	content, err := os.ReadFile(l.templateFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", l.templateFilePath, err)
	}
	return string(content), nil
}

// convertStoryToYAML converts the story struct to YAML format
func (l *TaskPromptLoader) convertStoryToYAML(storyObj *story.Story) (string, error) {
	yamlBytes, err := yaml.Marshal(storyObj)
	if err != nil {
		return "", fmt.Errorf("failed to marshal story to YAML: %w", err)
	}

	return string(yamlBytes), nil
}

// executeTemplate uses Go's text/template system to properly inject data
func (l *TaskPromptLoader) executeTemplate(templateContent, storyYAML string, architectureDocs map[string]string) (string, error) {
	// Create template data structure
	templateData := struct {
		StoryYAML            string
		Architecture         string
		FrontendArchitecture string
		CodingStandards      string
		SourceTree           string
		TechStack            string
	}{
		StoryYAML:            storyYAML,
		Architecture:         architectureDocs["Architecture"],
		FrontendArchitecture: architectureDocs["FrontendArchitecture"],
		CodingStandards:      architectureDocs["CodingStandards"],
		SourceTree:           architectureDocs["SourceTree"],
		TechStack:            architectureDocs["TechStack"],
	}

	// Parse the template
	tmpl, err := template.New("task-prompt").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
