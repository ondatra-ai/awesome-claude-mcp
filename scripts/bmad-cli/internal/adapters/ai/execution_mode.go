package ai

// ExecutionMode defines tool permissions for AI execution
type ExecutionMode struct {
	AllowedTools    []string
	DisallowedTools []string
}

// Predefined execution modes (similar to const block syntax)
var (
	ThinkMode = ExecutionMode{
		[]string{
			"Read(**)",
			"Write(./tmp/**)",
			"Glob(**)",
			"Grep(**)",
		},
		[]string{
			"Bash",
			"Edit(**.go)",
			"MultiEdit(**.go)",
		},
	}

	FullAccessMode = ExecutionMode{
		[]string{
			"Read(**)",
			"Write(**)",
			"Edit(**)",
			"MultiEdit(**)",
			"Glob(**)",
			"Grep(**)",
			"Bash",
			"WebFetch",
			"WebSearch",
			"Task",
		},
		[]string{}, // No disallowed tools
	}
)
