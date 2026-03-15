package commands

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"bmad-cli/internal/domain/models/checklist"
)

const (
	maxQuestionLen        = 40
	maxQuestionLenFixList = 80 // Longer question display for fix prompts list
	separatorLine         = "================================================================================"
	percentMultiplier     = 100
	minTruncateLen        = 3
)

// ScenarioSummary holds per-scenario summary data for the compact report table.
type ScenarioSummary struct {
	ScenarioID string
	Level      string
	Service    string
	TestFile   string
	PassCount  int
	FailCount  int
	Total      int
	Error      string
}

// TableRenderer renders checklist reports as ASCII tables.
type TableRenderer struct {
	writer io.Writer
}

// NewTableRenderer creates a new table renderer writing to stdout.
func NewTableRenderer() *TableRenderer {
	return &TableRenderer{
		writer: os.Stdout,
	}
}

// RenderReport renders a checklist report as an ASCII table.
// If showFixPrompts is true, fix prompts for failed checks are displayed.
func (r *TableRenderer) RenderReport(report *checklist.ChecklistReport, showFixPrompts bool) {
	r.renderHeader(report)
	r.renderTable(report)
	r.renderSummary(report)

	if showFixPrompts {
		r.renderFixPrompts(report)
	}
}

// RenderCompactSummary renders a single compact table summarizing all scenarios.
func (r *TableRenderer) RenderCompactSummary(summaries []ScenarioSummary) {
	_, _ = fmt.Fprintln(r.writer, separatorLine)
	_, _ = fmt.Fprintln(r.writer, "TEST VALIDATION REPORT")
	_, _ = fmt.Fprintln(r.writer, separatorLine)
	_, _ = fmt.Fprintln(r.writer)

	const columnPadding = 2

	tabWriter := tabwriter.NewWriter(r.writer, 0, 0, columnPadding, ' ', 0)

	_, _ = fmt.Fprintln(tabWriter, "SCENARIO\tLEVEL\tSERVICE\tTEST FILE\tPASS\tFAIL\tSTATUS")
	_, _ = fmt.Fprintln(tabWriter, "--------\t-----\t-------\t---------\t----\t----\t------")

	passedCount := 0

	for _, summary := range summaries {
		status := "FAIL"
		if summary.Error != "" {
			status = "ERROR"
		} else if summary.FailCount == 0 {
			status = "PASS"
			passedCount++
		}

		_, _ = fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\t%d/%d\t%d/%d\t%s\n",
			summary.ScenarioID, summary.Level, summary.Service, summary.TestFile,
			summary.PassCount, summary.Total, summary.FailCount, summary.Total, status)
	}

	_ = tabWriter.Flush()
	_, _ = fmt.Fprintln(r.writer)

	if passedCount == len(summaries) {
		_, _ = fmt.Fprintf(r.writer, "Total: %d/%d scenarios passed.\n", passedCount, len(summaries))
	} else {
		_, _ = fmt.Fprintf(r.writer, "Total: %d/%d scenarios passed. Use --fix to fix failures.\n",
			passedCount, len(summaries))
	}

	_, _ = fmt.Fprintln(r.writer, separatorLine)
}

// renderHeader renders the report header.
func (r *TableRenderer) renderHeader(report *checklist.ChecklistReport) {
	_, _ = fmt.Fprintln(r.writer, separatorLine)
	_, _ = fmt.Fprintf(r.writer, "CHECKLIST VALIDATION - %s: %s\n",
		report.SubjectID, report.SubjectTitle)
	_, _ = fmt.Fprintln(r.writer, separatorLine)
	_, _ = fmt.Fprintln(r.writer)
}

// renderTable renders the results table.
func (r *TableRenderer) renderTable(report *checklist.ChecklistReport) {
	const (
		columnPadding = 2
		answerMaxLen  = 12
	)

	tabWriter := tabwriter.NewWriter(r.writer, 0, 0, columnPadding, ' ', 0)

	// Header
	_, _ = fmt.Fprintln(tabWriter, "SECTION\tQUESTION\tEXPECTED\tACTUAL\tSTATUS")
	_, _ = fmt.Fprintln(tabWriter, "-------\t--------\t--------\t------\t------")

	// Results
	for _, result := range report.Results {
		question := truncateString(result.Question, maxQuestionLen)
		expected := truncateString(result.ExpectedAnswer, answerMaxLen)
		actual := truncateString(result.ActualAnswer, answerMaxLen)
		status := r.formatStatus(result.Status)

		_, _ = fmt.Fprintf(tabWriter, "%s\t%s\t%s\t%s\t%s\n",
			result.SectionPath,
			question,
			expected,
			actual,
			status,
		)
	}

	_ = tabWriter.Flush()
	_, _ = fmt.Fprintln(r.writer)
}

