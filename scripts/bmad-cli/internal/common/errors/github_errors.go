package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNoPRFound               = errors.New("no PR found")
	ErrUnexpectedRepoOutput    = errors.New("unexpected repo output")
	ErrNoEligibleReviewThreads = errors.New("no eligible review threads found")
	ErrUnexpectedOutput        = errors.New("unexpected output format")
)

func ErrNoPRFoundForBranch(branch string) error {
	return fmt.Errorf("no PR found for branch: %s: %w", branch, ErrNoPRFound)
}

func ErrUnexpectedRepoOutputWithDetails(output string) error {
	return fmt.Errorf("unexpected repo view output: %s: %w", output, ErrUnexpectedRepoOutput)
}

func ErrNoEligibleThreads(prNumber int) error {
	return fmt.Errorf("no eligible review threads found for PR %d: %w", prNumber, ErrNoEligibleReviewThreads)
}
