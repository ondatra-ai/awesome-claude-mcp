package cmd

import (
	"log"
	"os"

	"bmad-cli/internal/app/bootstrap"
	"github.com/spf13/cobra"
)

func Execute() {
	container, err := bootstrap.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "bmad-cli",
		Short: "BMAD CLI tool",
	}

	rootCmd.AddCommand(NewUSCommand(container))
	rootCmd.AddCommand(NewReqCommand(container))

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
