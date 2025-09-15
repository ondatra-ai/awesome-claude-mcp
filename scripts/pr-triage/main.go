package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
)

// pr-triage CLI entrypoint
//
// Flags:
//   -engine string AI engine: "claude" (default) or "codex"
func main() {
    log.SetFlags(0)

    engine := flag.String("engine", "claude", "AI engine: claude|codex")
    flag.Parse()

    // Log selected engine to stderr for visibility without affecting stdout blocks
    fmt.Fprintf(os.Stderr, "pr-triage engine: %s\n", *engine)

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    var codexClient CodexClient
    switch *engine {
    case "claude":
        codexClient = NewClaudeClient()
    case "codex":
        codexClient = NewStubCodex()
    default:
        log.Fatalf("unsupported engine: %s (supported: claude, codex)", *engine)
    }

    runner := NewRunner(NewGitHubCLIClient(), codexClient)
    if err := runner.Run(ctx); err != nil {
        log.Fatalf("triage: %v", err)
    }
}
