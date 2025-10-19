package subprocess

import (
	"errors"
	"testing"

	"bmad-cli/claudecode/internal/shared"
)

// TestErrorSenderWithRealChannels tests ErrorSender with actual Transport channels.
// This test is in the subprocess package (not subprocess_test) to access unexported fields.
func TestErrorSenderWithRealChannels(t *testing.T) {
	t.Run("error is sent to transport error channel", func(t *testing.T) {
		// Create a minimal Transport with channels
		transport := &Transport{
			errChan: make(chan error, 1),
			msgChan: make(chan shared.Message, 10),
			done:    make(chan struct{}),
		}
		defer close(transport.done)
		defer close(transport.errChan)
		defer close(transport.msgChan)

		sender := &ErrorSender{}
		testErr := errors.New("test error from handler")
		ctx := &ProcessContext{Error: testErr}

		// ErrorSender should send the error to the transport
		result := sender.Handle(ctx, transport)

		// Verify the handler returned true (sendError succeeded)
		if !result {
			t.Error("expected handler to return true when error is successfully sent")
		}

		// Verify error was actually sent to channel
		select {
		case receivedErr := <-transport.errChan:
			if receivedErr.Error() != testErr.Error() {
				t.Errorf("expected error %q, got %q", testErr, receivedErr)
			}
		default:
			t.Error("expected error in channel, got nothing")
		}
	})

	t.Run("no error continues to next handler", func(t *testing.T) {
		transport := &Transport{
			errChan: make(chan error, 1),
			msgChan: make(chan shared.Message, 10),
			done:    make(chan struct{}),
		}
		defer close(transport.done)
		defer close(transport.errChan)
		defer close(transport.msgChan)

		sender := &ErrorSender{}
		mockNext := &mockTransportHandler{shouldReturn: true}
		sender.SetNext(mockNext)

		ctx := &ProcessContext{Error: nil}

		sender.Handle(ctx, transport)

		if !mockNext.wasCalled {
			t.Error("expected next handler to be called when no error")
		}

		// Verify no error was sent
		select {
		case err := <-transport.errChan:
			t.Errorf("unexpected error in channel: %v", err)
		default:
			// Expected: no error sent
		}
	})
}

// TestMessageSenderSingleMessage tests MessageSender sends one message correctly.
func TestMessageSenderSingleMessage(t *testing.T) {
	transport := &Transport{
		errChan: make(chan error, 1),
		msgChan: make(chan shared.Message, 10),
		done:    make(chan struct{}),
	}
	defer close(transport.done)
	defer close(transport.errChan)
	defer close(transport.msgChan)

	mockMsg := &shared.SystemMessage{Subtype: "test"}
	sender := &MessageSender{}
	ctx := &ProcessContext{
		Messages: []shared.Message{mockMsg},
		Error:    nil,
	}

	result := sender.Handle(ctx, transport)

	if !result {
		t.Error("expected handler to return true after sending messages")
	}

	// Verify message was sent
	select {
	case receivedMsg := <-transport.msgChan:
		if receivedMsg == nil {
			t.Error("expected message in channel, got nil")
		}

		sysMsg, ok := receivedMsg.(*shared.SystemMessage)
		if !ok {
			t.Error("expected SystemMessage type")
		} else if sysMsg.Subtype != "test" {
			t.Errorf("expected subtype 'test', got %q", sysMsg.Subtype)
		}
	default:
		t.Error("expected message in channel, got nothing")
	}
}

// TestMessageSenderWithError tests that errors skip message sending.
func TestMessageSenderWithError(t *testing.T) {
	transport := &Transport{
		errChan: make(chan error, 1),
		msgChan: make(chan shared.Message, 10),
		done:    make(chan struct{}),
	}
	defer close(transport.done)
	defer close(transport.errChan)
	defer close(transport.msgChan)

	mockMsg := &shared.SystemMessage{Subtype: "test"}
	sender := &MessageSender{}
	ctx := &ProcessContext{
		Messages: []shared.Message{mockMsg},
		Error:    errors.New("context has error"),
	}

	result := sender.Handle(ctx, transport)

	if !result {
		t.Error("expected handler to return true even with error")
	}

	// Verify no message was sent
	select {
	case msg := <-transport.msgChan:
		t.Errorf("unexpected message in channel: %v", msg)
	default:
		// Expected: no message sent
	}
}

// TestMessageSenderMultipleMessages tests sending multiple messages.
func TestMessageSenderMultipleMessages(t *testing.T) {
	transport := &Transport{
		errChan: make(chan error, 1),
		msgChan: make(chan shared.Message, 10),
		done:    make(chan struct{}),
	}
	defer close(transport.done)
	defer close(transport.errChan)
	defer close(transport.msgChan)

	msg1 := &shared.SystemMessage{Subtype: "msg1"}
	msg2 := &shared.SystemMessage{Subtype: "msg2"}

	sender := &MessageSender{}
	ctx := &ProcessContext{
		Messages: []shared.Message{msg1, msg2},
		Error:    nil,
	}

	sender.Handle(ctx, transport)

	// Count messages received
	count := 0

	for range 2 {
		select {
		case <-transport.msgChan:
			count++
		default:
			// Channel empty
		}
	}

	if count != 2 {
		t.Errorf("expected 2 messages, got %d", count)
	}
}

// mockTransportHandler is used for testing handler chaining with Transport.
type mockTransportHandler struct {
	BaseStdoutHandler

	wasCalled    bool
	shouldReturn bool
}

func (m *mockTransportHandler) Handle(_ *ProcessContext, _ *Transport) bool {
	m.wasCalled = true

	return m.shouldReturn
}
