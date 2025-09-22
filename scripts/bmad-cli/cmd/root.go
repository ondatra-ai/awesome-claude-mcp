package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "bmad-cli",
	Short: "BMAD CLI tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Always use ./scripts/bmad-cli/bmad-cli.yml
	viper.SetConfigFile("./scripts/bmad-cli/bmad-cli.yml")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	// Defaults
	viper.SetDefault("engine.type", "claude")
	viper.SetDefault("templates.prompts.heuristic", "./templates/heuristic.prompt.tpl")
	viper.SetDefault("templates.prompts.apply", "./templates/apply.prompt.tpl")
}
