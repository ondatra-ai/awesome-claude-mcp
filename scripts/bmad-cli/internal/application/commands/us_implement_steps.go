package commands

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidStep = errors.New("invalid step")

// Step constants define available implementation steps.
const (
	StepValidateStory    = "validate_story"
	StepCreateBranch     = "create_branch"
	StepMergeScenarios   = "merge_scenarios"
	StepGenerateTests    = "generate_tests"
	StepImplementFeature = "implement_feature"
	StepAll              = "all"
)

// ExecutionSteps represents which steps should be executed.
type ExecutionSteps struct {
	ValidateStory    bool
	CreateBranch     bool
	MergeScenarios   bool
	GenerateTests    bool
	ImplementFeature bool
}

// ParseSteps parses the steps string and returns ExecutionSteps.
func ParseSteps(stepsStr string) (*ExecutionSteps, error) {
	stepsStr = strings.TrimSpace(stepsStr)
	if stepsStr == "" {
		stepsStr = StepAll
	}

	// Parse comma-separated steps
	steps := &ExecutionSteps{}
	stepList := strings.Split(stepsStr, ",")

	for _, step := range stepList {
		step = strings.TrimSpace(step)

		switch step {
		case StepAll:
			// "all" enables all steps EXCEPT create_branch
			// must be explicitly requested
			steps.ValidateStory = true
			steps.MergeScenarios = true
			steps.GenerateTests = true
			steps.ImplementFeature = true
		case StepValidateStory:
			steps.ValidateStory = true
		case StepCreateBranch:
			steps.CreateBranch = true
		case StepMergeScenarios:
			steps.MergeScenarios = true
		case StepGenerateTests:
			steps.GenerateTests = true
		case StepImplementFeature:
			steps.ImplementFeature = true
		default:
			return nil, fmt.Errorf("%w: %s (valid steps: %s, %s, %s, %s, %s, %s)",
				ErrInvalidStep,
				step,
				StepValidateStory,
				StepCreateBranch,
				StepMergeScenarios,
				StepGenerateTests,
				StepImplementFeature,
				StepAll,
			)
		}
	}

	return steps, nil
}

// String returns a string representation of the execution steps.
func (e *ExecutionSteps) String() string {
	var enabled []string

	if e.ValidateStory {
		enabled = append(enabled, StepValidateStory)
	}

	if e.CreateBranch {
		enabled = append(enabled, StepCreateBranch)
	}

	if e.MergeScenarios {
		enabled = append(enabled, StepMergeScenarios)
	}

	if e.GenerateTests {
		enabled = append(enabled, StepGenerateTests)
	}

	if e.ImplementFeature {
		enabled = append(enabled, StepImplementFeature)
	}

	if len(enabled) == 0 {
		return "none"
	}

	return strings.Join(enabled, ", ")
}
