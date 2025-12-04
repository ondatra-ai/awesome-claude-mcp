package implement

import (
	"context"
	"log/slog"
	"time"

	"bmad-cli/internal/adapters/ai"
	storyModels "bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/template"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// MergeScenariosGenerator merges story scenarios into requirements file using Claude.
type MergeScenariosGenerator struct {
	claudeClient *ai.ClaudeClient
	config       *config.ViperConfig
}

// NewMergeScenariosGenerator creates a new MergeScenariosGenerator.
func NewMergeScenariosGenerator(
	claudeClient *ai.ClaudeClient,
	config *config.ViperConfig,
) *MergeScenariosGenerator {
	return &MergeScenariosGenerator{
		claudeClient: claudeClient,
		config:       config,
	}
}

// MergeScenarios merges all scenarios from the story into the requirements file.
func (g *MergeScenariosGenerator) MergeScenarios(
	ctx context.Context,
	storyDoc *storyModels.StoryDocument,
	outputFile string,
	tmpDir string,
) (GenerationStatus, error) {
	slog.Info("Starting scenario merge",
		"story_id", storyDoc.Story.ID,
		"scenario_count", len(storyDoc.Scenarios.TestScenarios),
		"output_file", outputFile,
	)

	// Create template loaders
	userPromptPath := g.config.GetString("templates.prompts.merge_scenarios")
	systemPromptPath := g.config.GetString("templates.prompts.merge_scenarios_system")
	userPromptLoader := template.NewTemplateLoader[*template.ScenarioMergeData](userPromptPath)
	systemPromptLoader := template.NewTemplateLoader[*template.ScenarioMergeData](systemPromptPath)

	processedCount := 0

	// Process each scenario individually
	for i, scenario := range storyDoc.Scenarios.TestScenarios {
		startTime := time.Now()

		slog.Info("Processing scenario",
			"index", i+1,
			"total", len(storyDoc.Scenarios.TestScenarios),
			"scenario_id", scenario.ID,
		)

		// Create merge data adapter
		mergeData := template.NewScenarioMergeData(storyDoc.Story.ID, scenario, outputFile)

		// Load templates
		userPrompt, err := userPromptLoader.LoadTemplate(mergeData)
		if err != nil {
			return NewFailureStatus("load user prompt failed"),
				pkgerrors.ErrLoadUserPromptForScenarioFailed(scenario.ID, err)
		}

		systemPrompt, err := systemPromptLoader.LoadTemplate(mergeData)
		if err != nil {
			return NewFailureStatus("load system prompt failed"),
				pkgerrors.ErrLoadSystemPromptForScenarioFailed(scenario.ID, err)
		}

		// Call Claude to merge scenario
		_, err = g.claudeClient.ExecutePromptWithSystem(
			ctx,
			systemPrompt,
			userPrompt,
			"sonnet",
			ai.ExecutionMode{AllowedTools: []string{"Read", "Edit"}},
		)
		if err != nil {
			return NewFailureStatus("merge scenario failed"),
				pkgerrors.ErrMergeScenarioFailed(scenario.ID, err)
		}

		processedCount++
		duration := time.Since(startTime)
		slog.Info("✓ Scenario merged successfully",
			"scenario_id", scenario.ID,
			"duration", duration.Round(time.Second),
		)
	}

	slog.Info("All scenarios merged successfully", "total_count", processedCount)

	return NewSuccessStatus(
		processedCount,
		[]string{outputFile},
		"All scenarios merged successfully",
	), nil
}
