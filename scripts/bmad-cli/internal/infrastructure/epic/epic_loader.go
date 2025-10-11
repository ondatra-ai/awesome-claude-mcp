package epic

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/domain/models/story"
	"bmad-cli/internal/infrastructure/config"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

type EpicLoader struct {
	basePath string
}

func NewEpicLoader(cfg *config.ViperConfig) *EpicLoader {
	// Get the epic path from configuration
	basePath := cfg.GetString("epics.path")

	return &EpicLoader{
		basePath: basePath,
	}
}

func (el *EpicLoader) LoadStoryFromEpic(storyNumber string) (*story.Story, error) {
	epicNum, storyIndex, err := el.parseStoryNumber(storyNumber)
	if err != nil {
		return nil, pkgerrors.ErrParseStoryNumberFailed(storyNumber, err)
	}

	epicDoc, err := el.loadEpicFile(epicNum)
	if err != nil {
		return nil, pkgerrors.ErrLoadEpicFailed(epicNum, err)
	}

	if storyIndex < 1 || storyIndex > len(epicDoc.Stories) {
		return nil, pkgerrors.ErrStoryIndexOutOfBoundsError(storyIndex, epicNum, len(epicDoc.Stories))
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
		return 0, 0, pkgerrors.ErrInvalidStoryNumberFormatError(storyNumber)
	}

	return epicNum, storyNum, nil
}

func (el *EpicLoader) loadEpicFile(epicNum int) (*models.EpicDocument, error) {
	filename := fmt.Sprintf("epic-%02d-*.yaml", epicNum)
	pattern := filepath.Join(el.basePath, filename)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, pkgerrors.ErrSearchEpicFilesFailed(err)
	}

	if len(matches) == 0 {
		return nil, pkgerrors.ErrNoEpicFileError(pattern)
	}

	if len(matches) > 1 {
		return nil, pkgerrors.ErrMultipleEpicFilesError(pattern, matches)
	}

	epicFilePath := matches[0]

	data, err := os.ReadFile(epicFilePath)
	if err != nil {
		return nil, pkgerrors.ErrReadEpicFileFailed(epicFilePath, err)
	}

	var epicDoc models.EpicDocument
	if err := yaml.Unmarshal(data, &epicDoc); err != nil {
		return nil, pkgerrors.ErrParseEpicYAMLFailed(epicFilePath, err)
	}

	return &epicDoc, nil
}
