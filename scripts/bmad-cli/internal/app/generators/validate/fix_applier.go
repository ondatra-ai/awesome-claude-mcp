package validate

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const fixApplierFilePermissions = 0o644

// FixApplierData represents data needed for fix applier templates.
type FixApplierData struct {
	Story      *story.Story // Current story to modify
	FixPrompt  string       // The fix prompt to apply
	ResultPath string       // Path for FILE_START/FILE_END output
}

// FixApplier applies fix prompts to stories using AI.
type FixApplier struct {
	aiClient     ports.AIPort
	config       *config.ViperConfig
	modeFactory  *ai.ModeFactory
	systemLoader *template.TemplateLoader[FixApplierData]
	userLoader   *template.TemplateLoader[FixApplierData]
	tmpDir       string
}

// NewFixApplier creates a new fix applier.
func NewFixApplier(aiClient ports.AIPort, cfg *config.ViperConfig) *FixApplier {
	systemTemplatePath := cfg.GetString("templates.prompts.fix_applier_system")
	userTemplatePath := cfg.GetString("templates.prompts.fix_applier")

	return &FixApplier{
		aiClient:     aiClient,
		config:       cfg,
		modeFactory:  ai.NewModeFactory(cfg),
		systemLoader: template.NewTemplateLoader[FixApplierData](systemTemplatePath),
		userLoader:   template.NewTemplateLoader[FixApplierData](userTemplatePath),
	}
}

// Apply applies a fix prompt to the story and returns the updated story.
func (a *FixApplier) Apply(
	ctx context.Context,
	storyData *story.Story,
	fixPrompt string,
	tmpDir string,
	iteration int,
) (*story.Story, error) {
	a.tmpDir = tmpDir

	slog.Info("Applying fix prompt to story",
		"storyID", storyData.ID,
		"iteration", iteration,
	)

	resultPath := fmt.Sprintf("%s/apply-%s-iter%d-result.yaml", tmpDir, storyData.ID, iteration)

	promptData := FixApplierData{
		Story:      storyData,
		FixPrompt:  fixPrompt,
		ResultPath: resultPath,
	}

	systemPrompt, err := a.systemLoader.LoadTemplate(promptData)
	if err != nil {
		return nil, pkgerrors.ErrLoadChecklistSystemPromptFailed(err)
	}

	userPrompt, err := a.userLoader.LoadTemplate(promptData)
	if err != nil {
		return nil, pkgerrors.ErrLoadChecklistUserPromptFailed(err)
	}

	// Save prompts for debugging
	a.savePromptFile(storyData.ID, iteration, "system", systemPrompt)
	a.savePromptFile(storyData.ID, iteration, "user", userPrompt)

	mode := a.modeFactory.GetThinkMode()

	response, err := a.aiClient.ExecutePromptWithSystem(ctx, systemPrompt, userPrompt, "", mode)
	if err != nil {
		return nil, pkgerrors.ErrChecklistAIEvaluationFailed(err)
	}

	// Save response for debugging
	a.savePromptFile(storyData.ID, iteration, "response", response)

	// Parse updated acceptance criteria from response
	updatedACs, err := a.parseUpdatedACs(response, resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated acceptance criteria: %w", err)
	}

	// Create updated story with new ACs
	updatedStory := *storyData // Copy story
	updatedStory.AcceptanceCriteria = updatedACs

	slog.Info("Fix applied successfully",
		"storyID", storyData.ID,
		"originalACCount", len(storyData.AcceptanceCriteria),
		"updatedACCount", len(updatedACs),
	)

	return &updatedStory, nil
}

// parseUpdatedACs extracts and parses updated acceptance criteria from AI response.
func (a *FixApplier) parseUpdatedACs(response, resultPath string) ([]story.AcceptanceCriterion, error) {
	content := a.extractFileContent(response, resultPath)
	if content == "" {
		return nil, pkgerrors.ErrFixApplierNoContentFound(resultPath)
	}

	// Save the extracted content
	err := os.WriteFile(resultPath, []byte(content), fixApplierFilePermissions)
	if err != nil {
		slog.Warn("Failed to save result file", "path", resultPath, "error", err)
	}

	// Parse YAML as acceptance criteria array
	var acs []story.AcceptanceCriterion

	err = yaml.Unmarshal([]byte(content), &acs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse acceptance criteria YAML: %w", err)
	}

	return acs, nil
}

// extractFileContent extracts content between FILE_START and FILE_END markers.
func (a *FixApplier) extractFileContent(response, path string) string {
	startMarker := fmt.Sprintf("=== FILE_START: %s ===", path)
	endMarker := fmt.Sprintf("=== FILE_END: %s ===", path)

	startIdx := strings.Index(response, startMarker)
	if startIdx == -1 {
		return ""
	}

	contentStart := startIdx + len(startMarker)
	endIdx := strings.Index(response[contentStart:], endMarker)

	if endIdx == -1 {
		return ""
	}

	return strings.TrimSpace(response[contentStart : contentStart+endIdx])
}

// savePromptFile saves a prompt file for debugging.
func (a *FixApplier) savePromptFile(storyID string, iteration int, suffix, content string) {
	if a.tmpDir == "" {
		return
	}

	filePath := fmt.Sprintf("%s/apply-%s-iter%d-%s.txt", a.tmpDir, storyID, iteration, suffix)

	err := os.WriteFile(filePath, []byte(content), fixApplierFilePermissions)
	if err != nil {
		slog.Warn("Failed to save prompt file", "error", err)
	} else {
		slog.Info("Prompt saved", "file", filePath)
	}
}
