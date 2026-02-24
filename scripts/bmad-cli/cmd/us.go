package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/app/bootstrap"
	"bmad-cli/internal/app/commands"
	"bmad-cli/internal/pkg/console"
	"github.com/spf13/cobra"
)

func NewUSCommand(container *bootstrap.Container) *cobra.Command {
	usCmd := &cobra.Command{
		Use:   "us",
		Short: "User story commands",
	}

	usCmd.AddCommand(newUSImplementCmd(container))
	usCmd.AddCommand(newUSCreateCmd(container))
	usCmd.AddCommand(newUSRefineCmd(container))
	usCmd.AddCommand(newUSArchitectureCmd())
	usCmd.AddCommand(newUSReadyCmd(container))

	return usCmd
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

func newUSCreateCmd(container *bootstrap.Container) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create [story-number]",
		Short: "Create and validate a user story (Stage 1: Story Creation)",
		Long: `Extract a story from its epic and validate against Stage 1 (Story Creation)
checklist prompts.

The story is saved to docs/stories/ upon passing all checks.

Example:
  bmad-cli us create 4.1
  bmad-cli us create 4.1 --fix`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			fix, _ := cmd.Flags().GetBool("fix")

			config := commands.StageConfig{
				StageID:      "story_creation",
				GateStageIDs: nil,
				LoadFromEpic: true,
				StageName:    "Story Creation",
				CommandName:  "us create",
			}

			err := container.USValidationCmd.Execute(ctx, args[0], fix, config)

			stop()

			if err != nil {
				return fmt.Errorf("us create command failed: %w", err)
			}

			return nil
		},
	}

	createCmd.Flags().Bool("fix", false,
		"Enable interactive fix mode to resolve failed checks")

	return createCmd
}

func newUSRefineCmd(container *bootstrap.Container) *cobra.Command {
	refineCmd := &cobra.Command{
		Use:   "refine [story-number]",
		Short: "Refine a user story (Stage 2: Refinement)",
		Long: `Load a story from docs/stories/ and validate against Stage 2 (Refinement)
checklist prompts. Gate-checks Stage 1 before proceeding.

Example:
  bmad-cli us refine 4.1
  bmad-cli us refine 4.1 --fix`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			fix, _ := cmd.Flags().GetBool("fix")

			config := commands.StageConfig{
				StageID:      "refinement",
				GateStageIDs: []string{"story_creation"},
				LoadFromEpic: false,
				StageName:    "Refinement",
				CommandName:  "us refine",
			}

			err := container.USValidationCmd.Execute(ctx, args[0], fix, config)

			stop()

			if err != nil {
				return fmt.Errorf("us refine command failed: %w", err)
			}

			return nil
		},
	}

	refineCmd.Flags().Bool("fix", false,
		"Enable interactive fix mode to resolve failed checks")

	return refineCmd
}

func newUSArchitectureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "architecture [story-number]",
		Short: "Architecture review (Stage 3: Architecture)",
		Long: `Architecture review stage. No automated checks are defined yet.

Example:
  bmad-cli us architecture 4.1`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			console.Println("No architecture checks defined yet.")

			return nil
		},
	}
}

func newUSReadyCmd(container *bootstrap.Container) *cobra.Command {
	readyCmd := &cobra.Command{
		Use:   "ready [story-number]",
		Short: "Ready gate validation (Stage 4: Ready Gate)",
		Long: `Load a story from docs/stories/ and validate against Stage 4 (Ready Gate)
checklist prompts. Gate-checks Stages 1 and 2 before proceeding.

Example:
  bmad-cli us ready 4.1
  bmad-cli us ready 4.1 --fix`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			fix, _ := cmd.Flags().GetBool("fix")

			config := commands.StageConfig{
				StageID:      "ready_gate",
				GateStageIDs: []string{"story_creation", "refinement"},
				LoadFromEpic: false,
				StageName:    "Ready Gate",
				CommandName:  "us ready",
			}

			err := container.USValidationCmd.Execute(ctx, args[0], fix, config)

			stop()

			if err != nil {
				return fmt.Errorf("us ready command failed: %w", err)
			}

			return nil
		},
	}

	readyCmd.Flags().Bool("fix", false,
		"Enable interactive fix mode to resolve failed checks")

	return readyCmd
}
