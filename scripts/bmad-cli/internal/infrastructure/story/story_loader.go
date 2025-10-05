package story

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"bmad-cli/internal/infrastructure/config"
)

// ValidationRule represents a validation check
type ValidationRule struct {
	Name      string
	Validator func(nameWithoutExt, slug string) error
}

// filenameValidations defines all validation rules
var filenameValidations = []ValidationRule{
	{
		Name: "DashSeparator",
		Validator: func(nameWithoutExt, slug string) error {
			if !strings.Contains(nameWithoutExt, "-") {
				return fmt.Errorf("invalid story filename: must have format '<story-number>-<slug>.yaml' (e.g., '3.1-my-feature.yaml')")
			}
			return nil
		},
	},
	{
		Name: "EmptySlug",
		Validator: func(nameWithoutExt, slug string) error {
			if slug == "" {
				return fmt.Errorf("invalid story filename: slug cannot be empty (format: '<story-number>-<slug>.yaml')")
			}
			return nil
		},
	},
	{
		Name: "NoSpaces",
		Validator: func(nameWithoutExt, slug string) error {
			if strings.Contains(slug, " ") {
				return fmt.Errorf("invalid story slug '%s': cannot contain spaces (use dashes instead)", slug)
			}
			return nil
		},
	},
}

// validateFilename executes all validation rules
func validateFilename(filename, nameWithoutExt, slug string) error {
	for _, rule := range filenameValidations {
		if err := rule.Validator(nameWithoutExt, slug); err != nil {
			slog.Error("Validation failed", "rule", rule.Name, "filename", filename, "error", err)
			return err
		}
	}
	return nil
}

// StoryLoader provides utilities for loading story information
type StoryLoader struct {
	storiesDir string
}

// NewStoryLoader creates a new story loader
func NewStoryLoader(cfg *config.ViperConfig) *StoryLoader {
	return &StoryLoader{
		storiesDir: cfg.GetString("paths.stories_dir"),
	}
}

// GetStorySlug extracts the slug from the story filename
func (l *StoryLoader) GetStorySlug(storyNumber string) (string, error) {
	slog.Debug("Getting story slug", "story_number", storyNumber)

	pattern := filepath.Join(l.storiesDir, fmt.Sprintf("%s-*.yaml", storyNumber))

	matches, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("Failed to glob story files", "error", err, "pattern", pattern)
		return "", fmt.Errorf("failed to find story file: %w", err)
	}

	if len(matches) == 0 {
		slog.Error("No story file found", "story_number", storyNumber, "pattern", pattern)
		return "", fmt.Errorf("no story file found for story %s in %s (expected format: %s-<slug>.yaml)", storyNumber, l.storiesDir, storyNumber)
	}

	if len(matches) > 1 {
		slog.Error("Multiple story files found", "story_number", storyNumber, "count", len(matches), "files", matches)
		return "", fmt.Errorf("multiple story files found for story %s: %v", storyNumber, matches)
	}

	// Verify file exists and is readable
	storyFile := matches[0]
	if _, err := os.Stat(storyFile); err != nil {
		slog.Error("Story file not accessible", "file", storyFile, "error", err)
		return "", fmt.Errorf("story file not accessible: %s: %w", storyFile, err)
	}

	// Extract slug from filename: "3.1-mcp-server-implementation.yaml" -> "mcp-server-implementation"
	filename := filepath.Base(storyFile)
	nameWithoutExt := strings.TrimSuffix(filename, ".yaml")

	// Parse slug from filename
	parts := strings.SplitN(nameWithoutExt, "-", 2)
	slug := ""
	if len(parts) == 2 {
		slug = parts[1]
	}

	// Execute all validation rules
	if err := validateFilename(filename, nameWithoutExt, slug); err != nil {
		return "", err
	}

	slog.Debug("Story slug extracted", "story_number", storyNumber, "slug", slug, "file", storyFile)

	return slug, nil
}
