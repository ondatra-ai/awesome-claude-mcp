package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// PromptLoader is a generic template loader for any prompt data type
type PromptLoader[T any] struct {
	templateFilePath string
}

// NewPromptLoader creates a new generic PromptLoader instance
func NewPromptLoader[T any](templateFilePath string) *PromptLoader[T] {
	return &PromptLoader[T]{
		templateFilePath: templateFilePath,
	}
}

// LoadPromptTemplate loads and processes the prompt template with the provided data
func (l *PromptLoader[T]) LoadPromptTemplate(inputData T) (string, error) {
	// Load the template file
	templateContent, err := l.loadTemplateFile()
	if err != nil {
		return "", fmt.Errorf("failed to load template file: %w", err)
	}

	// Execute template directly with input data
	prompt, err := l.executeTemplate(templateContent, inputData)
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
func (l *PromptLoader[T]) executeTemplate(templateContent string, data T) (string, error) {
	// Parse the template
	tmpl, err := template.New("prompt").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