// renderSummary renders the report summary.
func (r *TableRenderer) renderSummary(report *checklist.ChecklistReport) {
	_, _ = fmt.Fprintln(r.writer, separatorLine)
	_, _ = fmt.Fprintln(r.writer, "SUMMARY")
	_, _ = fmt.Fprintln(r.writer, separatorLine)
	_, _ = fmt.Fprintf(r.writer, "Total Prompts: %d\n", report.Summary.TotalPrompts)
	_, _ = fmt.Fprintf(r.writer, "PASS: %d (%.1f%%)\n",
		report.Summary.PassCount,
		r.calculatePercentage(report.Summary.PassCount, report.Summary.TotalPrompts))
	_, _ = fmt.Fprintf(r.writer, "WARN: %d (%.1f%%)\n",
		report.Summary.WarnCount,
		r.calculatePercentage(report.Summary.WarnCount, report.Summary.TotalPrompts))
	_, _ = fmt.Fprintf(r.writer, "FAIL: %d (%.1f%%)\n",
		report.Summary.FailCount,
		r.calculatePercentage(report.Summary.FailCount, report.Summary.TotalPrompts))
	_, _ = fmt.Fprintf(r.writer, "SKIP: %d\n", report.Summary.SkipCount)
	_, _ = fmt.Fprintln(r.writer)
	_, _ = fmt.Fprintf(r.writer, "Overall: %s\n", report.GetOverallStatus())
	_, _ = fmt.Fprintln(r.writer, separatorLine)
}

// formatStatus formats the status with visual indicators.
func (r *TableRenderer) formatStatus(status checklist.Status) string {
	switch status {
	case checklist.StatusPass:
		return "PASS"
	case checklist.StatusWarn:
		return "WARN"
	case checklist.StatusFail:
		return "FAIL"
	case checklist.StatusSkip:
		return "SKIP"
	default:
		return string(status)
	}
}

// calculatePercentage calculates percentage safely.
func (r *TableRenderer) calculatePercentage(count, total int) float64 {
	if total == 0 {
		return 0
	}

	return float64(count) / float64(total) * percentMultiplier
}

// renderFixPrompts renders fix prompts for failed validations.
func (r *TableRenderer) renderFixPrompts(report *checklist.ChecklistReport) {
	// Collect results with fix prompts
	var fixes []checklist.ValidationResult

	for _, result := range report.Results {
		if result.Status == checklist.StatusFail && result.FixPrompt != "" {
			fixes = append(fixes, result)
		}
	}

	if len(fixes) == 0 {
		return
	}

	_, _ = fmt.Fprintln(r.writer)
	_, _ = fmt.Fprintln(r.writer, separatorLine)
	_, _ = fmt.Fprintln(r.writer, "FIX PROMPTS")
	_, _ = fmt.Fprintln(r.writer, separatorLine)

	for i, fix := range fixes {
		_, _ = fmt.Fprintf(r.writer, "\n### Fix %d: %s\n", i+1, fix.SectionPath)
		_, _ = fmt.Fprintf(r.writer, "Question: %s\n", truncateString(fix.Question, maxQuestionLenFixList))
		_, _ = fmt.Fprintln(r.writer)
		_, _ = fmt.Fprintln(r.writer, fix.FixPrompt)
	}

	_, _ = fmt.Fprintln(r.writer, separatorLine)
}

// truncateString truncates a string to maxLen, adding "..." if needed.
func truncateString(input string, maxLen int) string {
	// Remove newlines and extra whitespace
	cleaned := strings.ReplaceAll(input, "\n", " ")
	cleaned = strings.Join(strings.Fields(cleaned), " ")

	if len(cleaned) <= maxLen {
		return cleaned
	}

	if maxLen <= minTruncateLen {
		return cleaned[:maxLen]
	}

	return cleaned[:maxLen-minTruncateLen] + "..."
}
