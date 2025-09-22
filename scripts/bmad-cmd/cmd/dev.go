package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"bmad-cmd/prtriage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Developer persona",
}

var prTriageCmd = &cobra.Command{
	Use:   "pr-triage",
	Short: "Run PR triage",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Configure slog for clean CLI output (reuse existing logic)
		log.SetFlags(0)
		opts := &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
					return slog.Attr{}
				}
				return a
			},
		}
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))

		ctx, stop := signal.NotifyContext(context.Background(),
			os.Interrupt, syscall.SIGTERM)
		defer stop()

		// Get config from viper
		engineType := viper.GetString("engine.type")

		// Log selected engine to stderr for visibility
		fmt.Fprintf(os.Stderr, "bmad-cmd dev pr-triage engine: %s\n", engineType)

		// Run pr-triage using the prtriage package
		err := prtriage.RunPRTriage(ctx, engineType)

		stop() // Always call stop before exit

		if err != nil {
			return fmt.Errorf("triage: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
	devCmd.AddCommand(prTriageCmd)
}
