package validation

import (
	"fmt"
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
		return fmt.Errorf("failed to get absolute schema path: %w", err)
	}

	// Verify schema file exists
	if _, err := os.Stat(absSchemaPath); os.IsNotExist(err) {
		return fmt.Errorf("schema file does not exist: %s", absSchemaPath)
	}

	// Create temporary file for YAML content
	tmpFile, err := os.CreateTemp("", "story-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write YAML content to temp file
	if _, err := tmpFile.WriteString(yamlContent); err != nil {
		return fmt.Errorf("failed to write YAML content: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Use yamale CLI tool
	return v.validateWithYamaleCLI(absSchemaPath, tmpFile.Name())
}

func (v *YamaleValidator) validateWithYamaleCLI(schemaPath, dataPath string) error {
	// Check if yamale command is available first
	if _, err := exec.LookPath("yamale"); err != nil {
		fmt.Println("Warning: yamale command not found - skipping validation")
		return nil
	}

	// Run yamale validation
	cmd := exec.Command("yamale", "-s", schemaPath, dataPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("yamale validation failed: %s\nOutput: %s", err, string(output))
	}

	fmt.Println("✅ YAML validation passed")
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
			return fmt.Errorf("python yamale validation failed: %s\nOutput: %s", err, string(output))
		}
	}

	return nil
}

func (v *YamaleValidator) ValidateFromStdin(yamlContent string) error {
	// Get absolute schema path
	absSchemaPath, err := filepath.Abs(v.schemaPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute schema path: %w", err)
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
			return fmt.Errorf("yamale validation failed: %s\nOutput: %s", err, string(output))
		}
	}

	return nil
}
