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

// PromptLoader is a generic template loader for any prompt data type
type PromptLoader[T any] struct {
	templateFilePath string
	templateBuilder  func(T) (map[string]interface{}, error)
}

// NewPromptLoader creates a new generic PromptLoader instance
func NewPromptLoader[T any](templateFilePath string, templateBuilder func(T) (map[string]interface{}, error)) *PromptLoader[T] {
	return &PromptLoader[T]{
		templateFilePath: templateFilePath,
		templateBuilder:  templateBuilder,
	}
}

// LoadPromptTemplate loads and processes the prompt template with the provided data
func (l *PromptLoader[T]) LoadPromptTemplate(inputData T) (string, error) {
	// Load the template file
	templateContent, err := l.loadTemplateFile()
	if err != nil {
		return "", fmt.Errorf("failed to load template file: %w", err)
	}

	// Build template data directly from input
	templateData, err := l.templateBuilder(inputData)
	if err != nil {
		return "", fmt.Errorf("failed to build template data: %w", err)
	}

	// Execute template
	prompt, err := l.executeTemplate(templateContent, templateData)
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
func (l *PromptLoader[T]) executeTemplate(templateContent string, templateData map[string]interface{}) (string, error) {
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

// BuildTaskTemplateData builds template data for task generation
func BuildTaskTemplateData(story *story.Story, docs map[string]docs.ArchitectureDoc) (map[string]interface{}, error) {
	storyYAML, err := utils.MarshalToYAML(story)
	if err != nil {
		return nil, fmt.Errorf("failed to convert story to YAML: %w", err)
	}

	return map[string]interface{}{
		"StoryYAML":                storyYAML,
		"StoryID":                  story.ID,
		"Architecture":             docs["Architecture"].Content,
		"FrontendArchitecture":     docs["FrontendArchitecture"].Content,
		"CodingStandards":          docs["CodingStandards"].Content,
		"SourceTree":               docs["SourceTree"].Content,
		"TechStack":                docs["TechStack"].Content,
		"ArchitecturePath":         docs["Architecture"].FilePath,
		"FrontendArchitecturePath": docs["FrontendArchitecture"].FilePath,
		"CodingStandardsPath":      docs["CodingStandards"].FilePath,
		"SourceTreePath":           docs["SourceTree"].FilePath,
		"TechStackPath":            docs["TechStack"].FilePath,
	}, nil
}

// BuildDevNotesTemplateData builds template data for dev notes generation
func BuildDevNotesTemplateData(story *story.Story, tasks []story.Task, docs map[string]docs.ArchitectureDoc) (map[string]interface{}, error) {
	storyYAML, err := utils.MarshalToYAML(story)
	if err != nil {
		return nil, fmt.Errorf("failed to convert story to YAML: %w", err)
	}

	tasksYAML, err := utils.MarshalWithWrapper(tasks, "tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to convert tasks to YAML: %w", err)
	}

	return map[string]interface{}{
		"StoryYAML":                storyYAML,
		"StoryID":                  story.ID,
		"TasksYAML":                tasksYAML,
		"Architecture":             docs["Architecture"].Content,
		"FrontendArchitecture":     docs["FrontendArchitecture"].Content,
		"CodingStandards":          docs["CodingStandards"].Content,
		"SourceTree":               docs["SourceTree"].Content,
		"TechStack":                docs["TechStack"].Content,
		"ArchitecturePath":         docs["Architecture"].FilePath,
		"FrontendArchitecturePath": docs["FrontendArchitecture"].FilePath,
		"CodingStandardsPath":      docs["CodingStandards"].FilePath,
		"SourceTreePath":           docs["SourceTree"].FilePath,
		"TechStackPath":            docs["TechStack"].FilePath,
	}, nil
}
