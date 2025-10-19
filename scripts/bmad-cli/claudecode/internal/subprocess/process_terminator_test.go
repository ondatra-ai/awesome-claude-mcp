package subprocess_test

import (
	"testing"

	"bmad-cli/claudecode/internal/subprocess"
)

func TestTimeoutTerminator(t *testing.T) {
	terminator := &subprocess.TimeoutTerminator{}

	if terminator.GetReason() != "timeout" {
		t.Errorf("GetReason() = %v, want timeout", terminator.GetReason())
	}

	// Note: Actual Kill testing requires a real process, which is complex
	// This is covered by integration tests
}

func TestCancellationTerminator(t *testing.T) {
	terminator := &subprocess.CancellationTerminator{}

	if terminator.GetReason() != "context cancellation" {
		t.Errorf("GetReason() = %v, want context cancellation", terminator.GetReason())
	}

	// Note: Actual Kill testing requires a real process, which is complex
	// This is covered by integration tests
}

func TestTerminatorInterface(t *testing.T) {
	// Verify both terminators implement the interface
	var (
		_ subprocess.ProcessTerminator = &subprocess.TimeoutTerminator{}
		_ subprocess.ProcessTerminator = &subprocess.CancellationTerminator{}
	)
}

func TestBaseTerminatorTemplateMethod(t *testing.T) {
	// This test verifies the template method pattern structure
	// The actual termination logic is tested in integration tests
	// since it requires real process management
	tests := []struct {
		name       string
		terminator subprocess.ProcessTerminator
		reason     string
	}{
		{
			name:       "timeout terminator",
			terminator: &subprocess.TimeoutTerminator{},
			reason:     "timeout",
		},
		{
			name:       "cancellation terminator",
			terminator: &subprocess.CancellationTerminator{},
			reason:     "context cancellation",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.terminator.GetReason() != testCase.reason {
				t.Errorf("GetReason() = %v, want %v",
					testCase.terminator.GetReason(), testCase.reason)
			}

			// Note: Actual Kill() testing is skipped as it requires
			// a real running process. This is covered by integration tests.
		})
	}
}
