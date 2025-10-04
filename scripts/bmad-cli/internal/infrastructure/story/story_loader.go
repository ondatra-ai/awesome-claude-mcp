package story

import (
	"fmt"
	"log/slog"
	"os"
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

	// Validate filename has dash separator
	if !strings.Contains(nameWithoutExt, "-") {
		slog.Error("Story filename missing dash separator", "filename", filename)
		return "", fmt.Errorf("invalid story filename '%s': must have format '<story-number>-<slug>.yaml' (e.g., '3.1-my-feature.yaml')", filename)
	}

	parts := strings.SplitN(nameWithoutExt, "-", 2)
	if len(parts) != 2 || parts[1] == "" {
		slog.Error("Invalid story filename format", "filename", filename)
		return "", fmt.Errorf("invalid story filename '%s': must have format '<story-number>-<slug>.yaml' where slug is not empty", filename)
	}

	slug := parts[1]

	// Validate slug doesn't contain invalid characters
	if strings.Contains(slug, " ") {
		slog.Error("Story slug contains spaces", "filename", filename, "slug", slug)
		return "", fmt.Errorf("invalid story slug '%s': slug cannot contain spaces (use dashes instead)", slug)
	}

	slog.Debug("Story slug extracted", "story_number", storyNumber, "slug", slug, "file", storyFile)

	return slug, nil
}
