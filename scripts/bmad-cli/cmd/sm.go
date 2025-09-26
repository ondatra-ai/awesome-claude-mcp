package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/app"
	"github.com/spf13/cobra"
)

func NewSMCommand(container *app.Container) *cobra.Command {
	smCmd := &cobra.Command{
		Use:   "sm",
		Short: "Story management",
	}

	usCreateCmd := &cobra.Command{
		Use:   "us-create [story-number]",
		Short: "Create user story",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			err := container.USCreateCmd.Execute(ctx, args[0])

			stop()

			return err
		},
	}

	smCmd.AddCommand(usCreateCmd)
	return smCmd
}
