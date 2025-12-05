package prompt_builders_test

import (
	"testing"

	"bmad-cli/internal/app/prompt_builders"
)

func TestItemsParserBuilderSuccess(t *testing.T) {
	yaml := `
items:
  tools_present: true
  pr_detected: false
  conversations_fetched: true
  auto_resolved_outdated: false
  relevance_classified: true
  human_approval_needed: false
`

	items, err := prompt_builders.NewItemsParser(yaml).
		Extract().
		Validate().
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %d", len(items))
	}

	if !items["tools_present"] {
		t.Error("expected tools_present to be true")
	}

	if items["pr_detected"] {
		t.Error("expected pr_detected to be false")
	}
}

func TestItemsParserBuilderMissingItems(t *testing.T) {
	yaml := `
items:
  tools_present: true
`

	_, err := prompt_builders.NewItemsParser(yaml).
		Extract().
		Validate().
		Build()
	if err == nil {
		t.Error("expected error for missing required items")
	}
}

func TestItemsParserBuilderInvalidBoolean(t *testing.T) {
	yaml := `
items:
  tools_present: invalid
`

	_, err := prompt_builders.NewItemsParser(yaml).
		Extract().
		Validate().
		Build()
	if err == nil {
		t.Error("expected error for invalid boolean value")
	}
}

func TestItemsParserBuilderNoItemsBlock(t *testing.T) {
	yaml := `
summary: test
`

	_, err := prompt_builders.NewItemsParser(yaml).
		Extract().
		Validate().
		Build()
	if err == nil {
		t.Error("expected error for missing items block")
	}
}

func TestItemsParserBuilderFluentInterface(t *testing.T) {
	yaml := `
items:
  tools_present: true
  pr_detected: false
  conversations_fetched: true
  auto_resolved_outdated: false
  relevance_classified: true
  human_approval_needed: false
`

	// Test fluent interface
	builder := prompt_builders.NewItemsParser(yaml)
	builder = builder.Extract()
	builder = builder.Validate()

	items, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 6 {
		t.Errorf("expected 6 items, got %d", len(items))
	}
}

func TestItemsParserBuilderErrorPropagation(t *testing.T) {
	yaml := `invalid yaml`

	// Error from Extract should propagate through Validate and Build
	items, err := prompt_builders.NewItemsParser(yaml).
		Extract().
		Validate().
		Build()
	if err == nil {
		t.Error("expected error to propagate through chain")
	}

	if items != nil {
		t.Error("expected nil items on error")
	}
}
