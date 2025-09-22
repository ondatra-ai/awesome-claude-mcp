package ai

import (
	"fmt"

	"bmad-cli/internal/infrastructure/config"
)

type ClientFactory struct {
	config config.ConfigProvider
}

func NewClientFactory(config config.ConfigProvider) *ClientFactory {
	return &ClientFactory{config: config}
}

func (f *ClientFactory) Create() (AIClient, error) {
	engine := f.config.GetString("engine.type")

	switch engine {
	case "claude":
		return NewClaudeClient()
	default:
		return nil, fmt.Errorf("unsupported engine: %s (only 'claude' is supported)", engine)
	}
}
