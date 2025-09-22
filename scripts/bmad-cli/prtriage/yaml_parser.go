package prtriage

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	minRegexMatchGroups = 2
	keyValueSplitLimit  = 2
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

// Error wrapping functions for context-specific errors.
func ErrItemsMustBeBoolean(key, val string) error {
	return fmt.Errorf("items[%s] must be boolean, got %q: %w", key, val, ErrInvalidBooleanValue)
}

func ErrMissingItems(key string) error {
	return fmt.Errorf("missing items.%s: %w", key, ErrMissingRequiredItem)
}

func ErrInvalidRiskScore(score int) error {
	return fmt.Errorf("invalid risk_score: %d: %w", score, ErrInvalidRiskScoreValue)
}

// extractFinalYAML attempts to locate the last YAML block starting with
// a top-level `risk_score:` line in AI output that may include logs/prose.
// It returns the substring from that line to the end.
func extractFinalYAML(input string) string {
	re := regexp.MustCompile(`(?m)^\s*risk_score\s*:\s*[0-9]+\b`)

	locs := re.FindAllStringIndex(input, -1)
	if len(locs) == 0 {
		return ""
	}

	start := locs[len(locs)-1][0]

	return input[start:]
}

// parseRiskFromYAML extracts risk_score: N robustly from YAML.
func parseRiskFromYAML(yaml string) (int, error) {
	re := regexp.MustCompile(`(?m)^\s*risk_score\s*:\s*([0-9]+)\b`)

	m := re.FindStringSubmatch(yaml)
	if len(m) < minRegexMatchGroups {
		return 0, ErrRiskScoreNotFound
	}

	score, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse risk score: %w", err)
	}

	return score, nil
}

// parseActionsFromYAML extracts a minimal proposed action list from the YAML (preferred_option or items map).
func parseActionsFromYAML(yaml string) ([]string, error) {
	var actions []string

	for _, line := range strings.Split(yaml, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "preferred_option:") {
			parts := strings.SplitN(l, ":", keyValueSplitLimit)
			if len(parts) == keyValueSplitLimit {
				v := strings.TrimSpace(parts[1])

				v = strings.Trim(v, "\"'")
				if v != "" {
					actions = append(actions, v)
				}
			}
		}
	}

	if len(actions) == 0 {
		return nil, ErrPreferredOptionNotFound
	}

	return actions, nil
}

// parseSummaryFromYAML extracts the summary field from YAML.
func parseSummaryFromYAML(yaml string) (string, error) {
	for _, line := range strings.Split(yaml, "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "summary:") {
			parts := strings.SplitN(l, ":", keyValueSplitLimit)
			if len(parts) == keyValueSplitLimit {
				v := strings.TrimSpace(parts[1])
				v = strings.Trim(v, "\"'")

				return v, nil
			}
		}
	}

	return "", ErrSummaryNotFound
}

// parseItemsFromYAML extracts items map from YAML with validation.
func parseItemsFromYAML(yaml string) (map[string]bool, error) {
	items := map[string]bool{}
	inItems := false

	for _, raw := range strings.Split(yaml, "\n") {
		line := strings.TrimRight(raw, "\r")

		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "items:") {
			inItems = true

			continue
		}

		if inItems {
			isEmpty := trimmedLine == ""
			hasNoColon := !strings.Contains(trimmedLine, ":")

			isNotIndented := !strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t")
			if isEmpty || hasNoColon || isNotIndented {
				break
			}

			parts := strings.SplitN(strings.TrimSpace(trimmedLine), ":", keyValueSplitLimit)
			if len(parts) == keyValueSplitLimit {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, "\"'")

				lv := strings.ToLower(val)
				if lv != "true" && lv != "false" {
					return nil, ErrItemsMustBeBoolean(key, val)
				}

				b := lv == "true"
				items[key] = b
			}
		}
	}
	// minimal sanity
	if len(items) == 0 {
		return nil, ErrItemsBlockNotFound
	}
	// Ensure all required keys are present.
	required := []string{
		"tools_present", "pr_detected", "conversations_fetched",
		"auto_resolved_outdated", "relevance_classified", "human_approval_needed",
	}
	for _, k := range required {
		if _, ok := items[k]; !ok {
			return nil, ErrMissingItems(k)
		}
	}

	return items, nil
}

