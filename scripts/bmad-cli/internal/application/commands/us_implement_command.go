package commands

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/story"
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

	// Clone requirements.yml to tmp folder for safe testing
	tmpDir := c.config.GetString("paths.tmp_dir")
	outputFile := filepath.Join(tmpDir, "requirements-merged.yml")
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

	// Process each scenario individually
	for i, scenario := range storyDoc.Scenarios.TestScenarios {
		slog.Info("Processing scenario", "index", i+1, "scenario_id", scenario.ID)
		fmt.Printf("\nMerging scenario %d/%d: %s\n", i+1, len(storyDoc.Scenarios.TestScenarios), scenario.ID)

		// Load templates
		userPrompt, systemPrompt, err := c.loadMergeTemplates()
		if err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}

		// Render template with scenario data and output file
		prompt, err := c.renderMergePrompt(userPrompt, storyNumber, scenario, outputFile)
		if err != nil {
			return fmt.Errorf("failed to render prompt for scenario %s: %w", scenario.ID, err)
		}

		// Call Claude Code API to analyze and merge
		slog.Debug("Calling Claude Code for scenario merge", "scenario_id", scenario.ID)
		_, err = c.claudeClient.ExecutePromptWithSystem(
			ctx,
			systemPrompt,
			prompt,
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

func (c *USImplementCommand) loadMergeTemplates() (string, string, error) {
	// Load user prompt template
	userPromptPath := c.config.GetString("templates.prompts.merge_scenarios")
	userPromptBytes, err := os.ReadFile(userPromptPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read user prompt template: %w", err)
	}

	// Load system prompt template
	systemPromptPath := c.config.GetString("templates.prompts.merge_scenarios_system")
	systemPromptBytes, err := os.ReadFile(systemPromptPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read system prompt template: %w", err)
	}

	return string(userPromptBytes), string(systemPromptBytes), nil
}

func (c *USImplementCommand) renderMergePrompt(templateStr string, storyNumber string, scenario storyModels.TestScenario, outputFile string) (string, error) {
	// Create template
	tmpl, err := template.New("merge").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Format steps for template
	stepsStr := c.formatScenarioSteps(scenario)

	// Format acceptance criteria for template
	acStr := c.formatAcceptanceCriteria(scenario.AcceptanceCriteria)

	// Prepare template data
	data := map[string]interface{}{
		"StoryNumber":         storyNumber,
		"ScenarioID":          scenario.ID,
		"Level":               scenario.Level,
		"Priority":            scenario.Priority,
		"AcceptanceCriteria":  acStr,
		"Steps":               stepsStr,
		"RequirementsFile":    outputFile,
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func (c *USImplementCommand) formatScenarioSteps(scenario storyModels.TestScenario) string {
	var result strings.Builder

	for _, step := range scenario.Steps {
		if len(step.Given) > 0 {
			result.WriteString("  Given:\n")
			for _, g := range step.Given {
				result.WriteString(fmt.Sprintf("    - %s\n", g))
			}
		}
		if len(step.When) > 0 {
			result.WriteString("  When:\n")
			for _, w := range step.When {
				result.WriteString(fmt.Sprintf("    - %s\n", w))
			}
		}
		if len(step.Then) > 0 {
			result.WriteString("  Then:\n")
			for _, t := range step.Then {
				result.WriteString(fmt.Sprintf("    - %s\n", t))
			}
		}
	}

	return result.String()
}

func (c *USImplementCommand) formatAcceptanceCriteria(criteria []string) string {
	if len(criteria) == 0 {
		return "[]"
	}
	return fmt.Sprintf(`["%s"]`, strings.Join(criteria, `", "`))
}
