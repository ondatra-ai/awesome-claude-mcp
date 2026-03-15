package factories

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/app/generators/implement"
	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/infrastructure/fs"
	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/shell"
	"bmad-cli/internal/infrastructure/story"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// StoryValidator validates that a story exists.
type StoryValidator interface {
	Execute(storyNumber string) error
}

// ScenarioMerger merges scenarios into requirements.
type ScenarioMerger interface {
	Execute(ctx context.Context, storyNumber string) error
}

// File permission constants.
const fileModeReadWrite = 0o644

// ImplementFactory orchestrates the user story implementation workflow.
// It coordinates generators for each step of the implementation process.
type ImplementFactory struct {
	branchManager *git.BranchManager
	storyLoader   *story.StoryLoader
	claudeClient  *ai.ClaudeClient
	config        *config.ViperConfig
	runDir        *fs.RunDirectory
	shellExec     *shell.CommandRunner

	// Delegated commands
	validateCmd       StoryValidator
	mergeScenariosCmd ScenarioMerger

	// Generators
	mergeScenariosGen *implement.MergeScenariosGenerator
	testCodeGen       *implement.TestCodeGenerator
	testValidatorGen  *implement.TestValidatorGenerator
	featureImplGen    *implement.FeatureImplementerGenerator
}

// NewImplementFactory creates a new ImplementFactory with all dependencies.
func NewImplementFactory(
	branchManager *git.BranchManager,
	storyLoader *story.StoryLoader,
	claudeClient *ai.ClaudeClient,
	cfg *config.ViperConfig,
	runDir *fs.RunDirectory,
	shellExec *shell.CommandRunner,
	validateCmd StoryValidator,
	mergeScenariosCmd ScenarioMerger,
) *ImplementFactory {
	return &ImplementFactory{
		branchManager:     branchManager,
		storyLoader:       storyLoader,
		claudeClient:      claudeClient,
		config:            cfg,
		runDir:            runDir,
		shellExec:         shellExec,
		validateCmd:       validateCmd,
		mergeScenariosCmd: mergeScenariosCmd,
		mergeScenariosGen: implement.NewMergeScenariosGenerator(claudeClient, cfg),
		testCodeGen:       implement.NewTestCodeGenerator(claudeClient, cfg),
		testValidatorGen:  implement.NewTestValidatorGenerator(claudeClient, cfg),
		featureImplGen:    implement.NewFeatureImplementerGenerator(claudeClient, cfg),
	}
}

// Execute runs the implementation workflow based on the specified steps.
func (f *ImplementFactory) Execute(
	ctx context.Context,
	storyNumber string,
	steps *implement.ExecutionSteps,
	force bool,
) error {
	slog.Info("Starting user story implementation",
		"story_number", storyNumber,
		"force", force,
		"steps", steps.String(),
	)

	err := f.executePreparationSteps(storyNumber, steps, force)
	if err != nil {
		return err
	}

	err = f.executeImplementationSteps(ctx, storyNumber, steps)
	if err != nil {
		return err
	}

	slog.Info("✅ User story implementation completed successfully")

	return nil
}

// executePreparationSteps runs non-AI preparation steps (validate, branch).
func (f *ImplementFactory) executePreparationSteps(
	storyNumber string,
	steps *implement.ExecutionSteps,
	force bool,
) error {
	if steps.ValidateStory {
		err := f.validateStory(storyNumber)
		if err != nil {
			return err
		}
	}

	if steps.CreateBranch {
		err := f.createBranch(storyNumber, force)
		if err != nil {
			return err
		}
	}

	return nil
}

// executeImplementationSteps runs AI-powered implementation steps.
func (f *ImplementFactory) executeImplementationSteps(
	ctx context.Context,
	storyNumber string,
	steps *implement.ExecutionSteps,
) error {
	tmpDir := f.runDir.GetTmpOutPath()

	err := f.runMergeScenariosIfEnabled(ctx, storyNumber, steps)
	if err != nil {
		return err
	}

	err = f.runGenerateTestsIfEnabled(ctx, tmpDir, steps)
	if err != nil {
		return err
	}

	err = f.runValidateTestsIfEnabled(ctx, tmpDir, steps)
	if err != nil {
		return err
	}

	err = f.runValidateScenariosIfEnabled(steps)
	if err != nil {
		return err
	}

	return f.runImplementFeatureIfEnabled(ctx, storyNumber, tmpDir, steps)
}

