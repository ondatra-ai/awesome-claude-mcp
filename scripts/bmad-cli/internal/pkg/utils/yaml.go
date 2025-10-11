package utils

import (
	"log/slog"

	"bmad-cli/internal/pkg/errors"
	"gopkg.in/yaml.v3"
)

// MarshalToYAML converts any data structure to YAML string.
func MarshalToYAML[T any](data T) (string, error) {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		slog.Error("Failed to marshal to YAML", "error", err)

		return "", errors.ErrMarshalToYAMLFailed(err)
	}

	return string(yamlBytes), nil
}

// UnmarshalFromYAML converts YAML string to specified type.
func UnmarshalFromYAML[T any](yamlString string) (T, error) {
	var result T

	err := yaml.Unmarshal([]byte(yamlString), &result)
	if err != nil {
		slog.Error("Failed to unmarshal from YAML", "error", err)

		return result, errors.ErrUnmarshalFromYAMLFailed(err)
	}

	return result, nil
}

// MarshalWithWrapper wraps data in a map with the specified key before marshaling.
func MarshalWithWrapper[T any](data T, key string) (string, error) {
	wrapper := map[string]interface{}{key: data}

	return MarshalToYAML(wrapper)
}
