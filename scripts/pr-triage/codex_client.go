package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type CodexClient interface {
	HeuristicAnalysis(ctx context.Context, ctxInput ThreadContext) (HeuristicAnalysisResult, error)
	ImplementCode(ctx context.Context, ctxInput ThreadContext) (string, error)
}

// stubCodex now attempts Codex CLI with a stub prompt; no fallbacks.
type stubCodex struct{}

func NewStubCodex() CodexClient { return &stubCodex{} }

func (s *stubCodex) HeuristicAnalysis(ctx context.Context, ctxInput ThreadContext) (HeuristicAnalysisResult, error) {
	prompt, perr := buildHeuristicPrompt(ctxInput)
	if perr != nil {
		return HeuristicAnalysisResult{}, perr
	}
	out, err := tryCodex(ctx, prompt, PlanMode)
	if err != nil {
		return HeuristicAnalysisResult{}, err
	}

	// Extract the most likely final YAML block from Codex output.
	cleaned := extractFinalYAML(out)
	if strings.TrimSpace(cleaned) == "" {
		cleaned = out
	}

	score, serr := parseRiskFromYAML(cleaned)
	if serr != nil || score < 1 || score > 10 {
		if serr != nil {
			return HeuristicAnalysisResult{}, serr
		}
		return HeuristicAnalysisResult{}, fmt.Errorf("invalid risk_score: %d", score)
	}
	actions, aerr := parseActionsFromYAML(cleaned)
	if aerr != nil {
		return HeuristicAnalysisResult{}, aerr
	}
	summary, serr2 := parseSummaryFromYAML(cleaned)
	if serr2 != nil || strings.TrimSpace(summary) == "" {
		if serr2 != nil {
			return HeuristicAnalysisResult{}, serr2
		}
		return HeuristicAnalysisResult{}, fmt.Errorf("missing summary in YAML")
	}
	items, ierr := parseItemsFromYAML(cleaned)
	if ierr != nil {
		return HeuristicAnalysisResult{}, ierr
	}
	alts, alterr := parseAlternativesFromYAML(cleaned)
	if alterr != nil {
		return HeuristicAnalysisResult{}, alterr
	}
	return HeuristicAnalysisResult{Score: score, Summary: summary, ProposedActions: actions, Items: items, Alternatives: alts}, nil
}

func (s *stubCodex) ImplementCode(ctx context.Context, ctxInput ThreadContext) (string, error) {
	prompt, err := buildImplementCodePrompt(ctxInput)
	fmt.Println("Codex prompt:")
	fmt.Println("--------------------------------")
	fmt.Println(prompt)
	fmt.Println("--------------------------------")
	if err != nil {
		return "", err
	}
	out, err := tryCodex(ctx, prompt, ApplyMode)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(firstLine(out)), nil
}

// Codex execution modes
type ExecMode string

const (
	PlanMode  ExecMode = "plan"
	ApplyMode ExecMode = "apply"
)

// tryCodex executes Codex in plan or apply mode. In apply mode, it disables
// approvals and grants workspace write access so Codex can apply changes.
func tryCodex(ctx context.Context, prompt string, mode ExecMode) (string, error) {
	args := []string{"codex", "exec", prompt}
	if mode == ApplyMode {
		// Auto-apply changes with no interactive approvals, limited to workspace writes
		args = append(args, "--ask-for-approval", "never", "--sandbox", "workspace-write")
	}
	out, err := runShell(ctx, args[0], args[1:]...)
	if err != nil {
		return "", err
	}
	return out, nil
}

func hasPrefix(s, p string) bool {
	if len(s) < len(p) {
		return false
	}
	return s[:len(p)] == p
}

func parseFirstInt(s string) (int, error) {
	s = strings.TrimSpace(firstLine(s))
	return strconv.Atoi(s)
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

// extractFinalYAML attempts to locate the last YAML block starting with
// a top-level `risk_score:` line in Codex output that may include logs/prose.
// It returns the substring from that line to the end.
func extractFinalYAML(s string) string {
	re := regexp.MustCompile(`(?m)^\s*risk_score\s*:\s*[0-9]+\b`)
	locs := re.FindAllStringIndex(s, -1)
	if len(locs) == 0 {
		return ""
	}
	start := locs[len(locs)-1][0]
	return s[start:]
}

// buildImplementCodePrompt loads scripts/pr-triage/apply.prompt.tpl and fills
// placeholders with the given thread context for implementation.
func buildImplementCodePrompt(tc ThreadContext) (string, error) {
	tplPath := filepath.FromSlash("scripts/pr-triage/apply.prompt.tpl")
	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		return "", err
	}
	prompt := string(tplBytes)
	prompt = strings.ReplaceAll(prompt, "{{PR_NUMBER}}", fmt.Sprintf("%d", tc.PRNumber))
	loc := fmt.Sprintf("%s:%d", tc.Comment.File, tc.Comment.Line)
	prompt = strings.ReplaceAll(prompt, "{{LOCATION}}", loc)
	prompt = strings.ReplaceAll(prompt, "{{URL}}", tc.Comment.URL)
	prompt = strings.ReplaceAll(prompt, "{{CONVERSATION_TEXT}}", joinAllComments(tc.Thread))
	return prompt, nil
}

// heuristicToYAML renders the HeuristicAnalysisResult to a compact YAML block
// for inclusion in follow-up prompts.
// heuristicToYAML removed per request; apply prompt no longer includes heuristic YAML.

