package commands

import (
	"context"
	"log/slog"

	"bmad-cli/internal/app/generators/implement"
	"bmad-cli/internal/infrastructure/fs"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// ReqGenerateTestsCommand generates test code for pending scenarios in requirements.yaml.
type ReqGenerateTestsCommand struct {
	testCodeGen *implement.TestCodeGenerator
	runDir      *fs.RunDirectory
}

// NewReqGenerateTestsCommand creates a new ReqGenerateTestsCommand.
func NewReqGenerateTestsCommand(
	testCodeGen *implement.TestCodeGenerator,
	runDir *fs.RunDirectory,
) *ReqGenerateTestsCommand {
	return &ReqGenerateTestsCommand{
		testCodeGen: testCodeGen,
		runDir:      runDir,
	}
}

// Execute generates tests for all pending scenarios in the given requirements file.
func (c *ReqGenerateTestsCommand) Execute(
	ctx context.Context,
	requirementsFile string,
) error {
	slog.Info("Generating tests from requirements",
		"requirements_file", requirementsFile,
	)

	tmpDir := c.runDir.GetTmpOutPath()

	status, err := c.testCodeGen.GenerateTests(ctx, requirementsFile, tmpDir)
	if err != nil {
		return pkgerrors.ErrGenerateTestsFailed(err)
	}

	slog.Info("Test generation completed",
		"implemented", status.ItemsProcessed,
	)

	return nil
}
