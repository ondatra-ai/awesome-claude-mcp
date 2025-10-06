package shared

import "fmt"

const (
	// DefaultMaxThinkingTokens is the default maximum number of thinking tokens.
	DefaultMaxThinkingTokens = 8000
)

// Validate checks the options for valid values and constraints.
func (o *Options) Validate() error {
	// Validate MaxThinkingTokens
	if o.MaxThinkingTokens < 0 {
		return fmt.Errorf("MaxThinkingTokens must be non-negative, got %d", o.MaxThinkingTokens)
	}

	// Validate MaxTurns
	if o.MaxTurns < 0 {
		return fmt.Errorf("MaxTurns must be non-negative, got %d", o.MaxTurns)
	}

	// Validate tool conflicts (same tool in both allowed and disallowed)
	allowedSet := make(map[string]bool)
	for _, tool := range o.AllowedTools {
		allowedSet[tool] = true
	}

	for _, tool := range o.DisallowedTools {
		if allowedSet[tool] {
			return fmt.Errorf("tool '%s' cannot be in both AllowedTools and DisallowedTools", tool)
		}
	}

	return nil
}

// NewOptions creates Options with default values.
func NewOptions() *Options {
	return &Options{
		AllowedTools:      []string{},
		DisallowedTools:   []string{},
		MaxThinkingTokens: DefaultMaxThinkingTokens,
		AddDirs:           []string{},
		McpServers:        make(map[string]McpServerConfig),
		ExtraArgs:         make(map[string]*string),
	}
}
