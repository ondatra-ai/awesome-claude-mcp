package errors_test

import (
	"errors"
	"testing"

	pkgerrors "bmad-cli/internal/pkg/errors"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *pkgerrors.AppError
		want string
	}{
		{
			name: "error with cause",
			err: &pkgerrors.AppError{
				Category: pkgerrors.CategoryAI,
				Code:     "TEST_ERROR",
				Message:  "test message",
				Cause:    errors.New("underlying error"),
			},
			want: "[ai:TEST_ERROR] test message: underlying error",
		},
		{
			name: "error without cause",
			err: &pkgerrors.AppError{
				Category: pkgerrors.CategoryGitHub,
				Code:     "TEST_ERROR",
				Message:  "test message",
				Cause:    nil,
			},
			want: "[github:TEST_ERROR] test message",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if got := testCase.err.Error(); got != testCase.want {
				t.Errorf("AppError.Error() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &pkgerrors.AppError{
		Category: pkgerrors.CategoryAI,
		Code:     "TEST_ERROR",
		Message:  "test message",
		Cause:    cause,
	}

	got := err.Unwrap()
	if !errors.Is(got, cause) {
		t.Errorf("AppError.Unwrap() = %v, want %v", got, cause)
	}
}

func TestIsCategory(t *testing.T) {
	aiErr := &pkgerrors.AppError{Category: pkgerrors.CategoryAI, Code: "TEST", Message: "test"}
	githubErr := &pkgerrors.AppError{Category: pkgerrors.CategoryGitHub, Code: "TEST", Message: "test"}
	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		category pkgerrors.Category
		want     bool
	}{
		{
			name:     "ai error matches ai category",
			err:      aiErr,
			category: pkgerrors.CategoryAI,
			want:     true,
		},
		{
			name:     "ai error doesn't match github category",
			err:      aiErr,
			category: pkgerrors.CategoryGitHub,
			want:     false,
		},
		{
			name:     "github error matches github category",
			err:      githubErr,
			category: pkgerrors.CategoryGitHub,
			want:     true,
		},
		{
			name:     "regular error doesn't match any category",
			err:      regularErr,
			category: pkgerrors.CategoryAI,
			want:     false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if got := pkgerrors.IsCategory(testCase.err, testCase.category); got != testCase.want {
				t.Errorf("IsCategory() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestIsAIError(t *testing.T) {
	aiErr := &pkgerrors.AppError{Category: pkgerrors.CategoryAI, Code: "TEST", Message: "test"}
	githubErr := &pkgerrors.AppError{Category: pkgerrors.CategoryGitHub, Code: "TEST", Message: "test"}

	if !pkgerrors.IsAIError(aiErr) {
		t.Errorf("IsAIError() should return true for AI error")
	}

	if pkgerrors.IsAIError(githubErr) {
		t.Errorf("IsAIError() should return false for non-AI error")
	}
}

func TestIsGitHubError(t *testing.T) {
	aiErr := &pkgerrors.AppError{Category: pkgerrors.CategoryAI, Code: "TEST", Message: "test"}
	githubErr := &pkgerrors.AppError{Category: pkgerrors.CategoryGitHub, Code: "TEST", Message: "test"}

	if !pkgerrors.IsGitHubError(githubErr) {
		t.Errorf("IsGitHubError() should return true for GitHub error")
	}

	if pkgerrors.IsGitHubError(aiErr) {
		t.Errorf("IsGitHubError() should return false for non-GitHub error")
	}
}

func TestIsParsingError(t *testing.T) {
	parsingErr := &pkgerrors.AppError{Category: pkgerrors.CategoryParsing, Code: "TEST", Message: "test"}
	aiErr := &pkgerrors.AppError{Category: pkgerrors.CategoryAI, Code: "TEST", Message: "test"}

	if !pkgerrors.IsParsingError(parsingErr) {
		t.Errorf("IsParsingError() should return true for parsing error")
	}

	if pkgerrors.IsParsingError(aiErr) {
		t.Errorf("IsParsingError() should return false for non-parsing error")
	}
}

func TestErrorConstructors(t *testing.T) {
	// Test that all error constructors return proper AppError types
	tests := []struct {
		name     string
		err      error
		category pkgerrors.Category
	}{
		{
			name:     "ErrEmptyClientOutput",
			err:      pkgerrors.ErrEmptyClientOutput("TestClient"),
			category: pkgerrors.CategoryAI,
		},
		{
			name:     "ErrNoPRFoundForBranch",
			err:      pkgerrors.ErrNoPRFoundForBranch("main"),
			category: pkgerrors.CategoryGitHub,
		},
		{
			name:     "ErrInvalidRiskScore",
			err:      pkgerrors.ErrInvalidRiskScore(15),
			category: pkgerrors.CategoryParsing,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var appErr *pkgerrors.AppError
			if !errors.As(testCase.err, &appErr) {
				t.Errorf("%s should return AppError type", testCase.name)

				return
			}

			if appErr.Category != testCase.category {
				t.Errorf("%s category = %v, want %v", testCase.name, appErr.Category, testCase.category)
			}

			if appErr.Code == "" {
				t.Errorf("%s should have non-empty code", testCase.name)
			}

			if appErr.Message == "" {
				t.Errorf("%s should have non-empty message", testCase.name)
			}
		})
	}
}
