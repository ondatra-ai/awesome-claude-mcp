package epic

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/domain/models/epic"
	"bmad-cli/internal/domain/models/story"
)

type EpicLoader struct {
	basePath string
}

func NewEpicLoader() *EpicLoader {
	// Get the base path relative to the current working directory
	basePath := filepath.Join("..", "..", "docs", "epics", "jsons")
	return &EpicLoader{
		basePath: basePath,
	}
}

func (el *EpicLoader) LoadStoryFromEpic(storyNumber string) (*story.Story, error) {
	epicNum, storyIndex, err := el.parseStoryNumber(storyNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to parse story number %s: %w", storyNumber, err)
	}

	epicDoc, err := el.loadEpicFile(epicNum)
	if err != nil {
		return nil, fmt.Errorf("failed to load epic %d: %w", epicNum, err)
	}

	if storyIndex < 1 || storyIndex > len(epicDoc.Stories) {
		return nil, fmt.Errorf("story index %d out of bounds for epic %d (has %d stories)",
			storyIndex, epicNum, len(epicDoc.Stories))
	}

	// Stories are 1-indexed in the story number, but 0-indexed in the slice
	targetStory := epicDoc.Stories[storyIndex-1]
	return &targetStory, nil
}

func (el *EpicLoader) parseStoryNumber(storyNumber string) (int, int, error) {
	// Parse format like "3.1" into epic=3, story=1
	var epicNum, storyNum int
	n, err := fmt.Sscanf(storyNumber, "%d.%d", &epicNum, &storyNum)
	if n != 2 || err != nil {
		return 0, 0, fmt.Errorf("invalid story number format, expected X.Y but got %s", storyNumber)
	}
	return epicNum, storyNum, nil
}

func (el *EpicLoader) loadEpicFile(epicNum int) (*epic.EpicDocument, error) {
	filename := fmt.Sprintf("epic-%02d-*.yaml", epicNum)
	pattern := filepath.Join(el.basePath, filename)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search for epic files: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no epic file found matching pattern %s", pattern)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple epic files found matching pattern %s: %v", pattern, matches)
	}

	epicFilePath := matches[0]
	data, err := os.ReadFile(epicFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read epic file %s: %w", epicFilePath, err)
	}

	var epicDoc epic.EpicDocument
	if err := yaml.Unmarshal(data, &epicDoc); err != nil {
		return nil, fmt.Errorf("failed to parse epic YAML from %s: %w", epicFilePath, err)
	}

	return &epicDoc, nil
}
