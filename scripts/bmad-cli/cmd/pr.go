package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/app"
	pkgerrors "bmad-cli/internal/pkg/errors"
	"github.com/spf13/cobra"
)

func NewPRCommand(container *app.Container) *cobra.Command {
	prCmd := &cobra.Command{
		Use:   "pr",
		Short: "Pull request commands",
	}

	triageCmd := &cobra.Command{
		Use:   "triage",
		Short: "Run PR triage",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			engineType := container.Config.GetString("engine.type")
			fmt.Fprintf(os.Stderr, "bmad-cli pr triage engine: %s\n", engineType)

			err := container.PRTriageCmd.Execute(ctx)

			stop()

			if err != nil {
				return pkgerrors.ErrTriageFailed(err)
			}

			return nil
		},
	}

	prCmd.AddCommand(triageCmd)

	return prCmd
}
