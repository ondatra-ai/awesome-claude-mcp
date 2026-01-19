package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/app/bootstrap"
	"github.com/spf13/cobra"
)

func NewUSCommand(container *bootstrap.Container) *cobra.Command {
	usCmd := &cobra.Command{
		Use:   "us",
		Short: "User story commands",
	}

	usCmd.AddCommand(newUSCreateCmd(container))
	usCmd.AddCommand(newUSImplementCmd(container))
	usCmd.AddCommand(newUSChecklistCmd(container))

	return usCmd
}

func newUSCreateCmd(container *bootstrap.Container) *cobra.Command {
	return &cobra.Command{
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
}

func newUSImplementCmd(container *bootstrap.Container) *cobra.Command {
	implementCmd := &cobra.Command{
		Use:   "implement [story-number]",
		Short: "Implement user story (placeholder)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			force, _ := cmd.Flags().GetBool("force")
			steps, _ := cmd.Flags().GetString("steps")
			err := container.USImplementCmd.Execute(ctx, args[0], force, steps)

			stop()

			if err != nil {
				return fmt.Errorf("us implement command failed: %w", err)
			}

			return nil
		},
	}

	implementCmd.Flags().BoolP("force", "f", false,
		"Force recreate the story branch even if it already exists")

	stepsHelp := "Comma-separated list of steps to execute " +
		"(validate_story,create_branch,merge_scenarios,generate_tests," +
		"validate_tests,validate_scenarios,implement_feature,all)"
	implementCmd.Flags().StringP("steps", "s", "all", stepsHelp)

	return implementCmd
}

func newUSChecklistCmd(container *bootstrap.Container) *cobra.Command {
	checklistCmd := &cobra.Command{
		Use:   "checklist [story-number]",
		Short: "Validate user story against checklist",
		Long: `Validate a user story against the validation checklist using AI.

Each validation prompt from the checklist will be evaluated against the story,
and results will be displayed as a table with PASS/WARN/FAIL status.

Example:
  bmad-cli us checklist 4.1
  bmad-cli us checklist 4.1 --fix`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			fix, _ := cmd.Flags().GetBool("fix")
			err := container.USChecklistCmd.Execute(ctx, args[0], fix)

			stop()

			if err != nil {
				return fmt.Errorf("us checklist command failed: %w", err)
			}

			return nil
		},
	}

	checklistCmd.Flags().Bool("fix", false,
		"Enable interactive fix mode to resolve failed checks")

	return checklistCmd
}
