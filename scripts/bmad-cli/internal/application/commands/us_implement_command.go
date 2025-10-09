package commands

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/story"
	"bmad-cli/internal/infrastructure/template"
	storyModels "bmad-cli/internal/domain/models/story"
)

type USImplementCommand struct {
	branchManager *git.BranchManager
	storyLoader   *story.StoryLoader
	claudeClient  *ai.ClaudeClient
	config        *config.ViperConfig
}

func NewUSImplementCommand(
	branchManager *git.BranchManager,
	storyLoader *story.StoryLoader,
	claudeClient *ai.ClaudeClient,
	cfg *config.ViperConfig,
) *USImplementCommand {
	return &USImplementCommand{
		branchManager: branchManager,
		storyLoader:   storyLoader,
		claudeClient:  claudeClient,
		config:        cfg,
	}
}

func (c *USImplementCommand) Execute(ctx context.Context, storyNumber string, force bool) error {
	slog.Info("Starting user story implementation", "story_number", storyNumber, "force", force)

	// Get story slug from file
	_, err := c.storyLoader.GetStorySlug(storyNumber)
	if err != nil {
		return fmt.Errorf("failed to get story slug: %w", err)
	}

	// TEMPORARILY COMMENTED OUT FOR TESTING
	// Ensure correct branch is checked out
	// if err := c.branchManager.EnsureBranch(ctx, storyNumber, storySlug, force); err != nil {
	// 	return fmt.Errorf("failed to ensure branch: %w", err)
	// }
	// slog.Info("Branch setup completed successfully")

	fmt.Println("⚠️  Branch management temporarily disabled for testing")

	// Create run directory for this execution
	tmpBasePath := c.config.GetString("paths.tmp_dir")
	runDir, err := fs.NewRunDirectory(tmpBasePath)
	if err != nil {
		return fmt.Errorf("failed to create run directory: %w", err)
	}

	// Clone requirements.yml to run directory for safe testing
	outputFile := filepath.Join(runDir.GetPath(), "requirements-merged.yml")
	if err := c.cloneRequirements(outputFile); err != nil {
		return fmt.Errorf("failed to clone requirements file: %w", err)
	}

	// Load story document
	storyDoc, err := c.storyLoader.Load(storyNumber)
	if err != nil {
		return fmt.Errorf("failed to load story: %w", err)
	}

	// Merge scenarios from story into requirements-merged.yml (test file)
	if err := c.mergeScenarios(ctx, storyNumber, storyDoc, outputFile); err != nil {
		return fmt.Errorf("failed to merge scenarios: %w", err)
	}

	slog.Info("User story implementation completed successfully")
	fmt.Println("\n✅ Scenario merge completed successfully!")
	fmt.Printf("Merged %d scenarios from story %s into %s\n", len(storyDoc.Scenarios.TestScenarios), storyNumber, outputFile)
	fmt.Printf("\nTo review changes: diff docs/requirements.yml %s\n", outputFile)
	return nil
}

func (c *USImplementCommand) cloneRequirements(outputFile string) error {
	slog.Info("Cloning requirements file for safe testing", "output", outputFile)

	// Ensure tmp directory exists
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Read original file
	data, err := os.ReadFile("docs/requirements.yml")
	if err != nil {
		return fmt.Errorf("failed to read requirements.yml: %w", err)
	}

	// Write to output file
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputFile, err)
	}

	fmt.Printf("✓ Cloned docs/requirements.yml → %s\n", outputFile)
	return nil
}

func (c *USImplementCommand) mergeScenarios(ctx context.Context, storyNumber string, storyDoc *storyModels.StoryDocument, outputFile string) error {
	slog.Info("Starting scenario merge", "story_id", storyNumber, "scenario_count", len(storyDoc.Scenarios.TestScenarios), "output_file", outputFile)

	// Create template loaders
	userPromptPath := c.config.GetString("templates.prompts.merge_scenarios")
	systemPromptPath := c.config.GetString("templates.prompts.merge_scenarios_system")
	userPromptLoader := template.NewTemplateLoader[*template.ScenarioMergeData](userPromptPath)
	systemPromptLoader := template.NewTemplateLoader[*template.ScenarioMergeData](systemPromptPath)

	// Process each scenario individually
	for i, scenario := range storyDoc.Scenarios.TestScenarios {
		slog.Info("Processing scenario", "index", i+1, "scenario_id", scenario.ID)
		fmt.Printf("\nMerging scenario %d/%d: %s\n", i+1, len(storyDoc.Scenarios.TestScenarios), scenario.ID)

		// Create merge data adapter
		mergeData := template.NewScenarioMergeData(storyNumber, scenario, outputFile)

		// Load templates using TemplateLoader
		userPrompt, err := userPromptLoader.LoadTemplate(mergeData)
		if err != nil {
			return fmt.Errorf("failed to load user prompt for scenario %s: %w", scenario.ID, err)
		}

		systemPrompt, err := systemPromptLoader.LoadTemplate(mergeData)
		if err != nil {
			return fmt.Errorf("failed to load system prompt for scenario %s: %w", scenario.ID, err)
		}

		// Call Claude Code API to analyze and merge
		slog.Debug("Calling Claude Code for scenario merge", "scenario_id", scenario.ID)
		_, err = c.claudeClient.ExecutePromptWithSystem(
			ctx,
			systemPrompt,
			userPrompt,
			"sonnet",
			ai.ExecutionMode{
				AllowedTools: []string{"Read", "Edit"},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to merge scenario %s: %w", scenario.ID, err)
		}

		slog.Info("Scenario merged successfully", "scenario_id", scenario.ID)
		fmt.Printf("✓ Merged scenario: %s\n", scenario.ID)
	}

	slog.Info("All scenarios merged successfully", "total_count", len(storyDoc.Scenarios.TestScenarios))
	return nil
}
