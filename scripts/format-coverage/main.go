// format-coverage parses Go coverage files (e.g., coverage/unit.out) and displays
// file-level coverage percentages in a clean table format. Unlike package-level
// reporting, this gives granular visibility into which individual files need more test coverage.
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

type FileCoverage struct {
	SmartLinesSet            map[int]bool // Smart filtered lines that actually need coverage
	SmartCoveredLinesSet     map[int]bool // Covered lines from the smart set
	SmartLineCoveragePercent float64
}

// repoPrefix caches the computed GitHub repo import prefix, e.g.,
// "github.com/owner/repo/". Used to shorten paths for display.
var repoPrefix string

func getRepoPrefix() string {
    if repoPrefix != "" {
        return repoPrefix
    }
    out, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
    if err != nil {
        return ""
    }
    url := strings.TrimSpace(string(out))
    ownerRepo := ""
    if strings.HasPrefix(url, "git@") {
        // git@github.com:owner/repo.git
        parts := strings.SplitN(url, ":", 2)
        if len(parts) == 2 {
            ownerRepo = strings.TrimSuffix(parts[1], ".git")
        }
    } else {
        // https://github.com/owner/repo.git or similar
        idx := strings.Index(url, "github.com/")
        if idx != -1 {
            ownerRepo = strings.TrimSuffix(url[idx+len("github.com/"):], ".git")
        }
    }
    if ownerRepo != "" {
        repoPrefix = "github.com/" + ownerRepo + "/"
    }
    return repoPrefix
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <coverage-file>")
	}

	coverageFile := os.Args[1]

	// Parse coverage file
	fileCoverage, err := parseCoverageFile(coverageFile)
	if err != nil {
		log.Fatalf("Error parsing coverage file: %v", err)
	}

	// Display results
	displayCoverage(fileCoverage)
}

func parseCoverageFile(filename string) (map[string]*FileCoverage, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open coverage file: %w", err)
	}

	defer func() { _ = file.Close() }()

	fileCoverage := make(map[string]*FileCoverage)
	scanner := bufio.NewScanner(file)

	// Skip the first line (mode: set)
	_ = scanner.Scan() // Skip mode line

	for scanner.Scan() {
		if err := processCoverageLine(scanner.Text(), fileCoverage); err != nil {
			continue // Skip malformed lines
		}
	}

	calculateCoveragePercentages(fileCoverage)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return fileCoverage, nil
}

func processCoverageLine(line string, fileCoverage map[string]*FileCoverage) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return fmt.Errorf("empty line")
	}

	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("invalid line format")
	}

	filePath, startLine, endLine, err := parseLineInfo(parts[0])
	if err != nil {
		return err
	}

	coverageFlag, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("invalid coverage flag: %w", err)
	}

	initializeFileCoverage(fileCoverage, filePath)
	processLineRange(fileCoverage[filePath], filePath, startLine, endLine, coverageFlag)

	return nil
}

func parseLineInfo(pathInfo string) (string, int, int, error) {
	colonIndex := strings.Index(pathInfo, ":")
	if colonIndex == -1 {
		return "", 0, 0, fmt.Errorf("invalid path format")
	}

	filePath := pathInfo[:colonIndex]
	lineRange := pathInfo[colonIndex+1:]

	startLine, endLine, err := parseLineRange(lineRange)
	if err != nil {
		return "", 0, 0, err
	}

	return filePath, startLine, endLine, nil
}

func parseLineRange(lineRange string) (int, int, error) {
	commaIndex := strings.Index(lineRange, ",")
	if commaIndex == -1 {
		return 0, 0, fmt.Errorf("invalid range format")
	}

	startPart := lineRange[:commaIndex]
	endPart := lineRange[commaIndex+1:]

	startLine, err := extractLineNumber(startPart)
	if err != nil {
		return 0, 0, err
	}

	endLine, err := extractLineNumber(endPart)
	if err != nil {
		return 0, 0, err
	}

	return startLine, endLine, nil
}

func extractLineNumber(part string) (int, error) {
	dotIndex := strings.Index(part, ".")
	if dotIndex == -1 {
		return 0, fmt.Errorf("invalid line format")
	}

	lineNum, err := strconv.Atoi(part[:dotIndex])
	if err != nil {
		return 0, fmt.Errorf("invalid line number: %w", err)
	}

	return lineNum, nil
}

