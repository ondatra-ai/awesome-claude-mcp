package commands

import (
	"os"

	pkgerrors "bmad-cli/internal/pkg/errors"

	"gopkg.in/yaml.v3"
)

// ValidateTestsResult represents the result of test validation.
type ValidateTestsResult struct {
	Result string                  `yaml:"result"`
	Data   ValidateTestsResultData `yaml:"data"`
}

// ValidateTestsResultData contains the detailed validation data.
type ValidateTestsResultData struct {
	FilesScanned  int            `yaml:"files_scanned"`
	IssuesFound   int            `yaml:"issues_found"`
	IssuesFixed   int            `yaml:"issues_fixed"`
	UnfixedIssues []UnfixedIssue `yaml:"unfixed_issues"`
}

// UnfixedIssue represents an issue that could not be automatically fixed.
type UnfixedIssue struct {
	File         string `yaml:"file"`
	Line         int    `yaml:"line"`
	Description  string `yaml:"description"`
	SuggestedFix string `yaml:"suggested_fix"`
}

// parseValidateTestsResult reads and parses the validation result YAML file.
func parseValidateTestsResult(filePath string) (*ValidateTestsResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, pkgerrors.ErrReadResultFileFailed(err)
	}

	var result ValidateTestsResult

	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, pkgerrors.ErrParseResultYAMLFailed(err)
	}

	return &result, nil
}

// IsSuccess returns true if the validation result is "ok".
func (r *ValidateTestsResult) IsSuccess() bool {
	return r.Result == "ok"
}

// HasUnfixedIssues returns true if there are unfixed issues.
func (r *ValidateTestsResult) HasUnfixedIssues() bool {
	return len(r.Data.UnfixedIssues) > 0
}
