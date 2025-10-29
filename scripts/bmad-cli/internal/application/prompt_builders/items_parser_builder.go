package prompt_builders

import (
	"strings"

	"bmad-cli/internal/pkg/errors"
)

const (
	keyValueSplitLimit = 2
)

// ItemsParserBuilder provides a fluent interface for parsing YAML items.
type ItemsParserBuilder struct {
	yaml  string
	items map[string]bool
	err   error
}

// NewItemsParser creates a new items parser builder.
func NewItemsParser(yaml string) *ItemsParserBuilder {
	return &ItemsParserBuilder{
		yaml:  yaml,
		items: make(map[string]bool),
	}
}

// Extract extracts items from the YAML content.
func (b *ItemsParserBuilder) Extract() *ItemsParserBuilder {
	if b.err != nil {
		return b
	}

	b.items, b.err = b.extractItemsMap()

	return b
}

// Validate validates that all required items are present.
func (b *ItemsParserBuilder) Validate() *ItemsParserBuilder {
	if b.err != nil {
		return b
	}

	b.err = b.validateRequiredItems()

	return b
}

// Build returns the parsed items map and any error.
func (b *ItemsParserBuilder) Build() (map[string]bool, error) {
	return b.items, b.err
}

// extractItemsMap extracts the items map from YAML content.
func (b *ItemsParserBuilder) extractItemsMap() (map[string]bool, error) {
	items := map[string]bool{}
	inItems := false

	for _, raw := range strings.Split(b.yaml, "\n") {
		line := strings.TrimRight(raw, "\r")
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "items:") {
			inItems = true

			continue
		}

		if inItems {
			if b.shouldStopParsing(line, trimmedLine) {
				break
			}

			err := b.parseItemLine(trimmedLine, items)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(items) == 0 {
		return nil, errors.ErrItemsBlockNotFound
	}

	return items, nil
}

// shouldStopParsing checks if we should stop parsing the items block.
func (b *ItemsParserBuilder) shouldStopParsing(line, trimmedLine string) bool {
	isEmpty := trimmedLine == ""
	hasNoColon := !strings.Contains(trimmedLine, ":")
	isNotIndented := !strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t")

	return isEmpty || hasNoColon || isNotIndented
}

// parseItemLine parses a single item line and adds it to the items map.
func (b *ItemsParserBuilder) parseItemLine(trimmedLine string, items map[string]bool) error {
	parts := strings.SplitN(strings.TrimSpace(trimmedLine), ":", keyValueSplitLimit)
	if len(parts) != keyValueSplitLimit {
		return nil
	}

	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	val = strings.Trim(val, "\"'")

	lv := strings.ToLower(val)
	if lv != "true" && lv != "false" {
		return errors.ErrItemsMustBeBoolean(key, val)
	}

	items[key] = lv == "true"

	return nil
}

// validateRequiredItems checks that all required items are present.
func (b *ItemsParserBuilder) validateRequiredItems() error {
	required := []string{
		"tools_present", "pr_detected", "conversations_fetched",
		"auto_resolved_outdated", "relevance_classified", "human_approval_needed",
	}

	for _, k := range required {
		if _, ok := b.items[k]; !ok {
			return errors.ErrMissingItems(k)
		}
	}

	return nil
}
