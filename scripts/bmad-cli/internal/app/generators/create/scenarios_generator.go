package create

import (
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/pkg/ai"
	pkgerrors "bmad-cli/internal/pkg/errors"
	"context"
	"strings"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// AIScenariosGenerator generates test scenarios for stories using AI.
type AIScenariosGenerator struct {
	aiClient ports.AIPort
	config   *config.ViperConfig
}

// ScenariosData contains all data needed for test scenarios generation.
type ScenariosData struct {
	Story            *story.Story
	Tasks            []story.Task
	DevNotes         story.DevNotes
	Testing          story.Testing
	ArchitectureDocs *docs.ArchitectureDocs
	TmpDir           string // Path to run-specific tmp directory
}

// NewAIScenariosGenerator creates a new test scenarios generator.
func NewAIScenariosGenerator(aiClient ports.AIPort, config *config.ViperConfig) *AIScenariosGenerator {
	return &AIScenariosGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateScenarios generates comprehensive test scenarios in Given-When-Then
// format.
func (g *AIScenariosGenerator) GenerateScenarios(
	ctx context.Context,
	storyDoc *story.StoryDocument,
	tmpDir string,
) (story.Scenarios, error) {
	// Create AI generator for test scenarios
	generator := ai.NewAIGenerator[ScenariosData, story.Scenarios](
		ctx,
		g.aiClient,
		g.config,
		storyDoc.Story.ID,
		"scenarios",
	).
		WithTmpDir(tmpDir).
		WithData(func() (ScenariosData, error) {
			return ScenariosData{
				Story:            &storyDoc.Story,
				Tasks:            storyDoc.Tasks,
				DevNotes:         storyDoc.DevNotes,
				Testing:          storyDoc.Testing,
				ArchitectureDocs: storyDoc.ArchitectureDocs,
				TmpDir:           tmpDir,
			}, nil
		}).
		WithPrompt(func(data ScenariosData) (string, string, error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.scenarios_system")
			systemLoader := template.NewTemplateLoader[ScenariosData](systemTemplatePath)

			systemPrompt, err := systemLoader.LoadTemplate(ScenariosData{})
			if err != nil {
				return "", "", pkgerrors.ErrLoadScenariosSystemPromptFailed(err)
			}

			// Load user prompt
			userPrompt, err := g.loadScenariosPrompt(data)
			if err != nil {
				return "", "", pkgerrors.ErrLoadScenariosUserPromptFailed(err)
			}

			return systemPrompt, userPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[story.Scenarios](
			g.config,
			storyDoc.Story.ID,
			"scenarios",
			"scenarios",
			tmpDir,
		)).
		WithValidator(g.validateScenarios(storyDoc.Story.AcceptanceCriteria))

	// Generate test scenarios
	scenarios, err := generator.Generate(ctx)
	if err != nil {
		return story.Scenarios{}, pkgerrors.ErrGenerateTestScenariosFailed(err)
	}

	return scenarios, nil
}

// loadScenariosPrompt loads the test scenarios prompt template.
func (g *AIScenariosGenerator) loadScenariosPrompt(data ScenariosData) (string, error) {
	templatePath := g.config.GetString("templates.prompts.scenarios")

	promptLoader := template.NewTemplateLoader[ScenariosData](templatePath)

	prompt, err := promptLoader.LoadTemplate(data)
	if err != nil {
		return "", pkgerrors.ErrLoadScenariosPromptFailed(err)
	}

	return prompt, nil
}

// validateScenarios validates the generated test scenarios.
func (g *AIScenariosGenerator) validateScenarios(
	acceptanceCriteria []story.AcceptanceCriterion,
) func(story.Scenarios) error {
	return func(scenarios story.Scenarios) error {
		if len(scenarios.TestScenarios) == 0 {
			return pkgerrors.ErrAtLeastOneTestScenario
		}

		coveredACs := make(map[string]bool)

		for i, scenario := range scenarios.TestScenarios {
			err := g.validateSingleScenario(i, scenario, coveredACs)
			if err != nil {
				return err
			}
		}

		err := g.verifyACCoverage(acceptanceCriteria, coveredACs)
		if err != nil {
			return err
		}

		return nil
	}
}

// validateSingleScenario validates a single test scenario.
func (g *AIScenariosGenerator) validateSingleScenario(
	index int,
	scenario story.TestScenario,
	coveredACs map[string]bool,
) error {
	err := g.validateScenarioBasicFields(index, scenario)
	if err != nil {
		return err
	}

	err = g.validateScenarioSteps(scenario)
	if err != nil {
		return err
	}

	err = g.validateScenarioMetadata(scenario)
	if err != nil {
		return err
	}

	for _, ac := range scenario.AcceptanceCriteria {
		coveredACs[ac] = true
	}

	return nil
}

// validateScenarioBasicFields validates basic scenario fields.
func (g *AIScenariosGenerator) validateScenarioBasicFields(
	index int,
	scenario story.TestScenario,
) error {
	if scenario.ID == "" {
		return pkgerrors.ErrEmptyScenarioIDError(index)
	}

	if len(scenario.AcceptanceCriteria) == 0 {
		return pkgerrors.ErrNoCriteriaError(scenario.ID)
	}

	if len(scenario.Steps) == 0 {
		return pkgerrors.ErrNoStepsError(scenario.ID)
	}

	return nil
}

// validateScenarioSteps validates all steps in a scenario.
func (g *AIScenariosGenerator) validateScenarioSteps(scenario story.TestScenario) error {
	hasGiven, hasWhen, hasThen := false, false, false

	for stepIdx, step := range scenario.Steps {
		givenPresent, whenPresent, thenPresent, err := g.validateStep(scenario.ID, stepIdx, step)
		if err != nil {
			return err
		}

		hasGiven = hasGiven || givenPresent
		hasWhen = hasWhen || whenPresent
		hasThen = hasThen || thenPresent
	}

	return g.validateScenarioHasAllKeywords(scenario.ID, hasGiven, hasWhen, hasThen)
}

// validateStep validates a single step and returns which keywords are present.
func (g *AIScenariosGenerator) validateStep(
	scenarioID string,
	stepIdx int,
	step story.ScenarioStep,
) (bool, bool, bool, error) {
	nonEmptyCount := 0
	hasGiven := false
	hasWhen := false
	hasThen := false

	if len(step.Given) > 0 {
		nonEmptyCount++
		hasGiven = true

		err := validateStepStatements(scenarioID, stepIdx, "Given", step.Given)
		if err != nil {
			return false, false, false, err
		}
	}

	if len(step.When) > 0 {
		nonEmptyCount++
		hasWhen = true

		err := validateStepStatements(scenarioID, stepIdx, "When", step.When)
		if err != nil {
			return false, false, false, err
		}
	}

	if len(step.Then) > 0 {
		nonEmptyCount++
		hasThen = true

		err := validateStepStatements(scenarioID, stepIdx, "Then", step.Then)
		if err != nil {
			return false, false, false, err
		}
	}

	if nonEmptyCount == 0 {
		return false, false, false, pkgerrors.ErrNoKeywordSetError(scenarioID, stepIdx)
	}

	if nonEmptyCount > 1 {
		return false, false, false, pkgerrors.ErrMultipleKeywordsError(scenarioID, stepIdx)
	}

	return hasGiven, hasWhen, hasThen, nil
}

// validateScenarioHasAllKeywords ensures scenario has Given, When, and Then.
func (g *AIScenariosGenerator) validateScenarioHasAllKeywords(
	scenarioID string,
	hasGiven bool,
	hasWhen bool,
	hasThen bool,
) error {
	if !hasGiven {
		return pkgerrors.ErrNoGivenStepError(scenarioID)
	}

	if !hasWhen {
		return pkgerrors.ErrNoWhenStepError(scenarioID)
	}

	if !hasThen {
		return pkgerrors.ErrNoThenStepError(scenarioID)
	}

	return nil
}

// validateScenarioMetadata validates scenario metadata fields.
func (g *AIScenariosGenerator) validateScenarioMetadata(scenario story.TestScenario) error {
	if scenario.ScenarioOutline {
		if len(scenario.Examples) == 0 {
			return pkgerrors.ErrNoExamplesError(scenario.ID)
		}
	}

	validLevels := map[string]bool{"integration": true, "e2e": true}
	if !validLevels[scenario.Level] {
		return pkgerrors.ErrInvalidLevelError(scenario.ID)
	}

	validPriorities := map[string]bool{"P0": true, "P1": true, "P2": true, "P3": true}
	if !validPriorities[scenario.Priority] {
		return pkgerrors.ErrInvalidPriorityError(scenario.ID)
	}

	return nil
}

// verifyACCoverage verifies all acceptance criteria are covered.
func (g *AIScenariosGenerator) verifyACCoverage(
	acceptanceCriteria []story.AcceptanceCriterion,
	coveredACs map[string]bool,
) error {
	for _, ac := range acceptanceCriteria {
		if !coveredACs[ac.ID] {
			return pkgerrors.ErrUncoveredCriterionError(ac.ID)
		}
	}

	return nil
}

// validateStepStatements validates an array of step statements.
func validateStepStatements(scenarioID string, stepIdx int, keyword string, statements []story.StepStatement) error {
	if len(statements) == 0 {
		return pkgerrors.ErrNoStatementsError(scenarioID, stepIdx, keyword)
	}

	for stmtIdx, stmt := range statements {
		// Check statement is non-empty
		if strings.TrimSpace(stmt.Statement) == "" {
			return pkgerrors.ErrEmptyStatementError(scenarioID, stepIdx, keyword, stmtIdx)
		}

		// First statement must be main (no type)
		if stmtIdx == 0 && stmt.Type != "" {
			return pkgerrors.ErrInvalidFirstStmtError(scenarioID, stepIdx, keyword)
		}

		// Additional statements must be and/but
		if stmtIdx > 0 {
			if stmt.Type != story.ModifierTypeAnd && stmt.Type != story.ModifierTypeBut {
				return pkgerrors.ErrInvalidFollowingStmtError(scenarioID, stepIdx, keyword, stmtIdx)
			}
		}
	}

	return nil
}
