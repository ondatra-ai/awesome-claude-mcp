package prtriage

import (
	"context"
	"os/exec"
	"strings"
)

// runShell executes a command and returns combined stdout/stderr as string.
func runShell(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	// Always return captured output, even on error, to aid debugging
	return string(out), err
}

// runShellWithStdin executes a command with stdin input and returns combined stdout/stderr as string.
func runShellWithStdin(ctx context.Context, name, stdin string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.CombinedOutput()
	// Always return captured output, even on error, to aid debugging
	return string(out), err
}
