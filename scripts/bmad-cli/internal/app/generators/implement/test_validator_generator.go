package implement

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"

	"gopkg.in/yaml.v3"
)

// TestValidatorGenerator validates test quality using Claude.
type TestValidatorGenerator struct {
	claudeClient *ai.ClaudeClient
	config       *config.ViperConfig
}

// ValidateTestsData holds the data for the validate tests prompt.
type ValidateTestsData struct {
	RequirementsFile string
	TestFilesGlob    string
	TmpDir           string
}

// ValidateTestsResult represents the result of test validation.
type ValidateTestsResult struct {
	FilesScanned  int            `yaml:"files_scanned"`
	IssuesFound   int            `yaml:"issues_found"`
	IssuesFixed   int            `yaml:"issues_fixed"`
	UnfixedIssues []UnfixedIssue `yaml:"unfixed_issues"`
}

// UnfixedIssue represents an issue that could not be fixed.
type UnfixedIssue struct {
	File         string `yaml:"file"`
	Line         int    `yaml:"line"`
	Description  string `yaml:"description"`
	SuggestedFix string `yaml:"suggested_fix"`
}

// NewTestValidatorGenerator creates a new TestValidatorGenerator.
func NewTestValidatorGenerator(
	claudeClient *ai.ClaudeClient,
	config *config.ViperConfig,
) *TestValidatorGenerator {
	return &TestValidatorGenerator{
		claudeClient: claudeClient,
		config:       config,
	}
}

// ValidateTests validates test quality using Claude and returns the result.
func (g *TestValidatorGenerator) ValidateTests(
	ctx context.Context,
	tmpDir string,
) (GenerationStatus, error) {
	// Load template paths from config
	userPromptPath := g.config.GetString("templates.prompts.validate_tests")
	systemPromptPath := g.config.GetString("templates.prompts.validate_tests_system")

	// Create template loaders
	userPromptLoader := template.NewTemplateLoader[*ValidateTestsData](userPromptPath)
	systemPromptLoader := template.NewTemplateLoader[*ValidateTestsData](systemPromptPath)

	// Create validation data
	validateData := &ValidateTestsData{
		RequirementsFile: "docs/requirements.yaml",
		TestFilesGlob:    "tests/**/*.spec.ts",
		TmpDir:           tmpDir,
	}

	// Load prompts
	userPrompt, err := userPromptLoader.LoadTemplate(validateData)
	if err != nil {
		return NewFailureStatus("load user prompt failed"),
			pkgerrors.ErrLoadPromptsFailed(err)
	}

	g.savePromptFile(tmpDir, "validate-tests-user-prompt.txt", userPrompt)

	systemPrompt, err := systemPromptLoader.LoadTemplate(validateData)
	if err != nil {
		return NewFailureStatus("load system prompt failed"),
			pkgerrors.ErrLoadPromptsFailed(err)
	}

	g.savePromptFile(tmpDir, "validate-tests-system-prompt.txt", systemPrompt)

	slog.Info("🤖 Calling Claude to validate tests")

	response, err := g.claudeClient.ExecutePromptWithSystem(
		ctx,
		systemPrompt,
		userPrompt,
		"sonnet",
		ai.ExecutionMode{AllowedTools: []string{"Read", "Edit", "Glob", "Grep", "Write"}},
	)

	// Save response
	if response != "" {
		g.savePromptFile(tmpDir, "validate-tests-response.txt", response)
	}

	if err != nil {
		return NewFailureStatus("validate tests failed"),
			pkgerrors.ErrValidateTestsFailed(err)
	}

	// Parse and check validation result
	result, parseErr := g.parseValidateTestsResult(tmpDir)
	if parseErr != nil {
		slog.Warn("⚠️  Could not parse validation result file", "error", parseErr)
		slog.Info("✅ Test validation completed (no structured result)")

		return NewSuccessStatus(0, nil, "Test validation completed"), nil
	}

	// Log validation summary
	slog.Info("Test validation summary",
		"files_scanned", result.FilesScanned,
		"issues_found", result.IssuesFound,
		"issues_fixed", result.IssuesFixed,
		"unfixed_count", len(result.UnfixedIssues),
	)

	// Check if there are unfixed issues
	if len(result.UnfixedIssues) > 0 {
		g.logUnfixedIssues(result.UnfixedIssues)

		return NewFailureStatus("unfixed issues remain"),
			pkgerrors.ErrUnfixedTestIssuesError(len(result.UnfixedIssues))
	}

	slog.Info("✅ Test validation completed successfully")

	return NewSuccessStatus(
		result.IssuesFixed,
		nil,
		"Test validation completed successfully",
	), nil
}

// parseValidateTestsResult parses the validation result YAML file.
func (g *TestValidatorGenerator) parseValidateTestsResult(tmpDir string) (*ValidateTestsResult, error) {
	resultPath := filepath.Join(tmpDir, "validate-tests-result.yaml")

	data, err := os.ReadFile(resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read validation result file: %w", err)
	}

	var wrapper struct {
		Data ValidateTestsResult `yaml:"data"`
	}

	err = yaml.Unmarshal(data, &wrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to parse validation result YAML: %w", err)
	}

	return &wrapper.Data, nil
}

// logUnfixedIssues logs all unfixed issues.
func (g *TestValidatorGenerator) logUnfixedIssues(issues []UnfixedIssue) {
	slog.Error("❌ Some issues could not be automatically fixed:")

	for _, issue := range issues {
		slog.Error("  Unfixed issue",
			"file", issue.File,
			"line", issue.Line,
			"description", issue.Description,
			"suggested_fix", issue.SuggestedFix,
		)
	}
}

// savePromptFile saves content to a file in the tmp directory.
func (g *TestValidatorGenerator) savePromptFile(tmpDir, filename, content string) {
	filePath := filepath.Join(tmpDir, filename)

	err := os.WriteFile(filePath, []byte(content), fileModeReadWrite)
	if err != nil {
		slog.Warn("Failed to save file", "file", filePath, "error", err)
	} else {
		slog.Info("💾 Prompt saved", "file", filePath)
	}
}