// buildHeuristicPrompt loads the template and checklist, fills placeholders, and returns the final prompt.
func buildHeuristicPrompt(tc ThreadContext) (string, error) {
	tplPath := filepath.FromSlash("scripts/pr-triage/heuristic.prompt.tpl")
	chkPath := filepath.FromSlash(".bmad-core/checklists/triage-heuristic-checklist.md")
	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		return "", err
	}
	chkBytes, err := os.ReadFile(chkPath)
	if err != nil {
		return "", err
	}
	prompt := string(tplBytes)
	prompt = strings.ReplaceAll(prompt, "{{PR_NUMBER}}", fmt.Sprintf("%d", tc.PRNumber))
	loc := fmt.Sprintf("%s:%d", tc.Comment.File, tc.Comment.Line)
	prompt = strings.ReplaceAll(prompt, "{{LOCATION}}", loc)
	prompt = strings.ReplaceAll(prompt, "{{URL}}", tc.Comment.URL)
	prompt = strings.ReplaceAll(prompt, "{{CONVERSATION_TEXT}}", joinAllComments(tc.Thread))
	prompt = strings.ReplaceAll(prompt, "{{CHECKLIST_MD}}", string(chkBytes))
	return prompt, nil
}

func joinAllComments(t Thread) string {
	var b strings.Builder
	for i, c := range t.Comments {
		if i > 0 {
			b.WriteString("\n---\n")
		}
		b.WriteString(c.Body)
	}
	return b.String()
}

// parseRiskFromYAML extracts risk_score: N robustly from YAML
func parseRiskFromYAML(yaml string) (int, error) {
	re := regexp.MustCompile(`(?m)^\s*risk_score\s*:\s*([0-9]+)\b`)
	m := re.FindStringSubmatch(yaml)
	if len(m) < 2 {
		return 0, fmt.Errorf("risk_score not found")
	}
	return strconv.Atoi(m[1])
}

// parseActionsFromYAML extracts a minimal proposed action list from the YAML (preferred_option or items map)
func parseActionsFromYAML(yaml string) ([]string, error) {
	var actions []string
	for _, line := range strings.Split(yaml, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "preferred_option:") {
			parts := strings.SplitN(l, ":", 2)
			if len(parts) == 2 {
				v := strings.TrimSpace(parts[1])
				v = strings.Trim(v, "\"'")
				if v != "" {
					actions = append(actions, v)
				}
			}
		}
	}
	if len(actions) == 0 {
		return nil, fmt.Errorf("preferred_option not found in YAML")
	}
	return actions, nil
}

func parseSummaryFromYAML(yaml string) (string, error) {
	for _, line := range strings.Split(yaml, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "summary:") {
			parts := strings.SplitN(l, ":", 2)
			if len(parts) == 2 {
				v := strings.TrimSpace(parts[1])
				v = strings.Trim(v, "\"'")
				return v, nil
			}
		}
	}
	return "", fmt.Errorf("summary not found in YAML")
}

func parseItemsFromYAML(yaml string) (map[string]bool, error) {
	items := map[string]bool{}
	inItems := false
	for _, raw := range strings.Split(yaml, "\n") {
		line := strings.TrimRight(raw, "\r")
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "items:") {
			inItems = true
			continue
		}
		if inItems {
			if l == "" || !strings.Contains(l, ":") || (!strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t")) {
				break
			}
			parts := strings.SplitN(strings.TrimSpace(l), ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, "\"'")
				b := val == "true"
				items[key] = b
			}
		}
	}
	// minimal sanity
	if len(items) == 0 {
		return nil, fmt.Errorf("items block not found or empty")
	}
	return items, nil
}

func parseAlternativesFromYAML(yaml string) ([]map[string]string, error) {
	var alts []map[string]string
	inAlts := false
	current := map[string]string{}
	for _, raw := range strings.Split(yaml, "\n") {
		line := strings.TrimRight(raw, "\r")
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "alternatives:") {
			inAlts = true
			continue
		}
		if inAlts {
			if strings.HasPrefix(l, "- ") {
				// New alt item. Flush previous if populated.
				if len(current) > 0 {
					alts = append(alts, current)
					current = map[string]string{}
				}
				// Handle inline form: "- option: value"
				rest := strings.TrimSpace(strings.TrimPrefix(l, "- "))
				if strings.Contains(rest, ":") {
					parts := strings.SplitN(rest, ":", 2)
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					val = strings.Trim(val, "\"'")
					current[key] = val
				}
				continue
			}
			if l == "" || (!strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t")) {
				if len(current) > 0 {
					alts = append(alts, current)
				}
				break
			}
			if strings.Contains(l, ":") {
				parts := strings.SplitN(l, ":", 2)
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, "\"'")
				current[key] = val
			}
		}
	}
	if len(current) > 0 {
		alts = append(alts, current)
	}
	// Deduplicate by option name while preserving order
	if len(alts) == 0 {
		return nil, fmt.Errorf("alternatives not found in YAML")
	}
	seen := map[string]bool{}
	var uniq []map[string]string
	for _, a := range alts {
		opt := strings.TrimSpace(a["option"])
		if opt == "" {
			// Keep entries lacking option to avoid data loss, but don't duplicate empty.
			if !seen["__empty__"] {
				uniq = append(uniq, a)
				seen["__empty__"] = true
			}
			continue
		}
		if !seen[opt] {
			uniq = append(uniq, a)
			seen[opt] = true
		}
	}
	return uniq, nil
}
