package ai

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"

	"gopkg.in/yaml.v3"
)

// AIGenerator is a generic AI content generator with builder pattern
type AIGenerator[T1 any, T2 any] struct {
	ctx            context.Context
	aiClient       ports.AIPort
	config         *config.ViperConfig
	storyID        string
	filePrefix     string
	dataLoader     func() (T1, error)
	promptLoader   func(T1) (systemPrompt string, userPrompt string, err error)
	responseParser func(aiResponse string) (T2, error)
	validator      func(T2) error
	model          string
	mode           ai.ExecutionMode
}

// NewAIGenerator creates a new generator instance
func NewAIGenerator[T1 any, T2 any](ctx context.Context, aiClient ports.AIPort, config *config.ViperConfig, storyID string, filePrefix string) *AIGenerator[T1, T2] {
	modeFactory := ai.NewModeFactory(config)
	return &AIGenerator[T1, T2]{
		ctx:        ctx,
		aiClient:   aiClient,
		config:     config,
		storyID:    storyID,
		filePrefix: filePrefix,
		model:      "sonnet",                        // default model
		mode:       modeFactory.GetFullAccessMode(), // TEMP: test with full write access
	}
}

// WithData sets the data loader functor
func (g *AIGenerator[T1, T2]) WithData(loader func() (T1, error)) *AIGenerator[T1, T2] {
	g.dataLoader = loader
	return g
}

// WithPrompt sets the prompt loader functor - can return either single prompt or dual prompts (system, user)
func (g *AIGenerator[T1, T2]) WithPrompt(loader interface{}) *AIGenerator[T1, T2] {
	switch l := loader.(type) {
	case func(T1) (string, error):
		// Convert single prompt to dual prompt format with empty system prompt
		g.promptLoader = func(data T1) (string, string, error) {
			userPrompt, err := l(data)
			return "", userPrompt, err
		}
	case func(T1) (string, string, error):
		// Use dual prompt directly
		g.promptLoader = l
	}
	return g
}

// WithResponseParser sets the response parser functor
func (g *AIGenerator[T1, T2]) WithResponseParser(parser func(string) (T2, error)) *AIGenerator[T1, T2] {
	g.responseParser = parser
	return g
}

// WithValidator sets the validation functor
func (g *AIGenerator[T1, T2]) WithValidator(validator func(T2) error) *AIGenerator[T1, T2] {
	g.validator = validator
	return g
}

// WithModel sets the AI model to use ("sonnet" or "opus")
func (g *AIGenerator[T1, T2]) WithModel(model string) *AIGenerator[T1, T2] {
	g.model = model
	return g
}

// WithMode sets the execution mode (ThinkMode or FullAccessMode)
func (g *AIGenerator[T1, T2]) WithMode(mode ai.ExecutionMode) *AIGenerator[T1, T2] {
	g.mode = mode
	return g
}

// Generate executes the generation pipeline
func (g *AIGenerator[T1, T2]) Generate() (T2, error) {
	var zero T2

	// 0. Create tmp directory for debugging early
	tmpDir := g.config.GetString("paths.tmp_dir")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return zero, fmt.Errorf("failed to create tmp directory: %w", err)
	}

	// 1. Load input data
	data, err := g.dataLoader()
	if err != nil {
		return zero, fmt.Errorf("failed to load data: %w", err)
	}

	// 2. Generate prompts and call AI
	systemPrompt, userPrompt, err := g.promptLoader(data)
	if err != nil {
		return zero, fmt.Errorf("failed to load prompts: %w", err)
	}

	var response string

	if systemPrompt != "" {
		// Dual prompt mode - save both prompts for debugging
		systemPromptFile := fmt.Sprintf("%s/%s-%s-system-prompt.txt", tmpDir, g.storyID, g.filePrefix)
		if err := os.WriteFile(systemPromptFile, []byte(systemPrompt), 0644); err != nil {
			slog.Warn("Failed to save system prompt file", "error", err)
		} else {
			slog.Info("💾 System prompt saved", "file", systemPromptFile)
		}

		userPromptFile := fmt.Sprintf("%s/%s-%s-user-prompt.txt", tmpDir, g.storyID, g.filePrefix)
		if err := os.WriteFile(userPromptFile, []byte(userPrompt), 0644); err != nil {
			slog.Warn("Failed to save user prompt file", "error", err)
		} else {
			slog.Info("💾 User prompt saved", "file", userPromptFile)
		}

		// Use system + user prompt
		response, err = g.aiClient.ExecutePromptWithSystem(g.ctx, systemPrompt, userPrompt, g.model, g.mode)
		if err != nil {
			return zero, fmt.Errorf("failed to generate content with system prompt: %w", err)
		}
	} else {
		// Single prompt mode - save single prompt for debugging
		promptFile := fmt.Sprintf("%s/%s-%s-prompt.txt", tmpDir, g.storyID, g.filePrefix)
		if err := os.WriteFile(promptFile, []byte(userPrompt), 0644); err != nil {
			slog.Warn("Failed to save prompt file", "error", err)
		} else {
			slog.Info("💾 Prompt saved", "file", promptFile)
		}

		// Use single prompt (empty system prompt)
		response, err = g.aiClient.ExecutePromptWithSystem(g.ctx, "", userPrompt, g.model, g.mode)
		if err != nil {
			return zero, fmt.Errorf("failed to generate content: %w", err)
		}
	}

	// 4. Save AI response for debugging
	responseFile := fmt.Sprintf("%s/%s-%s-full-response.txt", tmpDir, g.storyID, g.filePrefix)
	if err := os.WriteFile(responseFile, []byte(response), 0644); err != nil {
		return zero, fmt.Errorf("failed to write response file: %w", err)
	}
	slog.Info("💾 AI response saved", "file", responseFile)

	// 5. Parse response
	result, err := g.responseParser(response)
	if err != nil {
		return zero, fmt.Errorf("failed to parse response: %w", err)
	}

	// 6. Log success
	slog.Info("Content generated successfully", "type", g.filePrefix)

	// 7. Validate if validator is set
	if g.validator != nil {
		if err := g.validator(result); err != nil {
			return zero, fmt.Errorf("validation failed: %w", err)
		}
	}

	return result, nil
}

// CreateYAMLFileParser creates a parser function for reading YAML files
// This is a higher-order function that returns a closure configured with the given parameters
func CreateYAMLFileParser[T any](config *config.ViperConfig, storyID, filePrefix, yamlKey string) func(string) (T, error) {
	return func(aiResponse string) (T, error) {
		var zero T

		// Construct file path
		tmpDir := config.GetString("paths.tmp_dir")
		filePath := fmt.Sprintf("%s/%s-%s.yaml", tmpDir, storyID, filePrefix)

		// Read file
		content, err := os.ReadFile(filePath)
		if err != nil {
			return zero, fmt.Errorf("%s file not found: %s", filePrefix, filePath)
		}

		// Parse YAML based on key
		data := make(map[string]T)
		if err := yaml.Unmarshal(content, &data); err != nil {
			return zero, fmt.Errorf("failed to parse %s YAML: %w", filePrefix, err)
		}

		result, exists := data[yamlKey]
		if !exists {
			return zero, fmt.Errorf("%s key not found in YAML", yamlKey)
		}

		return result, nil
	}
}
