package validation

import (
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
	if _, err := os.Stat(absSchemaPath); os.IsNotExist(err) {
		slog.Error("schema file does not exist", "path", absSchemaPath)

		return errors.New("schema file does not exist: " + absSchemaPath)
	}

	// Create temporary file for YAML content
	tmpFile, err := os.CreateTemp("", "story-*.yaml")
	if err != nil {
		slog.Error("failed to create temporary file", "error", err)

		return errors.Join(errors.New("failed to create temporary file"), err)
	}
	defer os.Remove(tmpFile.Name())

	// Write YAML content to temp file
	if _, err := tmpFile.WriteString(yamlContent); err != nil {
		slog.Error("failed to write YAML content", "error", err)

		return errors.Join(errors.New("failed to write YAML content"), err)
	}

	if err := tmpFile.Close(); err != nil {
		slog.Error("failed to close temporary file", "error", err)

		return errors.Join(errors.New("failed to close temporary file"), err)
	}

	// Use yamale CLI tool
	return v.validateWithYamaleCLI(absSchemaPath, tmpFile.Name())
}

func (v *YamaleValidator) validateWithYamaleCLI(schemaPath, dataPath string) error {
	// Check if yamale command is available first
	if _, err := exec.LookPath("yamale"); err != nil {
		slog.Warn("yamale command not found - skipping validation")

		return nil
	}

	// Run yamale validation
	cmd := exec.Command("yamale", "-s", schemaPath, dataPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("yamale validation failed", "error", err, "output", string(output))

		return errors.Join(errors.New("yamale validation failed"), err)
	}

	slog.Info("✅ YAML validation passed")

	return nil
}

func (v *YamaleValidator) validateWithPythonModule(schemaPath, dataPath string) error {
	// Try python -m yamale
	cmd := exec.Command("python", "-m", "yamale", "--schema", schemaPath, dataPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try python3 if python failed
		cmd = exec.Command("python3", "-m", "yamale", "--schema", schemaPath, dataPath)

		output, err = cmd.CombinedOutput()
		if err != nil {
			slog.Error("python yamale validation failed", "error", err, "output", string(output))

			return errors.Join(errors.New("python yamale validation failed"), err)
		}
	}

	return nil
}

func (v *YamaleValidator) ValidateFromStdin(yamlContent string) error {
	// Get absolute schema path
	absSchemaPath, err := filepath.Abs(v.schemaPath)
	if err != nil {
		slog.Error("failed to get absolute schema path", "error", err, "path", v.schemaPath)

		return errors.Join(errors.New("failed to get absolute schema path"), err)
	}

	// Try yamale with stdin
	cmd := exec.Command("yamale", "--schema", absSchemaPath, "-")
	cmd.Stdin = strings.NewReader(yamlContent)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try with python module if CLI fails
		cmd = exec.Command("python", "-m", "yamale", "--schema", absSchemaPath, "-")
		cmd.Stdin = strings.NewReader(yamlContent)

		output, err = cmd.CombinedOutput()
		if err != nil {
			slog.Error("yamale validation failed", "error", err, "output", string(output))

			return errors.Join(errors.New("yamale validation failed"), err)
		}
	}

	return nil
}
