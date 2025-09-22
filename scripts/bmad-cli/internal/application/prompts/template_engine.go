package prompts

import (
	"fmt"
	"strconv"
	"strings"

	"bmad-cli/internal/domain/models"
)

type TemplateEngine struct {
	loader *TemplateLoader
}

func NewTemplateEngine() *TemplateEngine {
	loader := NewTemplateLoader("")
	return &TemplateEngine{loader: loader}
}

func (e *TemplateEngine) BuildFromTemplate(threadCtx models.ThreadContext, templatePath, checklistPath string) (string, error) {
	template, err := e.loader.LoadTemplate(templatePath)
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

	return e.fillPlaceholders(template, threadCtx, checklist), nil
}

func (e *TemplateEngine) fillPlaceholders(template string, threadCtx models.ThreadContext, checklist string) string {
	prompt := template

	prompt = strings.ReplaceAll(prompt, "{{PR_NUMBER}}", strconv.Itoa(threadCtx.PRNumber))

	loc := fmt.Sprintf("%s:%d", threadCtx.Comment.File, threadCtx.Comment.Line)
	prompt = strings.ReplaceAll(prompt, "{{LOCATION}}", loc)
	prompt = strings.ReplaceAll(prompt, "{{URL}}", threadCtx.Comment.URL)
	prompt = strings.ReplaceAll(prompt, "{{CONVERSATION_TEXT}}", e.joinAllComments(threadCtx.Thread))

	if checklist != "" {
		prompt = strings.ReplaceAll(prompt, "{{CHECKLIST_MD}}", checklist)
	}

	return prompt
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
