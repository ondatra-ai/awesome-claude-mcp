package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"bmad-cli/internal/domain/models/story"
	"gopkg.in/yaml.v3"
)

// DevNotesPromptLoader loads and processes the dev notes generation prompt template
type DevNotesPromptLoader struct {
	templateFilePath string
}

// NewDevNotesPromptLoader creates a new DevNotesPromptLoader instance
func NewDevNotesPromptLoader(templateFilePath string) *DevNotesPromptLoader {
	return &DevNotesPromptLoader{
		templateFilePath: templateFilePath,
	}
}

// LoadDevNotesPromptTemplate loads the dev notes prompt template and injects story and architecture data
func (l *DevNotesPromptLoader) LoadDevNotesPromptTemplate(story *story.Story, tasks []story.Task, architectureDocs map[string]string) (string, error) {
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

	// Convert tasks to YAML for injection
	tasksYAML, err := l.convertTasksToYAML(tasks)
	if err != nil {
		return "", fmt.Errorf("failed to convert tasks to YAML: %w", err)
	}

	// Use proper template system to inject data
	prompt, err := l.executeTemplate(templateContent, storyYAML, tasksYAML, architectureDocs)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return prompt, nil
}

// loadTemplateFile loads the template file from disk
func (l *DevNotesPromptLoader) loadTemplateFile() (string, error) {
	content, err := os.ReadFile(l.templateFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", l.templateFilePath, err)
	}
	return string(content), nil
}

// convertStoryToYAML converts the story struct to YAML format
func (l *DevNotesPromptLoader) convertStoryToYAML(storyObj *story.Story) (string, error) {
	yamlBytes, err := yaml.Marshal(storyObj)
	if err != nil {
		return "", fmt.Errorf("failed to marshal story to YAML: %w", err)
	}

	return string(yamlBytes), nil
}

// convertTasksToYAML converts the tasks slice to YAML format
func (l *DevNotesPromptLoader) convertTasksToYAML(tasks []story.Task) (string, error) {
	taskMap := map[string]interface{}{"tasks": tasks}
	yamlBytes, err := yaml.Marshal(taskMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tasks to YAML: %w", err)
	}

	return string(yamlBytes), nil
}

// executeTemplate uses Go's text/template system to properly inject data
func (l *DevNotesPromptLoader) executeTemplate(templateContent, storyYAML, tasksYAML string, architectureDocs map[string]string) (string, error) {
	// Create template data structure
	templateData := struct {
		StoryYAML            string
		TasksYAML            string
		Architecture         string
		FrontendArchitecture string
		CodingStandards      string
		SourceTree           string
		TechStack            string
	}{
		StoryYAML:            storyYAML,
		TasksYAML:            tasksYAML,
		Architecture:         architectureDocs["Architecture"],
		FrontendArchitecture: architectureDocs["FrontendArchitecture"],
		CodingStandards:      architectureDocs["CodingStandards"],
		SourceTree:           architectureDocs["SourceTree"],
		TechStack:            architectureDocs["TechStack"],
	}

	// Parse the template
	tmpl, err := template.New("dev-notes-prompt").Parse(templateContent)
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
