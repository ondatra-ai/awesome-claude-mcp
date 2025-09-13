// usr/bin/env go run "$0" "$@"; exit

/*
SonarCloud Coverage Analysis Tool

Analyzes test coverage for files with coverage data in the current PR using SonarCloud API.
Shows exact coverage data matching the SonarCloud UI.

Usage: go run scripts/sonar-uncovered-lines/main.go

Requirements:
- SONAR_TOKEN environment variable
- GitHub CLI (gh) installed and authenticated
- sonar-project.properties file with project configuration
*/

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Core types
type Config struct {
	ProjectKey   string
	Organization string
	Token        string
	HostURL      string
}

type Component struct {
	Key      string `json:"key"`
	Path     string `json:"path"`
	Measures []struct {
		Metric  string `json:"metric"`
		Periods []struct {
			Value string `json:"value"`
		} `json:"periods"`
	} `json:"measures"`
}

func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout:       30 * time.Second,
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
	}
}

// Configuration loading
func loadConfig() (*Config, error) {
	loadEnvFile()

	config := &Config{
		ProjectKey:   "",
		Organization: "",
		Token:        "",
		HostURL:      "https://sonarcloud.io",
	}

	// Load from sonar-project.properties
	if err := parsePropertiesFile("sonar-project.properties", map[string]*string{
		"sonar.projectKey":   &config.ProjectKey,
		"sonar.organization": &config.Organization,
	}); err != nil {
		return nil, fmt.Errorf("failed to load sonar properties: %w", err)
	}

	// Load from environment
	config.Token = os.Getenv("SONAR_TOKEN")
	if config.Token == "" {
		return nil, fmt.Errorf("SONAR_TOKEN environment variable required")
	}

	if hostURL := os.Getenv("SONAR_HOST_URL"); hostURL != "" {
		config.HostURL = hostURL
	}

	// Validate
	if config.ProjectKey == "" || config.Organization == "" {
		return nil, fmt.Errorf("missing project configuration in sonar-project.properties")
	}

	return config, nil
}

func parsePropertiesFile(filename string, targets map[string]*string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open properties file: %w", err)
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if target, exists := targets[key]; exists {
				*target = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading properties file: %w", err)
	}

	return nil
}

func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}

	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if parts := strings.SplitN(line, "=", 2); len(parts) == 2 {
			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if len(value) >= 2 && ((strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'"))) {
				value = value[1 : len(value)-1]
			}

			if os.Getenv(key) == "" {
				_ = os.Setenv(key, value)
			}
		}
	}
}

// Git operations
func getCurrentPRNumber() (int, error) {
	branch, err := execCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return 0, fmt.Errorf("failed to get current branch: %w", err)
	}

	if branch == "main" || branch == "master" {
		return 0, fmt.Errorf("not on a feature branch")
	}

	output, err := execCommand("gh", "pr", "list", "--head", branch, "--json", "number", "--limit", "1")
	if err != nil {
		return 0, fmt.Errorf("failed to get PR info: %w", err)
	}

	var prs []struct {
		Number int `json:"number"`
	}
	if err := json.Unmarshal([]byte(output), &prs); err != nil || len(prs) == 0 {
		return 0, fmt.Errorf("no PR found for branch: %s", branch)
	}

	return prs[0].Number, nil
}

func execCommand(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).Output()
	return strings.TrimSpace(string(output)), err
}

