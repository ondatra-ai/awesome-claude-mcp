package fs

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/domain/models/story"
)

const (
	filePermissions = 0o644 // Standard file permissions for YAML files
)

// StoryVersionManager manages versioned copies of stories in tmp directory.
// Each operation creates a new version (v01, v02, v03...) for audit trail.
type StoryVersionManager struct {
	runDir   *RunDirectory
	storyID  string
	currentV int
}

// NewStoryVersionManager creates a new version manager for the specified story.
func NewStoryVersionManager(runDir *RunDirectory, storyID string) *StoryVersionManager {
	return &StoryVersionManager{
		runDir:   runDir,
		storyID:  storyID,
		currentV: 0, // Will be set to 1 after SaveInitialVersion
	}
}

// SaveInitialVersion saves the original story as v01.
func (m *StoryVersionManager) SaveInitialVersion(storyData *story.Story) error {
	m.currentV = 1

	return m.saveVersion(storyData)
}

// SaveNextVersion increments version and saves the story.
func (m *StoryVersionManager) SaveNextVersion(storyData *story.Story) (string, error) {
	m.currentV++

	err := m.saveVersion(storyData)
	if err != nil {
		return "", err
	}

	return m.GetLatestPath(), nil
}

// LoadLatest loads the most recent version of the story.
func (m *StoryVersionManager) LoadLatest() (*story.Story, error) {
	path := m.GetLatestPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read story file %s: %w", path, err)
	}

	var storyData story.Story

	err = yaml.Unmarshal(data, &storyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse story YAML %s: %w", path, err)
	}

	slog.Info("Loaded story version", "path", path, "version", m.currentV)

	return &storyData, nil
}

// GetLatestPath returns the file path for the current version.
func (m *StoryVersionManager) GetLatestPath() string {
	return m.getVersionPath(m.currentV)
}

// GetCurrentVersion returns the current version number.
func (m *StoryVersionManager) GetCurrentVersion() int {
	return m.currentV
}

// saveVersion saves the story to the versioned file path.
func (m *StoryVersionManager) saveVersion(storyData *story.Story) error {
	path := m.GetLatestPath()

	data, err := yaml.Marshal(storyData)
	if err != nil {
		return fmt.Errorf("failed to marshal story to YAML: %w", err)
	}

	err = os.WriteFile(path, data, filePermissions)
	if err != nil {
		return fmt.Errorf("failed to write story file %s: %w", path, err)
	}

	slog.Info("Saved story version", "path", path, "version", m.currentV)

	return nil
}

// getVersionPath returns the file path for a specific version number.
func (m *StoryVersionManager) getVersionPath(version int) string {
	filename := fmt.Sprintf("story-%s-v%02d.yaml", m.storyID, version)

	return filepath.Join(m.runDir.GetTmpOutPath(), filename)
}
