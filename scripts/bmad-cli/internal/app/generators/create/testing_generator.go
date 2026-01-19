package create

import (
	"bmad-cli/internal/domain/ports"
	"context"
	"strings"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
	"bmad-cli/internal/pkg/ai"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// AITestingGenerator generates testing requirements for stories using AI.
type AITestingGenerator struct {
	aiClient ports.AIPort
	config   *config.ViperConfig
}

// TestingData contains all data needed for testing requirements generation.
type TestingData struct {
	Story            *story.Story
	Tasks            []story.Task
	DevNotes         story.DevNotes
	ArchitectureDocs *docs.ArchitectureDocs
	TmpDir           string // Path to run-specific tmp directory
}

// NewAITestingGenerator creates a new testing requirements generator.
func NewAITestingGenerator(aiClient ports.AIPort, config *config.ViperConfig) *AITestingGenerator {
	return &AITestingGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateTesting generates comprehensive testing requirements based on story
// analysis.
func (g *AITestingGenerator) GenerateTesting(
	ctx context.Context,
	storyDoc *story.StoryDocument,
	tmpDir string,
) (story.Testing, error) {
	// Create AI generator for testing requirements
	generator := g.createTestingGenerator(ctx, storyDoc, tmpDir)

	// Generate testing requirements
	testing, err := generator.Generate(ctx)
	if err != nil {
		return story.Testing{}, pkgerrors.ErrGenerateTestingFailed(err)
	}

	return testing, nil
}

// loadPrompts loads system and user prompts for testing requirements.
func (g *AITestingGenerator) loadPrompts(
	data TestingData,
) (string, string, error) {
	// Load system prompt (doesn't need data)
	systemTemplatePath := g.config.GetString("templates.prompts.testing_system")
	systemLoader := template.NewTemplateLoader[TestingData](systemTemplatePath)

	sysPrompt, err := systemLoader.LoadTemplate(TestingData{})
	if err != nil {
		return "", "", pkgerrors.ErrLoadTestingSystemPromptFailed(err)
	}

	// Load user prompt
	usrPrompt, err := g.loadTestingPrompt(data)
	if err != nil {
		return "", "", pkgerrors.ErrLoadTestingUserPromptFailed(err)
	}

	return sysPrompt, usrPrompt, nil
}

// createTestingGenerator creates and configures the AI generator for testing requirements.
func (g *AITestingGenerator) createTestingGenerator(
	ctx context.Context,
	storyDoc *story.StoryDocument,
	tmpDir string,
) *ai.AIGenerator[TestingData, story.Testing] {
	return ai.NewAIGenerator[TestingData, story.Testing](
		ctx,
		g.aiClient,
		g.config,
		storyDoc.Story.ID,
		"testing",
	).
		WithTmpDir(tmpDir).
		WithData(func() (TestingData, error) {
			return TestingData{
				Story:            &storyDoc.Story,
				Tasks:            storyDoc.Tasks,
				DevNotes:         storyDoc.DevNotes,
				ArchitectureDocs: storyDoc.ArchitectureDocs,
				TmpDir:           tmpDir,
			}, nil
		}).
		WithPrompt(g.loadPrompts).
		WithResponseParser(ai.CreateYAMLFileParser[story.Testing](
			g.config,
			storyDoc.Story.ID,
			"testing",
			"testing",
			tmpDir,
		)).
		WithValidator(g.validateTesting)
}

// loadTestingPrompt loads the testing requirements prompt template.
func (g *AITestingGenerator) loadTestingPrompt(data TestingData) (string, error) {
	templatePath := g.config.GetString("templates.prompts.testing")

	promptLoader := template.NewTemplateLoader[TestingData](templatePath)

	prompt, err := promptLoader.LoadTemplate(data)
	if err != nil {
		return "", pkgerrors.ErrLoadTestingPromptFailed(err)
	}

	return prompt, nil
}

// validateTesting validates the generated testing requirements.
func (g *AITestingGenerator) validateTesting(testing story.Testing) error {
	if testing.TestLocation == "" {
		return pkgerrors.ErrTestLocationEmpty
	}

	if len(testing.Frameworks) == 0 {
		return pkgerrors.ErrAtLeastOneFramework
	}

	if len(testing.Requirements) == 0 {
		return pkgerrors.ErrAtLeastOneTestingReq
	}

	if len(testing.Coverage) == 0 {
		return pkgerrors.ErrCoverageTargetsMustBeSpecified
	}

	// Validate coverage values are percentages
	for key, value := range testing.Coverage {
		if value == "" {
			return pkgerrors.ErrEmptyCoverageError(key)
		}
		// Simple validation that value contains % sign
		if !strings.Contains(value, "%") {
			return pkgerrors.ErrInvalidCoverageError(key)
		}
	}

	return nil
}
