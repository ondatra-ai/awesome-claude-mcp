package claudecode

import (
	"bmad-cli/claudecode/internal/shared"
)

// SDKError represents the base interface for all SDK errors.
type SDKError = shared.SDKError

// BaseError provides common error functionality across the SDK.
type BaseError = shared.BaseError

// ConnectionError represents errors that occur during CLI connection.
type ConnectionError = shared.ConnectionError

// CLINotFoundError indicates that the Claude Code CLI was not found.
type CLINotFoundError = shared.CLINotFoundError

// ProcessError represents errors from the CLI process execution.
type ProcessError = shared.ProcessError

// JSONDecodeError represents JSON parsing errors from CLI responses.
type JSONDecodeError = shared.JSONDecodeError

// MessageParseError represents errors parsing message content.
type MessageParseError = shared.MessageParseError

// NewConnectionError creates a new connection error.
func NewConnectionError(message string, cause error) *ConnectionError {
	return shared.NewConnectionError(message, cause)
}

// NewCLINotFoundError creates a new CLI not found error.
func NewCLINotFoundError(path, message string) *CLINotFoundError {
	return shared.NewCLINotFoundError(path, message)
}

// NewProcessError creates a new process error.
func NewProcessError(message string, exitCode int, stderr string) *ProcessError {
	return shared.NewProcessError(message, exitCode, stderr)
}

// NewJSONDecodeError creates a new JSON decode error.
func NewJSONDecodeError(line string, position int, cause error) *JSONDecodeError {
	return shared.NewJSONDecodeError(line, position, cause)
}

// NewMessageParseError creates a new message parse error.
func NewMessageParseError(message string, data any) *MessageParseError {
	return shared.NewMessageParseError(message, data)
}
