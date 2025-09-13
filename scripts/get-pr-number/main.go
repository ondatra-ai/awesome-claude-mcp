// usr/bin/env go run "$0" "$@"; exit

/*
Get PR Number Script

This script reveals the Pull Request number associated with
the current git branch. It uses the GitHub CLI (gh) to search
for open PRs linked to the current branch.

Usage:
  Direct execution: go run scripts/get-pr-number/main.go
  Build and run: go build -o get-pr-number scripts/get-pr-number/main.go && ./get-pr-number

Requirements:
  - GitHub CLI (gh) must be installed: https://cli.github.com/
  - Must be run from within a git repository
  - Repository must be hosted on GitHub

Output:
  - Shows PR number, title, state, and URL
  - If no open PR is found, shows closed/merged PRs
  - Displays just the PR number at the end for easy scripting
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
}

// getCurrentBranch gets the current git branch name
func getCurrentBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error getting current branch: %v", err)
	}

	return strings.TrimSpace(string(output))
}

// handleError handles errors with appropriate messages
func handleError(err error) {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "gh: command not found"):
		fmt.Fprintln(os.Stderr, "GitHub CLI (gh) is not installed. Please install it first:")
		fmt.Fprintln(os.Stderr, "https://cli.github.com/manual/installation")
	case strings.Contains(errMsg, "Could not resolve to a Repository"):
		fmt.Fprintln(os.Stderr, "This directory does not appear to be a GitHub repository.")
	default:
		fmt.Fprintf(os.Stderr, "Error fetching PR information: %v\n", err)
	}

	os.Exit(1)
}

// displayPR displays PR information
func displayPR(pullRequest PullRequest, branch string) {
	fmt.Printf("\nCurrent branch: %s\n", branch)
	fmt.Printf("PR #%d: %s\n", pullRequest.Number, pullRequest.Title)
	fmt.Printf("State: %s\n", pullRequest.State)
	fmt.Printf("URL: %s\n", pullRequest.URL)

	// Output just the PR number for easy scripting
	fmt.Printf("\nPR Number: %d\n", pullRequest.Number)
}

// checkForClosedPRs checks for closed/merged PRs on the branch
func checkForClosedPRs(branch string) {
	cmd := exec.Command("gh", "pr", "list", "--head", branch, "--state", "all",
		"--json", "number,title,state,url", "--limit", "5")

	output, err := cmd.Output()
	if err != nil {
		return // Ignore errors for closed PRs check
	}

	var allPRs []PullRequest
	if err := json.Unmarshal(output, &allPRs); err != nil {
		return // Ignore JSON parse errors
	}

	if len(allPRs) > 0 {
		fmt.Println("\nFound closed/merged PRs:")

		for _, pr := range allPRs {
			fmt.Printf("  PR #%d: %s (%s)\n", pr.Number, pr.Title, pr.State)
			fmt.Printf("  URL: %s\n", pr.URL)
		}
	}
}

// getPRForBranch gets PR information for the given branch
func getPRForBranch(branch string) {
	if branch == "main" || branch == "master" {
		fmt.Println("Current branch is the default branch. No PR associated.")
		return
	}

	// Try to get PR using GitHub CLI search
	cmd := exec.Command("gh", "pr", "list", "--head", branch,
		"--json", "number,title,state,url", "--limit", "1")

	output, err := cmd.Output()
	if err != nil {
		handleError(err)
	}

	var prs []PullRequest
	if err := json.Unmarshal(output, &prs); err != nil {
		handleError(fmt.Errorf("failed to parse PR data: %w", err))
	}

	if len(prs) == 0 {
		fmt.Printf("No PR found for branch: %s\n", branch)
		// Try to find PRs that might have been merged or closed
		checkForClosedPRs(branch)

		return
	}

	displayPR(prs[0], branch)
}

func main() {
	currentBranch := getCurrentBranch()
	fmt.Printf("Checking for PR associated with branch: %s\n", currentBranch)
	getPRForBranch(currentBranch)
}
