package utils

import (
	"testing"
)

func TestMarshalToYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "simple struct",
			input:    struct{ Name string }{Name: "test"},
			expected: "name: test\n",
			wantErr:  false,
		},
		{
			name:     "map",
			input:    map[string]string{"key": "value"},
			expected: "key: value\n",
			wantErr:  false,
		},
		{
			name:     "slice",
			input:    []string{"one", "two"},
			expected: "- one\n- two\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalToYAML(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalToYAML() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if result != tt.expected {
				t.Errorf("MarshalToYAML() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnmarshalFromYAML(t *testing.T) {
	type testStruct struct {
		Name string `yaml:"name"`
	}

	tests := []struct {
		name     string
		input    string
		expected testStruct
		wantErr  bool
	}{
		{
			name:     "valid yaml",
			input:    "name: test",
			expected: testStruct{Name: "test"},
			wantErr:  false,
		},
		{
			name:     "invalid yaml",
			input:    "invalid: yaml: content:",
			expected: testStruct{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnmarshalFromYAML[testStruct](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalFromYAML() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("UnmarshalFromYAML() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMarshalWithWrapper(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		key      string
		expected string
		wantErr  bool
	}{
		{
			name:     "wrap string",
			input:    "test",
			key:      "data",
			expected: "data: test\n",
			wantErr:  false,
		},
		{
			name:     "wrap slice",
			input:    []string{"one", "two"},
			key:      "items",
			expected: "items:\n    - one\n    - two\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalWithWrapper(tt.input, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalWithWrapper() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if result != tt.expected {
				t.Errorf("MarshalWithWrapper() = %v, want %v", result, tt.expected)
			}
		})
	}
}
