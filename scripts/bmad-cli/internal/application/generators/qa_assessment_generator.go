package generators

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/pkg/ai"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// AIQAAssessmentGenerator generates QA results for stories using AI.
type AIQAAssessmentGenerator struct {
	aiClient ports.AIPort
	config   *config.ViperConfig
}

// QAAssessmentData contains all data needed for QA assessment generation.
type QAAssessmentData struct {
	Story            *story.Story
	Tasks            []story.Task
	DevNotes         story.DevNotes
	ArchitectureDocs *docs.ArchitectureDocs
	TmpDir           string // Path to run-specific tmp directory
}

// NewAIQAAssessmentGenerator creates a new QA assessment generator.
func NewAIQAAssessmentGenerator(aiClient ports.AIPort, config *config.ViperConfig) *AIQAAssessmentGenerator {
	return &AIQAAssessmentGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateQAResults generates comprehensive QA results following Quinn persona.
func (g *AIQAAssessmentGenerator) GenerateQAResults(
	ctx context.Context,
	storyDoc *story.StoryDocument,
	tmpDir string,
) (story.QAResults, error) {
	// Create AI generator for QA assessment
	generator := g.createQAGenerator(ctx, storyDoc, tmpDir)

	// Generate QA results
	qaResults, err := generator.Generate(ctx)
	if err != nil {
		return story.QAResults{}, pkgerrors.ErrGenerateQAResultsFailed(err)
	}

	// Set review metadata
	qaResults.ReviewDate = time.Now().Format("2006-01-02")
	qaResults.ReviewedBy = "Quinn (Test Architect)"

	// Generate gate reference path
	slug := slugifyTitle(storyDoc.Story.Title)
	qaGatesPath := g.config.GetString("paths.qa_gates")
	qaResults.GateReference = fmt.Sprintf("%s/%s-%s.yml", qaGatesPath, storyDoc.Story.ID, slug)

	return qaResults, nil
}

// loadPrompts loads system and user prompts for QA assessment.
func (g *AIQAAssessmentGenerator) loadPrompts(
	data QAAssessmentData,
) (string, string, error) {
	// Load system prompt (doesn't need data)
	systemTemplatePath := g.config.GetString("templates.prompts.qa_system")
	systemLoader := template.NewTemplateLoader[QAAssessmentData](systemTemplatePath)

	systemPrompt, err := systemLoader.LoadTemplate(QAAssessmentData{})
	if err != nil {
		return "", "", pkgerrors.ErrLoadQASystemPromptFailed(err)
	}

	// Load user prompt
	userPrompt, err := g.loadQAPrompt(data)
	if err != nil {
		return "", "", pkgerrors.ErrLoadQAUserPromptFailed(err)
	}

	return systemPrompt, userPrompt, nil
}

// createQAGenerator creates and configures the AI generator for QA assessment.
func (g *AIQAAssessmentGenerator) createQAGenerator(
	ctx context.Context,
	storyDoc *story.StoryDocument,
	tmpDir string,
) *ai.AIGenerator[QAAssessmentData, story.QAResults] {
	return ai.NewAIGenerator[QAAssessmentData, story.QAResults](
		ctx,
		g.aiClient,
		g.config,
		storyDoc.Story.ID,
		"qa-assessment",
	).
		WithTmpDir(tmpDir).
		WithData(func() (QAAssessmentData, error) {
			return QAAssessmentData{
				Story:            &storyDoc.Story,
				Tasks:            storyDoc.Tasks,
				DevNotes:         storyDoc.DevNotes,
				ArchitectureDocs: storyDoc.ArchitectureDocs,
				TmpDir:           tmpDir,
			}, nil
		}).
		WithPrompt(g.loadPrompts).
		WithResponseParser(ai.CreateYAMLFileParser[story.QAResults](
			g.config,
			storyDoc.Story.ID,
			"qa-assessment",
			"qa_results",
			tmpDir,
		)).
		WithValidator(g.validateQAResults)
}

// loadQAPrompt loads the QA assessment prompt template.
func (g *AIQAAssessmentGenerator) loadQAPrompt(data QAAssessmentData) (string, error) {
	templatePath := g.config.GetString("templates.prompts.qa")

	promptLoader := template.NewTemplateLoader[QAAssessmentData](templatePath)

	prompt, err := promptLoader.LoadTemplate(data)
	if err != nil {
		return "", pkgerrors.ErrLoadQAPromptFailed(err)
	}

	return prompt, nil
}

// validateQAResults validates the generated QA results.
func (g *AIQAAssessmentGenerator) validateQAResults(qaResults story.QAResults) error {
	if qaResults.Assessment.Summary == "" {
		return errors.New("assessment summary cannot be empty")
	}

	if len(qaResults.Assessment.Strengths) == 0 {
		return errors.New("at least one strength must be identified")
	}

	if qaResults.Assessment.RiskLevel == "" {
		return errors.New("risk level must be specified")
	}

	validRiskLevels := map[string]bool{
		"Low":    true,
		"Medium": true,
		"High":   true,
	}
	if !validRiskLevels[qaResults.Assessment.RiskLevel] {
		return pkgerrors.ErrInvalidRiskLevelError(qaResults.Assessment.RiskLevel)
	}

	if qaResults.Assessment.TestabilityScore < 1 || qaResults.Assessment.TestabilityScore > 10 {
		return errors.New("testability score must be between 1 and 10")
	}

	if qaResults.Assessment.ImplementationReadiness < 1 || qaResults.Assessment.ImplementationReadiness > 10 {
		return errors.New("implementation readiness must be between 1 and 10")
	}

	validGateStatuses := map[string]bool{
		"PASS":     true,
		"CONCERNS": true,
		"FAIL":     true,
		"WAIVED":   true,
	}
	if !validGateStatuses[qaResults.GateStatus] {
		return pkgerrors.ErrInvalidGateStatusError(qaResults.GateStatus)
	}

	return nil
}

// slugifyTitle converts a title to a URL-friendly slug.
func slugifyTitle(title string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}
