package checklist

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// mappingNodeEntryStride is the number of yaml.Node children per map entry
// (one for the key, one for the value).
const mappingNodeEntryStride = 2

// ParseAnswerMap parses a raw answer string as a YAML mapping node.
// Returns the mapping node and true when the answer is a valid YAML map
// (possibly empty). Returns nil, false for scalars, sequences, or parse errors.
func ParseAnswerMap(raw string) (*yaml.Node, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, false
	}

	var root yaml.Node

	err := yaml.Unmarshal([]byte(trimmed), &root)
	if err != nil {
		return nil, false
	}

	// Document nodes wrap the actual content in Content[0].
	node := &root
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		node = node.Content[0]
	}

	if node.Kind != yaml.MappingNode {
		return nil, false
	}

	return node, true
}

// AnswerMapEntryCount returns the number of top-level keys in a YAML map node.
// Returns 0 for nil or non-mapping nodes.
func AnswerMapEntryCount(node *yaml.Node) int {
	if node == nil || node.Kind != yaml.MappingNode {
		return 0
	}

	return len(node.Content) / mappingNodeEntryStride
}