func (f *ImplementFactory) runMergeScenariosIfEnabled(
	ctx context.Context,
	storyNumber string,
	steps *implement.ExecutionSteps,
) error {
	if !steps.MergeScenarios {
		return nil
	}

	return f.executeMergeScenarios(ctx, storyNumber)
}

func (f *ImplementFactory) runGenerateTestsIfEnabled(
	ctx context.Context,
	tmpDir string,
	steps *implement.ExecutionSteps,
) error {
	if !steps.GenerateTests {
		return nil
	}

	return f.executeGenerateTests(ctx, tmpDir)
}

func (f *ImplementFactory) runValidateTestsIfEnabled(
	ctx context.Context,
	tmpDir string,
	steps *implement.ExecutionSteps,
) error {
	if !steps.ValidateTests {
		return nil
	}

	return f.executeValidateTests(ctx, tmpDir)
}

func (f *ImplementFactory) runValidateScenariosIfEnabled(steps *implement.ExecutionSteps) error {
	if !steps.ValidateScenarios {
		return nil
	}

	return f.executeValidateScenarios()
}

func (f *ImplementFactory) runImplementFeatureIfEnabled(
	ctx context.Context,
	storyNumber, tmpDir string,
	steps *implement.ExecutionSteps,
) error {
	if !steps.ImplementFeature {
		return nil
	}

	return f.executeImplementFeature(ctx, storyNumber, tmpDir)
}

func (f *ImplementFactory) validateStory(storyNumber string) error {
	slog.Info("Step 1: Validating story file")

	err := f.validateCmd.Execute(storyNumber)
	if err != nil {
		return fmt.Errorf("validate story: %w", err)
	}

	return nil
}

func (f *ImplementFactory) createBranch(storyNumber string, force bool) error {
	slog.Info("Step 2: Creating story branch")

	storySlug, err := f.storyLoader.GetStorySlug(storyNumber)
	if err != nil {
		return pkgerrors.ErrGetStorySlugFailed(err)
	}

	// TEMPORARILY COMMENTED OUT FOR TESTING
	// if err := f.branchManager.EnsureBranch(ctx, storyNumber, storySlug, force); err != nil {
	// 	return fmt.Errorf("failed to ensure branch: %w", err)
	// }
	_ = storySlug
	_ = force

	slog.Warn("⚠️  Branch management temporarily disabled for testing")

	return nil
}

func (f *ImplementFactory) executeMergeScenarios(
	ctx context.Context,
	storyNumber string,
) error {
	slog.Info("Step 3: Merging scenarios into requirements")

	err := f.mergeScenariosCmd.Execute(ctx, storyNumber)
	if err != nil {
		return fmt.Errorf("merge scenarios: %w", err)
	}

	return nil
}

func (f *ImplementFactory) executeGenerateTests(ctx context.Context, tmpDir string) error {
	slog.Info("Step 4: Generating test code")

	err := f.validateBaselineTests(ctx)
	if err != nil {
		return err
	}

	status, err := f.testCodeGen.GenerateTests(ctx, "docs/requirements.yaml", tmpDir)
	if err != nil {
		return pkgerrors.ErrGenerateTestsFailed(err)
	}

	slog.Info("✅ Test generation completed", "implemented", status.ItemsProcessed)

	return f.validateGeneratedTests(ctx)
}

func (f *ImplementFactory) executeValidateTests(ctx context.Context, tmpDir string) error {
	slog.Info("Step 5: Validating test quality (Claude-based)")

	result, err := f.testValidatorGen.ValidateTests(ctx, tmpDir)
	if err != nil {
		return pkgerrors.ErrValidateTestsFailed(err)
	}

	if !result.Success {
		return pkgerrors.ErrUnfixedTestIssuesError(result.ItemsProcessed)
	}

	slog.Info("✅ Test validation completed successfully")

	return nil
}

