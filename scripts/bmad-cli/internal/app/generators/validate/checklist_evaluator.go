package validate

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/domain/models/checklist"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const (
	minRegexMatchLen = 2     // Minimum matches for regex capture groups
	rangeParts       = 2     // Number of parts in a range like "3-7"
	warnThresholdPct = 10    // Warning threshold for percentage comparisons
	filePermissions  = 0o644 // File permissions for saved prompts
)

// getDocKeyToConfigPath returns the mapping of doc keys to bmad-cli.yaml config paths.
func getDocKeyToConfigPath() map[string]string {
	return map[string]string{
		"architecture":          "documents.architecture",
		"frontend_architecture": "documents.frontend_architecture",
		"coding_standards":      "documents.coding_standards",
		"source_tree":           "documents.source_tree",
		"tech_stack":            "documents.tech_stack",
		"prd":                   "documents.prd",
		"user_roles":            "documents.user_roles",
		"architecture_yaml":     "documents.architecture_yaml",
	}
}

// ChecklistPromptData represents data needed for checklist validation prompts.
type ChecklistPromptData struct {
	Story      *story.Story
	Question   string
	Rationale  string
	ResultPath string
	Docs       map[string]*docs.ArchitectureDoc
}

// ChecklistEvaluator evaluates user stories against validation prompts using AI.
type ChecklistEvaluator struct {
	aiClient     ports.AIPort
	config       *config.ViperConfig
	modeFactory  *ai.ModeFactory
	systemLoader *template.TemplateLoader[ChecklistPromptData]
	userLoader   *template.TemplateLoader[ChecklistPromptData]
	tmpDir       string
	storyID      string
}

// NewChecklistEvaluator creates a new checklist evaluator.
func NewChecklistEvaluator(aiClient ports.AIPort, cfg *config.ViperConfig) *ChecklistEvaluator {
	systemTemplatePath := cfg.GetString("templates.prompts.checklist_system")
	userTemplatePath := cfg.GetString("templates.prompts.checklist")

	return &ChecklistEvaluator{
		aiClient:     aiClient,
		config:       cfg,
		modeFactory:  ai.NewModeFactory(cfg),
		systemLoader: template.NewTemplateLoader[ChecklistPromptData](systemTemplatePath),
		userLoader:   template.NewTemplateLoader[ChecklistPromptData](userTemplatePath),
	}
}

// Evaluate evaluates all prompts against the given story.
func (e *ChecklistEvaluator) Evaluate(
	ctx context.Context,
	storyData *story.Story,
	prompts []checklist.PromptWithContext,
	tmpDir string,
) (*checklist.ChecklistReport, error) {
	e.tmpDir = tmpDir
	e.storyID = storyData.ID

	report := &checklist.ChecklistReport{
		StoryNumber: storyData.ID,
		StoryTitle:  storyData.Title,
		Results:     make([]checklist.ValidationResult, 0, len(prompts)),
	}

	for promptIndex, promptCtx := range prompts {
		slog.Info("Evaluating prompt",
			"index", promptIndex+1,
			"total", len(prompts),
			"section", promptCtx.GetFullSectionPath(),
		)

		result, err := e.evaluatePrompt(ctx, storyData, promptCtx, promptIndex+1)
		if err != nil {
			slog.Error("Failed to evaluate prompt", "error", err)
			// Continue with other prompts, mark this one as failed
			result = checklist.ValidationResult{
				SectionPath:    promptCtx.GetFullSectionPath(),
				Question:       promptCtx.Prompt.Question,
				ExpectedAnswer: promptCtx.Prompt.Answer,
				ActualAnswer:   "ERROR: " + err.Error(),
				Status:         checklist.StatusFail,
				Rationale:      promptCtx.Prompt.Rationale,
			}
		}

		report.Results = append(report.Results, result)
	}

	report.CalculateSummary()

	return report, nil
}