func initializeFileCoverage(fileCoverage map[string]*FileCoverage, filePath string) {
	if fileCoverage[filePath] == nil {
		fileCoverage[filePath] = &FileCoverage{
			SmartLinesSet:            make(map[int]bool),
			SmartCoveredLinesSet:     make(map[int]bool),
			SmartLineCoveragePercent: 0.0,
		}
	}
}

func processLineRange(coverage *FileCoverage, filePath string, startLine, endLine, coverageFlag int) {
	sourceLines, err := loadSourceFile(filePath)
	if err != nil {
		addBasicLineRange(coverage, startLine, endLine, coverageFlag)
		return
	}

	addSmartFilteredLines(coverage, sourceLines, startLine, endLine, coverageFlag)
}

func addBasicLineRange(coverage *FileCoverage, startLine, endLine, coverageFlag int) {
	for line := startLine; line <= endLine; line++ {
		coverage.SmartLinesSet[line] = true
		if coverageFlag == 1 {
			coverage.SmartCoveredLinesSet[line] = true
		}
	}
}

func addSmartFilteredLines(coverage *FileCoverage, sourceLines []string, startLine, endLine, coverageFlag int) {
	for line := startLine; line <= endLine; line++ {
		if isExecutableLine(sourceLines, line) {
			coverage.SmartLinesSet[line] = true
			if coverageFlag == 1 {
				coverage.SmartCoveredLinesSet[line] = true
			}
		}
	}
}

func calculateCoveragePercentages(fileCoverage map[string]*FileCoverage) {
	for _, coverage := range fileCoverage {
		smartLines := len(coverage.SmartLinesSet)

		smartCoveredLines := len(coverage.SmartCoveredLinesSet)
		if smartLines > 0 {
			coverage.SmartLineCoveragePercent = float64(smartCoveredLines) / float64(smartLines) * 100
		}
	}
}

func displayCoverage(fileCoverage map[string]*FileCoverage) {
	if len(fileCoverage) == 0 {
		fmt.Println("No coverage data found")
		return
	}

	maxFilePathLen := calculateMaxPathLength(fileCoverage)
	printHeader(maxFilePathLen)

	sortedFiles := getSortedFiles(fileCoverage)
	printCoverageRows(fileCoverage, sortedFiles, maxFilePathLen)
	printFooter(maxFilePathLen)
	printSummary(fileCoverage, sortedFiles)
}

func calculateMaxPathLength(fileCoverage map[string]*FileCoverage) int {
	maxLen := 0

	for filePath := range fileCoverage {
        rp := getRepoPrefix()
        shortPath := filePath
        if rp != "" {
            shortPath = strings.TrimPrefix(filePath, rp)
        }
		if len(shortPath) > maxLen {
			maxLen = len(shortPath)
		}
	}

	if maxLen < 40 {
		maxLen = 40
	}

	return maxLen
}

func printHeader(maxLen int) {
	fmt.Printf("╭%s┬%s╮\n", strings.Repeat("─", maxLen+2), "────────")
	fmt.Printf("│ %-*s │ Cover  │\n", maxLen, "File")
	fmt.Printf("├%s┼%s┤\n", strings.Repeat("─", maxLen+2), "────────")
}

func getSortedFiles(fileCoverage map[string]*FileCoverage) []string {
	sortedFiles := make([]string, 0, len(fileCoverage))

	for filePath := range fileCoverage {
		sortedFiles = append(sortedFiles, filePath)
	}

	// Simple sort by file path
	for i := 0; i < len(sortedFiles); i++ {
		for j := i + 1; j < len(sortedFiles); j++ {
			if sortedFiles[i] > sortedFiles[j] {
				sortedFiles[i], sortedFiles[j] = sortedFiles[j], sortedFiles[i]
			}
		}
	}

	return sortedFiles
}

func printCoverageRows(fileCoverage map[string]*FileCoverage, sortedFiles []string, maxLen int) {
	for _, filePath := range sortedFiles {
		coverage := fileCoverage[filePath]
        rp := getRepoPrefix()
        shortPath := filePath
        if rp != "" {
            shortPath = strings.TrimPrefix(filePath, rp)
        }

		if len(coverage.SmartLinesSet) == 0 {
			fmt.Printf("│ %-*s │   --   │\n", maxLen, shortPath)
		} else {
			fmt.Printf("│ %-*s │ %5.1f%% │\n", maxLen, shortPath, coverage.SmartLineCoveragePercent)
		}
	}
}

