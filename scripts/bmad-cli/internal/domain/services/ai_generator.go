package services

import (
	"context"
	"fmt"
	"os"
	"strings"

	"bmad-cli/internal/adapters/ai"

	"gopkg.in/yaml.v3"
)

// AIClient defines the interface for AI communication
type AIClient interface {
	ExecutePrompt(ctx context.Context, prompt string, mode ai.ExecutionMode) (string, error)
}

// AIGenerator is a generic AI content generator with builder pattern
type AIGenerator[T1 any, T2 any] struct {
	ctx            context.Context
	aiClient       AIClient
	storyID        string
	filePrefix     string
	dataLoader     func() (T1, error)
	promptLoader   func(T1) (string, error)
	responseParser func(aiResponse string) (T2, error)
	validator      func(T2) error
}

// NewAIGenerator creates a new generator instance
func NewAIGenerator[T1 any, T2 any](ctx context.Context, aiClient AIClient, storyID string, filePrefix string) *AIGenerator[T1, T2] {
	return &AIGenerator[T1, T2]{
		ctx:        ctx,
		aiClient:   aiClient,
		storyID:    storyID,
		filePrefix: filePrefix,
	}
}

// WithData sets the data loader functor
func (g *AIGenerator[T1, T2]) WithData(loader func() (T1, error)) *AIGenerator[T1, T2] {
	g.dataLoader = loader
	return g
}

// WithPrompt sets the prompt loader functor
func (g *AIGenerator[T1, T2]) WithPrompt(loader func(T1) (string, error)) *AIGenerator[T1, T2] {
	g.promptLoader = loader
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

// Generate executes the generation pipeline
func (g *AIGenerator[T1, T2]) Generate() (T2, error) {
	var zero T2

	// 0. Create tmp directory for debugging early
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		return zero, fmt.Errorf("failed to create tmp directory: %w", err)
	}

	// 1. Load input data
	data, err := g.dataLoader()
	if err != nil {
		return zero, fmt.Errorf("failed to load data: %w", err)
	}

	// 2. Generate prompt
	prompt, err := g.promptLoader(data)
	if err != nil {
		return zero, fmt.Errorf("failed to load prompt: %w", err)
	}

	// 3. Call AI
	// Save prompt for debugging
	promptFile := fmt.Sprintf("./tmp/%s-%s-prompt.txt", g.storyID, g.filePrefix)
	if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
		fmt.Printf("⚠️ Failed to save prompt file: %v\n", err)
	} else {
		fmt.Printf("💾 Prompt saved to: %s\n", promptFile)
	}

	// Use PlanMode as required
	response, err := g.aiClient.ExecutePrompt(g.ctx, prompt, ai.PlanMode)
	if err != nil {
		return zero, fmt.Errorf("failed to generate content: %w", err)
	}

	// 4. Save AI response for debugging
	responseFile := fmt.Sprintf("./tmp/%s-%s-full-response.txt", g.storyID, g.filePrefix)
	if err := os.WriteFile(responseFile, []byte(response), 0644); err != nil {
		return zero, fmt.Errorf("failed to write response file: %w", err)
	}
	fmt.Printf("💾 Full AI response saved to: %s\n", responseFile)

	// 5. Parse response
	result, err := g.responseParser(response)
	if err != nil {
		return zero, fmt.Errorf("failed to parse response: %w", err)
	}

	// 6. Log success
	fmt.Printf("✅ %s generated successfully\n", g.filePrefix)

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
func CreateYAMLFileParser[T any](storyID, filePrefix, yamlKey string) func(string) (T, error) {
	return func(aiResponse string) (T, error) {
		var zero T

		// Construct file path
		filePath := fmt.Sprintf("./tmp/%s-%s.yaml", storyID, filePrefix)

		// First, try to extract and save file from AI response if it contains file markers
		if err := extractAndSaveFileFromResponse(aiResponse, filePath, storyID, filePrefix); err != nil {
			// If extraction fails, continue with existing file reading logic
			fmt.Printf("⚠️ Failed to extract file from response: %v\n", err)
		}

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

// extractAndSaveFileFromResponse extracts file content from AI response using file markers
func extractAndSaveFileFromResponse(aiResponse, filePath, storyID, filePrefix string) error {
	startMarker := fmt.Sprintf("=== FILE_START: ./tmp/%s-%s.yaml ===", storyID, filePrefix)
	endMarker := fmt.Sprintf("=== FILE_END: ./tmp/%s-%s.yaml ===", storyID, filePrefix)

	startIdx := strings.Index(aiResponse, startMarker)
	if startIdx == -1 {
		return fmt.Errorf("start marker not found")
	}

	endIdx := strings.Index(aiResponse, endMarker)
	if endIdx == -1 {
		return fmt.Errorf("end marker not found")
	}

	if endIdx <= startIdx {
		return fmt.Errorf("invalid marker positions")
	}

	// Extract content between markers
	startIdx += len(startMarker)
	fileContent := strings.TrimSpace(aiResponse[startIdx:endIdx])

	// Save to file
	if err := os.WriteFile(filePath, []byte(fileContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Extracted and saved file: %s\n", filePath)
	return nil
}
