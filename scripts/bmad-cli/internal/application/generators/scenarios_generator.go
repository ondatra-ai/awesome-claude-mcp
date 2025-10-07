package generators

import (
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/pkg/ai"
	"context"
	"fmt"
	"strings"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/docs"
	"bmad-cli/internal/infrastructure/template"
)

// AIScenariosGenerator generates test scenarios for stories using AI
type AIScenariosGenerator struct {
	aiClient ports.AIPort
	config   *config.ViperConfig
}

// ScenariosData contains all data needed for test scenarios generation
type ScenariosData struct {
	Story            *story.Story
	Tasks            []story.Task
	DevNotes         story.DevNotes
	Testing          story.Testing
	ArchitectureDocs *docs.ArchitectureDocs
}

// NewAIScenariosGenerator creates a new test scenarios generator
func NewAIScenariosGenerator(aiClient ports.AIPort, config *config.ViperConfig) *AIScenariosGenerator {
	return &AIScenariosGenerator{
		aiClient: aiClient,
		config:   config,
	}
}

// GenerateScenarios generates comprehensive test scenarios in Given-When-Then format
func (g *AIScenariosGenerator) GenerateScenarios(ctx context.Context, storyDoc *story.StoryDocument) (story.Scenarios, error) {
	// Create AI generator for test scenarios
	generator := ai.NewAIGenerator[ScenariosData, story.Scenarios](ctx, g.aiClient, g.config, storyDoc.Story.ID, "scenarios").
		WithData(func() (ScenariosData, error) {
			return ScenariosData{
				Story:            &storyDoc.Story,
				Tasks:            storyDoc.Tasks,
				DevNotes:         storyDoc.DevNotes,
				Testing:          storyDoc.Testing,
				ArchitectureDocs: storyDoc.ArchitectureDocs,
			}, nil
		}).
		WithPrompt(func(data ScenariosData) (systemPrompt string, userPrompt string, err error) {
			// Load system prompt (doesn't need data)
			systemTemplatePath := g.config.GetString("templates.prompts.scenarios_system")
			systemLoader := template.NewTemplateLoader[ScenariosData](systemTemplatePath)
			systemPrompt, err = systemLoader.LoadTemplate(ScenariosData{})
			if err != nil {
				return "", "", fmt.Errorf("failed to load scenarios system prompt: %w", err)
			}

			// Load user prompt
			userPrompt, err = g.loadScenariosPrompt(data)
			if err != nil {
				return "", "", fmt.Errorf("failed to load scenarios user prompt: %w", err)
			}

			return systemPrompt, userPrompt, nil
		}).
		WithResponseParser(ai.CreateYAMLFileParser[story.Scenarios](g.config, storyDoc.Story.ID, "scenarios", "scenarios")).
		WithValidator(g.validateScenarios(storyDoc.Story.AcceptanceCriteria))

	// Generate test scenarios
	scenarios, err := generator.Generate()
	if err != nil {
		return story.Scenarios{}, fmt.Errorf("failed to generate test scenarios: %w", err)
	}

	return scenarios, nil
}

// loadScenariosPrompt loads the test scenarios prompt template
func (g *AIScenariosGenerator) loadScenariosPrompt(data ScenariosData) (string, error) {
	templatePath := g.config.GetString("templates.prompts.scenarios")

	promptLoader := template.NewTemplateLoader[ScenariosData](templatePath)
	prompt, err := promptLoader.LoadTemplate(data)
	if err != nil {
		return "", fmt.Errorf("failed to load scenarios prompt: %w", err)
	}

	return prompt, nil
}

// validateScenarios validates the generated test scenarios
func (g *AIScenariosGenerator) validateScenarios(acceptanceCriteria []story.AcceptanceCriterion) func(story.Scenarios) error {
	return func(scenarios story.Scenarios) error {
		if len(scenarios.TestScenarios) == 0 {
			return fmt.Errorf("at least one test scenario must be specified")
		}

		// Track which ACs are covered
		coveredACs := make(map[string]bool)

		// Validate each scenario
		for i, scenario := range scenarios.TestScenarios {
			// Validate required fields
			if scenario.ID == "" {
				return fmt.Errorf("scenario %d: ID cannot be empty", i)
			}
			if len(scenario.AcceptanceCriteria) == 0 {
				return fmt.Errorf("scenario %s: must reference at least one acceptance criterion", scenario.ID)
			}

			// Validate steps array
			if len(scenario.Steps) == 0 {
				return fmt.Errorf("scenario %s: must have at least one step", scenario.ID)
			}

			// Validate each step has exactly one keyword set
			hasGiven, hasWhen, hasThen := false, false, false
			for stepIdx, step := range scenario.Steps {
				nonEmptyCount := 0
				if step.Given != "" {
					nonEmptyCount++
					hasGiven = true
				}
				if step.When != "" {
					nonEmptyCount++
					hasWhen = true
				}
				if step.Then != "" {
					nonEmptyCount++
					hasThen = true
				}
				if step.And != "" {
					nonEmptyCount++
				}
				if step.But != "" {
					nonEmptyCount++
				}

				if nonEmptyCount == 0 {
					return fmt.Errorf("scenario %s, step %d: step must have at least one keyword set", scenario.ID, stepIdx)
				}
				if nonEmptyCount > 1 {
					return fmt.Errorf("scenario %s, step %d: step must have exactly one keyword set", scenario.ID, stepIdx)
				}
			}

			// Ensure scenario has at least Given, When, and Then (And/But are optional)
			if !hasGiven {
				return fmt.Errorf("scenario %s: must have at least one 'Given' step", scenario.ID)
			}
			if !hasWhen {
				return fmt.Errorf("scenario %s: must have at least one 'When' step", scenario.ID)
			}
			if !hasThen {
				return fmt.Errorf("scenario %s: must have at least one 'Then' step", scenario.ID)
			}

			// Validate scenario outline has examples
			if scenario.ScenarioOutline {
				if len(scenario.Examples) == 0 {
					return fmt.Errorf("scenario %s: scenario outline must have at least one example", scenario.ID)
				}
			}

			// Validate level (only integration and e2e allowed)
			validLevels := map[string]bool{"integration": true, "e2e": true}
			if !validLevels[scenario.Level] {
				return fmt.Errorf("scenario %s: level must be integration or e2e (unit scenarios are not allowed in BDD)", scenario.ID)
			}

			// Validate priority
			validPriorities := map[string]bool{"P0": true, "P1": true, "P2": true, "P3": true}
			if !validPriorities[scenario.Priority] {
				return fmt.Errorf("scenario %s: priority must be P0, P1, P2, or P3", scenario.ID)
			}

			// Track covered ACs
			for _, ac := range scenario.AcceptanceCriteria {
				coveredACs[ac] = true
			}
		}

		// Verify all ACs are covered by at least one scenario
		for _, ac := range acceptanceCriteria {
			if !coveredACs[ac.ID] {
				return fmt.Errorf("acceptance criterion %s is not covered by any test scenario", ac.ID)
			}
		}

		return nil
	}
}

// validateScenarioID validates scenario ID format (e.g., "3.1-INT-001")
func validateScenarioID(id string) error {
	parts := strings.Split(id, "-")
	if len(parts) != 3 {
		return fmt.Errorf("scenario ID must be in format {epic}.{story}-{LEVEL}-{SEQ}")
	}

	level := strings.ToUpper(parts[1])
	if level != "INT" && level != "E2E" {
		return fmt.Errorf("scenario ID level must be INT or E2E (UNIT is not allowed in BDD scenarios)")
	}

	return nil
}
