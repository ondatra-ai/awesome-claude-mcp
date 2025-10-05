package errors

import (
	"errors"
	"fmt"
)

// Category represents the type of error
type Category string

const (
	CategoryAI      Category = "ai"
	CategoryGitHub  Category = "github"
	CategoryParsing Category = "parsing"
)

// AppError represents a structured application error
type AppError struct {
	Category Category
	Code     string
	Message  string
	Cause    error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Category, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Category, e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// AI Errors
var (
	ErrEmptyOutput = errors.New("client returned empty output")
)

func ErrEmptyClientOutput(clientName string) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "EMPTY_OUTPUT",
		Message:  fmt.Sprintf("%s returned empty output", clientName),
		Cause:    ErrEmptyOutput,
	}
}

// GitHub Errors
var (
	ErrNoPRFound               = errors.New("no PR found")
	ErrUnexpectedRepoOutput    = errors.New("unexpected repo output")
	ErrNoEligibleReviewThreads = errors.New("no eligible review threads found")
	ErrUnexpectedOutput        = errors.New("unexpected output format")
	ErrNoComments              = errors.New("no comments found in thread")
)

func ErrNoPRFoundForBranch(branch string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "NO_PR_FOUND",
		Message:  fmt.Sprintf("no PR found for branch: %s", branch),
		Cause:    ErrNoPRFound,
	}
}

func ErrUnexpectedRepoOutputWithDetails(output string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "UNEXPECTED_REPO_OUTPUT",
		Message:  fmt.Sprintf("unexpected repo view output: %s", output),
		Cause:    ErrUnexpectedRepoOutput,
	}
}

func ErrNoEligibleThreads(prNumber int) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "NO_ELIGIBLE_THREADS",
		Message:  fmt.Sprintf("no eligible review threads found for PR %d", prNumber),
		Cause:    ErrNoEligibleReviewThreads,
	}
}

func ErrNoCommentsInThread(threadID string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "NO_COMMENTS",
		Message:  fmt.Sprintf("no comments found in thread %s", threadID),
		Cause:    ErrNoComments,
	}
}

// Parsing Errors
var (
	ErrRiskScoreNotFound       = errors.New("risk_score not found")
	ErrPreferredOptionNotFound = errors.New("preferred_option not found in YAML")
	ErrSummaryNotFound         = errors.New("summary not found in YAML")
	ErrItemsBlockNotFound      = errors.New("items block not found or empty")
	ErrInvalidBooleanValue     = errors.New("invalid boolean value")
	ErrMissingRequiredItem     = errors.New("missing required item")
	ErrInvalidRiskScoreValue   = errors.New("invalid risk score value")
)

func ErrItemsMustBeBoolean(key, val string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_BOOLEAN_VALUE",
		Message:  fmt.Sprintf("items[%s] must be boolean, got %q", key, val),
		Cause:    ErrInvalidBooleanValue,
	}
}

func ErrMissingItems(key string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "MISSING_REQUIRED_ITEM",
		Message:  fmt.Sprintf("missing items.%s", key),
		Cause:    ErrMissingRequiredItem,
	}
}

func ErrInvalidRiskScore(score int) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_RISK_SCORE",
		Message:  fmt.Sprintf("invalid risk_score: %d", score),
		Cause:    ErrInvalidRiskScoreValue,
	}
}

// Helper functions for error checking
func IsCategory(err error, category Category) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Category == category
	}
	return false
}

func IsAIError(err error) bool {
	return IsCategory(err, CategoryAI)
}

func IsGitHubError(err error) bool {
	return IsCategory(err, CategoryGitHub)
}

func IsParsingError(err error) bool {
	return IsCategory(err, CategoryParsing)
}
