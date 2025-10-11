package shared

import (
	"encoding/json"
)

// UserMessage represents a message from the user.
type UserMessage struct {
	MessageType string      `json:"type"`
	Content     interface{} `json:"content"` // string or []ContentBlock
}

// Type returns the message type for UserMessage.
func (m *UserMessage) Type() string {
	return MessageTypeUser
}

// MarshalJSON implements custom JSON marshaling for UserMessage.
func (m *UserMessage) MarshalJSON() ([]byte, error) {
	type userMessage UserMessage

	temp := struct {
		Type string `json:"type"`
		*userMessage
	}{
		Type:        MessageTypeUser,
		userMessage: (*userMessage)(m),
	}

	return json.Marshal(temp)
}
