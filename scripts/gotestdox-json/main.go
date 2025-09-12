// gotestdox-json converts Go test JSON output into readable documentation format
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"

	StatusPass = "pass"
	StatusFail = "fail"
	StatusSkip = "skip"
)

type TestEvent struct {
	Time    time.Time `json:"time"`
	Action  string    `json:"action"`
	Package string    `json:"package"`
	Test    string    `json:"test,omitempty"`
	Output  string    `json:"output,omitempty"`
	Elapsed float64   `json:"elapsed,omitempty"`
}

type TestResult struct {
	Name    string
	Package string
	Status  string
	Elapsed float64
	Output  []string
}

func main() {
	var (
		testResults map[string]*TestResult
		err         error
	)

	if len(os.Args) >= 2 {
		testResults, err = parseTestJSON(os.Args[1])
	} else {
		testResults = parseTestJSONFromReader(os.Stdin)
		err = nil
	}

	if err != nil {
		log.Fatalf("Error parsing test JSON: %v", err)
	}

	displayGotestdoxReport(testResults)
}

func parseTestJSON(filename string) (map[string]*TestResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() { _ = file.Close() }()

	return parseTestJSONFromReader(file), nil
}

func parseTestJSONFromReader(reader io.Reader) map[string]*TestResult {
	testResults := make(map[string]*TestResult)
	scanner := bufio.NewScanner(reader)
	// Increase buffer in case of long output lines
	const maxScannerTokenSize = 1024 * 1024

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxScannerTokenSize)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}

		var event TestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		if event.Test == "" {
			continue
		}

		testKey := event.Package + "::" + event.Test
		result := getOrCreateTestResult(testResults, testKey, event)
		updateTestResult(result, event)
	}
	// Ignore scanner errors and return what we parsed
	return testResults
}

func getOrCreateTestResult(testResults map[string]*TestResult, testKey string, event TestEvent) *TestResult {
	if testResults[testKey] == nil {
		testResults[testKey] = &TestResult{
			Name:    event.Test,
			Package: event.Package,
			Status:  "unknown",
			Output:  []string{},
			Elapsed: 0.0,
		}
	}

	return testResults[testKey]
}

func updateTestResult(result *TestResult, event TestEvent) {
	switch event.Action {
	case StatusPass:
		result.Status = StatusPass
		result.Elapsed = event.Elapsed
	case StatusFail:
		result.Status = StatusFail
		result.Elapsed = event.Elapsed
	case StatusSkip:
		result.Status = StatusSkip
		result.Elapsed = event.Elapsed
	case "output":
		if isFailureOutput(event.Output) {
			result.Output = append(result.Output, strings.TrimSpace(event.Output))
		}
	}
}

func isFailureOutput(output string) bool {
	return strings.Contains(output, ": want ") ||
		strings.Contains(output, "FAIL:") ||
		strings.Contains(output, "panic:") ||
		(strings.Contains(output, ".go:") && strings.Contains(output, ":"))
}

func displayGotestdoxReport(testResults map[string]*TestResult) {
	if len(testResults) == 0 {
		fmt.Println("No test results found")
		return
	}

	packageTests := groupTestsByPackage(testResults)
	packages := getSortedPackages(packageTests)

	stats := displayPackageResults(packageTests, packages)
	printSummary(stats)
}

func groupTestsByPackage(testResults map[string]*TestResult) map[string][]*TestResult {
	packageTests := make(map[string][]*TestResult)
	for _, result := range testResults {
		packageTests[result.Package] = append(packageTests[result.Package], result)
	}

	return packageTests
}

func getSortedPackages(packageTests map[string][]*TestResult) []string {
	packages := make([]string, 0, len(packageTests))
	for pkg := range packageTests {
		packages = append(packages, pkg)
	}

	sort.Strings(packages)

	return packages
}

type TestStats struct {
	totalPassed  int
	totalFailed  int
	totalSkipped int
	totalElapsed float64
	hadFailures  bool
}

func displayPackageResults(packageTests map[string][]*TestResult, packages []string) TestStats {
	stats := TestStats{
		totalPassed:  0,
		totalFailed:  0,
		totalSkipped: 0,
		totalElapsed: 0.0,
		hadFailures:  false,
	}

	for idx, pkg := range packages {
		if idx > 0 {
			fmt.Println()
		}

		shortPkg := strings.TrimPrefix(pkg, "github.com/ondatra-ai/flow-test-go/")
		fmt.Printf("%s:\n", shortPkg)

		tests := packageTests[pkg]
		sort.Slice(tests, func(i, j int) bool {
			return tests[i].Name < tests[j].Name
		})

		displayTestResults(tests, &stats)
	}

	return stats
}

func displayTestResults(tests []*TestResult, stats *TestStats) {
	for _, test := range tests {
		if test.Status == "unknown" {
			continue
		}

		symbol, color := getTestSymbolAndColor(test.Status)
		updateStats(stats, test)

		readableName := convertTestName(test.Name)

		elapsedStr := fmt.Sprintf("%.2fs", test.Elapsed)
		if test.Elapsed < 0.01 {
			elapsedStr = "0.00s"
		}

		coloredSymbol := colorize(symbol, color)
		fmt.Printf(" %s %s (%s)\n", coloredSymbol, readableName, elapsedStr)

		if test.Status == StatusFail && len(test.Output) > 0 {
			for _, output := range test.Output {
				if strings.TrimSpace(output) != "" {
					fmt.Printf("   %s\n", colorize(output, ColorRed))
				}
			}
		}
	}
}

func getTestSymbolAndColor(status string) (string, string) {
	switch status {
	case StatusPass:
		return "✔", ColorGreen
	case StatusFail:
		return "✗", ColorRed
	case StatusSkip:
		return "○", ColorYellow
	default:
		return "?", ""
	}
}

func updateStats(stats *TestStats, test *TestResult) {
	switch test.Status {
	case StatusPass:
		stats.totalPassed++
	case StatusFail:
		stats.totalFailed++
		stats.hadFailures = true
	case StatusSkip:
		stats.totalSkipped++
	}

	stats.totalElapsed += test.Elapsed
}

func printSummary(stats TestStats) {
	fmt.Println()

	totalTests := stats.totalPassed + stats.totalFailed + stats.totalSkipped
	fmt.Printf("Summary: %d tests, %d passed, %d failed, %d skipped (%.2fs)\n",
		totalTests, stats.totalPassed, stats.totalFailed, stats.totalSkipped, stats.totalElapsed)

	if stats.hadFailures {
		os.Exit(1)
	}
}

func colorize(text, color string) string {
	if !shouldUseColor() {
		return text
	}

	return color + text + ColorReset
}

func shouldUseColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}

func convertTestName(testName string) string {
	name := strings.TrimPrefix(testName, "Test")

	underscoreIndex := strings.Index(name, "_")
	if underscoreIndex != -1 {
		functionName := name[:underscoreIndex]
		description := name[underscoreIndex+1:]

		return functionName + " " + camelCaseToSentence(description)
	}

	return camelCaseToSentence(name)
}

func camelCaseToSentence(camelCase string) string {
	if len(camelCase) == 0 {
		return ""
	}

	// Simple conversion: insert space before capitals
	var result strings.Builder

	for idx, char := range camelCase {
		if idx > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune(' ')
		}

		if idx == 0 {
			result.WriteRune(char)
		} else {
			result.WriteRune(char)
		}
	}

	sentence := strings.ToLower(result.String())
	if len(sentence) > 0 {
		sentence = strings.ToUpper(sentence[:1]) + sentence[1:]
	}

	return sentence
}
