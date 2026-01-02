package validate

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	Story    *story.Story
	Question string
	Docs     map[string]*docs.ArchitectureDoc
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

	for i, promptCtx := range prompts {
		slog.Info("Evaluating prompt",
			"index", i+1,
			"total", len(prompts),
			"section", promptCtx.GetFullSectionPath(),
		)

		result, err := e.evaluatePrompt(ctx, storyData, promptCtx)
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
) (checklist.ValidationResult, error) {
	// Load requested documents for this prompt (uses prompt-specific or defaults)
	requestedDocs := e.loadRequestedDocs(promptCtx.GetEffectiveDocs())

	// Load system prompt template (uses cached loader)
	systemPrompt, err := e.systemLoader.LoadTemplate(ChecklistPromptData{})
	if err != nil {
		return checklist.ValidationResult{}, pkgerrors.ErrLoadChecklistSystemPromptFailed(err)
	}

	// Load user prompt template with data (uses cached loader)
	promptData := ChecklistPromptData{
		Story:    storyData,
		Question: promptCtx.Prompt.Question,
		Docs:     requestedDocs,
	}

	userPrompt, err := e.userLoader.LoadTemplate(promptData)
	if err != nil {
		return checklist.ValidationResult{}, pkgerrors.ErrLoadChecklistUserPromptFailed(err)
	}

	// Save prompts to tmp for debugging
	sectionPath := promptCtx.GetFullSectionPath()
	e.savePromptFile(sectionPath, "system", systemPrompt)
	e.savePromptFile(sectionPath, "user", userPrompt)

	// Use think mode - allows Read, Glob, Grep tools for accessing reference docs
	mode := e.modeFactory.GetThinkMode()

	response, err := e.aiClient.ExecutePromptWithSystem(ctx, systemPrompt, userPrompt, "", mode)
	if err != nil {
		return checklist.ValidationResult{}, pkgerrors.ErrChecklistAIEvaluationFailed(err)
	}

	// Save response to tmp
	e.savePromptFile(sectionPath, "response", response)

	// Parse the answer from response
	actualAnswer := e.parseAnswer(response)

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

// parseAnswer extracts the answer from AI response.
func (e *ChecklistEvaluator) parseAnswer(response string) string {
	// Clean up the response - take first line, trim whitespace
	response = strings.TrimSpace(response)

	lines := strings.Split(response, "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}

	return response
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
func (e *ChecklistEvaluator) savePromptFile(sectionPath, suffix, content string) {
	if e.tmpDir == "" {
		return
	}

	// Replace slashes in section path with dashes for filename
	safeSectionPath := strings.ReplaceAll(sectionPath, "/", "-")
	filePath := fmt.Sprintf("%s/%s-checklist-%s-%s.txt", e.tmpDir, e.storyID, safeSectionPath, suffix)

	err := os.WriteFile(filePath, []byte(content), filePermissions)
	if err != nil {
		slog.Warn("Failed to save prompt file", "error", err)
	} else {
		slog.Info("Prompt saved", "file", filePath)
	}
}
