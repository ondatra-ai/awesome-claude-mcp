package story

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"bmad-cli/internal/infrastructure/config"
)

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
		slog.Error("No story file found", "story_number", storyNumber)
		return "", fmt.Errorf("no story file found for story %s", storyNumber)
	}

	if len(matches) > 1 {
		slog.Error("Multiple story files found", "story_number", storyNumber, "count", len(matches))
		return "", fmt.Errorf("multiple story files found for story %s", storyNumber)
	}

	// Extract slug from filename: "3.1-mcp-server-implementation.yaml" -> "mcp-server-implementation"
	filename := filepath.Base(matches[0])
	nameWithoutExt := strings.TrimSuffix(filename, ".yaml")
	parts := strings.SplitN(nameWithoutExt, "-", 2)

	if len(parts) != 2 {
		slog.Error("Invalid story filename format", "filename", filename)
		return "", fmt.Errorf("invalid story filename format: %s", filename)
	}

	slug := parts[1]
	slog.Debug("Story slug extracted", "story_number", storyNumber, "slug", slug)

	return slug, nil
}
