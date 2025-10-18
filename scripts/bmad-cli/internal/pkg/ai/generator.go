package ai

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"
	pkgerrors "bmad-cli/internal/pkg/errors"

	"gopkg.in/yaml.v3"
)

const (
	// File permission constants.
	fileModeReadWrite = 0644 // Standard file permission for read/write files
	fileModeDirectory = 0755 // Standard directory permission
)

// AIGenerator is a generic AI content generator with builder pattern.
type AIGenerator[T1 any, T2 any] struct {
	aiClient       ports.AIPort
	config         *config.ViperConfig
	storyID        string
	filePrefix     string
	tmpDir         string // run-specific tmp directory
	dataLoader     func() (T1, error)
	promptLoader   func(T1) (systemPrompt string, userPrompt string, err error)
	responseParser func(aiResponse string) (T2, error)
	validator      func(T2) error
	model          string
	mode           ai.ExecutionMode
}

// NewAIGenerator creates a new generator instance.
func NewAIGenerator[T1 any, T2 any](
	_ context.Context,
	aiClient ports.AIPort,
	config *config.ViperConfig,
	storyID string,
	filePrefix string,
) *AIGenerator[T1, T2] {
	modeFactory := ai.NewModeFactory(config)

	return &AIGenerator[T1, T2]{
		aiClient:   aiClient,
		config:     config,
		storyID:    storyID,
		filePrefix: filePrefix,
		model:      "sonnet",                        // default model
		mode:       modeFactory.GetFullAccessMode(), // TEMP: test with full write access
	}
}

// WithData sets the data loader functor.
func (g *AIGenerator[T1, T2]) WithData(loader func() (T1, error)) *AIGenerator[T1, T2] {
	g.dataLoader = loader

	return g
}

// WithPrompt sets the prompt loader functor - can return either single prompt or dual prompts (system, user).
func (g *AIGenerator[T1, T2]) WithPrompt(loader interface{}) *AIGenerator[T1, T2] {
	switch loaderFunc := loader.(type) {
	case func(T1) (string, error):
		// Convert single prompt to dual prompt format with empty system prompt
		g.promptLoader = func(data T1) (string, string, error) {
			userPrompt, err := loaderFunc(data)

			return "", userPrompt, err
		}
	case func(T1) (string, string, error):
		// Use dual prompt directly
		g.promptLoader = loaderFunc
	}

	return g
}

// WithResponseParser sets the response parser functor.
func (g *AIGenerator[T1, T2]) WithResponseParser(parser func(string) (T2, error)) *AIGenerator[T1, T2] {
	g.responseParser = parser

	return g
}

// WithValidator sets the validation functor.
func (g *AIGenerator[T1, T2]) WithValidator(validator func(T2) error) *AIGenerator[T1, T2] {
	g.validator = validator

	return g
}

// WithModel sets the AI model to use ("sonnet" or "opus").
func (g *AIGenerator[T1, T2]) WithModel(model string) *AIGenerator[T1, T2] {
	g.model = model

	return g
}

// WithMode sets the execution mode (ThinkMode or FullAccessMode).
func (g *AIGenerator[T1, T2]) WithMode(mode ai.ExecutionMode) *AIGenerator[T1, T2] {
	g.mode = mode

	return g
}

// WithTmpDir sets the run-specific tmp directory path.
func (g *AIGenerator[T1, T2]) WithTmpDir(tmpDir string) *AIGenerator[T1, T2] {
	g.tmpDir = tmpDir

	return g
}

// Generate executes the generation pipeline.
func (g *AIGenerator[T1, T2]) Generate(ctx context.Context) (T2, error) {
	var zero T2

	tmpDir, err := g.prepareTmpDirectory()
	if err != nil {
		return zero, err
	}

	data, err := g.dataLoader()
	if err != nil {
		return zero, fmt.Errorf("load data failed: %w", pkgerrors.ErrLoadDataFailed(err))
	}

	systemPrompt, userPrompt, err := g.promptLoader(data)
	if err != nil {
		return zero, fmt.Errorf("load prompts failed: %w", pkgerrors.ErrLoadPromptsFailed(err))
	}

	response, err := g.executeAIPrompt(ctx, tmpDir, systemPrompt, userPrompt)
	if err != nil {
		return zero, err
	}

	err = g.saveResponseFile(tmpDir, response)
	if err != nil {
		return zero, err
	}

	result, err := g.responseParser(response)
	if err != nil {
		return zero, fmt.Errorf("parse response failed: %w", pkgerrors.ErrParseResponseFailed(err))
	}

	slog.Info("Content generated successfully", "type", g.filePrefix)

	if g.validator != nil {
		err := g.validator(result)
		if err != nil {
			return zero, fmt.Errorf("validation failed: %w", pkgerrors.ErrValidationFailed(err))
		}
	}

	return result, nil
}

