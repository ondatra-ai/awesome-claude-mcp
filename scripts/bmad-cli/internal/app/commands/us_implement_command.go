package commands

import (
	"context"

	"bmad-cli/internal/app/factories"
	"bmad-cli/internal/app/generators/implement"
	pkgerrors "bmad-cli/internal/pkg/errors"
)

// USImplementCommand is a thin wrapper that delegates to ImplementFactory.
type USImplementCommand struct {
	factory *factories.ImplementFactory
}

// NewUSImplementCommand creates a new USImplementCommand.
func NewUSImplementCommand(factory *factories.ImplementFactory) *USImplementCommand {
	return &USImplementCommand{
		factory: factory,
	}
}

// Execute runs the implementation workflow based on the specified steps.
func (c *USImplementCommand) Execute(
	ctx context.Context,
	storyNumber string,
	force bool,
	stepsStr string,
) error {
	steps, err := implement.ParseSteps(stepsStr)
	if err != nil {
		return pkgerrors.ErrInvalidSteps(err)
	}

	err = c.factory.Execute(ctx, storyNumber, steps, force)
	if err != nil {
		return pkgerrors.ErrImplementFailed(err)
	}

	return nil
}
