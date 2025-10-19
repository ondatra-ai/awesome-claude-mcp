package subprocess

import (
	"fmt"
	"os/exec"
)

// ProcessTerminator defines the interface for process termination strategies.
type ProcessTerminator interface {
	Kill(cmd *exec.Cmd, done chan error) error
	GetReason() string
}

// BaseTerminator implements the template method for process termination.
type BaseTerminator struct{}

// terminateProcess implements the template method algorithm.
func (b *BaseTerminator) terminateProcess(cmd *exec.Cmd, done chan error, reason string) error {
	killErr := cmd.Process.Kill()
	if killErr != nil && !isProcessAlreadyFinishedError(killErr) {
		return fmt.Errorf("kill process after %s: %w", reason, killErr)
	}

	<-done

	return nil
}

// TimeoutTerminator handles process termination after timeout.
type TimeoutTerminator struct {
	BaseTerminator
}

// Kill terminates the process after timeout.
func (t *TimeoutTerminator) Kill(cmd *exec.Cmd, done chan error) error {
	return t.terminateProcess(cmd, done, "timeout")
}

// GetReason returns the termination reason.
func (t *TimeoutTerminator) GetReason() string {
	return "timeout"
}

// CancellationTerminator handles process termination after context cancellation.
type CancellationTerminator struct {
	BaseTerminator
}

// Kill terminates the process after context cancellation.
func (c *CancellationTerminator) Kill(cmd *exec.Cmd, done chan error) error {
	return c.terminateProcess(cmd, done, "context cancellation")
}

// GetReason returns the termination reason.
func (c *CancellationTerminator) GetReason() string {
	return "context cancellation"
}
