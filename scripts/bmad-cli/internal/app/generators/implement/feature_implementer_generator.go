package implement

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"bmad-cli/internal/adapters/ai"
	storyModels "bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// Maximum implementation attempts.
const maxAttempts = 5

// FeatureImplementerGenerator implements features using Claude.
type FeatureImplementerGenerator struct {
	claudeClient *ai.ClaudeClient
	config       *config.ViperConfig
}

// ImplementFeatureData holds the data for the implement feature prompt.
type ImplementFeatureData struct {
	StoryID     string
	StoryTitle  string
	AsA         string
	IWant       string
	SoThat      string
	TestCommand string
	TestOutput  string
	Attempt     int
	MaxAttempts int
}

// NewFeatureImplementerGenerator creates a new FeatureImplementerGenerator.
func NewFeatureImplementerGenerator(
	claudeClient *ai.ClaudeClient,
	config *config.ViperConfig,
) *FeatureImplementerGenerator {
	return &FeatureImplementerGenerator{
		claudeClient: claudeClient,
		config:       config,
	}
}

// Implement attempts to implement the feature to make tests pass.
// This is called in a loop by the factory until tests pass or max attempts reached.
func (g *FeatureImplementerGenerator) Implement(
	ctx context.Context,
	storyDoc *storyModels.StoryDocument,
	attempt int,
	testOutput string,
	tmpDir string,
) (GenerationStatus, error) {
	promptData := g.buildPromptData(storyDoc, attempt, testOutput)

	userPrompt, systemPrompt, err := g.loadPrompts(promptData)
	if err != nil {
		return NewFailureStatus("load prompts failed"), err
	}

	g.saveAttemptPrompts(tmpDir, storyDoc.Story.ID, attempt, userPrompt, systemPrompt)

	slog.Info("🤖 Calling Claude to implement feature", "attempt", attempt)

	response, err := g.claudeClient.ExecutePromptWithSystem(
		ctx,
		systemPrompt,
		userPrompt,
		"sonnet",
		ai.ExecutionMode{AllowedTools: []string{"Read", "Write", "Edit", "Bash"}},
	)

	g.saveAttemptResponse(tmpDir, storyDoc.Story.ID, attempt, response)

	if err != nil {
		return NewFailureStatus("implement feature failed"),
			pkgerrors.ErrImplementFeaturesFailed(err)
	}

	slog.Info("✓ Claude finished attempt", "attempt", attempt)

	return NewSuccessStatus(1, nil, fmt.Sprintf("Implementation attempt %d completed", attempt)), nil
}

// buildPromptData creates the data structure for prompt templates.
func (g *FeatureImplementerGenerator) buildPromptData(
	storyDoc *storyModels.StoryDocument,
	attempt int,
	testOutput string,
) *ImplementFeatureData {
	testCommand := g.config.GetString("testing.command")
	if testCommand == "" {
		testCommand = "make test-e2e"
	}

	return &ImplementFeatureData{
		StoryID:     storyDoc.Story.ID,
		StoryTitle:  storyDoc.Story.Title,
		AsA:         storyDoc.Story.AsA,
		IWant:       storyDoc.Story.IWant,
		SoThat:      storyDoc.Story.SoThat,
		TestCommand: testCommand,
		TestOutput:  testOutput,
		Attempt:     attempt,
		MaxAttempts: maxAttempts,
	}
}

// loadPrompts loads user and system prompts from templates.
func (g *FeatureImplementerGenerator) loadPrompts(
	promptData *ImplementFeatureData,
) (string, string, error) {
	userPromptPath := g.config.GetString("templates.prompts.implement_feature")
	systemPromptPath := g.config.GetString("templates.prompts.implement_feature_system")

	userPromptLoader := template.NewTemplateLoader[*ImplementFeatureData](userPromptPath)
	systemPromptLoader := template.NewTemplateLoader[*ImplementFeatureData](systemPromptPath)

	userPrompt, err := userPromptLoader.LoadTemplate(promptData)
	if err != nil {
		return "", "", pkgerrors.ErrLoadPromptsFailed(err)
	}

	systemPrompt, err := systemPromptLoader.LoadTemplate(promptData)
	if err != nil {
		return "", "", pkgerrors.ErrLoadPromptsFailed(err)
	}

	return userPrompt, systemPrompt, nil
}

// saveAttemptPrompts saves user and system prompts for an attempt.
func (g *FeatureImplementerGenerator) saveAttemptPrompts(
	tmpDir, storyID string, attempt int, userPrompt, systemPrompt string,
) {
	userFileName := fmt.Sprintf("%s-implement-feature-attempt-%d-user-prompt.txt", storyID, attempt)
	systemFileName := fmt.Sprintf("%s-implement-feature-attempt-%d-system-prompt.txt", storyID, attempt)

	g.savePromptFile(tmpDir, userFileName, userPrompt)
	g.savePromptFile(tmpDir, systemFileName, systemPrompt)
}

// saveAttemptResponse saves the response for an attempt if not empty.
func (g *FeatureImplementerGenerator) saveAttemptResponse(tmpDir, storyID string, attempt int, response string) {
	if response != "" {
		g.savePromptFile(tmpDir, fmt.Sprintf("%s-implement-feature-attempt-%d-response.txt", storyID, attempt), response)
	}
}

// savePromptFile saves content to a file in the tmp directory.
func (g *FeatureImplementerGenerator) savePromptFile(tmpDir, filename, content string) {
	filePath := filepath.Join(tmpDir, filename)

	err := os.WriteFile(filePath, []byte(content), fileModeReadWrite)
	if err != nil {
		slog.Warn("Failed to save file", "file", filePath, "error", err)
	} else {
		slog.Info("💾 Prompt saved", "file", filePath)
	}
}
