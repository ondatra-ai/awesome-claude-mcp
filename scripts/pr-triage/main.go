package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
)

// pr-triage CLI entrypoint
//
// Flags:
//   -mode string   triage mode: "heuristic" (default) or "apply"
//                  Note: current runner prints heuristic for first thread.
//                        "apply" mode is reserved for future use.
func main() {
    log.SetFlags(0)

    mode := flag.String("mode", "heuristic", "triage mode: heuristic|apply")
    flag.Parse()

    // Log selected mode to stderr for visibility without affecting stdout blocks
    fmt.Fprintf(os.Stderr, "pr-triage mode: %s\n", *mode)

    ctx := context.Background()
    runner := NewRunner(NewGitHubCLIClient(), NewStubCodex())
    if err := runner.Run(ctx); err != nil {
        log.Fatalf("triage: %v", err)
    }
}
