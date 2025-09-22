package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/app"
	"github.com/spf13/cobra"
)

func NewDevCommand(container *app.Container) *cobra.Command {
	devCmd := &cobra.Command{
		Use:   "dev",
		Short: "Developer persona",
	}

	prTriageCmd := &cobra.Command{
		Use:   "pr-triage",
		Short: "Run PR triage",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			engineType := container.Config.GetString("engine.type")
			fmt.Fprintf(os.Stderr, "bmad-cli dev pr-triage engine: %s\n", engineType)

			err := container.PRTriageCmd.Execute(ctx)

			stop()

			if err != nil {
				return fmt.Errorf("triage: %w", err)
			}

			return nil
		},
	}

	devCmd.AddCommand(prTriageCmd)
	return devCmd
}
