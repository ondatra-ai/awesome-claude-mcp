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

			if err != nil {
				return fmt.Errorf("us create command failed: %w", err)
			}

			return nil
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

			force, _ := cmd.Flags().GetBool("force")
			err := container.USImplementCmd.Execute(ctx, args[0], force)

			stop()

			if err != nil {
				return fmt.Errorf("us implement command failed: %w", err)
			}

			return nil
		},
	}

	implementCmd.Flags().BoolP("force", "f", false, "Force recreate the story branch even if it already exists")

	usCmd.AddCommand(createCmd)
	usCmd.AddCommand(implementCmd)

	return usCmd
}
