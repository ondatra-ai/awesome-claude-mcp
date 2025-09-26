package prompts

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"bmad-cli/internal/domain/models"
)

type TemplateEngine struct {
	loader *TemplateLoader
}

func NewTemplateEngine() *TemplateEngine {
	loader := NewTemplateLoader("")
	return &TemplateEngine{loader: loader}
}

type TemplateData struct {
	PRNumber         int
	Location         string
	URL              string
	ConversationText string
	ChecklistMD      string
}

func (e *TemplateEngine) BuildFromTemplate(threadCtx models.ThreadContext, templatePath, checklistPath string) (string, error) {
	templateContent, err := e.loader.LoadTemplate(templatePath)
	if err != nil {
		return "", err
	}

	checklist := ""
	if checklistPath != "" {
		checklist, err = e.loader.LoadChecklist(checklistPath)
		if err != nil {
			return "", err
		}
	}

	tmpl, err := template.New("prompt").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := TemplateData{
		PRNumber:         threadCtx.PRNumber,
		Location:         fmt.Sprintf("%s:%d", threadCtx.Comment.File, threadCtx.Comment.Line),
		URL:              threadCtx.Comment.URL,
		ConversationText: e.joinAllComments(threadCtx.Thread),
		ChecklistMD:      checklist,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (e *TemplateEngine) joinAllComments(thread models.Thread) string {
	var builder strings.Builder

	for i, comment := range thread.Comments {
		if i > 0 {
			builder.WriteString("\n---\n")
		}
		builder.WriteString(comment.Body)
	}

	return builder.String()
}
