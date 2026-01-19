package checklist

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/infrastructure/config"
)

// ChecklistLoader loads and parses the validation checklist YAML.
type ChecklistLoader struct {
	checklistPath string
}

// NewChecklistLoader creates a new checklist loader.
func NewChecklistLoader(cfg *config.ViperConfig) *ChecklistLoader {
	return &ChecklistLoader{
		checklistPath: cfg.GetString("paths.checklist"),
	}
}

// Load loads and parses the checklist YAML file.
func (l *ChecklistLoader) Load() (*checklist.Checklist, error) {
	slog.Debug("Loading checklist", "path", l.checklistPath)

	data, err := os.ReadFile(l.checklistPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read checklist file: %w", err)
	}

	var parsedChecklist checklist.Checklist

	err = yaml.Unmarshal(data, &parsedChecklist)
	if err != nil {
		return nil, fmt.Errorf("failed to parse checklist YAML: %w", err)
	}

	slog.Debug("Checklist loaded successfully",
		"version", parsedChecklist.Version,
		"stages", len(parsedChecklist.Stages),
	)

	return &parsedChecklist, nil
}

// ExtractAllPrompts extracts all prompts from the checklist with their context.
// Iterates through stages → sections → prompts.
// Prompts with Skip field set are excluded.
func (l *ChecklistLoader) ExtractAllPrompts(chkList *checklist.Checklist) []checklist.PromptWithContext {
	prompts := make([]checklist.PromptWithContext, 0)

	for _, stage := range chkList.Stages {
		for _, section := range stage.Sections {
			for _, prompt := range section.ValidationPrompts {
				if prompt.ShouldSkip() {
					continue
				}

				prompts = append(prompts, checklist.PromptWithContext{
					SectionID:     stage.ID,
					SectionName:   stage.Name,
					CriterionID:   section.ID,
					CriterionName: section.Name,
					DefaultDocs:   chkList.DefaultDocs,
					Prompt:        prompt,
				})
			}
		}
	}

	slog.Debug("Extracted prompts from checklist", "count", len(prompts))

	return prompts
}
