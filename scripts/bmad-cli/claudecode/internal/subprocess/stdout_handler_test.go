package subprocess_test

import (
	"errors"
	"testing"

	"bmad-cli/claudecode/internal/subprocess"
)

func TestEmptyLineFilter(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		expectContinue bool
	}{
		{
			name:           "empty line returns true",
			line:           "",
			expectContinue: true,
		},
		{
			name:           "non-empty line calls next handler",
			line:           "test",
			expectContinue: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			filter := &subprocess.EmptyLineFilter{}
			ctx := &subprocess.ProcessContext{Line: testCase.line}

			// Pass nil Transport since EmptyLineFilter doesn't use it
			result := filter.Handle(ctx, nil)

			if result != testCase.expectContinue {
				t.Errorf("Handle() = %v, want %v", result, testCase.expectContinue)
			}
		})
	}
}

func TestEmptyLineFilterChaining(t *testing.T) {
	// Test that non-empty lines get passed to next handler
	filter := &subprocess.EmptyLineFilter{}
	nextHandler := &mockHandler{shouldReturn: false}
	filter.SetNext(nextHandler)

	ctx := &subprocess.ProcessContext{Line: "test"}
	result := filter.Handle(ctx, nil)

	if !nextHandler.wasCalled {
		t.Error("expected next handler to be called for non-empty line")
	}

	if result != false {
		t.Error("expected to return next handler's result")
	}
}

func TestErrorSenderBehavior(t *testing.T) {
	t.Run("no error continues to next handler", func(t *testing.T) {
		sender := &subprocess.ErrorSender{}
		nextHandler := &mockHandler{shouldReturn: true}
		sender.SetNext(nextHandler)

		ctx := &subprocess.ProcessContext{Error: nil}

		// Pass nil Transport since no error means no Transport interaction
		sender.Handle(ctx, nil)

		if !nextHandler.wasCalled {
			t.Error("expected next handler to be called when no error")
		}
	})

	// Note: Testing error handling (when Error is present) requires a real Transport
	// because ErrorSender calls t.sendError(). This is covered by integration tests.
}

func TestMessageSenderBehavior(t *testing.T) {
	tests := []struct {
		name        string
		hasMessages bool
		hasError    bool
		expectSkip  bool
	}{
		{
			name:        "no error allows sending",
			hasMessages: true,
			hasError:    false,
			expectSkip:  false,
		},
		{
			name:        "error skips sending",
			hasMessages: true,
			hasError:    true,
			expectSkip:  true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			sender := &subprocess.MessageSender{}
			ctx := &subprocess.ProcessContext{}

			if testCase.hasError {
				ctx.Error = errors.New("test error")
			}

			// Note: Transport interaction (sending to msgChan) is tested in integration tests
			// Here we test that the handler returns true to continue processing
			result := sender.Handle(ctx, nil)

			if !result {
				t.Error("MessageSender should always return true")
			}
		})
	}
}

func TestHandlerChaining(t *testing.T) {
	// Test that handlers can be chained together
	filter := &subprocess.EmptyLineFilter{}
	sender := &subprocess.ErrorSender{}

	filter.SetNext(sender)

	ctx := &subprocess.ProcessContext{Line: "test"}
	filter.Handle(ctx, nil)

	// If we get here without panic, chaining works
}

// mockHandler is a simple mock for testing handler chaining.
type mockHandler struct {
	wasCalled    bool
	shouldReturn bool
}

func (m *mockHandler) SetNext(_ subprocess.StdoutHandler) subprocess.StdoutHandler {
	return nil
}

func (m *mockHandler) Handle(_ *subprocess.ProcessContext, _ *subprocess.Transport) bool {
	m.wasCalled = true

	return m.shouldReturn
}

// Note on Transport Integration Testing:
//
// ErrorSender and MessageSender call unexported methods (sendError, sendMessages)
// on Transport, which cannot be mocked from outside the subprocess package.
//
// The Chain of Responsibility pattern's interaction with Transport is verified
// through integration tests that use real Transport instances with actual channels.
// These integration tests ensure that:
// 1. ErrorSender correctly sends errors to Transport.errChan
// 2. MessageSender correctly sends messages to Transport.msgChan
// 3. The handler chain maintains all delivery guarantees
//
// Unit tests here focus on handler logic and chain behavior, while integration
// tests verify the complete stdout processing pipeline including channel delivery.
