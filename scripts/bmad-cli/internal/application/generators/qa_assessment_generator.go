package generators

import (
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/pkg/ai"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// AIQAAssessmentGenerator generates QA results for stories using AI
type AIQAAssessmentGenerator struct {
	aiClient ports.AIClient
	config   *config.ViperConfig
}

// QAAssessmentData contains all data needed for QA assessment generation
type QAAssessmentData struct {
	Story            *story.Story
	Tasks            []story.Task
	DevNotes         story.DevNotes
	ArchitectureDocs *docs.ArchitectureDocs
}

// NewAIQAAssessmentGenerator creates a new QA assessment generator
func NewAIQAAssessmentGenerator(aiClient ports.AIClient, config *config.ViperConfig) *AIQAAssessmentGenerator {
	return &AIQAAssessmentGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateQAResults generates comprehensive QA results following Quinn persona
func (g *AIQAAssessmentGenerator) GenerateQAResults(ctx context.Context, storyDoc *story.StoryDocument) (story.QAResults, error) {
	// Create AI generator for QA assessment
	generator := ai.NewAIGenerator[QAAssessmentData, story.QAResults](ctx, g.aiClient, g.config, storyDoc.Story.ID, "qa-assessment").
		WithData(func() (QAAssessmentData, error) {
			return QAAssessmentData{
				Story:            &storyDoc.Story,
				Tasks:            storyDoc.Tasks,
				DevNotes:         storyDoc.DevNotes,
				ArchitectureDocs: storyDoc.ArchitectureDocs,
			}, nil
		}).
		WithPrompt(func(data QAAssessmentData) (systemPrompt string, userPrompt string, err error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.qa_system")
			systemLoader := template.NewTemplateLoader[QAAssessmentData](systemTemplatePath)
			systemPrompt, err = systemLoader.LoadTemplate(QAAssessmentData{})
			if err != nil {
				return "", "", fmt.Errorf("failed to load QA system prompt: %w", err)
			}

			// Load user prompt
			userPrompt, err = g.loadQAPrompt(data)
			if err != nil {
				return "", "", fmt.Errorf("failed to load QA user prompt: %w", err)
			}

			return systemPrompt, userPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[story.QAResults](g.config, storyDoc.Story.ID, "qa-assessment", "qa_results")).
		WithValidator(g.validateQAResults)

	// Generate QA results
	qaResults, err := generator.Generate()
	if err != nil {
		return story.QAResults{}, fmt.Errorf("failed to generate QA results: %w", err)
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

// loadQAPrompt loads the QA assessment prompt template
func (g *AIQAAssessmentGenerator) loadQAPrompt(data QAAssessmentData) (string, error) {
	templatePath := g.config.GetString("templates.prompts.qa")

	promptLoader := template.NewTemplateLoader[QAAssessmentData](templatePath)
	prompt, err := promptLoader.LoadTemplate(data)
	if err != nil {
		return "", fmt.Errorf("failed to load QA prompt: %w", err)
	}

	return prompt, nil
}

// validateQAResults validates the generated QA results
func (g *AIQAAssessmentGenerator) validateQAResults(qaResults story.QAResults) error {
	if qaResults.Assessment.Summary == "" {
		return fmt.Errorf("assessment summary cannot be empty")
	}

	if len(qaResults.Assessment.Strengths) == 0 {
		return fmt.Errorf("at least one strength must be identified")
	}

	if qaResults.Assessment.RiskLevel == "" {
		return fmt.Errorf("risk level must be specified")
	}

	validRiskLevels := map[string]bool{
		"Low":    true,
		"Medium": true,
		"High":   true,
	}
	if !validRiskLevels[qaResults.Assessment.RiskLevel] {
		return fmt.Errorf("invalid risk level: %s (must be Low, Medium, or High)", qaResults.Assessment.RiskLevel)
	}

	if qaResults.Assessment.TestabilityScore < 1 || qaResults.Assessment.TestabilityScore > 10 {
		return fmt.Errorf("testability score must be between 1 and 10")
	}

	if qaResults.Assessment.ImplementationReadiness < 1 || qaResults.Assessment.ImplementationReadiness > 10 {
		return fmt.Errorf("implementation readiness must be between 1 and 10")
	}

	validGateStatuses := map[string]bool{
		"PASS":     true,
		"CONCERNS": true,
		"FAIL":     true,
		"WAIVED":   true,
	}
	if !validGateStatuses[qaResults.GateStatus] {
		return fmt.Errorf("invalid gate status: %s", qaResults.GateStatus)
	}

	return nil
}

// slugifyTitle converts a title to a URL-friendly slug
func slugifyTitle(title string) string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(title)
	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`[\s_-]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
