package shared

import (
	"encoding/json"
)

// AssistantMessage represents a message from the assistant.
type AssistantMessage struct {
	MessageType string         `json:"type"`
	Content     []ContentBlock `json:"content"`
	Model       string         `json:"model"`
}

// Type returns the message type for AssistantMessage.
func (m *AssistantMessage) Type() string {
	return MessageTypeAssistant
}

// MarshalJSON implements custom JSON marshaling for AssistantMessage.
func (m *AssistantMessage) MarshalJSON() ([]byte, error) {
	type assistantMessage AssistantMessage

	temp := struct {
		*assistantMessage

		Type string `json:"type"`
	}{
		assistantMessage: (*assistantMessage)(m),
		Type:             MessageTypeAssistant,
	}

	return json.Marshal(temp)
}