func (g *AIGenerator[T1, T2]) prepareTmpDirectory() (string, error) {
	tmpDir := g.tmpDir
	if tmpDir == "" {
		tmpDir = g.config.GetString("paths.tmp_dir")
	}

	err := os.MkdirAll(tmpDir, fileModeDirectory)
	if err != nil {
		return "", fmt.Errorf("create tmp directory failed: %w", pkgerrors.ErrCreateTmpDirectoryFailed(err))
	}

	return tmpDir, nil
}

func (g *AIGenerator[T1, T2]) executeAIPrompt(
	ctx context.Context, tmpDir, systemPrompt, userPrompt string,
) (string, error) {
	if systemPrompt != "" {
		return g.executeDualPrompt(ctx, tmpDir, systemPrompt, userPrompt)
	}

	return g.executeSinglePrompt(ctx, tmpDir, userPrompt)
}

func (g *AIGenerator[T1, T2]) executeDualPrompt(
	ctx context.Context, tmpDir, systemPrompt, userPrompt string,
) (string, error) {
	g.savePromptFile(tmpDir, "system-prompt", systemPrompt)
	g.savePromptFile(tmpDir, "user-prompt", userPrompt)

	response, err := g.aiClient.ExecutePromptWithSystem(ctx, systemPrompt, userPrompt, g.model, g.mode)
	if err != nil {
		return "", fmt.Errorf(
			"generate content with system prompt failed: %w",
			pkgerrors.ErrGenerateContentWithSystemPromptFailed(err),
		)
	}

	return response, nil
}

func (g *AIGenerator[T1, T2]) executeSinglePrompt(ctx context.Context, tmpDir, userPrompt string) (string, error) {
	g.savePromptFile(tmpDir, "prompt", userPrompt)

	response, err := g.aiClient.ExecutePromptWithSystem(ctx, "", userPrompt, g.model, g.mode)
	if err != nil {
		return "", fmt.Errorf("generate content failed: %w", pkgerrors.ErrGenerateContentFailed(err))
	}

	return response, nil
}

func (g *AIGenerator[T1, T2]) savePromptFile(tmpDir, suffix, content string) {
	filePath := fmt.Sprintf("%s/%s-%s-%s.txt", tmpDir, g.storyID, g.filePrefix, suffix)

	err := os.WriteFile(filePath, []byte(content), fileModeReadWrite)
	if err != nil {
		slog.Warn("Failed to save prompt file", "error", err)
	} else {
		slog.Info("💾 Prompt saved", "file", filePath)
	}
}

func (g *AIGenerator[T1, T2]) saveResponseFile(tmpDir, response string) error {
	responseFile := fmt.Sprintf("%s/%s-%s-full-response.txt", tmpDir, g.storyID, g.filePrefix)

	err := os.WriteFile(responseFile, []byte(response), fileModeReadWrite)
	if err != nil {
		return fmt.Errorf("write response file failed: %w", pkgerrors.ErrWriteResponseFileFailed(err))
	}

	slog.Info("💾 AI response saved", "file", responseFile)

	return nil
}

// CreateYAMLFileParser creates a parser function for reading YAML files
// This is a higher-order function that returns a closure configured with the
// given parameters.
func CreateYAMLFileParser[T any](
	config *config.ViperConfig,
	storyID, filePrefix, yamlKey string,
	tmpDir string,
) func(string) (T, error) {
	return func(aiResponse string) (T, error) {
		var zero T

		// Construct file path - use provided tmpDir or fallback to config
		dir := tmpDir
		if dir == "" {
			dir = config.GetString("paths.tmp_dir")
		}

		filePath := fmt.Sprintf("%s/%s-%s.yaml", dir, storyID, filePrefix)

		// Read file
		content, err := os.ReadFile(filePath)
		if err != nil {
			return zero, fmt.Errorf("YAML file not found: %w", pkgerrors.ErrYAMLFileNotFound(filePrefix, filePath))
		}

		// Parse YAML based on key
		data := make(map[string]T)

		err = yaml.Unmarshal(content, &data)
		if err != nil {
			return zero, fmt.Errorf("parse YAML failed: %w", pkgerrors.ErrParseYAMLFailed(filePrefix, err))
		}

		result, exists := data[yamlKey]
		if !exists {
			return zero, fmt.Errorf("YAML key not found: %w", pkgerrors.ErrYAMLKeyNotFound(yamlKey))
		}

		return result, nil
	}
}
