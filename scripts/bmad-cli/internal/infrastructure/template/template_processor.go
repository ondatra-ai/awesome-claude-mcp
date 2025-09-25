package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"bmad-cli/internal/domain/models/story"
)

type TemplateProcessor struct {
	templatePath string
}

func NewTemplateProcessor(templatePath string) *TemplateProcessor {
	return &TemplateProcessor{
		templatePath: templatePath,
	}
}

func (p *TemplateProcessor) ProcessTemplate(doc *story.StoryDocument) (string, error) {
	// Get absolute path for template
	absTemplatePath, err := filepath.Abs(p.templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute template path: %w", err)
	}

	// Load template content
	templateContent, err := os.ReadFile(absTemplatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", absTemplatePath, err)
	}

	// Parse template
	tmpl, err := template.New("story").Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Create template data that matches the expected template structure
	// Template expects story fields at root level and other nested structures
	templateData := map[string]interface{}{
		// Story fields at root level
		"ID":                 doc.Story.ID,
		"Title":              doc.Story.Title,
		"AsA":                doc.Story.AsA,
		"IWant":              doc.Story.IWant,
		"SoThat":             doc.Story.SoThat,
		"Status":             doc.Story.Status,
		"AcceptanceCriteria": doc.Story.AcceptanceCriteria,

		// Other sections as nested structures
		"Tasks":          doc.Tasks,
		"DevNotes":       doc.DevNotes,
		"Testing":        doc.Testing,
		"ChangeLog":      doc.ChangeLog,
		"QAResults":      doc.QAResults,
		"DevAgentRecord": doc.DevAgentRecord,
	}

	// Execute template with flattened data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (p *TemplateProcessor) GetTemplatePath() string {
	return p.templatePath
}

func (p *TemplateProcessor) SetTemplatePath(path string) {
	p.templatePath = path
}

func (p *TemplateProcessor) ValidateTemplate() error {
	absPath, err := filepath.Abs(p.templatePath)
	if err != nil {
		return fmt.Errorf("failed to resolve template path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("template file does not exist: %s", absPath)
	}

	// Try to parse the template to ensure it's valid
	content, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	_, err = template.New("validation").Parse(string(content))
	if err != nil {
		return fmt.Errorf("template syntax error: %w", err)
	}

	return nil
}
