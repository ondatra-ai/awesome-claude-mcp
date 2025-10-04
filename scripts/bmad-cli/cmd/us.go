package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/app"
	"github.com/spf13/cobra"
)

func NewUSCommand(container *app.Container) *cobra.Command {
	usCmd := &cobra.Command{
		Use:   "us",
		Short: "User story commands",
	}

	createCmd := &cobra.Command{
		Use:   "create [story-number]",
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

	implementCmd := &cobra.Command{
		Use:   "implement [story-number]",
		Short: "Implement user story (placeholder)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			err := container.USImplementCmd.Execute(ctx, args[0])

			stop()

			return err
		},
	}

	usCmd.AddCommand(createCmd)
	usCmd.AddCommand(implementCmd)
	return usCmd
}
