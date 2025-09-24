package commands

import (
	"context"
	"fmt"
)

type USCreateCommand struct{}

func NewUSCreateCommand() *USCreateCommand {
	return &USCreateCommand{}
}

func (c *USCreateCommand) Execute(ctx context.Context, storyNumber string) error {
	fmt.Printf("creating user story number %s\n", storyNumber)
	return nil
}
