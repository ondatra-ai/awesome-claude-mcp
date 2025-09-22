package commands

import (
	"context"
	"fmt"

	"bmad-cli/internal/domain/services"
)

type PRTriageCommand struct {
	orchestrator *services.PRTriageOrchestrator
}

func NewPRTriageCommand(orchestrator *services.PRTriageOrchestrator) *PRTriageCommand {
	return &PRTriageCommand{orchestrator: orchestrator}
}

func (c *PRTriageCommand) Execute(ctx context.Context) error {
	err := c.orchestrator.Run(ctx)
	if err != nil {
		return fmt.Errorf("triage: %w", err)
	}
	return nil
}
