package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// claudeClient implements CodexClient interface using Claude Code CLI
type claudeClient struct{}

func NewClaudeClient() CodexClient { return &claudeClient{} }

func (c *claudeClient) HeuristicAnalysis(ctx context.Context, ctxInput ThreadContext) (HeuristicAnalysisResult, error) {
	prompt, perr := buildHeuristicPrompt(ctxInput)
	if perr != nil {
		return HeuristicAnalysisResult{}, perr
	}
	out, err := tryClaude(ctx, prompt, PlanMode)
	if err != nil {
		return HeuristicAnalysisResult{}, err
	}

	// If Claude returns empty output, provide a default fallback response
	if strings.TrimSpace(out) == "" {
		fmt.Printf("Claude returned empty output, using fallback response\n")
		return HeuristicAnalysisResult{
			Score:           5, // Medium risk to require approval
			Summary:         "Claude analysis unavailable - manual review required",
			ProposedActions: []string{"manual-review"},
			Items: map[string]bool{
				"tools_present":            false,
				"pr_detected":              true,
				"conversations_fetched":    true,
				"auto_resolved_outdated":   false,
				"relevance_classified":     false,
			},
			Alternatives: []map[string]string{
				{"option": "manual-review", "why": "Claude analysis failed, requires human evaluation"},
			},
		}, nil
	}

	// Extract the most likely final YAML block from Claude output.
	cleaned := extractFinalYAML(out)
	if strings.TrimSpace(cleaned) == "" {
		cleaned = out
	}

	// Try to parse YAML, but provide fallbacks if parsing fails
	score, serr := parseRiskFromYAML(cleaned)
	if serr != nil {
		fmt.Printf("Failed to parse risk_score from Claude output, using fallback: %v\n", serr)
		score = 5 // Default to medium risk
	}
	if score < 1 || score > 10 {
		fmt.Printf("Invalid risk_score %d, using fallback\n", score)
		score = 5
	}

	actions, aerr := parseActionsFromYAML(cleaned)
	if aerr != nil {
		fmt.Printf("Failed to parse actions from Claude output: %v\n", aerr)
		actions = []string{"manual-review"}
	}

	summary, serr2 := parseSummaryFromYAML(cleaned)
	if serr2 != nil || strings.TrimSpace(summary) == "" {
		fmt.Printf("Failed to parse summary from Claude output: %v\n", serr2)
		summary = "Claude analysis parsing failed - see raw output above"
	}

	items, ierr := parseItemsFromYAML(cleaned)
	if ierr != nil {
		fmt.Printf("Failed to parse items from Claude output: %v\n", ierr)
		items = map[string]bool{
			"tools_present":            false,
			"pr_detected":              true,
			"conversations_fetched":    true,
			"auto_resolved_outdated":   false,
			"relevance_classified":     false,
		}
	}

	alts, alterr := parseAlternativesFromYAML(cleaned)
	if alterr != nil {
		fmt.Printf("Failed to parse alternatives from Claude output: %v\n", alterr)
		alts = []map[string]string{
			{"option": "manual-review", "why": "YAML parsing failed, requires human evaluation"},
		}
	}

	return HeuristicAnalysisResult{Score: score, Summary: summary, ProposedActions: actions, Items: items, Alternatives: alts}, nil
}

func (c *claudeClient) ImplementCode(ctx context.Context, ctxInput ThreadContext) (string, error) {
	prompt, err := buildImplementCodePrompt(ctxInput)
	fmt.Println("Claude prompt:")
	fmt.Println("--------------------------------")
	fmt.Println(prompt)
	fmt.Println("--------------------------------")
	if err != nil {
		return "", err
	}
	fmt.Println("Claude output:")
	out, err := tryClaude(ctx, prompt, ApplyMode)
	fmt.Println("--------------------------------")
	fmt.Println(out)
	fmt.Println("--------------------------------")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(firstLine(out)), nil
}

// tryClaude executes Claude Code CLI in plan or apply mode. In apply mode, it enables
// auto-approval and grants workspace write access so Claude can apply changes.
func tryClaude(ctx context.Context, prompt string, mode ExecMode) (string, error) {
	fmt.Printf("BEGIN_CLAUDE_RUN\n")
	fmt.Printf("mode: %s\n", mode)

	args := []string{"claude"}

	// Note: Claude Code CLI doesn't support --auto-approve flag
	// Apply mode will use the same command as plan mode

	// Execute Claude with prompt via stdin
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(prompt)

	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("exit: %d\n", exitCode(err))
	} else {
		fmt.Printf("exit: 0\n")
	}

	result := strings.TrimSpace(string(out))
	fmt.Printf("stdout:\n%s\n", result)
	fmt.Printf("END_CLAUDE_RUN\n")

	if err != nil {
		return "", err
	}
	return result, nil
}

// claudeJSONRequest represents a request to Claude via stdin JSON format
type claudeJSONRequest struct {
	Messages []claudeMessage `json:"messages"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// tryClaudeJSON uses Claude Code CLI with JSON input format for more control
func tryClaudeJSON(ctx context.Context, prompt string, mode ExecMode) (string, error) {
	fmt.Printf("BEGIN_CLAUDE_JSON_RUN\n")
	fmt.Printf("mode: %s\n", mode)

	// Prepare the JSON request
	req := claudeJSONRequest{
		Messages: []claudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	args := []string{"claude", "--json"}
	if mode == ApplyMode {
		args = append(args, "--auto-approve")
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(string(jsonBytes))

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("exit: %d\n", exitCode(err))
		if ee, ok := err.(*exec.ExitError); ok {
			fmt.Printf("stderr: %s\n", string(ee.Stderr))
		}
	} else {
		fmt.Printf("exit: 0\n")
	}

	result := strings.TrimSpace(string(out))
	fmt.Printf("stdout:\n%s\n", result)
	fmt.Printf("END_CLAUDE_JSON_RUN\n")

	if err != nil {
		return "", err
	}
	return result, nil
}