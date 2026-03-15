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

			err := container.ReqGenerateTestsCmd.Execute(ctx, requirements)

			stop()

			if err != nil {
				return fmt.Errorf("req generate_tests command failed: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("requirements", "r", defaultRequirementsFile,
		"Path to requirements.yaml file")

	return cmd
}
