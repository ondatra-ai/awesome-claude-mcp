package services

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// QAAssessmentGenerator generates QA results for stories using AI
type QAAssessmentGenerator struct {
	aiClient AIClient
}

// QAAssessmentData contains all data needed for QA assessment generation
type QAAssessmentData struct {
	Story            *story.Story
	Tasks            []story.Task
	DevNotes         story.DevNotes
	ArchitectureDocs *docs.ArchitectureDocs
}

// NewQAAssessmentGenerator creates a new QA assessment generator
func NewQAAssessmentGenerator(aiClient AIClient) *QAAssessmentGenerator {
	return &QAAssessmentGenerator{
		aiClient: aiClient,
	}
}

// GenerateQAResults generates comprehensive QA results following Quinn persona
func (g *QAAssessmentGenerator) GenerateQAResults(ctx context.Context, storyObj *story.Story, tasks []story.Task, devNotes story.DevNotes, architectureDocs *docs.ArchitectureDocs) (story.QAResults, error) {
	storyID := storyObj.ID

	// Create AI generator for QA assessment
	generator := NewAIGenerator[QAAssessmentData, story.QAResults](ctx, g.aiClient, storyID, "qa-assessment").
		WithData(func() (QAAssessmentData, error) {
			return QAAssessmentData{
				Story:            storyObj,
				Tasks:            tasks,
				DevNotes:         devNotes,
				ArchitectureDocs: architectureDocs,
			}, nil
		}).
		WithPrompt(g.loadQAPrompt).
		WithResponseParser(CreateYAMLFileParser[story.QAResults](storyID, "qa-assessment", "qa_results")).
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
	slug := slugifyTitle(storyObj.Title)
	qaResults.GateReference = fmt.Sprintf("docs/qa/gates/%s-%s.yml", storyID, slug)

	return qaResults, nil
}

// loadQAPrompt loads the QA assessment prompt template
func (g *QAAssessmentGenerator) loadQAPrompt(data QAAssessmentData) (string, error) {
	templatePath := filepath.Join("templates", "us-create.qa.prompt.tpl")

	promptLoader := template.NewTemplateLoader[QAAssessmentData](templatePath)
	prompt, err := promptLoader.LoadTemplate(data)
	if err != nil {
		return "", fmt.Errorf("failed to load QA prompt: %w", err)
	}

	return prompt, nil
}

// validateQAResults validates the generated QA results
func (g *QAAssessmentGenerator) validateQAResults(qaResults story.QAResults) error {
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
