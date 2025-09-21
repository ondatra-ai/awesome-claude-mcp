package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// createAIClient creates an AI client based on the engine name.
func createAIClient(engine string) AIClient {
	switch engine {
	case "claude":
		return NewClaudeStrategy()
	case "codex":
		return NewCodexStrategy()
	default:
		return nil
	}
}

// pr-triage CLI entrypoint
//
// Flags:
//
//	-engine string AI engine: "claude" (default) or "codex"
func main() {
	log.SetFlags(0)

	// Configure slog for clean CLI output (no timestamps/levels for Info messages)
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
				return slog.Attr{}
			}

			return a
		},
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))

	engine := flag.String("engine", "claude", "AI engine: claude|codex")
	flag.Parse()

	// Log selected engine to stderr for visibility without affecting stdout blocks
	fmt.Fprintf(os.Stderr, "pr-triage engine: %s\n", *engine)

	// Validate engine before setting up context
	aiClient := createAIClient(*engine)
	if aiClient == nil {
		log.Printf("unsupported engine: %s (supported: claude, codex)", *engine)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// Create GitHub client
	ghClient := NewGitHubCLIClient()

	// Create GitHub operation components
	prFetcher := NewPRNumberFetcher(ghClient)
	threadsFetcher := NewThreadsFetcher(ghClient)
	resolver := NewThreadResolver(ghClient)

	// Create AI operation components
	analyzer := NewThreadAnalyzer(aiClient)
	implementer := NewThreadImplementer(aiClient)

	// Create runner with all components
	runner := NewRunner(prFetcher, threadsFetcher, resolver, analyzer, implementer)

	err := runner.Run(ctx)

	stop() // Always call stop before exit

	if err != nil {
		log.Printf("triage: %v", err)
		os.Exit(1)
	}
}
