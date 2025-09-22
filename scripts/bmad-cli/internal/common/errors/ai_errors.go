package errors

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyOutput = errors.New("client returned empty output")
)

func ErrEmptyClientOutput(clientName string) error {
	return fmt.Errorf("%s returned empty output: %w", clientName, ErrEmptyOutput)
}
