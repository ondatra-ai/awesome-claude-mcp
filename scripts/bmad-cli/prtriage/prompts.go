package prtriage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// buildPromptFromTemplate loads a template file, optionally loads a checklist,
// and fills placeholders with the given thread context.
func buildPromptFromTemplate(threadCtx ThreadContext, templatePath, checklistPath string) (string, error) {
	tplPath := filepath.FromSlash(templatePath)

	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", tplPath, err)
	}

	prompt := string(tplBytes)

	// Replace common placeholders
	prompt = strings.ReplaceAll(prompt, "{{PR_NUMBER}}", strconv.Itoa(threadCtx.PRNumber))
	loc := fmt.Sprintf("%s:%d", threadCtx.Comment.File, threadCtx.Comment.Line)
	prompt = strings.ReplaceAll(prompt, "{{LOCATION}}", loc)
	prompt = strings.ReplaceAll(prompt, "{{URL}}", threadCtx.Comment.URL)
	prompt = strings.ReplaceAll(prompt, "{{CONVERSATION_TEXT}}", joinAllComments(threadCtx.Thread))

	// If checklist path is provided, load and replace it
	if checklistPath != "" {
		chkPath := filepath.FromSlash(checklistPath)

		chkBytes, err := os.ReadFile(chkPath)
		if err != nil {
			return "", fmt.Errorf("failed to read checklist file %s: %w", chkPath, err)
		}

		prompt = strings.ReplaceAll(prompt, "{{CHECKLIST_MD}}", string(chkBytes))
	}

	return prompt, nil
}

// buildHeuristicPrompt loads the template and checklist, fills placeholders, and returns the final prompt.
func buildHeuristicPrompt(threadCtx ThreadContext) (string, error) {
	return buildPromptFromTemplate(
		threadCtx,
		"scripts/pr-triage/heuristic.prompt.tpl",
		".bmad-core/checklists/triage-heuristic-checklist.md",
	)
}

// buildImplementCodePrompt loads scripts/pr-triage/apply.prompt.tpl and fills
// placeholders with the given thread context for implementation.
func buildImplementCodePrompt(threadCtx ThreadContext) (string, error) {
	return buildPromptFromTemplate(threadCtx, "scripts/pr-triage/apply.prompt.tpl", "")
}

// joinAllComments concatenates all comments in a thread with separators.
func joinAllComments(thread Thread) string {
	var builder strings.Builder

	for i, comment := range thread.Comments {
		if i > 0 {
			builder.WriteString("\n---\n")
		}

		builder.WriteString(comment.Body)
	}

	return builder.String()
}
