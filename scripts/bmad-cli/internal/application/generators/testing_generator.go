package generators

import (
	"bmad-cli/internal/domain/ports"
	"context"
	"fmt"
	"strings"

	"bmad-cli/internal/pkg/ai"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// AITestingGenerator generates testing requirements for stories using AI
type AITestingGenerator struct {
	aiClient ports.AIPort
	config   *config.ViperConfig
}

// TestingData contains all data needed for testing requirements generation
type TestingData struct {
	Story            *story.Story
	Tasks            []story.Task
	DevNotes         story.DevNotes
	ArchitectureDocs *docs.ArchitectureDocs
}

// NewAITestingGenerator creates a new testing requirements generator
func NewAITestingGenerator(aiClient ports.AIPort, config *config.ViperConfig) *AITestingGenerator {
	return &AITestingGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateTesting generates comprehensive testing requirements based on story analysis
func (g *AITestingGenerator) GenerateTesting(ctx context.Context, storyDoc *story.StoryDocument, tmpDir string) (story.Testing, error) {
	// Create AI generator for testing requirements
	generator := ai.NewAIGenerator[TestingData, story.Testing](ctx, g.aiClient, g.config, storyDoc.Story.ID, "testing").
		WithTmpDir(tmpDir).
		WithData(func() (TestingData, error) {
			return TestingData{
				Story:            &storyDoc.Story,
				Tasks:            storyDoc.Tasks,
				DevNotes:         storyDoc.DevNotes,
				ArchitectureDocs: storyDoc.ArchitectureDocs,
			}, nil
		}).
		WithPrompt(func(data TestingData) (systemPrompt string, userPrompt string, err error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.testing_system")
			systemLoader := template.NewTemplateLoader[TestingData](systemTemplatePath)
			systemPrompt, err = systemLoader.LoadTemplate(TestingData{})
			if err != nil {
				return "", "", fmt.Errorf("failed to load testing system prompt: %w", err)
			}

			// Load user prompt
			userPrompt, err = g.loadTestingPrompt(data)
			if err != nil {
				return "", "", fmt.Errorf("failed to load testing user prompt: %w", err)
			}

			return systemPrompt, userPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[story.Testing](g.config, storyDoc.Story.ID, "testing", "testing", tmpDir)).
		WithValidator(g.validateTesting)

	// Generate testing requirements
	testing, err := generator.Generate()
	if err != nil {
		return story.Testing{}, fmt.Errorf("failed to generate testing requirements: %w", err)
	}

	return testing, nil
}

// loadTestingPrompt loads the testing requirements prompt template
func (g *AITestingGenerator) loadTestingPrompt(data TestingData) (string, error) {
	templatePath := g.config.GetString("templates.prompts.testing")

	promptLoader := template.NewTemplateLoader[TestingData](templatePath)
	prompt, err := promptLoader.LoadTemplate(data)
	if err != nil {
		return "", fmt.Errorf("failed to load testing prompt: %w", err)
	}

	return prompt, nil
}

// validateTesting validates the generated testing requirements
func (g *AITestingGenerator) validateTesting(testing story.Testing) error {
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