func printFooter(maxLen int) {
	fmt.Printf("╰%s┴%s╯\n", strings.Repeat("─", maxLen+2), "────────")
}

func printSummary(fileCoverage map[string]*FileCoverage, sortedFiles []string) {
	stats := calculateStats(fileCoverage, sortedFiles)

	fmt.Println()

	if stats.filesWithCoverage > 0 {
		fmt.Printf("Summary: %d files, %.1f%% average coverage (%.1f%% - %.1f%% range)\n",
			len(sortedFiles), stats.avgCoverage, stats.minCoverage, stats.maxCoverage)
		fmt.Printf("Total: %d/%d lines covered (%.1f%% overall coverage)\n",
			stats.totalCovered, stats.totalLines, stats.overallCoverage)
		fmt.Printf("Note: Smart filtering excludes braces, comments, and structural elements\n")
	} else {
		fmt.Printf("Summary: %d files, no coverage data available\n", len(sortedFiles))
	}
}

type CoverageStats struct {
	filesWithCoverage int
	avgCoverage       float64
	minCoverage       float64
	maxCoverage       float64
	totalLines        int
	totalCovered      int
	overallCoverage   float64
}

func calculateStats(fileCoverage map[string]*FileCoverage, sortedFiles []string) CoverageStats {
	stats := CoverageStats{
		filesWithCoverage: 0,
		avgCoverage:       0.0,
		minCoverage:       100.0,
		maxCoverage:       0.0,
		totalLines:        0,
		totalCovered:      0,
		overallCoverage:   0.0,
	}

	totalSmartCovered := 0.0

	for _, filePath := range sortedFiles {
		coverage := fileCoverage[filePath]
		smartLines := len(coverage.SmartLinesSet)
		smartCovered := len(coverage.SmartCoveredLinesSet)

		stats.totalLines += smartLines
		stats.totalCovered += smartCovered

		if smartLines > 0 {
			stats.filesWithCoverage++
			smartPercent := coverage.SmartLineCoveragePercent
			totalSmartCovered += smartPercent

			if smartPercent < stats.minCoverage {
				stats.minCoverage = smartPercent
			}

			if smartPercent > stats.maxCoverage {
				stats.maxCoverage = smartPercent
			}
		}
	}

	if stats.filesWithCoverage > 0 {
		stats.avgCoverage = totalSmartCovered / float64(stats.filesWithCoverage)
	}

	if stats.totalLines > 0 {
		stats.overallCoverage = float64(stats.totalCovered) / float64(stats.totalLines) * 100
	}

	return stats
}

// Helper functions for smart filtering
func loadSourceFile(modulePath string) ([]string, error) {
	// Accept both module-qualified and repo-relative paths.
    rp := getRepoPrefix()
    filePath := modulePath
    if rp != "" {
        filePath = strings.TrimPrefix(modulePath, rp)
    }
	// If the path still doesn't exist relative to CWD (repo root), bail and let caller fall back.
	if _, statErr := os.Stat(filePath); statErr != nil {
		return nil, fmt.Errorf("failed to locate source file %s: %w", filePath, statErr)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file %s: %w", filePath, err)
	}

	defer func() { _ = file.Close() }()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return lines, nil
}

func isExecutableLine(sourceLines []string, lineNum int) bool {
	if lineNum < 1 || lineNum > len(sourceLines) {
		return true // Default to true if we can't check
	}

	line := strings.TrimSpace(sourceLines[lineNum-1]) // sourceLines is 0-indexed, lineNum is 1-indexed

	// Simple, reliable filters for obvious non-executable lines
	if line == "" {
		return false // Empty lines
	}

	if line == "{" || line == "}" {
		return false // Standalone braces
	}

	if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
		return false // Comments
	}

	// Package declarations
	if strings.HasPrefix(line, "package ") {
		return false
	}

	// Import lines
	if strings.HasPrefix(line, "import ") || (strings.HasPrefix(line, "\"") && strings.HasSuffix(line, "\"")) {
		return false
	}

	return true // Conservative: assume executable if not obviously structural
}
