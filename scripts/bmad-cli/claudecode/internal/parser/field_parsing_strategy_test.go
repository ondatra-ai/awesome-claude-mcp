package parser_test

import (
	"testing"

	"bmad-cli/claudecode/internal/parser"
	"bmad-cli/claudecode/internal/shared"
)

func TestRequiredFieldsStrategy(t *testing.T) {
	strategy := &parser.RequiredFieldsStrategy{}

	tests := []struct {
		name        string
		data        map[string]any
		expectError bool
	}{
		{
			name: "all required fields present",
			data: map[string]any{
				"subtype":         "test",
				"duration_ms":     float64(100),
				"duration_api_ms": float64(50),
				"is_error":        false,
				"num_turns":       float64(1),
				"session_id":      "session123",
			},
			expectError: false,
		},
		{
			name: "missing subtype",
			data: map[string]any{
				"duration_ms":     float64(100),
				"duration_api_ms": float64(50),
				"is_error":        false,
				"num_turns":       float64(1),
				"session_id":      "session123",
			},
			expectError: true,
		},
		{
			name: "missing duration_ms",
			data: map[string]any{
				"subtype":         "test",
				"duration_api_ms": float64(50),
				"is_error":        false,
				"num_turns":       float64(1),
				"session_id":      "session123",
			},
			expectError: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := &shared.ResultMessage{}
			err := strategy.ParseFields(testCase.data, result)

			if testCase.expectError && err == nil {
				t.Error("expected error, got nil")
			}

			if !testCase.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestOptionalFieldsStrategyAllFields(t *testing.T) {
	strategy := &parser.OptionalFieldsStrategy{}
	data := map[string]any{
		"total_cost_usd": float64(0.5),
		"usage":          map[string]any{"tokens": 100},
		"result":         map[string]any{"status": "ok"},
	}

	result := &shared.ResultMessage{}

	err := strategy.ParseFields(data, result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.TotalCostUSD == nil {
		t.Error("expected TotalCostUSD to be set")
	}

	if result.Usage == nil {
		t.Error("expected Usage to be set")
	}

	if result.Result == nil {
		t.Error("expected Result to be set")
	}
}

func TestOptionalFieldsStrategyNoFields(t *testing.T) {
	strategy := &parser.OptionalFieldsStrategy{}
	data := map[string]any{}

	result := &shared.ResultMessage{}

	err := strategy.ParseFields(data, result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.TotalCostUSD != nil {
		t.Error("expected TotalCostUSD to be nil")
	}

	if result.Usage != nil {
		t.Error("expected Usage to be nil")
	}

	if result.Result != nil {
		t.Error("expected Result to be nil")
	}
}

func TestOptionalFieldsStrategyPartialFields(t *testing.T) {
	strategy := &parser.OptionalFieldsStrategy{}
	data := map[string]any{
		"total_cost_usd": float64(0.5),
	}

	result := &shared.ResultMessage{}

	err := strategy.ParseFields(data, result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.TotalCostUSD == nil {
		t.Error("expected TotalCostUSD to be set")
	}

	if result.Usage != nil {
		t.Error("expected Usage to be nil")
	}

	if result.Result != nil {
		t.Error("expected Result to be nil")
	}
}
func TestStrategyInjection(t *testing.T) {
	// Test that parser accepts custom strategies (Dependency Inversion)
	customRequired := &parser.RequiredFieldsStrategy{}
	customOptional := &parser.OptionalFieldsStrategy{}

	parserWithStrategies := parser.NewWithStrategies(customRequired, customOptional)

	// Verify parser can parse with injected strategies
	line := `{"type":"result","subtype":"test","duration_ms":100,` +
		`"duration_api_ms":50,"is_error":false,"num_turns":1,"session_id":"session123"}`

	messages, err := parserWithStrategies.ProcessLine(line)
	if err != nil {
		t.Fatalf("unexpected error with custom strategies: %v", err)
	}

	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}
}
