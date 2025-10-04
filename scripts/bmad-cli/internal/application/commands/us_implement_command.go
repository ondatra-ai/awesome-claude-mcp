package commands

import (
	"context"
	"fmt"
	"log/slog"

	"bmad-cli/internal/infrastructure/git"
	"bmad-cli/internal/infrastructure/story"
)

type USImplementCommand struct {
	branchManager *git.BranchManager
	storyLoader   *story.StoryLoader
}

func NewUSImplementCommand(branchManager *git.BranchManager, storyLoader *story.StoryLoader) *USImplementCommand {
	return &USImplementCommand{
		branchManager: branchManager,
		storyLoader:   storyLoader,
	}
}

func (c *USImplementCommand) Execute(ctx context.Context, storyNumber string, force bool) error {
	slog.Info("Starting user story implementation", "story_number", storyNumber, "force", force)

	// Get story slug from file
	storySlug, err := c.storyLoader.GetStorySlug(storyNumber)
	if err != nil {
		return fmt.Errorf("failed to get story slug: %w", err)
	}

	// Ensure correct branch is checked out
	if err := c.branchManager.EnsureBranch(ctx, storyNumber, storySlug, force); err != nil {
		return fmt.Errorf("failed to ensure branch: %w", err)
	}

	slog.Info("Branch setup completed successfully")
	fmt.Println("Implementation not yet available - placeholder command")
	return nil
}
