package main

import (
	"context"
	"os/exec"
)

// runShell executes a command and returns combined stdout/stderr as string.
func runShell(ctx context.Context, name string, args ...string) (string, error) {
    cmd := exec.CommandContext(ctx, name, args...)
    out, err := cmd.CombinedOutput()
    // Always return captured output, even on error, to aid debugging
    return string(out), err
}
