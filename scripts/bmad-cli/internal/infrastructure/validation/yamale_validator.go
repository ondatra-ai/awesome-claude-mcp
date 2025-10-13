package validation

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type YamaleValidator struct {
	schemaPath string
}

func NewYamaleValidator(schemaPath string) *YamaleValidator {
	return &YamaleValidator{
		schemaPath: schemaPath,
	}
}

func (v *YamaleValidator) Validate(yamlContent string) error {
	// Get absolute schema path
	absSchemaPath, err := filepath.Abs(v.schemaPath)
	if err != nil {
		slog.Error("failed to get absolute schema path", "error", err, "path", v.schemaPath)

		return errors.Join(errors.New("failed to get absolute schema path"), err)
	}

	// Verify schema file exists
	_, err = os.Stat(absSchemaPath)
	if os.IsNotExist(err) {
		slog.Error("schema file does not exist", "path", absSchemaPath)

		return errors.New("schema file does not exist: " + absSchemaPath)
	}

	// Create temporary file for YAML content
	tmpFile, err := os.CreateTemp("", "story-*.yaml")
	if err != nil {
		slog.Error("failed to create temporary file", "error", err)

		return errors.Join(errors.New("failed to create temporary file"), err)
	}

	defer func() {
		err := os.Remove(tmpFile.Name())
		if err != nil {
			slog.Warn("failed to remove temporary file", "path", tmpFile.Name(), "error", err)
		}
	}()

	// Write YAML content to temp file
	_, err = tmpFile.WriteString(yamlContent)
	if err != nil {
		slog.Error("failed to write YAML content", "error", err)

		return errors.Join(errors.New("failed to write YAML content"), err)
	}

	err = tmpFile.Close()
	if err != nil {
		slog.Error("failed to close temporary file", "error", err)

		return errors.Join(errors.New("failed to close temporary file"), err)
	}

	// Use yamale CLI tool
	return v.validateWithYamaleCLI(absSchemaPath, tmpFile.Name())
}

func (v *YamaleValidator) ValidateFromStdin(yamlContent string) error {
	// Get absolute schema path
	absSchemaPath, err := filepath.Abs(v.schemaPath)
	if err != nil {
		slog.Error("failed to get absolute schema path", "error", err, "path", v.schemaPath)

		return errors.Join(errors.New("failed to get absolute schema path"), err)
	}

	// Try yamale with stdin
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "yamale", "--schema", absSchemaPath, "-")
	cmd.Stdin = strings.NewReader(yamlContent)

	_, err = cmd.CombinedOutput()
	if err != nil {
		// Try with python module if CLI fails
		cmd = exec.CommandContext(ctx, "python", "-m", "yamale", "--schema", absSchemaPath, "-")
		cmd.Stdin = strings.NewReader(yamlContent)

		output, err := cmd.CombinedOutput()
		if err != nil {
			slog.Error("yamale validation failed", "error", err, "output", string(output))

			return errors.Join(errors.New("yamale validation failed"), err)
		}
	}

	return nil
}

func (v *YamaleValidator) validateWithYamaleCLI(schemaPath, dataPath string) error {
	// Check if yamale command is available first
	_, err := exec.LookPath("yamale")
	if err != nil {
		slog.Warn("yamale command not found - skipping validation")

		return err
	}

	// Run yamale validation
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "yamale", "-s", schemaPath, dataPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("yamale validation failed", "error", err, "output", string(output))

		return errors.Join(errors.New("yamale validation failed"), err)
	}

	slog.Info("✅ YAML validation passed")

	return nil
}
