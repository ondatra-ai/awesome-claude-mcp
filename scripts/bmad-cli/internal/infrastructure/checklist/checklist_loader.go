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

	slog.Debug("Checklist loaded successfully", "version", parsedChecklist.Version)

	return &parsedChecklist, nil
}

// ExtractAllPrompts extracts all prompts from the checklist with their context.
// Prompts with Skip field set are excluded.
func (l *ChecklistLoader) ExtractAllPrompts(chkList *checklist.Checklist) []checklist.PromptWithContext {
	prompts := make([]checklist.PromptWithContext, 0)

	prompts = l.extractTemplatePrompts(chkList, prompts)
	prompts = l.extractInvestPrompts(chkList, prompts)
	prompts = l.extractDependenciesPrompts(chkList, prompts)
	prompts = l.extractAcceptancePrompts(chkList, prompts)
	prompts = l.extractAntiPatternPrompts(chkList, prompts)
	prompts = l.extractDefinitionReadyPrompts(chkList, prompts)
	prompts = l.extractSplittingPrompts(chkList, prompts)

	slog.Debug("Extracted prompts from checklist", "count", len(prompts))

	return prompts
}

func (l *ChecklistLoader) extractTemplatePrompts(
	chkList *checklist.Checklist,
	prompts []checklist.PromptWithContext,
) []checklist.PromptWithContext {
	for _, criterion := range chkList.Template.Criteria {
		for _, prompt := range criterion.ValidationPrompts {
			if prompt.ShouldSkip() {
				continue
			}

			prompts = append(prompts, checklist.PromptWithContext{
				SectionID:     "template",
				SectionName:   "Template",
				CriterionID:   criterion.ID,
				CriterionName: criterion.Name,
				DefaultDocs:   chkList.DefaultDocs,
				Prompt:        prompt,
			})
		}
	}

	return prompts
}

func (l *ChecklistLoader) extractInvestPrompts(
	chkList *checklist.Checklist,
	prompts []checklist.PromptWithContext,
) []checklist.PromptWithContext {
	for _, criterion := range chkList.Invest.Criteria {
		for _, prompt := range criterion.ValidationPrompts {
			if prompt.ShouldSkip() {
				continue
			}

			prompts = append(prompts, checklist.PromptWithContext{
				SectionID:     "invest",
				SectionName:   "INVEST",
				CriterionID:   criterion.ID,
				CriterionName: criterion.Name,
				DefaultDocs:   chkList.DefaultDocs,
				Prompt:        prompt,
			})
		}
	}

	return prompts
}

func (l *ChecklistLoader) extractDependenciesPrompts(
	chkList *checklist.Checklist,
	prompts []checklist.PromptWithContext,
) []checklist.PromptWithContext {
	for _, prompt := range chkList.DependenciesRisk.ValidationPrompts {
		if prompt.ShouldSkip() {
			continue
		}

		prompts = append(prompts, checklist.PromptWithContext{
			SectionID:   "dependencies",
			SectionName: "Dependencies & Risks",
			DefaultDocs: chkList.DefaultDocs,
			Prompt:      prompt,
		})
	}

	if chkList.DependenciesRisk.RiskScoring != nil {
		for _, prompt := range chkList.DependenciesRisk.RiskScoring.ValidationPrompts {
			if prompt.ShouldSkip() {
				continue
			}

			prompts = append(prompts, checklist.PromptWithContext{
				SectionID:     "dependencies",
				SectionName:   "Dependencies & Risks",
				CriterionID:   "risk_scoring",
				CriterionName: "Risk Scoring",
				DefaultDocs:   chkList.DefaultDocs,
				Prompt:        prompt,
			})
		}
	}

	return prompts
}

func (l *ChecklistLoader) extractAcceptancePrompts(
	chkList *checklist.Checklist,
	prompts []checklist.PromptWithContext,
) []checklist.PromptWithContext {
	for _, prompt := range chkList.AcceptanceCrit.ValidationPrompts {
		if prompt.ShouldSkip() {
			continue
		}

		prompts = append(prompts, checklist.PromptWithContext{
			SectionID:   "acceptance",
			SectionName: "Acceptance Criteria",
			DefaultDocs: chkList.DefaultDocs,
			Prompt:      prompt,
		})
	}

	return prompts
}

func (l *ChecklistLoader) extractAntiPatternPrompts(
	chkList *checklist.Checklist,
	prompts []checklist.PromptWithContext,
) []checklist.PromptWithContext {
	for _, antiPattern := range chkList.AntiPatterns {
		for _, prompt := range antiPattern.ValidationPrompts {
			if prompt.ShouldSkip() {
				continue
			}

			prompts = append(prompts, checklist.PromptWithContext{
				SectionID:     "anti_patterns",
				SectionName:   "Anti-Patterns",
				CriterionID:   antiPattern.ID,
				CriterionName: antiPattern.Name,
				DefaultDocs:   chkList.DefaultDocs,
				Prompt:        prompt,
			})
		}
	}

	return prompts
}

func (l *ChecklistLoader) extractDefinitionReadyPrompts(
	chkList *checklist.Checklist,
	prompts []checklist.PromptWithContext,
) []checklist.PromptWithContext {
	for checklistKey, checklistPrompts := range chkList.DefinitionReady.Checklist {
		for _, prompt := range checklistPrompts {
			if prompt.ShouldSkip() {
				continue
			}

			prompts = append(prompts, checklist.PromptWithContext{
				SectionID:     "ready",
				SectionName:   "Definition of Ready",
				CriterionID:   checklistKey,
				CriterionName: checklistKey,
				DefaultDocs:   chkList.DefaultDocs,
				Prompt:        prompt,
			})
		}
	}

	return prompts
}

func (l *ChecklistLoader) extractSplittingPrompts(
	chkList *checklist.Checklist,
	prompts []checklist.PromptWithContext,
) []checklist.PromptWithContext {
	// When to split prompts
	for _, prompt := range chkList.Splitting.WhenToSplit {
		if prompt.ShouldSkip() {
			continue
		}

		prompts = append(prompts, checklist.PromptWithContext{
			SectionID:     "splitting",
			SectionName:   "Splitting",
			CriterionID:   "when_to_split",
			CriterionName: "When to Split",
			DefaultDocs:   chkList.DefaultDocs,
			Prompt:        prompt,
		})
	}

	// SPIDR techniques
	for techniqueID, technique := range chkList.Splitting.SPIDRTechniques {
		for _, prompt := range technique.Validation {
			if prompt.ShouldSkip() {
				continue
			}

			prompts = append(prompts, checklist.PromptWithContext{
				SectionID:     "splitting",
				SectionName:   "Splitting",
				CriterionID:   techniqueID,
				CriterionName: technique.Description,
				DefaultDocs:   chkList.DefaultDocs,
				Prompt:        prompt,
			})
		}
	}

	return prompts
}