// SonarCloud API
func makeRequest(config *Config, endpoint string, params url.Values) (*http.Response, error) {
	apiURL := fmt.Sprintf("%s%s?%s", config.HostURL, endpoint, params.Encode())

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.Token)
	req.Header.Set("Accept", "application/json")

	client := getHTTPClient()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func getComponentsWithCoverage(config *Config, prNumber int) ([]Component, error) {
	params := url.Values{
		"pullRequest":      {strconv.Itoa(prNumber)},
		"additionalFields": {"metrics"},
		"ps":               {"500"},
		"s":                {"metricPeriod"},
		"metricSort":       {"new_coverage"},
		"metricSortFilter": {"withMeasuresOnly"},
		"metricPeriodSort": {"1"},
		"component":        {config.ProjectKey},
		"metricKeys":       {"new_coverage,new_uncovered_lines"},
		"strategy":         {"leaves"},
	}

	resp, err := makeRequest(config, "/api/measures/component_tree", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get components: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, body)
	}

	var response struct {
		Components []Component `json:"components"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Components, nil
}

type UncoveredLines struct {
	NewLines []int
	OldLines []int
}

func getUncoveredLineNumbers(config *Config, componentKey string, prNumber int) (*UncoveredLines, error) {
	params := url.Values{
		"key":         {componentKey},
		"from":        {"1"},
		"to":          {"10000"},
		"pullRequest": {strconv.Itoa(prNumber)},
	}

	resp, err := makeRequest(config, "/api/sources/lines", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get source lines: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return &UncoveredLines{
			NewLines: []int{},
			OldLines: []int{},
		}, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, body)
	}

	var response struct {
		Sources []struct {
			Line     int  `json:"line"`
			LineHits *int `json:"lineHits,omitempty"`
			IsNew    bool `json:"isNew,omitempty"`
		} `json:"sources"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode sources response: %w", err)
	}

	result := &UncoveredLines{
		NewLines: []int{},
		OldLines: []int{},
	}

	for _, source := range response.Sources {
		// Only include lines that have lineHits field present AND lineHits == 0
		if source.LineHits != nil && *source.LineHits == 0 {
			if source.IsNew {
				result.NewLines = append(result.NewLines, source.Line)
			} else {
				result.OldLines = append(result.OldLines, source.Line)
			}
		}
	}

	return result, nil
}

func analyzeFiles(config *Config, prNumber int) error {
	components, err := getComponentsWithCoverage(config, prNumber)
	if err != nil {
		return fmt.Errorf("failed to get components with coverage: %w", err)
	}

	if len(components) == 0 {
		fmt.Println("No files with coverage data found in this PR")
		return nil
	}

	for _, component := range components {
		fmt.Printf("📁 %s\n", component.Path)

		// Display official metrics from component tree
		for _, measure := range component.Measures {
			if measure.Metric == "new_coverage" && len(measure.Periods) > 0 {
				if coverageFloat, err := strconv.ParseFloat(measure.Periods[0].Value, 64); err == nil {
					fmt.Printf("   📊 New Code Coverage: %.1f%%\n", coverageFloat)
				}
			}

			if measure.Metric == "new_uncovered_lines" && len(measure.Periods) > 0 {
				fmt.Printf("   ❌ New Uncovered Lines: %s\n", measure.Periods[0].Value)
			}
		}

		// Get specific line numbers
		uncoveredLines, err := getUncoveredLineNumbers(config, component.Key, prNumber)
		if err != nil {
			fmt.Printf("   Warning: Failed to get line numbers: %v\n", err)
		} else {
			if len(uncoveredLines.NewLines) > 0 {
				fmt.Printf("   🆕 New Code Uncovered Lines: %v\n", uncoveredLines.NewLines)
			}

			if len(uncoveredLines.OldLines) > 0 {
				fmt.Printf("   📜 Existing Code Uncovered Lines: %v\n", uncoveredLines.OldLines)
			}
		}

		fmt.Println()
	}

	return nil
}

// Main execution
func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("SonarCloud Coverage Analysis Tool")
		fmt.Println("\nUsage: go run scripts/sonar-uncovered-lines/main.go")
		fmt.Println("\nRequirements:")
		fmt.Println("  - SONAR_TOKEN environment variable")
		fmt.Println("  - GitHub CLI (gh) installed and authenticated")
		fmt.Println("  - sonar-project.properties file with project configuration")

		return
	}

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	prNumber, err := getCurrentPRNumber()
	if err != nil {
		log.Fatalf("Failed to get PR number: %v", err)
	}

	fmt.Printf("📊 SonarCloud Coverage Analysis for PR #%d\n", prNumber)
	fmt.Printf("🏢 Project: %s/%s\n\n", config.Organization, config.ProjectKey)

	if err := analyzeFiles(config, prNumber); err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}
}
