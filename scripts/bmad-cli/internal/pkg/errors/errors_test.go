package errors

import (
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		want string
	}{
		{
			name: "error with cause",
			err: &AppError{
				Category: CategoryAI,
				Code:     "TEST_ERROR",
				Message:  "test message",
				Cause:    errors.New("underlying error"),
			},
			want: "[ai:TEST_ERROR] test message: underlying error",
		},
		{
			name: "error without cause",
			err: &AppError{
				Category: CategoryGitHub,
				Code:     "TEST_ERROR",
				Message:  "test message",
				Cause:    nil,
			},
			want: "[github:TEST_ERROR] test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &AppError{
		Category: CategoryAI,
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
	aiErr := &AppError{Category: CategoryAI, Code: "TEST", Message: "test"}
	githubErr := &AppError{Category: CategoryGitHub, Code: "TEST", Message: "test"}
	regularErr := errors.New("regular error")

	tests := []struct {
		name     string
		err      error
		category Category
		want     bool
	}{
		{
			name:     "ai error matches ai category",
			err:      aiErr,
			category: CategoryAI,
			want:     true,
		},
		{
			name:     "ai error doesn't match github category",
			err:      aiErr,
			category: CategoryGitHub,
			want:     false,
		},
		{
			name:     "github error matches github category",
			err:      githubErr,
			category: CategoryGitHub,
			want:     true,
		},
		{
			name:     "regular error doesn't match any category",
			err:      regularErr,
			category: CategoryAI,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCategory(tt.err, tt.category); got != tt.want {
				t.Errorf("IsCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAIError(t *testing.T) {
	aiErr := &AppError{Category: CategoryAI, Code: "TEST", Message: "test"}
	githubErr := &AppError{Category: CategoryGitHub, Code: "TEST", Message: "test"}

	if !IsAIError(aiErr) {
		t.Errorf("IsAIError() should return true for AI error")
	}

	if IsAIError(githubErr) {
		t.Errorf("IsAIError() should return false for non-AI error")
	}
}

func TestIsGitHubError(t *testing.T) {
	aiErr := &AppError{Category: CategoryAI, Code: "TEST", Message: "test"}
	githubErr := &AppError{Category: CategoryGitHub, Code: "TEST", Message: "test"}

	if !IsGitHubError(githubErr) {
		t.Errorf("IsGitHubError() should return true for GitHub error")
	}

	if IsGitHubError(aiErr) {
		t.Errorf("IsGitHubError() should return false for non-GitHub error")
	}
}

func TestIsParsingError(t *testing.T) {
	parsingErr := &AppError{Category: CategoryParsing, Code: "TEST", Message: "test"}
	aiErr := &AppError{Category: CategoryAI, Code: "TEST", Message: "test"}

	if !IsParsingError(parsingErr) {
		t.Errorf("IsParsingError() should return true for parsing error")
	}

	if IsParsingError(aiErr) {
		t.Errorf("IsParsingError() should return false for non-parsing error")
	}
}

func TestErrorConstructors(t *testing.T) {
	// Test that all error constructors return proper AppError types
	tests := []struct {
		name     string
		err      error
		category Category
	}{
		{
			name:     "ErrEmptyClientOutput",
			err:      ErrEmptyClientOutput("TestClient"),
			category: CategoryAI,
		},
		{
			name:     "ErrNoPRFoundForBranch",
			err:      ErrNoPRFoundForBranch("main"),
			category: CategoryGitHub,
		},
		{
			name:     "ErrInvalidRiskScore",
			err:      ErrInvalidRiskScore(15),
			category: CategoryParsing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var appErr *AppError
			if !errors.As(tt.err, &appErr) {
				t.Errorf("%s should return AppError type", tt.name)

				return
			}

			if appErr.Category != tt.category {
				t.Errorf("%s category = %v, want %v", tt.name, appErr.Category, tt.category)
			}

			if appErr.Code == "" {
				t.Errorf("%s should have non-empty code", tt.name)
			}

			if appErr.Message == "" {
				t.Errorf("%s should have non-empty message", tt.name)
			}
		})
	}
}
