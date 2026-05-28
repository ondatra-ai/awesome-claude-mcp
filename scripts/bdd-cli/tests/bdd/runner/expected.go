package runner

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrJudgeSpecRequired is returned when expected.yaml has no `judge`
// field. Every fixture must declare a judge rubric — without one the
// Claude verdict step has nothing to compare against.
var ErrJudgeSpecRequired = errors.New("expected.yaml: judge is required")

// Expected mirrors the assertion strategies declared in a fixture's
// expected.yaml. Each top-level field corresponds to one strategy
// applied by checks.go / judge.go.
type Expected struct {
	// ExitCode is the exit status the CLI must return. Defaults to 0
	// when absent.
	ExitCode int `yaml:"exit_code"`

	// StdoutRegex is a list of Go regexp patterns. Each pattern is
	// asserted to match somewhere in captured stdout. Absent or empty
	// means no stdout assertions.
	StdoutRegex []string `yaml:"stdout_regex"`

	// Judge is the markdown rubric handed to the Claude judge. Required.
	Judge string `yaml:"judge"`
}

// LoadExpected reads and validates a fixture's expected.yaml. It also
// compiles each StdoutRegex entry so a bad pattern fails at load time
// rather than at assertion time.
func LoadExpected(path string) (*Expected, []*regexp.Regexp, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read expected.yaml: %w", err)
	}

	var expected Expected

	err = yaml.Unmarshal(data, &expected)
	if err != nil {
		return nil, nil, fmt.Errorf("parse expected.yaml: %w", err)
	}

	if strings.TrimSpace(expected.Judge) == "" {
		return nil, nil, ErrJudgeSpecRequired
	}

	regexes, err := compileStdoutRegexes(expected.StdoutRegex)
	if err != nil {
		return nil, nil, err
	}

	return &expected, regexes, nil
}

func compileStdoutRegexes(patterns []string) ([]*regexp.Regexp, error) {
	var regexes []*regexp.Regexp

	for _, raw := range patterns {
		pattern := strings.TrimSpace(raw)
		if pattern == "" {
			continue
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("compile regex %q: %w", pattern, err)
		}

		regexes = append(regexes, re)
	}

	return regexes, nil
}