// evaluatePrompt evaluates a single prompt against the story.
func (e *ChecklistEvaluator) evaluatePrompt(
	ctx context.Context,
	storyData *story.Story,
	promptCtx checklist.PromptWithContext,
	promptIndex int,
) (checklist.ValidationResult, error) {
	// Load requested documents for this prompt (uses prompt-specific or defaults)
	requestedDocs := e.loadRequestedDocs(promptCtx.GetEffectiveDocs())

	// Build result file path for FILE_START/FILE_END pattern
	sectionPath := promptCtx.GetFullSectionPath()
	safeSectionPath := strings.ReplaceAll(sectionPath, "/", "-")
	resultPath := fmt.Sprintf("%s/%02d-%s-checklist-%s-result.yaml",
		e.tmpDir, promptIndex, e.storyID, safeSectionPath)

	// Load system prompt template (uses cached loader)
	systemPrompt, err := e.systemLoader.LoadTemplate(ChecklistPromptData{})
	if err != nil {
		return checklist.ValidationResult{}, pkgerrors.ErrLoadChecklistSystemPromptFailed(err)
	}

	// Load user prompt template with data (uses cached loader)
	promptData := ChecklistPromptData{
		Story:      storyData,
		Question:   promptCtx.Prompt.Question,
		Rationale:  promptCtx.Prompt.Rationale,
		ResultPath: resultPath,
		Docs:       requestedDocs,
	}

	userPrompt, err := e.userLoader.LoadTemplate(promptData)
	if err != nil {
		return checklist.ValidationResult{}, pkgerrors.ErrLoadChecklistUserPromptFailed(err)
	}

	// Save prompts to tmp for debugging
	e.savePromptFile(sectionPath, promptIndex, "system", systemPrompt)
	e.savePromptFile(sectionPath, promptIndex, "user", userPrompt)

	// Use think mode - allows Read, Glob, Grep tools for accessing reference docs
	mode := e.modeFactory.GetThinkMode()

	response, err := e.aiClient.ExecutePromptWithSystem(ctx, systemPrompt, userPrompt, "", mode)
	if err != nil {
		return checklist.ValidationResult{}, pkgerrors.ErrChecklistAIEvaluationFailed(err)
	}

	// Save response to tmp
	e.savePromptFile(sectionPath, promptIndex, "response", response)

	// Parse the answer from result file (extracted from FILE_START/FILE_END in response)
	actualAnswer := e.parseResultFile(response, resultPath)

	// Compare with expected
	status := e.compareAnswers(promptCtx.Prompt.Answer, actualAnswer, storyData)

	return checklist.ValidationResult{
		SectionPath:    promptCtx.GetFullSectionPath(),
		Question:       promptCtx.Prompt.Question,
		ExpectedAnswer: promptCtx.Prompt.Answer,
		ActualAnswer:   actualAnswer,
		Status:         status,
		Rationale:      promptCtx.Prompt.Rationale,
	}, nil
}

// resultYAML represents the structure of the result file.
type resultYAML struct {
	Answer string `yaml:"answer"`
}

// parseResultFile extracts FILE_START/FILE_END content from response, saves to file, and parses.
func (e *ChecklistEvaluator) parseResultFile(response, path string) string {
	// Extract content between FILE_START and FILE_END markers
	content := e.extractFileContent(response, path)
	if content == "" {
		slog.Warn("No FILE_START/FILE_END content found in response", "path", path)

		return ""
	}

	// Save the extracted content to file
	err := os.WriteFile(path, []byte(content), filePermissions)
	if err != nil {
		slog.Warn("Failed to save result file", "path", path, "error", err)
	} else {
		slog.Info("Result file saved", "file", path)
	}

	// Parse the YAML
	var result resultYAML

	err = yaml.Unmarshal([]byte(content), &result)
	if err != nil {
		slog.Warn("Failed to parse result YAML", "path", path, "error", err)

		return ""
	}

	return strings.TrimSpace(result.Answer)
}

// extractFileContent extracts content between FILE_START and FILE_END markers.
func (e *ChecklistEvaluator) extractFileContent(response, path string) string {
	startMarker := fmt.Sprintf("=== FILE_START: %s ===", path)
	endMarker := fmt.Sprintf("=== FILE_END: %s ===", path)

	startIdx := strings.Index(response, startMarker)
	if startIdx == -1 {
		return ""
	}

	contentStart := startIdx + len(startMarker)
	endIdx := strings.Index(response[contentStart:], endMarker)

	if endIdx == -1 {
		return ""
	}

	return strings.TrimSpace(response[contentStart : contentStart+endIdx])
}

