package story

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

const filenamePartsCount = 2 // Story filename has 2 parts: story-number and slug

// ValidationRule represents a validation check.
type ValidationRule struct {
	Name      string
	Validator func(nameWithoutExt, slug string) error
}

// getFilenameValidations returns all validation rules.
func getFilenameValidations() []ValidationRule {
	return []ValidationRule{
		{
			Name: "DashSeparator",
			Validator: func(nameWithoutExt, slug string) error {
				if !strings.Contains(nameWithoutExt, "-") {
					return errors.New(
						"invalid story filename: must have format " +
							"'<story-number>-<slug>.yaml' (e.g., '3.1-my-feature.yaml')",
					)
				}

				return nil
			},
		},
		{
			Name: "EmptySlug",
			Validator: func(nameWithoutExt, slug string) error {
				if slug == "" {
					return errors.New("invalid story filename: slug cannot be empty (format: '<story-number>-<slug>.yaml')")
				}

				return nil
			},
		},
		{
			Name: "NoSpaces",
			Validator: func(nameWithoutExt, slug string) error {
				if strings.Contains(slug, " ") {
					return pkgerrors.ErrInvalidStorySlugError(slug)
				}

				return nil
			},
		},
	}
}

// validateFilename executes all validation rules.
func validateFilename(filename, nameWithoutExt, slug string) error {
	for _, rule := range getFilenameValidations() {
		err := rule.Validator(nameWithoutExt, slug)
		if err != nil {
			slog.Error("Validation failed", "rule", rule.Name, "filename", filename, "error", err)

			return err
		}
	}

	return nil
}

// StoryLoader provides utilities for loading story information.
type StoryLoader struct {
	storiesDir string
}

// NewStoryLoader creates a new story loader.
func NewStoryLoader(cfg *config.ViperConfig) *StoryLoader {
	return &StoryLoader{
		storiesDir: cfg.GetString("paths.stories_dir"),
	}
}

// GetStorySlug extracts the slug from the story filename.
func (l *StoryLoader) GetStorySlug(storyNumber string) (string, error) {
	slog.Debug("Getting story slug", "story_number", storyNumber)

	pattern := filepath.Join(l.storiesDir, storyNumber+"-*.yaml")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("Failed to glob story files", "error", err, "pattern", pattern)

		return "", pkgerrors.ErrFindStoryFileFailed(err)
	}

	if len(matches) == 0 {
		slog.Error("No story file found", "story_number", storyNumber, "pattern", pattern)

		return "", pkgerrors.ErrStoryFileNotFoundError(storyNumber, l.storiesDir, storyNumber)
	}

	if len(matches) > 1 {
		slog.Error("Multiple story files found", "story_number", storyNumber, "count", len(matches), "files", matches)

		return "", pkgerrors.ErrMultipleStoryFilesError(storyNumber, matches)
	}

	// Verify file exists and is readable
	storyFile := matches[0]

	_, err = os.Stat(storyFile)
	if err != nil {
		slog.Error("Story file not accessible", "file", storyFile, "error", err)

		return "", pkgerrors.ErrStoryFileNotAccessibleError(storyFile, err)
	}

	// Extract slug from filename: "3.1-mcp-server-implementation.yaml" -> "mcp-server-implementation"
	filename := filepath.Base(storyFile)
	nameWithoutExt := strings.TrimSuffix(filename, ".yaml")

	// Parse slug from filename
	parts := strings.SplitN(nameWithoutExt, "-", filenamePartsCount)

	slug := ""
	if len(parts) == filenamePartsCount {
		slug = parts[1]
	}

	// Execute all validation rules
	err = validateFilename(filename, nameWithoutExt, slug)
	if err != nil {
		return "", err
	}

	slog.Debug("Story slug extracted", "story_number", storyNumber, "slug", slug, "file", storyFile)

	return slug, nil
}

// Load loads the complete story document from YAML file.
func (l *StoryLoader) Load(storyNumber string) (*story.StoryDocument, error) {
	slog.Debug("Loading story document", "story_number", storyNumber)

	// Find story file using same logic as GetStorySlug
	pattern := filepath.Join(l.storiesDir, storyNumber+"-*.yaml")

	matches, err := filepath.Glob(pattern)
	if err != nil {
		slog.Error("Failed to glob story files", "error", err, "pattern", pattern)

		return nil, pkgerrors.ErrFindStoryFileFailed(err)
	}

	if len(matches) == 0 {
		slog.Error("No story file found", "story_number", storyNumber, "pattern", pattern)

		return nil, pkgerrors.ErrStoryFileNotFoundError(storyNumber, l.storiesDir, storyNumber)
	}

	if len(matches) > 1 {
		slog.Error("Multiple story files found", "story_number", storyNumber, "count", len(matches), "files", matches)

		return nil, pkgerrors.ErrMultipleStoryFilesError(storyNumber, matches)
	}

	// Read story file
	storyFile := matches[0]

	data, err := os.ReadFile(storyFile)
	if err != nil {
		slog.Error("Failed to read story file", "file", storyFile, "error", err)

		return nil, pkgerrors.ErrReadStoryFileFailed(storyFile, err)
	}

	// Unmarshal YAML
	var storyDoc story.StoryDocument

	err = yaml.Unmarshal(data, &storyDoc)
	if err != nil {
		slog.Error("Failed to unmarshal story YAML", "file", storyFile, "error", err)

		return nil, pkgerrors.ErrParseStoryYAMLFailed(storyFile, err)
	}

	slog.Debug(
		"Story document loaded successfully",
		"story_number", storyNumber,
		"file", storyFile,
		"scenario_count", len(storyDoc.Scenarios.TestScenarios),
	)

	return &storyDoc, nil
}
