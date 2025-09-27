package cmd

import (
	"log"
	"os"

	"bmad-cli/internal/app"
	"github.com/spf13/cobra"
)

func Execute() {
	container, err := app.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "bmad-cli",
		Short: "BMAD CLI tool",
	}

	rootCmd.AddCommand(NewDevCommand(container))
	rootCmd.AddCommand(NewSMCommand(container))
	rootCmd.AddCommand(NewMCPCommand(container))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
