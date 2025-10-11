package generators

import (
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/pkg/ai"
	pkgerrors "bmad-cli/internal/pkg/errors"
	"context"
	"errors"
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

// GenerateScenarios generates comprehensive test scenarios in Given-When-Then format.
func (g *AIScenariosGenerator) GenerateScenarios(ctx context.Context, storyDoc *story.StoryDocument, tmpDir string) (story.Scenarios, error) {
	// Create AI generator for test scenarios
	generator := ai.NewAIGenerator[ScenariosData, story.Scenarios](ctx, g.aiClient, g.config, storyDoc.Story.ID, "scenarios").
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
		WithPrompt(func(data ScenariosData) (systemPrompt string, userPrompt string, err error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.scenarios_system")
			systemLoader := template.NewTemplateLoader[ScenariosData](systemTemplatePath)

			systemPrompt, err = systemLoader.LoadTemplate(ScenariosData{})
			if err != nil {
				return "", "", pkgerrors.ErrLoadScenariosSystemPromptFailed(err)
			}

			// Load user prompt
			userPrompt, err = g.loadScenariosPrompt(data)
			if err != nil {
				return "", "", pkgerrors.ErrLoadScenariosUserPromptFailed(err)
			}

			return systemPrompt, userPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[story.Scenarios](g.config, storyDoc.Story.ID, "scenarios", "scenarios", tmpDir)).
		WithValidator(g.validateScenarios(storyDoc.Story.AcceptanceCriteria))

	// Generate test scenarios
	scenarios, err := generator.Generate()
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
func (g *AIScenariosGenerator) validateScenarios(acceptanceCriteria []story.AcceptanceCriterion) func(story.Scenarios) error {
	return func(scenarios story.Scenarios) error {
		if len(scenarios.TestScenarios) == 0 {
			return errors.New("at least one test scenario must be specified")
		}

		// Track which ACs are covered
		coveredACs := make(map[string]bool)

		// Validate each scenario
		for i, scenario := range scenarios.TestScenarios {
			// Validate required fields
			if scenario.ID == "" {
				return pkgerrors.ErrEmptyScenarioIDError(i)
			}

			if len(scenario.AcceptanceCriteria) == 0 {
				return pkgerrors.ErrNoCriteriaError(scenario.ID)
			}

			// Validate steps array
			if len(scenario.Steps) == 0 {
				return pkgerrors.ErrNoStepsError(scenario.ID)
			}

			// Validate each step has exactly one keyword set and statements are valid
			hasGiven, hasWhen, hasThen := false, false, false

			for stepIdx, step := range scenario.Steps {
				nonEmptyCount := 0

				// Validate Given
				if len(step.Given) > 0 {
					nonEmptyCount++
					hasGiven = true

					err := validateStepStatements(scenario.ID, stepIdx, "Given", step.Given)
					if err != nil {
						return err
					}
				}

				// Validate When
				if len(step.When) > 0 {
					nonEmptyCount++
					hasWhen = true

					err := validateStepStatements(scenario.ID, stepIdx, "When", step.When)
					if err != nil {
						return err
					}
				}

				// Validate Then
				if len(step.Then) > 0 {
					nonEmptyCount++
					hasThen = true

					err := validateStepStatements(scenario.ID, stepIdx, "Then", step.Then)
					if err != nil {
						return err
					}
				}

				if nonEmptyCount == 0 {
					return pkgerrors.ErrNoKeywordSetError(scenario.ID, stepIdx)
				}

				if nonEmptyCount > 1 {
					return pkgerrors.ErrMultipleKeywordsError(scenario.ID, stepIdx)
				}
			}

			// Ensure scenario has at least Given, When, and Then
			if !hasGiven {
				return pkgerrors.ErrNoGivenStepError(scenario.ID)
			}

			if !hasWhen {
				return pkgerrors.ErrNoWhenStepError(scenario.ID)
			}

			if !hasThen {
				return pkgerrors.ErrNoThenStepError(scenario.ID)
			}

			// Validate scenario outline has examples
			if scenario.ScenarioOutline {
				if len(scenario.Examples) == 0 {
					return pkgerrors.ErrNoExamplesError(scenario.ID)
				}
			}

			// Validate level (only integration and e2e allowed)
			validLevels := map[string]bool{"integration": true, "e2e": true}
			if !validLevels[scenario.Level] {
				return pkgerrors.ErrInvalidLevelError(scenario.ID)
			}

			// Validate priority
			validPriorities := map[string]bool{"P0": true, "P1": true, "P2": true, "P3": true}
			if !validPriorities[scenario.Priority] {
				return pkgerrors.ErrInvalidPriorityError(scenario.ID)
			}

			// Track covered ACs
			for _, ac := range scenario.AcceptanceCriteria {
				coveredACs[ac] = true
			}
		}

		// Verify all ACs are covered by at least one scenario
		for _, ac := range acceptanceCriteria {
			if !coveredACs[ac.ID] {
				return pkgerrors.ErrUncoveredCriterionError(ac.ID)
			}
		}

		return nil
	}
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

// validateScenarioID validates scenario ID format (e.g., "3.1-INT-001").
func validateScenarioID(id string) error {
	parts := strings.Split(id, "-")
	if len(parts) != 3 {
		return errors.New("scenario ID must be in format {epic}.{story}-{LEVEL}-{SEQ}")
	}

	level := strings.ToUpper(parts[1])
	if level != "INT" && level != "E2E" {
		return errors.New("scenario ID level must be INT or E2E (UNIT is not allowed in BDD scenarios)")
	}

	return nil
}
