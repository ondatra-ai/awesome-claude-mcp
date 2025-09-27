package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"bmad-cli/internal/common/utils"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
)

// PromptData represents the base data structure for template execution
type PromptData struct {
	Story *story.Story
	Docs  map[string]docs.ArchitectureDoc
	Extra map[string]interface{}
}

// PromptLoader is a generic template loader for any prompt data type
type PromptLoader[T any] struct {
	templateFilePath string
	dataConverter    func(T) (*PromptData, error)
}

// NewPromptLoader creates a new generic PromptLoader instance
func NewPromptLoader[T any](templateFilePath string, dataConverter func(T) (*PromptData, error)) *PromptLoader[T] {
	return &PromptLoader[T]{
		templateFilePath: templateFilePath,
		dataConverter:    dataConverter,
	}
}

// LoadPromptTemplate loads and processes the prompt template with the provided data
func (l *PromptLoader[T]) LoadPromptTemplate(inputData T) (string, error) {
	// Convert input data to standard prompt data structure
	promptData, err := l.dataConverter(inputData)
	if err != nil {
		return "", fmt.Errorf("failed to convert input data: %w", err)
	}

	// Load the template file
	templateContent, err := l.loadTemplateFile()
	if err != nil {
		return "", fmt.Errorf("failed to load template file: %w", err)
	}

	// Convert story to YAML for injection
	storyYAML, err := utils.MarshalToYAML(promptData.Story)
	if err != nil {
		return "", fmt.Errorf("failed to convert story to YAML: %w", err)
	}

	// Execute template with unified data structure
	prompt, err := l.executeTemplate(templateContent, storyYAML, promptData)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return prompt, nil
}

// loadTemplateFile loads the template file from disk
func (l *PromptLoader[T]) loadTemplateFile() (string, error) {
	content, err := os.ReadFile(l.templateFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", l.templateFilePath, err)
	}
	return string(content), nil
}

// executeTemplate uses Go's text/template system to properly inject data
func (l *PromptLoader[T]) executeTemplate(templateContent, storyYAML string, promptData *PromptData) (string, error) {
	// Create base template data structure
	templateData := map[string]interface{}{
		"StoryYAML":                storyYAML,
		"StoryID":                  promptData.Story.ID,
		"Architecture":             promptData.Docs["Architecture"].Content,
		"FrontendArchitecture":     promptData.Docs["FrontendArchitecture"].Content,
		"CodingStandards":          promptData.Docs["CodingStandards"].Content,
		"SourceTree":               promptData.Docs["SourceTree"].Content,
		"TechStack":                promptData.Docs["TechStack"].Content,
		"ArchitecturePath":         promptData.Docs["Architecture"].FilePath,
		"FrontendArchitecturePath": promptData.Docs["FrontendArchitecture"].FilePath,
		"CodingStandardsPath":      promptData.Docs["CodingStandards"].FilePath,
		"SourceTreePath":           promptData.Docs["SourceTree"].FilePath,
		"TechStackPath":            promptData.Docs["TechStack"].FilePath,
	}

	// Add any extra data from the converter
	for key, value := range promptData.Extra {
		templateData[key] = value
	}

	// Parse the template
	tmpl, err := template.New("prompt").Parse(templateContent)
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

// TaskPromptDataConverter converts task-specific data to PromptData
func TaskPromptDataConverter(story *story.Story, docs map[string]docs.ArchitectureDoc) (*PromptData, error) {
	return &PromptData{
		Story: story,
		Docs:  docs,
		Extra: map[string]interface{}{},
	}, nil
}

// DevNotesPromptDataConverter converts dev notes-specific data to PromptData
func DevNotesPromptDataConverter(story *story.Story, tasks []story.Task, docs map[string]docs.ArchitectureDoc) (*PromptData, error) {
	tasksYAML, err := utils.MarshalWithWrapper(tasks, "tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to convert tasks to YAML: %w", err)
	}

	return &PromptData{
		Story: story,
		Docs:  docs,
		Extra: map[string]interface{}{
			"TasksYAML": tasksYAML,
		},
	}, nil
}
