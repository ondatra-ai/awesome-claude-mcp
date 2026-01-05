package validate

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const fixPromptFilePermissions = 0o644

// FixPromptData represents data needed for fix prompt generation templates.
type FixPromptData struct {
	Story       *story.Story
	FailedCheck checklist.ValidationResult
	ResultPath  string
}

// GenerateParams contains parameters for fix prompt generation.
type GenerateParams struct {
	StoryData   *story.Story
	FailedCheck checklist.ValidationResult // Single failed check to generate fix for
	TmpDir      string
	// Note: PromptIndex is taken from FailedCheck.PromptIndex
}

// FixPromptGenerator generates complete, actionable fix prompts for failed checklist validations.
type FixPromptGenerator struct {
	aiClient     ports.AIPort
	config       *config.ViperConfig
	modeFactory  *ai.ModeFactory
	systemLoader *template.TemplateLoader[FixPromptData]
	userLoader   *template.TemplateLoader[FixPromptData]
}

// NewFixPromptGenerator creates a new fix prompt generator.
func NewFixPromptGenerator(aiClient ports.AIPort, cfg *config.ViperConfig) *FixPromptGenerator {
	systemTemplatePath := cfg.GetString("templates.prompts.fix_generator_system")
	userTemplatePath := cfg.GetString("templates.prompts.fix_generator")

	return &FixPromptGenerator{
		aiClient:     aiClient,
		config:       cfg,
		modeFactory:  ai.NewModeFactory(cfg),
		systemLoader: template.NewTemplateLoader[FixPromptData](systemTemplatePath),
		userLoader:   template.NewTemplateLoader[FixPromptData](userTemplatePath),
	}
}

// Generate creates a complete, actionable fix prompt for a single failed validation.
func (g *FixPromptGenerator) Generate(
	ctx context.Context,
	params GenerateParams,
) (string, error) {
	// Use PromptIndex from the failed check
	promptIndex := params.FailedCheck.PromptIndex
	if promptIndex == 0 {
		slog.Warn("FailedCheck has no PromptIndex, skipping fix generation")

		return "", nil
	}

	slog.Info("Generating fix prompt",
		"story", params.StoryData.ID,
		"promptIndex", promptIndex,
		"section", params.FailedCheck.SectionPath,
	)

	// Build result file path with naming convention: XX-<storyID>-fix-prompts.md
	resultPath := fmt.Sprintf("%s/%02d-%s-fix-prompts.md",
		params.TmpDir, promptIndex, params.StoryData.ID)

	// Build prompt data with full story context
	promptData := FixPromptData{
		Story:       params.StoryData,
		FailedCheck: params.FailedCheck,
		ResultPath:  resultPath,
	}

	// Load system prompt template
	systemPrompt, err := g.systemLoader.LoadTemplate(promptData)
	if err != nil {
		return "", pkgerrors.ErrLoadChecklistSystemPromptFailed(err)
	}

	// Load user prompt template with data
	userPrompt, err := g.userLoader.LoadTemplate(promptData)
	if err != nil {
		return "", pkgerrors.ErrLoadChecklistUserPromptFailed(err)
	}

	// Save prompts to tmp for debugging with naming convention
	g.savePromptFile(params.TmpDir, params.StoryData.ID, promptIndex, "fix-system", systemPrompt)
	g.savePromptFile(params.TmpDir, params.StoryData.ID, promptIndex, "fix-user", userPrompt)

	// Use think mode - allows Read tool for accessing BDD guidelines
	mode := g.modeFactory.GetThinkMode()

	response, err := g.aiClient.ExecutePromptWithSystem(ctx, systemPrompt, userPrompt, "", mode)
	if err != nil {
		return "", pkgerrors.ErrChecklistAIEvaluationFailed(err)
	}

	// Save response to tmp
	g.savePromptFile(params.TmpDir, params.StoryData.ID, promptIndex, "fix-response", response)

	// Extract fix prompt from FILE_START/FILE_END markers
	fixPrompt := g.extractFixPrompt(response, resultPath)
	if fixPrompt == "" {
		slog.Warn("No fix prompt content found in response")

		return "", nil
	}

	// Save the extracted fix prompt to file
	err = os.WriteFile(resultPath, []byte(fixPrompt), fixPromptFilePermissions)
	if err != nil {
		slog.Warn("Failed to save fix prompt file", "path", resultPath, "error", err)
	} else {
		slog.Info("Fix prompt saved", "file", resultPath)
	}

	return fixPrompt, nil
}

// extractFixPrompt extracts content between FILE_START and FILE_END markers.
func (g *FixPromptGenerator) extractFixPrompt(response, path string) string {
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

// savePromptFile saves a prompt to a file in the tmp directory with naming convention.
func (g *FixPromptGenerator) savePromptFile(tmpDir, storyID string, promptIndex int, suffix, content string) {
	if tmpDir == "" {
		return
	}

	// Follow naming convention: XX-<storyID>-<suffix>.txt
	filePath := fmt.Sprintf("%s/%02d-%s-%s.txt", tmpDir, promptIndex, storyID, suffix)

	err := os.WriteFile(filePath, []byte(content), fixPromptFilePermissions)
	if err != nil {
		slog.Warn("Failed to save prompt file", "error", err)
	} else {
		slog.Info("Prompt saved", "file", filePath)
	}
}
