package prompt_builders

import (
	"fmt"
	"os"
	"path/filepath"
)

// PromptFileLoader loads raw prompt template files and checklists from disk
// This is distinct from infrastructure/template.TemplateLoader which executes Go templates
type PromptFileLoader struct {
	basePath string
}

func NewPromptFileLoader(basePath string) *PromptFileLoader {
	return &PromptFileLoader{basePath: basePath}
}

func (l *PromptFileLoader) LoadTemplate(templatePath string) (string, error) {
	tplPath := filepath.FromSlash(templatePath)
	if !filepath.IsAbs(tplPath) && l.basePath != "" {
		tplPath = filepath.Join(l.basePath, tplPath)
	}

	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", tplPath, err)
	}

	return string(tplBytes), nil
}

func (l *PromptFileLoader) LoadChecklist(checklistPath string) (string, error) {
	if checklistPath == "" {
		return "", nil
	}

	chkPath := filepath.FromSlash(checklistPath)
	if !filepath.IsAbs(chkPath) && l.basePath != "" {
		chkPath = filepath.Join(l.basePath, chkPath)
	}

	chkBytes, err := os.ReadFile(chkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read checklist file %s: %w", chkPath, err)
	}

	return string(chkBytes), nil
}