// compareAnswers compares actual answer to expected and returns status.
func (e *ChecklistEvaluator) compareAnswers(
	expected, actual string,
	storyData *story.Story,
) checklist.Status {
	expected = strings.TrimSpace(strings.ToLower(expected))
	actual = strings.TrimSpace(strings.ToLower(actual))

	// Try specialized comparisons in order
	if status, matched := e.trySpecializedComparison(expected, actual, storyData); matched {
		return status
	}

	// Exact match fallback
	if expected == actual {
		return checklist.StatusPass
	}

	return checklist.StatusFail
}

// trySpecializedComparison attempts to match specialized comparison patterns.
// Returns the status and whether a pattern was matched.
func (e *ChecklistEvaluator) trySpecializedComparison(
	expected, actual string,
	storyData *story.Story,
) (checklist.Status, bool) {
	// Handle special case: "= total AC count" or similar
	if e.isACCountComparison(expected) {
		return e.compareToACCount(expected, actual, storyData), true
	}

	// Handle comparison operators
	if e.isGreaterOrEqualComparison(expected) {
		return e.compareGreaterOrEqual(expected, actual), true
	}

	if e.isLessOrEqualComparison(expected) {
		return e.compareLessOrEqual(expected, actual), true
	}

	// Handle ranges like "3-7" or "0-2"
	if e.isRangeComparison(expected) {
		return e.compareRange(expected, actual), true
	}

	// Handle percentage comparisons
	if strings.Contains(expected, "%") {
		return e.comparePercentage(expected, actual), true
	}

	return checklist.StatusFail, false
}

// isACCountComparison checks if the expected value is an AC count comparison.
func (e *ChecklistEvaluator) isACCountComparison(expected string) bool {
	return strings.Contains(expected, "total") && strings.Contains(expected, "ac")
}

// isGreaterOrEqualComparison checks for >= or ≥ prefix.
func (e *ChecklistEvaluator) isGreaterOrEqualComparison(expected string) bool {
	return strings.HasPrefix(expected, ">=") || strings.HasPrefix(expected, "≥")
}

// isLessOrEqualComparison checks for <= or ≤ prefix.
func (e *ChecklistEvaluator) isLessOrEqualComparison(expected string) bool {
	return strings.HasPrefix(expected, "<=") || strings.HasPrefix(expected, "≤")
}

// isRangeComparison checks for range pattern like "3-7".
func (e *ChecklistEvaluator) isRangeComparison(expected string) bool {
	return strings.Contains(expected, "-") && !strings.HasPrefix(expected, "-")
}

// compareGreaterOrEqual handles >=N comparisons.
func (e *ChecklistEvaluator) compareGreaterOrEqual(expected, actual string) checklist.Status {
	// Extract number from expected (e.g., ">=2" -> 2)
	re := regexp.MustCompile(`[>=≥]+\s*(\d+)`)

	matches := re.FindStringSubmatch(expected)
	if len(matches) < minRegexMatchLen {
		return checklist.StatusFail
	}

	expectedNum, err := strconv.Atoi(matches[1])
	if err != nil {
		return checklist.StatusFail
	}

	actualNum, err := strconv.Atoi(actual)
	if err != nil {
		return checklist.StatusFail
	}

	if actualNum >= expectedNum {
		return checklist.StatusPass
	}

	return checklist.StatusFail
}

// compareLessOrEqual handles <=N comparisons.
func (e *ChecklistEvaluator) compareLessOrEqual(expected, actual string) checklist.Status {
	// Extract number from expected (e.g., "<=10" -> 10)
	re := regexp.MustCompile(`[<=≤]+\s*(\d+)`)

	matches := re.FindStringSubmatch(expected)
	if len(matches) < minRegexMatchLen {
		return checklist.StatusFail
	}

	expectedNum, err := strconv.Atoi(matches[1])
	if err != nil {
		return checklist.StatusFail
	}

	actualNum, err := strconv.Atoi(actual)
	if err != nil {
		return checklist.StatusFail
	}

	if actualNum <= expectedNum {
		return checklist.StatusPass
	}

	return checklist.StatusFail
}

