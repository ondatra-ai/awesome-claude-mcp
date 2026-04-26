package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/app/bootstrap"
	"bmad-cli/internal/app/commands"
	"github.com/spf13/cobra"
)

const defaultRequirementsFile = "docs/requirements.yaml"

// errUSImplementTakesNoArgs is returned when `us implement` is given a
// positional argument (e.g. a story id) — the command now walks every
// scenario in requirements.yaml rather than targeting one story.
var errUSImplementTakesNoArgs = errors.New(
	"us implement takes no arguments; it walks all scenarios in docs/requirements.yaml",
)

func NewUSCommand(container *bootstrap.Container) *cobra.Command {
	usCmd := &cobra.Command{
		Use:   "us",
		Short: "User story commands",
	}

	usCmd.AddCommand(newUSCreateCmd(container))
	usCmd.AddCommand(newUSRefineCmd(container))
	usCmd.AddCommand(newUSApplyCmd(container))
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

// newUSScenarioChecklistCmd builds a subcommand that walks every scenario in
// docs/requirements.yaml and runs the named checklist against each.
func newUSScenarioChecklistCmd(
	container *bootstrap.Container,
	use string,
	short string,
	long string,
	checklistName string,
	commandLabel string,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}

			return errUSImplementTakesNoArgs
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			fix, _ := cmd.Flags().GetBool("fix")

			err := container.USValidationCmd.ExecuteScenarioChecklist(
				ctx,
				defaultRequirementsFile,
				checklistName,
				fix,
				container.ScenarioParser,
				container.ScenarioEvaluator,
				container.ScenarioFixPromptGenerator,
				container.ScenarioFixApplier,
			)

			stop()

			if err != nil {
				return fmt.Errorf("%s command failed: %w", commandLabel, err)
			}

			return nil
		},
	}

	cmd.Flags().Bool("fix", false,
		"Enable interactive fix mode with checklist-based validation")

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

func newUSApplyCmd(container *bootstrap.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply [story-number]",
		Short: "Apply scenarios from a refined user story into the registry",
		Long: `Walk every acceptance criterion in docs/stories/<story-number>-*.yaml and
validate each one against the us-apply checklist. With --fix, every failed
(AC, prompt) cell drives a Claude-mediated edit on a scratch copy of
docs/requirements.yaml. The canonical registry file is replaced atomically
only when every AC passes every prompt; otherwise it is left untouched.

Stories that still use the deprecated scenarios.test_scenarios[] format are
rejected — convert them to acceptance_criteria with embedded steps first.

Example:
  bmad-cli us apply 4.1
  bmad-cli us apply 4.1 --fix`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			fix, _ := cmd.Flags().GetBool("fix")

			err := container.USValidationCmd.ExecuteStoryScenarioChecklist(
				ctx,
				args[0],
				defaultRequirementsFile,
				"us-apply",
				fix,
				container.StoryScenarioParser,
				container.ApplyEvaluator,
				container.ApplyFixPromptGenerator,
				container.ApplyFixApplier,
			)

			stop()

			if err != nil {
				return fmt.Errorf("us apply command failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().Bool("fix", false,
		"Enable interactive fix mode to merge scenarios into the registry")

	return cmd
}

func newUSImplementCmd(container *bootstrap.Container) *cobra.Command {
	return newUSScenarioChecklistCmd(
		container,
		"implement",
		"Walk every scenario in requirements.yaml and run the us-implement checklist",
		`Walk all scenarios in docs/requirements.yaml and validate them against
the us-implement checklist. With --fix, the interactive loop drives
feature-code fixes for scenarios whose checks fail. The checklist is
currently empty; the command exists as a slot for future prompts.

Example:
  bmad-cli us implement
  bmad-cli us implement --fix`,
		"us-implement",
		"us implement",
	)
}

func newUSGenerateTestsCmd(container *bootstrap.Container) *cobra.Command {
	return newUSScenarioChecklistCmd(
		container,
		"generate_tests",
		"Generate and validate tests for every scenario in requirements.yaml",
		`Walk all scenarios in docs/requirements.yaml and validate their test
files against the us-generate_tests checklist. With --fix, missing test files
are created and existing ones updated in place.

Example:
  bmad-cli us generate_tests
  bmad-cli us generate_tests --fix`,
		"us-generate_tests",
		"us generate_tests",
	)
}
