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

const defaultRequirementsFile = "docs/requirements.yaml"

func NewReqCommand(container *bootstrap.Container) *cobra.Command {
	reqCmd := &cobra.Command{
		Use:   "req",
		Short: "Requirements commands",
	}

	reqCmd.AddCommand(newReqGenerateTestsCmd(container))

	return reqCmd
}

func newReqGenerateTestsCmd(container *bootstrap.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate_tests",
		Short: "Generate tests for pending scenarios in requirements.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(),
				os.Interrupt, syscall.SIGTERM)
			defer stop()

			requirements, _ := cmd.Flags().GetString("requirements")
			fix, _ := cmd.Flags().GetBool("fix")
			all, _ := cmd.Flags().GetBool("all")

			err := container.ReqValidationCmd.Execute(ctx, requirements, fix, all)

			stop()

			if err != nil {
				return fmt.Errorf("req generate_tests command failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("requirements", "r", defaultRequirementsFile,
		"Path to requirements.yaml file")
	cmd.Flags().Bool("fix", false,
		"Enable interactive fix mode with checklist-based validation")
	cmd.Flags().Bool("all", false,
		"Validate all scenarios (not just pending)")

	return cmd
}