// compareRange handles range comparisons like "3-7".
func (e *ChecklistEvaluator) compareRange(expected, actual string) checklist.Status {
	parts := strings.Split(expected, "-")
	if len(parts) != rangeParts {
		return checklist.StatusFail
	}

	minVal, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return checklist.StatusFail
	}

	maxVal, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return checklist.StatusFail
	}

	actualNum, err := strconv.Atoi(actual)
	if err != nil {
		return checklist.StatusFail
	}

	if actualNum >= minVal && actualNum <= maxVal {
		return checklist.StatusPass
	}

	// Close to range = warning
	if actualNum == minVal-1 || actualNum == maxVal+1 {
		return checklist.StatusWarn
	}

	return checklist.StatusFail
}

// comparePercentage handles percentage comparisons.
func (e *ChecklistEvaluator) comparePercentage(expected, actual string) checklist.Status {
	// Extract percentage from expected (e.g., ">=80%" -> 80)
	re := regexp.MustCompile(`[>=≥]*\s*(\d+)%`)

	expectedMatches := re.FindStringSubmatch(expected)
	if len(expectedMatches) < minRegexMatchLen {
		return checklist.StatusFail
	}

	expectedPct, err := strconv.Atoi(expectedMatches[1])
	if err != nil {
		return checklist.StatusFail
	}

	// Extract percentage from actual
	actualPctRe := regexp.MustCompile(`(\d+)%?`)

	actualMatches := actualPctRe.FindStringSubmatch(actual)
	if len(actualMatches) < minRegexMatchLen {
		return checklist.StatusFail
	}

	actualPct, err := strconv.Atoi(actualMatches[1])
	if err != nil {
		return checklist.StatusFail
	}

	if strings.Contains(expected, ">=") || strings.Contains(expected, "≥") {
		if actualPct >= expectedPct {
			return checklist.StatusPass
		}
		// Within threshold = warning
		if actualPct >= expectedPct-warnThresholdPct {
			return checklist.StatusWarn
		}
	}

	return checklist.StatusFail
}

// compareToACCount handles "= total AC count" comparisons.
func (e *ChecklistEvaluator) compareToACCount(
	_, actual string,
	storyData *story.Story,
) checklist.Status {
	acCount := len(storyData.AcceptanceCriteria)

	actualNum, err := strconv.Atoi(actual)
	if err != nil {
		return checklist.StatusFail
	}

	if actualNum == acCount {
		return checklist.StatusPass
	}

	return checklist.StatusFail
}

// loadRequestedDocs resolves document keys to file paths.
func (e *ChecklistEvaluator) loadRequestedDocs(keys []string) map[string]*docs.ArchitectureDoc {
	result := make(map[string]*docs.ArchitectureDoc, len(keys))
	docKeyMapping := getDocKeyToConfigPath()

	for _, key := range keys {
		configPath, ok := docKeyMapping[key]
		if !ok {
			slog.Warn("Unknown document key, skipping", "key", key)

			continue
		}

		filePath := e.config.GetString(configPath)
		if filePath == "" {
			slog.Warn("Document path not configured, skipping", "key", key, "configPath", configPath)

			continue
		}

		result[key] = &docs.ArchitectureDoc{
			FilePath: filePath,
		}
	}

	return result
}

// savePromptFile saves a prompt to a file in the tmp directory.
func (e *ChecklistEvaluator) savePromptFile(sectionPath string, promptIndex int, suffix, content string) {
	if e.tmpDir == "" {
		return
	}

	// Replace slashes in section path with dashes for filename
	safeSectionPath := strings.ReplaceAll(sectionPath, "/", "-")
	filePath := fmt.Sprintf("%s/%02d-%s-checklist-%s-%s.txt", e.tmpDir, promptIndex, e.storyID, safeSectionPath, suffix)

	err := os.WriteFile(filePath, []byte(content), filePermissions)
	if err != nil {
		slog.Warn("Failed to save prompt file", "error", err)
	} else {
		slog.Info("Prompt saved", "file", filePath)
	}
}
