package prtriage

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyOutput = errors.New("client returned empty output")
)

// ErrEmptyClientOutput returns an error when AI client returns empty output.
func ErrEmptyClientOutput(clientName string) error {
	return fmt.Errorf("%s returned empty output: %w", clientName, ErrEmptyOutput)
}

// ThreadAnalyzer handles analyzing PR review threads for risk assessment.
type ThreadAnalyzer struct {
	client AIClient
}

// NewThreadAnalyzer creates a new ThreadAnalyzer.
func NewThreadAnalyzer(client AIClient) *ThreadAnalyzer {
	return &ThreadAnalyzer{client: client}
}

// Analyze performs heuristic analysis of a PR thread and returns risk assessment.
func (ta *ThreadAnalyzer) Analyze(ctx context.Context, threadContext ThreadContext) (HeuristicAnalysisResult, error) {
	// Build the heuristic prompt
	prompt, err := buildHeuristicPrompt(threadContext)
	if err != nil {
		return HeuristicAnalysisResult{}, fmt.Errorf("failed to build heuristic prompt: %w", err)
	}

	// Execute the prompt using the AI client strategy
	rawOutput, err := ta.client.ExecutePrompt(ctx, prompt, PlanMode)
	if err != nil {
		return HeuristicAnalysisResult{}, fmt.Errorf("AI client execution failed: %w", err)
	}

	// Handle empty output
	if strings.TrimSpace(rawOutput) == "" {
		return HeuristicAnalysisResult{}, ErrEmptyClientOutput(ta.client.Name())
	}

	// Parse the YAML response
	result, err := parseHeuristicResult(rawOutput)
	if err != nil {
		return HeuristicAnalysisResult{}, fmt.Errorf("failed to parse %s output: %w", ta.client.Name(), err)
	}

	return result, nil
}
