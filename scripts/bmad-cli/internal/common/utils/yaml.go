package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// MarshalToYAML converts any data structure to YAML string
func MarshalToYAML[T any](data T) (string, error) {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	return string(yamlBytes), nil
}

// UnmarshalFromYAML converts YAML string to specified type
func UnmarshalFromYAML[T any](yamlString string) (T, error) {
	var result T
	if err := yaml.Unmarshal([]byte(yamlString), &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal from YAML: %w", err)
	}
	return result, nil
}

// MarshalWithWrapper wraps data in a map with the specified key before marshaling
func MarshalWithWrapper[T any](data T, key string) (string, error) {
	wrapper := map[string]interface{}{key: data}
	return MarshalToYAML(wrapper)
}
