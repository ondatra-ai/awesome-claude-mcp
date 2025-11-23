package commands

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidStep = errors.New("invalid step")

// Step constants define available implementation steps.
const (
	StepValidateStory     = "validate_story"
	StepCreateBranch      = "create_branch"
	StepMergeScenarios    = "merge_scenarios"
	StepGenerateTests     = "generate_tests"
	StepValidateTests     = "validate_tests"
	StepValidateScenarios = "validate_scenarios"
	StepImplementFeature  = "implement_feature"
	StepAll               = "all"
)

// ExecutionSteps represents which steps should be executed.
type ExecutionSteps struct {
	ValidateStory     bool
	CreateBranch      bool
	MergeScenarios    bool
	GenerateTests     bool
	ValidateTests     bool
	ValidateScenarios bool
	ImplementFeature  bool
}

// getStepSetters returns a map of step names to functions that enable them.
func getStepSetters() map[string]func(*ExecutionSteps) {
	return map[string]func(*ExecutionSteps){
		StepValidateStory:     func(s *ExecutionSteps) { s.ValidateStory = true },
		StepCreateBranch:      func(s *ExecutionSteps) { s.CreateBranch = true },
		StepMergeScenarios:    func(s *ExecutionSteps) { s.MergeScenarios = true },
		StepGenerateTests:     func(s *ExecutionSteps) { s.GenerateTests = true },
		StepValidateTests:     func(s *ExecutionSteps) { s.ValidateTests = true },
		StepValidateScenarios: func(s *ExecutionSteps) { s.ValidateScenarios = true },
		StepImplementFeature:  func(s *ExecutionSteps) { s.ImplementFeature = true },
	}
}

// enableAllSteps enables all steps except create_branch.
func enableAllSteps(steps *ExecutionSteps) {
	steps.ValidateStory = true
	steps.MergeScenarios = true
	steps.GenerateTests = true
	steps.ValidateTests = true
	steps.ValidateScenarios = true
	steps.ImplementFeature = true
}

// ParseSteps parses the steps string and returns ExecutionSteps.
func ParseSteps(stepsStr string) (*ExecutionSteps, error) {
	stepsStr = strings.TrimSpace(stepsStr)
	if stepsStr == "" {
		stepsStr = StepAll
	}

	steps := &ExecutionSteps{}
	stepList := strings.Split(stepsStr, ",")

	for _, step := range stepList {
		step = strings.TrimSpace(step)

		err := applyStep(steps, step)
		if err != nil {
			return nil, err
		}
	}

	return steps, nil
}

// applyStep applies a single step to the execution steps.
func applyStep(steps *ExecutionSteps, step string) error {
	if step == StepAll {
		enableAllSteps(steps)

		return nil
	}

	setter, ok := getStepSetters()[step]
	if !ok {
		return fmt.Errorf("%w: %s (valid steps: %s)",
			ErrInvalidStep,
			step,
			strings.Join(validStepNames(), ", "),
		)
	}

	setter(steps)

	return nil
}

// validStepNames returns all valid step names.
func validStepNames() []string {
	return []string{
		StepValidateStory,
		StepCreateBranch,
		StepMergeScenarios,
		StepGenerateTests,
		StepValidateTests,
		StepValidateScenarios,
		StepImplementFeature,
		StepAll,
	}
}

// String returns a string representation of the execution steps.
func (e *ExecutionSteps) String() string {
	enabled := e.enabledSteps()

	if len(enabled) == 0 {
		return "none"
	}

	return strings.Join(enabled, ", ")
}

// enabledSteps returns a slice of enabled step names.
func (e *ExecutionSteps) enabledSteps() []string {
	var enabled []string

	stepChecks := []struct {
		enabled bool
		name    string
	}{
		{e.ValidateStory, StepValidateStory},
		{e.CreateBranch, StepCreateBranch},
		{e.MergeScenarios, StepMergeScenarios},
		{e.GenerateTests, StepGenerateTests},
		{e.ValidateTests, StepValidateTests},
		{e.ValidateScenarios, StepValidateScenarios},
		{e.ImplementFeature, StepImplementFeature},
	}

	for _, check := range stepChecks {
		if check.enabled {
			enabled = append(enabled, check.name)
		}
	}

	return enabled
}
