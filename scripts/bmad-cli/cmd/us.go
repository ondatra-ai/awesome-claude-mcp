package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/app/bootstrap"
	"bmad-cli/internal/app/commands"
	"github.com/spf13/cobra"
)

const defaultRequirementsFile = "docs/requirements.yaml"

func NewUSCommand(container *bootstrap.Container) *cobra.Command {
	usCmd := &cobra.Command{
		Use:   "us",
		Short: "User story commands",
	}

	usCmd.AddCommand(newUSCreateCmd(container))
	usCmd.AddCommand(newUSRefineCmd(container))
	usCmd.AddCommand(newUSGenerateTestsCmd(container))
	usCmd.AddCommand(newUSImplementCmd(container))

	return usCmd
}

// newUSChecklistCmd builds a subcommand that drives a per-story checklist.
func newUSChecklistCmd(
	container *bootstrap.Container,
	use string,
	short string,
	long string,
	config commands.CommandConfig,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			fix, _ := cmd.Flags().GetBool("fix")

			err := container.USValidationCmd.Execute(ctx, args[0], fix, config)

			stop()

			if err != nil {
				return fmt.Errorf("%s command failed: %w", config.CommandName, err)
			}

			return nil
		},
	}

	cmd.Flags().Bool("fix", false,
		"Enable interactive fix mode to resolve failed checks")

	return cmd
}

func newUSCreateCmd(container *bootstrap.Container) *cobra.Command {
	return newUSChecklistCmd(
		container,
		"create [story-number]",
		"Create and validate a user story",
		`Extract a story from its epic and validate it against the us-create
checklist. The story is saved to docs/stories/ upon passing all checks.

Example:
  bmad-cli us create 4.1
  bmad-cli us create 4.1 --fix`,
		commands.CommandConfig{
			CommandName:   "us create",
			ChecklistName: "us-create",
			LoadFromEpic:  true,
		},
	)
}

func newUSRefineCmd(container *bootstrap.Container) *cobra.Command {
	return newUSChecklistCmd(
		container,
		"refine [story-number]",
		"Refine a user story",
		`Load a story from docs/stories/ and validate it against the us-refine
checklist. The story file is updated in place upon passing all checks.

Example:
  bmad-cli us refine 4.1
  bmad-cli us refine 4.1 --fix`,
		commands.CommandConfig{
			CommandName:   "us refine",
			ChecklistName: "us-refine",
			LoadFromEpic:  false,
		},
	)
}

func newUSImplementCmd(container *bootstrap.Container) *cobra.Command {
	return newUSChecklistCmd(
		container,
		"implement [story-number]",
		"Run the us-implement checklist against a user story",
		`Load a story from docs/stories/ and validate it against the us-implement
checklist. The checklist is currently empty; the command exists as a slot
for future validation prompts.

Example:
  bmad-cli us implement 4.1
  bmad-cli us implement 4.1 --fix`,
		commands.CommandConfig{
			CommandName:   "us implement",
			ChecklistName: "us-implement",
			LoadFromEpic:  false,
		},
	)
}

func newUSGenerateTestsCmd(container *bootstrap.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate_tests",
		Short: "Generate and validate tests for every scenario in requirements.yaml",
		Long: `Walk all scenarios in docs/requirements.yaml and validate their test
files against the us-generate_tests checklist. With --fix, missing test files
are created and existing ones updated in place.

Example:
  bmad-cli us generate_tests
  bmad-cli us generate_tests --fix`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			fix, _ := cmd.Flags().GetBool("fix")

			err := container.USValidationCmd.ExecuteTestValidation(
				ctx,
				defaultRequirementsFile,
				fix,
				container.ScenarioParser,
				container.TestChecklistEvaluator,
				container.TestFixPromptGenerator,
				container.TestFixApplier,
			)

			stop()

			if err != nil {
				return fmt.Errorf("us generate_tests command failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().Bool("fix", false,
		"Enable interactive fix mode with checklist-based validation")

	return cmd
}
