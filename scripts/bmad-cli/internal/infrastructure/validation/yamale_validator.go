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

	// Try different validation methods
	if err := v.validateWithYamaleCLI(absSchemaPath, tmpFile.Name()); err != nil {
		if err := v.validateWithPythonModule(absSchemaPath, tmpFile.Name()); err != nil {
			return fmt.Errorf("yamale validation failed: %w", err)
		}
	}

	return nil
}

func (v *YamaleValidator) validateWithYamaleCLI(schemaPath, dataPath string) error {
	// Try yamale CLI first
	cmd := exec.Command("yamale", "--schema", schemaPath, dataPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("yamale CLI validation failed: %s\nOutput: %s", err, string(output))
	}

	return nil
}

func (v *YamaleValidator) validateWithPythonModule(schemaPath, dataPath string) error {
	// Try python -m yamale
	cmd := exec.Command("python", "-m", "yamale", "--schema", schemaPath, dataPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("python yamale validation failed: %s\nOutput: %s", err, string(output))
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

func (v *YamaleValidator) GetSchemaPath() string {
	return v.schemaPath
}

func (v *YamaleValidator) SetSchemaPath(path string) {
	v.schemaPath = path
}

func (v *YamaleValidator) ValidateSchema() error {
	absPath, err := filepath.Abs(v.schemaPath)
	if err != nil {
		return fmt.Errorf("failed to resolve schema path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("schema file does not exist: %s", absPath)
	}

	return nil
}

func (v *YamaleValidator) CheckYamaleAvailable() error {
	// Check if yamale CLI is available
	if err := exec.Command("yamale", "--version").Run(); err == nil {
		return nil
	}

	// Check if python yamale module is available
	if err := exec.Command("python", "-c", "import yamale").Run(); err == nil {
		return nil
	}

	// Check if python3 yamale module is available
	if err := exec.Command("python3", "-c", "import yamale").Run(); err == nil {
		return nil
	}

	return fmt.Errorf("yamale is not available - please install with: pip install yamale")
}
