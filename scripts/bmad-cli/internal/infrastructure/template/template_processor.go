package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"bmad-cli/internal/domain/models/story"
	"gopkg.in/yaml.v3"
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

	// Create custom template functions
	funcMap := template.FuncMap{
		"toYaml": func(v interface{}) string {
			data, err := yaml.Marshal(v)
			if err != nil {
				return fmt.Sprintf("# Error marshaling to YAML: %v", err)
			}
			return string(data)
		},
		"nindent": func(spaces int, text string) string {
			lines := strings.Split(text, "\n")
			indent := strings.Repeat(" ", spaces)
			for i, line := range lines {
				if line != "" {
					lines[i] = indent + line
				}
			}
			return strings.Join(lines, "\n")
		},
	}

	// Parse template with custom functions
	tmpl, err := template.New("story").Funcs(funcMap).Parse(string(templateContent))
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
