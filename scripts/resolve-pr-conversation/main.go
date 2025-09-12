// usr/bin/env go run "$0" "$@"; exit
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// GitHubThread represents a GitHub review thread
type GitHubThread struct {
	ID         string `json:"id"`
	IsResolved bool   `json:"isResolved"`
}

// GitHubError represents a GitHub API error
type GitHubError struct {
	Message string `json:"message"`
}

// GitHubResponse represents the GitHub GraphQL API response
type GitHubResponse struct {
	Data *struct {
		ResolveReviewThread *struct {
			Thread GitHubThread `json:"thread"`
		} `json:"resolveReviewThread,omitempty"`
		AddPullRequestReviewThreadReply *struct {
			Comment struct {
				ID string `json:"id"`
			} `json:"comment"`
		} `json:"addPullRequestReviewThreadReply,omitempty"`
	} `json:"data,omitempty"`
	Errors []GitHubError `json:"errors,omitempty"`
}

// addCommentToThread adds a comment to a specific review thread
func addCommentToThread(threadID, comment string) error {
	fmt.Printf("Adding comment to thread: %s\n", threadID)

	// Escape the comment for GraphQL
	escapedComment := strings.ReplaceAll(comment, `\`, `\\`)
	escapedComment = strings.ReplaceAll(escapedComment, `"`, `\"`)
	escapedComment = strings.ReplaceAll(escapedComment, "\n", `\n`)

	query := fmt.Sprintf(`mutation {
		addPullRequestReviewThreadReply(input: {
			pullRequestReviewThreadId: "%s",
			body: "%s"
		}) {
			comment {
				id
			}
		}
	}`, threadID, escapedComment)

	// #nosec G204 - query is constructed internally and not from user input
	cmd := exec.Command("gh", "api", "graphql", "-f", "query="+query)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute GraphQL query: %w", err)
	}

	var response GitHubResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Data != nil && response.Data.AddPullRequestReviewThreadReply != nil {
		fmt.Println("✅ Successfully added comment to thread")
		return nil
	}

	if len(response.Errors) > 0 {
		fmt.Fprintln(os.Stderr, "❌ GitHub API errors:")

		for _, err := range response.Errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", err.Message)
		}

		return fmt.Errorf("GitHub API returned errors")
	}

	return fmt.Errorf("unexpected response format")
}

// resolveConversation resolves a GitHub review thread (conversation) using the provided thread ID
func resolveConversation(conversationID, comment string) error {
	if conversationID == "" {
		fmt.Fprintln(os.Stderr, "Error: conversationId is required")
		fmt.Fprintln(os.Stderr, "Usage: go run resolve-pr-conversation.go <conversationId> [comment]")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "If a comment is provided, it will be added to the conversation thread before resolving.")

		return fmt.Errorf("conversationId is required")
	}

	// Add comment to thread if provided
	if comment != "" {
		if err := addCommentToThread(conversationID, comment); err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}
	}

	// Resolve the thread
	query := fmt.Sprintf(`mutation {
		resolveReviewThread(input: {threadId: "%s"}) {
			thread { id isResolved }
		}
	}`, conversationID)

	fmt.Printf("Resolving conversation: %s\n", conversationID)

	// #nosec G204 - query is constructed internally and not from user input
	cmd := exec.Command("gh", "api", "graphql", "-f", "query="+query)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to execute GraphQL query: %w", err)
	}

	var response GitHubResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Data != nil && response.Data.ResolveReviewThread != nil {
		thread := response.Data.ResolveReviewThread.Thread
		fmt.Printf("✅ Successfully resolved conversation %s\n", thread.ID)

		status := "Not resolved"
		if thread.IsResolved {
			status = "Resolved"
		}

		fmt.Printf("Status: %s\n", status)

		return nil
	}

	if len(response.Errors) > 0 {
		fmt.Fprintln(os.Stderr, "❌ GitHub API errors:")

		for _, err := range response.Errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", err.Message)
		}

		return fmt.Errorf("GitHub API returned errors")
	}

	return fmt.Errorf("unexpected response format")
}

func main() {
	// Get conversationId and optional comment from command line arguments
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run resolve-pr-conversation.go <conversationId> [comment]")
		os.Exit(1)
	}

	conversationID := os.Args[1]

	var comment string
	if len(os.Args) > 2 {
		comment = os.Args[2]
	}

	// Call the function and handle it properly
	if err := resolveConversation(conversationID, comment); err != nil {
		log.Fatalf("Failed to resolve conversation: %v", err)
	}
}
