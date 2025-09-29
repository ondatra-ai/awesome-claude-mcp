package services

import (
	"context"
	"fmt"
	"strings"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// TestingGenerator generates testing requirements for stories using AI
type TestingGenerator struct {
	aiClient AIClient
	config   *config.ViperConfig
}

// TestingData contains all data needed for testing requirements generation
type TestingData struct {
	Story            *story.Story
	Tasks            []story.Task
	DevNotes         story.DevNotes
	ArchitectureDocs *docs.ArchitectureDocs
}

// NewTestingGenerator creates a new testing requirements generator
func NewTestingGenerator(aiClient AIClient, config *config.ViperConfig) *TestingGenerator {
	return &TestingGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateTesting generates comprehensive testing requirements based on story analysis
func (g *TestingGenerator) GenerateTesting(ctx context.Context, storyDoc *story.StoryDocument) (story.Testing, error) {
	// Create AI generator for testing requirements
	generator := NewAIGenerator[TestingData, story.Testing](ctx, g.aiClient, g.config, storyDoc.Story.ID, "testing").
		WithData(func() (TestingData, error) {
			return TestingData{
				Story:            &storyDoc.Story,
				Tasks:            storyDoc.Tasks,
				DevNotes:         storyDoc.DevNotes,
				ArchitectureDocs: storyDoc.ArchitectureDocs,
			}, nil
		}).
		WithPrompt(g.loadTestingPrompt).
		WithResponseParser(CreateYAMLFileParser[story.Testing](g.config, storyDoc.Story.ID, "testing", "testing")).
		WithValidator(g.validateTesting)

	// Generate testing requirements
	testing, err := generator.Generate()
	if err != nil {
		return story.Testing{}, fmt.Errorf("failed to generate testing requirements: %w", err)
	}

	return testing, nil
}

// loadTestingPrompt loads the testing requirements prompt template
func (g *TestingGenerator) loadTestingPrompt(data TestingData) (string, error) {
	templatePath := g.config.GetString("templates.prompts.testing")

	promptLoader := template.NewTemplateLoader[TestingData](templatePath)
	prompt, err := promptLoader.LoadTemplate(data)
	if err != nil {
		return "", fmt.Errorf("failed to load testing prompt: %w", err)
	}

	return prompt, nil
}

// validateTesting validates the generated testing requirements
func (g *TestingGenerator) validateTesting(testing story.Testing) error {
	if testing.TestLocation == "" {
		return fmt.Errorf("test location cannot be empty")
	}

	if len(testing.Frameworks) == 0 {
		return fmt.Errorf("at least one testing framework must be specified")
	}

	if len(testing.Requirements) == 0 {
		return fmt.Errorf("at least one testing requirement must be specified")
	}

	if len(testing.Coverage) == 0 {
		return fmt.Errorf("coverage targets must be specified")
	}

	// Validate coverage values are percentages
	for key, value := range testing.Coverage {
		if value == "" {
			return fmt.Errorf("coverage value for %s cannot be empty", key)
		}
		// Simple validation that value contains % sign
		if !strings.Contains(value, "%") {
			return fmt.Errorf("coverage value for %s should be a percentage", key)
		}
	}

	return nil
}