func (f *ImplementFactory) executeValidateScenarios() error {
	slog.Info("Step 6: Validating scenario coverage (Go-based)")

	validator := implement.NewScenarioValidator("docs/requirements.yaml", "tests")

	result, err := validator.Validate()
	if err != nil {
		return pkgerrors.ErrValidateScenariosFailed(err)
	}

	slog.Info("Scenario validation results",
		"total_scenarios", result.TotalScenarios,
		"covered", result.CoveredCount,
		"missing", len(result.MissingScenarios),
	)

	if len(result.MissingScenarios) > 0 {
		slog.Warn("⚠️  Missing scenario coverage:")

		for _, missing := range result.MissingScenarios {
			slog.Warn("  - " + missing)
		}

		return pkgerrors.ErrMissingScenarioCoverageError(result.MissingScenarios)
	}

	slog.Info("✅ All scenarios have test coverage")

	return nil
}

func (f *ImplementFactory) executeImplementFeature(
	ctx context.Context,
	storyNumber string,
	tmpDir string,
) error {
	slog.Info("Step 7: Implementing feature")

	storyDoc, err := f.storyLoader.Load(storyNumber)
	if err != nil {
		return pkgerrors.ErrLoadStoryFailed(err)
	}

	const maxAttempts = 5

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		slog.Info("🔄 Implementation attempt", "attempt", attempt, "max", maxAttempts)

		testOutput, testErr := f.runTests(ctx, "implement-feature")
		if testErr == nil {
			slog.Info("✅ All tests passing - feature implementation complete!", "attempts_used", attempt)

			return nil
		}

		slog.Info("❌ Tests failing - calling Claude to fix", "attempt", attempt)

		status, err := f.featureImplGen.Implement(ctx, storyDoc, attempt, testOutput, tmpDir)
		if err != nil {
			return fmt.Errorf("feature implementation attempt %d failed: %w", attempt, err)
		}

		slog.Info("✓ Claude finished attempt", "attempt", attempt, "status", status.Message)
	}

	return pkgerrors.ErrImplementFeaturesMaxAttemptsExceeded(maxAttempts)
}

func (f *ImplementFactory) runTests(ctx context.Context, phase string) (string, error) {
	slog.Info("🧪 Running tests", "phase", phase)

	testCommand := f.config.GetString("testing.command")
	if testCommand == "" {
		testCommand = "make test-e2e"
	}

	output, err := f.shellExec.Run(ctx, "sh", "-c", testCommand)
	if err != nil {
		return output, pkgerrors.ErrRunTestsFailed(phase, err)
	}

	outputFile := filepath.Join(f.runDir.GetTmpOutPath(), "test-output-"+phase+".txt")

	writeErr := os.WriteFile(outputFile, []byte(output), fileModeReadWrite)
	if writeErr != nil {
		slog.Warn("Failed to write test output", "file", outputFile, "error", writeErr)
	}

	return output, nil
}

func (f *ImplementFactory) validateBaselineTests(ctx context.Context) error {
	slog.Info("📋 Validating baseline tests (must pass)")

	output, err := f.runTests(ctx, "before")
	if err != nil {
		slog.Error("❌ Baseline tests failed", "error", err)

		return pkgerrors.ErrBaselineTestsFailedError(output)
	}

	slog.Info("✅ Baseline tests passed - ready for test generation")

	return nil
}

func (f *ImplementFactory) validateGeneratedTests(ctx context.Context) error {
	slog.Info("🔴 Validating generated tests (must fail - TDD red phase)")

	output, err := f.runTests(ctx, "after")
	if err == nil {
		slog.Error("❌ Generated tests are passing but should fail (TDD red phase)")

		return pkgerrors.ErrGeneratedTestsPassError(output)
	}

	slog.Info("✅ Generated tests are failing as expected (TDD red phase)")

	return nil
}