// parseAlternativesFromYAML extracts alternatives list from YAML.
func parseAlternativesFromYAML(yaml string) []map[string]string {
	var alts []map[string]string

	inAlts := false
	current := map[string]string{}

	for _, raw := range strings.Split(yaml, "\n") {
		line := strings.TrimRight(raw, "\r")
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "alternatives:") {
			inAlts = true

			continue
		}

		if !inAlts {
			continue
		}

		// Process alternatives section
		if shouldBreakFromAlternatives(trimmedLine, line) {
			if len(current) > 0 {
				alts = append(alts, current)
			}

			break
		}

		if strings.HasPrefix(trimmedLine, "- ") {
			alts, current = processNewAlternativeItem(alts, current, trimmedLine)

			continue
		}

		if strings.Contains(trimmedLine, ":") {
			processKeyValuePair(current, trimmedLine)
		}
	}

	if len(current) > 0 {
		alts = append(alts, current)
	}

	return deduplicateAlternatives(alts)
}

// shouldBreakFromAlternatives determines if we should stop processing alternatives.
func shouldBreakFromAlternatives(trimmedLine, line string) bool {
	isEmpty := trimmedLine == ""
	isNotIndented := !strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t")

	return isEmpty || isNotIndented
}

// processNewAlternativeItem handles a new alternative item starting with "- ".
func processNewAlternativeItem(
	alts []map[string]string,
	current map[string]string,
	trimmedLine string,
) ([]map[string]string, map[string]string) {
	// Flush previous item if populated
	if len(current) > 0 {
		alts = append(alts, current)
	}

	newCurrent := map[string]string{}

	// Handle inline form: "- option: value"
	rest := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- "))
	if strings.Contains(rest, ":") {
		processKeyValuePair(newCurrent, rest)
	}

	return alts, newCurrent
}

// processKeyValuePair extracts and stores a key-value pair.
func processKeyValuePair(target map[string]string, line string) {
	parts := strings.SplitN(line, ":", keyValueSplitLimit)
	if len(parts) == keyValueSplitLimit {
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "\"'")
		target[key] = val
	}
}

// deduplicateAlternatives removes duplicate alternatives while preserving order.
func deduplicateAlternatives(alts []map[string]string) []map[string]string {
	if len(alts) == 0 {
		return []map[string]string{}
	}

	seen := map[string]bool{}

	var uniq []map[string]string

	for _, alternative := range alts {
		opt := strings.TrimSpace(alternative["option"])
		if opt == "" {
			// Keep entries lacking option to avoid data loss, but don't duplicate empty.
			if !seen["__empty__"] {
				uniq = append(uniq, alternative)
				seen["__empty__"] = true
			}

			continue
		}

		if !seen[opt] {
			uniq = append(uniq, alternative)
			seen[opt] = true
		}
	}

	return uniq
}

// parseHeuristicResult parses the complete YAML response into HeuristicAnalysisResult.
func parseHeuristicResult(rawOutput string) (HeuristicAnalysisResult, error) {
	// Extract the most likely final YAML block from AI output.
	cleaned := extractFinalYAML(rawOutput)
	if strings.TrimSpace(cleaned) == "" {
		cleaned = rawOutput
	}

	score, serr := parseRiskFromYAML(cleaned)
	if serr != nil {
		return HeuristicAnalysisResult{}, serr
	}

	if score < 1 || score > 10 {
		return HeuristicAnalysisResult{}, ErrInvalidRiskScore(score)
	}

	actions, aerr := parseActionsFromYAML(cleaned)
	if aerr != nil {
		return HeuristicAnalysisResult{}, aerr
	}

	summary, serr2 := parseSummaryFromYAML(cleaned)
	if serr2 != nil || strings.TrimSpace(summary) == "" {
		return HeuristicAnalysisResult{}, fmt.Errorf("missing summary in YAML: %w", serr2)
	}

	items, ierr := parseItemsFromYAML(cleaned)
	if ierr != nil {
		return HeuristicAnalysisResult{}, ierr
	}

	alts := parseAlternativesFromYAML(cleaned)

	return HeuristicAnalysisResult{
		Score:           score,
		Summary:         summary,
		ProposedActions: actions,
		Items:           items,
		Alternatives:    alts,
	}, nil
}
