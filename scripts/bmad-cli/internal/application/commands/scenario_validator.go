package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// minRegexMatchGroups is the minimum number of match groups expected from the test pattern.
	minRegexMatchGroups = 2
)

// ScenarioValidator validates that all scenarios from requirements.yml
// have corresponding tests in the test files.
type ScenarioValidator struct {
	requirementsFile string
	testsDir         string
}

// ValidationResult holds the result of scenario validation.
type ValidationResult struct {
	TotalScenarios   int
	CoveredCount     int
	MissingScenarios []string
	ExtraTests       []string
}

// NewScenarioValidator creates a new scenario validator.
func NewScenarioValidator(requirementsFile, testsDir string) *ScenarioValidator {
	return &ScenarioValidator{
		requirementsFile: requirementsFile,
		testsDir:         testsDir,
	}
}

// Validate checks that all scenarios have corresponding tests.
func (v *ScenarioValidator) Validate() (*ValidationResult, error) {
	// Get scenario IDs from requirements.yml
	scenarioIDs, err := v.getScenarioIDs()
	if err != nil {
		return nil, err
	}

	// Get test IDs from test files
	testIDs, err := v.getTestIDs()
	if err != nil {
		return nil, err
	}

	// Compare and build result
	result := &ValidationResult{
		TotalScenarios:   len(scenarioIDs),
		MissingScenarios: make([]string, 0),
		ExtraTests:       make([]string, 0),
	}

	// Create a map for quick lookup
	testIDMap := make(map[string]bool)
	for _, id := range testIDs {
		testIDMap[id] = true
	}

	scenarioIDMap := make(map[string]bool)
	for _, id := range scenarioIDs {
		scenarioIDMap[id] = true
	}

	// Find missing scenarios (in requirements but not in tests)
	for _, id := range scenarioIDs {
		if testIDMap[id] {
			result.CoveredCount++
		} else {
			result.MissingScenarios = append(result.MissingScenarios, id)
		}
	}

	// Find extra tests (in tests but not in requirements)
	for _, id := range testIDs {
		if !scenarioIDMap[id] {
			result.ExtraTests = append(result.ExtraTests, id)
		}
	}

	return result, nil
}

// getScenarioIDs extracts scenario IDs from requirements.yml.
func (v *ScenarioValidator) getScenarioIDs() ([]string, error) {
	data, err := os.ReadFile(v.requirementsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read requirements file: %w", err)
	}

	var requirements struct {
		Scenarios map[string]struct {
			ImplementationStatus struct {
				Status string `yaml:"status"`
			} `yaml:"implementation_status"`
		} `yaml:"scenarios"`
	}

	err = yaml.Unmarshal(data, &requirements)
	if err != nil {
		return nil, fmt.Errorf("failed to parse requirements YAML: %w", err)
	}

	// Collect all scenario IDs (regardless of status)
	ids := make([]string, 0, len(requirements.Scenarios))
	for id := range requirements.Scenarios {
		ids = append(ids, id)
	}

	return ids, nil
}

// getTestIDs scans test files for test IDs matching the pattern 'ID: description'.
func (v *ScenarioValidator) getTestIDs() ([]string, error) {
	ids := make([]string, 0)

	// Pattern to match test IDs like: 'INT-012: description' or "E2E-001: description"
	testPattern := regexp.MustCompile(`['"]([A-Z0-9]+-\d+):`)

	// Walk through all .ts files in the tests directory
	err := filepath.Walk(v.testsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		// Skip directories and non-ts files
		if info.IsDir() || !strings.HasSuffix(path, ".ts") {
			return nil
		}

		// Read and scan the file
		fileIDs, scanErr := v.scanFileForTestIDs(path, testPattern)
		if scanErr != nil {
			return scanErr
		}

		ids = append(ids, fileIDs...)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk tests directory: %w", err)
	}

	return ids, nil
}

// scanFileForTestIDs scans a single file for test IDs.
func (v *ScenarioValidator) scanFileForTestIDs(path string, pattern *regexp.Regexp) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	defer func() {
		_ = file.Close()
	}()

	var ids []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		matches := pattern.FindStringSubmatch(line)
		if len(matches) >= minRegexMatchGroups {
			ids = append(ids, matches[1])
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		return nil, fmt.Errorf("error scanning file %s: %w", path, scanErr)
	}

	return ids, nil
}
