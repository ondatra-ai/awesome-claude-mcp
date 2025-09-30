package claudecode

import (
	"context"

	"bmad-cli/claudecode/internal/shared"
)

type Message = shared.Message
type ContentBlock = shared.ContentBlock
type UserMessage = shared.UserMessage
type AssistantMessage = shared.AssistantMessage
type SystemMessage = shared.SystemMessage
type ResultMessage = shared.ResultMessage
type TextBlock = shared.TextBlock
type ThinkingBlock = shared.ThinkingBlock
type ToolUseBlock = shared.ToolUseBlock
type ToolResultBlock = shared.ToolResultBlock
type StreamMessage = shared.StreamMessage
type MessageIterator = shared.MessageIterator

const (
	MessageTypeUser      = shared.MessageTypeUser
	MessageTypeAssistant = shared.MessageTypeAssistant
	MessageTypeSystem    = shared.MessageTypeSystem
	MessageTypeResult    = shared.MessageTypeResult
)

const (
	ContentBlockTypeText       = shared.ContentBlockTypeText
	ContentBlockTypeThinking   = shared.ContentBlockTypeThinking
	ContentBlockTypeToolUse    = shared.ContentBlockTypeToolUse
	ContentBlockTypeToolResult = shared.ContentBlockTypeToolResult
)

type Transport interface {
	Connect(ctx context.Context) error
	SendMessage(ctx context.Context, message StreamMessage) error
	ReceiveMessages(ctx context.Context) (<-chan Message, <-chan error)
	Interrupt(ctx context.Context) error
	Close() error
}
