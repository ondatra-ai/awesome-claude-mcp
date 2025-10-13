package utils_test

import (
	"testing"

	"bmad-cli/internal/pkg/utils"
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := utils.MarshalToYAML(testCase.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("MarshalToYAML() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if result != testCase.expected {
				t.Errorf("MarshalToYAML() = %v, want %v", result, testCase.expected)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := utils.UnmarshalFromYAML[testStruct](testCase.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("UnmarshalFromYAML() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !testCase.wantErr && result != testCase.expected {
				t.Errorf("UnmarshalFromYAML() = %v, want %v", result, testCase.expected)
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := utils.MarshalWithWrapper(testCase.input, testCase.key)
			if (err != nil) != testCase.wantErr {
				t.Errorf("MarshalWithWrapper() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if result != testCase.expected {
				t.Errorf("MarshalWithWrapper() = %v, want %v", result, testCase.expected)
			}
		})
	}
}
