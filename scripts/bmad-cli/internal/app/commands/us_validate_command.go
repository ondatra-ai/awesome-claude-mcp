package commands

import (
	"log/slog"

	"bmad-cli/internal/infrastructure/story"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// USValidateCommand validates that a story file exists and is loadable.
type USValidateCommand struct {
	storyLoader *story.StoryLoader
}

// NewUSValidateCommand creates a new USValidateCommand.
func NewUSValidateCommand(storyLoader *story.StoryLoader) *USValidateCommand {
	return &USValidateCommand{
		storyLoader: storyLoader,
	}
}

// Execute validates the story file for the given story number.
func (c *USValidateCommand) Execute(storyNumber string) error {
	slog.Info("Validating story file", "story_number", storyNumber)

	storySlug, err := c.storyLoader.GetStorySlug(storyNumber)
	if err != nil {
		return pkgerrors.ErrGetStorySlugFailed(err)
	}

	slog.Info("✓ Story validated successfully", "slug", storySlug)

	return nil
}
