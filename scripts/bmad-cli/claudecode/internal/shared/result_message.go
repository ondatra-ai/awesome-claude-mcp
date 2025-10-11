package shared

import (
	"encoding/json"
)

// ResultMessage represents the final result of a conversation turn.
type ResultMessage struct {
	MessageType   string          `json:"type"`
	Subtype       string          `json:"subtype"`
	DurationMs    int             `json:"duration_ms"`
	DurationAPIMs int             `json:"duration_api_ms"`
	IsError       bool            `json:"is_error"`
	NumTurns      int             `json:"num_turns"`
	SessionID     string          `json:"session_id"`
	TotalCostUSD  *float64        `json:"total_cost_usd,omitempty"`
	Usage         *map[string]any `json:"usage,omitempty"`
	Result        *map[string]any `json:"result,omitempty"`
}

// Type returns the message type for ResultMessage.
func (m *ResultMessage) Type() string {
	return MessageTypeResult
}

// MarshalJSON implements custom JSON marshaling for ResultMessage.
func (m *ResultMessage) MarshalJSON() ([]byte, error) {
	type resultMessage ResultMessage

	temp := struct {
		Type string `json:"type"`
		*resultMessage
	}{
		Type:          MessageTypeResult,
		resultMessage: (*resultMessage)(m),
	}

	return json.Marshal(temp)
}
