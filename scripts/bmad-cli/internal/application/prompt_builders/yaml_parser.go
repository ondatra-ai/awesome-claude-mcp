package prompt_builders

import (
	"regexp"
	"strconv"
	"strings"

	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/pkg/errors"
)

const (
	minRegexMatchGroups = 2
	keyValueSplitLimit  = 2
)

type YAMLParser struct{}

func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

func (p *YAMLParser) ParseHeuristicResult(rawOutput string) (models.HeuristicAnalysisResult, error) {
	cleaned := p.extractFinalYAML(rawOutput)
	if strings.TrimSpace(cleaned) == "" {
		cleaned = rawOutput
	}

	score, err := p.parseRiskFromYAML(cleaned)
	if err != nil {
		return models.HeuristicAnalysisResult{}, err
	}

	if score < 1 || score > 10 {
		return models.HeuristicAnalysisResult{}, errors.ErrInvalidRiskScore(score)
	}

	actions, err := p.parseActionsFromYAML(cleaned)
	if err != nil {
		return models.HeuristicAnalysisResult{}, err
	}

	summary, err := p.parseSummaryFromYAML(cleaned)
	if err != nil || strings.TrimSpace(summary) == "" {
		return models.HeuristicAnalysisResult{}, errors.ErrParseSummaryYAMLFailed(err)
	}

	items, err := p.parseItemsFromYAML(cleaned)
	if err != nil {
		return models.HeuristicAnalysisResult{}, err
	}

	alts := p.parseAlternativesFromYAML(cleaned)

	return models.HeuristicAnalysisResult{
		Score:           score,
		Summary:         summary,
		ProposedActions: actions,
		Items:           items,
		Alternatives:    alts,
	}, nil
}

func (p *YAMLParser) extractFinalYAML(input string) string {
	re := regexp.MustCompile(`(?m)^\s*risk_score\s*:\s*[0-9]+\b`)

	locs := re.FindAllStringIndex(input, -1)
	if len(locs) == 0 {
		return ""
	}

	start := locs[len(locs)-1][0]

	return input[start:]
}

func (p *YAMLParser) parseRiskFromYAML(yaml string) (int, error) {
	re := regexp.MustCompile(`(?m)^\s*risk_score\s*:\s*([0-9]+)\b`)

	m := re.FindStringSubmatch(yaml)
	if len(m) < minRegexMatchGroups {
		return 0, errors.ErrRiskScoreNotFound
	}

	score, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, errors.ErrParseRiskScoreFailed(err)
	}

	return score, nil
}

func (p *YAMLParser) parseActionsFromYAML(yaml string) ([]string, error) {
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
		return nil, errors.ErrPreferredOptionNotFound
	}

	return actions, nil
}

func (p *YAMLParser) parseSummaryFromYAML(yaml string) (string, error) {
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

	return "", errors.ErrSummaryNotFound
}

func (p *YAMLParser) parseItemsFromYAML(yaml string) (map[string]bool, error) {
	items, err := p.extractItemsMap(yaml)
	if err != nil {
		return nil, err
	}

	err = p.validateRequiredItems(items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// extractItemsMap extracts the items map from YAML content.
func (p *YAMLParser) extractItemsMap(yaml string) (map[string]bool, error) {
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
			if p.shouldStopParsing(line, trimmedLine) {
				break
			}

			err := p.parseItemLine(trimmedLine, items)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(items) == 0 {
		return nil, errors.ErrItemsBlockNotFound
	}

	return items, nil
}

// shouldStopParsing checks if we should stop parsing the items block.
func (p *YAMLParser) shouldStopParsing(line, trimmedLine string) bool {
	isEmpty := trimmedLine == ""
	hasNoColon := !strings.Contains(trimmedLine, ":")
	isNotIndented := !strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t")

	return isEmpty || hasNoColon || isNotIndented
}

// parseItemLine parses a single item line and adds it to the items map.
func (p *YAMLParser) parseItemLine(trimmedLine string, items map[string]bool) error {
	parts := strings.SplitN(strings.TrimSpace(trimmedLine), ":", keyValueSplitLimit)
	if len(parts) != keyValueSplitLimit {
		return nil
	}

	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	val = strings.Trim(val, "\"'")

	lv := strings.ToLower(val)
	if lv != "true" && lv != "false" {
		return errors.ErrItemsMustBeBoolean(key, val)
	}

	items[key] = lv == "true"

	return nil
}

// validateRequiredItems checks that all required items are present.
func (p *YAMLParser) validateRequiredItems(items map[string]bool) error {
	required := []string{
		"tools_present", "pr_detected", "conversations_fetched",
		"auto_resolved_outdated", "relevance_classified", "human_approval_needed",
	}

	for _, k := range required {
		if _, ok := items[k]; !ok {
			return errors.ErrMissingItems(k)
		}
	}

	return nil
}

func (p *YAMLParser) parseAlternativesFromYAML(yaml string) []map[string]string {
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

		if p.shouldBreakFromAlternatives(trimmedLine, line) {
			if len(current) > 0 {
				alts = append(alts, current)
			}

			break
		}

		if strings.HasPrefix(trimmedLine, "- ") {
			alts, current = p.processNewAlternativeItem(alts, current, trimmedLine)

			continue
		}

		if strings.Contains(trimmedLine, ":") {
			p.processKeyValuePair(current, trimmedLine)
		}
	}

	if len(current) > 0 {
		alts = append(alts, current)
	}

	return p.deduplicateAlternatives(alts)
}

func (p *YAMLParser) shouldBreakFromAlternatives(trimmedLine, line string) bool {
	isEmpty := trimmedLine == ""
	isNotIndented := !strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "\t")

	return isEmpty || isNotIndented
}

func (p *YAMLParser) processNewAlternativeItem(
	alts []map[string]string,
	current map[string]string,
	trimmedLine string,
) ([]map[string]string, map[string]string) {
	if len(current) > 0 {
		alts = append(alts, current)
	}

	newCurrent := map[string]string{}

	rest := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "- "))
	if strings.Contains(rest, ":") {
		p.processKeyValuePair(newCurrent, rest)
	}

	return alts, newCurrent
}

func (p *YAMLParser) processKeyValuePair(target map[string]string, line string) {
	parts := strings.SplitN(line, ":", keyValueSplitLimit)
	if len(parts) == keyValueSplitLimit {
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "\"'")
		target[key] = val
	}
}

func (p *YAMLParser) deduplicateAlternatives(alts []map[string]string) []map[string]string {
	if len(alts) == 0 {
		return []map[string]string{}
	}

	seen := map[string]bool{}

	var uniq []map[string]string

	for _, alternative := range alts {
		opt := strings.TrimSpace(alternative["option"])
		if opt == "" {
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
