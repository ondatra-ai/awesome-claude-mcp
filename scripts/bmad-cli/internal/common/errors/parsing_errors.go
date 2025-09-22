package errors

import (
	"errors"
	"fmt"
)

var (
	ErrRiskScoreNotFound       = errors.New("risk_score not found")
	ErrPreferredOptionNotFound = errors.New("preferred_option not found in YAML")
	ErrSummaryNotFound         = errors.New("summary not found in YAML")
	ErrItemsBlockNotFound      = errors.New("items block not found or empty")
	ErrInvalidBooleanValue     = errors.New("invalid boolean value")
	ErrMissingRequiredItem     = errors.New("missing required item")
	ErrInvalidRiskScoreValue   = errors.New("invalid risk score value")
)

func ErrItemsMustBeBoolean(key, val string) error {
	return fmt.Errorf("items[%s] must be boolean, got %q: %w", key, val, ErrInvalidBooleanValue)
}

func ErrMissingItems(key string) error {
	return fmt.Errorf("missing items.%s: %w", key, ErrMissingRequiredItem)
}

func ErrInvalidRiskScore(score int) error {
	return fmt.Errorf("invalid risk_score: %d: %w", score, ErrInvalidRiskScoreValue)
}
