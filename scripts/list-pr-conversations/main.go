// usr/bin/env go run "$0" "$@"; exit
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"
)

// Comment represents a PR review comment
type Comment struct {
	File      *string `json:"file"`
	Line      *int    `json:"line"`
	Author    string  `json:"author"`
	Body      string  `json:"body"`
	CreatedAt string  `json:"createdAt"`
	Outdated  bool    `json:"outdated"`
	Resolved  bool    `json:"resolved"`
	DiffHunk  string  `json:"diffHunk"`
	URL       string  `json:"url"`
}

// Conversation represents a PR review conversation
type Conversation struct {
	ID         string    `json:"id"`
	IsResolved bool      `json:"isResolved"`
	Comments   []Comment `json:"comments"`
}

// CommentNode represents a comment node from GitHub GraphQL API
type CommentNode struct {
	Path      *string `json:"path"`
	Line      *int    `json:"line"`
	Body      string  `json:"body"`
	CreatedAt string  `json:"createdAt"`
	Outdated  bool    `json:"outdated"`
	DiffHunk  string  `json:"diffHunk"`
	URL       string  `json:"url"`
	Author    struct {
		Login string `json:"login"`
	} `json:"author"`
}

// ThreadNode represents a review thread from GitHub GraphQL API
type ThreadNode struct {
	ID         string `json:"id"`
	IsResolved bool   `json:"isResolved"`
	Comments   struct {
		Nodes []CommentNode `json:"nodes"`
	} `json:"comments"`
}

// GraphQLResponse represents the GitHub GraphQL API response
type GraphQLResponse struct {
	Data struct {
		Repository struct {
			PullRequest struct {
				ReviewThreads struct {
					Nodes []ThreadNode `json:"nodes"`
				} `json:"reviewThreads"`
			} `json:"pullRequest"`
		} `json:"repository"`
	} `json:"data"`
}

// getPRNumber gets the PR number from command line arguments
func getPRNumber() string {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: list-pr-conversations.go [PR_NUMBER]")
		os.Exit(1)
	}

	return os.Args[1]
}

// buildGraphQLQuery builds the GraphQL query for fetching PR review threads
func buildGraphQLQuery() string {
	return `
query($prNumber: Int!) {
  repository(owner: "ondatra-ai", name: "flow-test-go") {
    pullRequest(number: $prNumber) {
      reviewThreads(first: 100) {
        nodes {
          id
          isResolved
          comments(first: 10) {
            nodes {
              path
              line
              body
              createdAt
              outdated
              diffHunk
              url
              author {
                login
              }
            }
          }
        }
      }
    }
  }
}`
}

// parseConversations parses the GraphQL response into conversations
func parseConversations(data GraphQLResponse) []Conversation {
	var conversations []Conversation

	for _, thread := range data.Data.Repository.PullRequest.ReviewThreads.Nodes {
		// Only include unresolved threads
		if thread.IsResolved {
			continue
		}

		var comments []Comment
		for _, comment := range thread.Comments.Nodes {
			comments = append(comments, Comment{
				File:      comment.Path,
				Line:      comment.Line,
				Author:    comment.Author.Login,
				Body:      comment.Body,
				CreatedAt: comment.CreatedAt,
				Outdated:  comment.Outdated,
				Resolved:  thread.IsResolved,
				DiffHunk:  comment.DiffHunk,
				URL:       comment.URL,
			})
		}

		conversations = append(conversations, Conversation{
			ID:         thread.ID,
			IsResolved: thread.IsResolved,
			Comments:   comments,
		})
	}

	return conversations
}

// getPRComments fetches and displays PR comments
func getPRComments(prNumber string) {
	query := buildGraphQLQuery()

	// Convert PR number to int to validate it
	prNum, err := strconv.Atoi(prNumber)
	if err != nil {
		log.Fatalf("Invalid PR number: %s", prNumber)
	}

	// #nosec G204 - query is constructed internally and not from user input
	cmd := exec.Command("gh", "api", "graphql", "-f", "query="+query, "-F", fmt.Sprintf("prNumber=%d", prNum))

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error fetching PR comments: %v", err)
	}

	var data GraphQLResponse
	if err := json.Unmarshal(output, &data); err != nil {
		log.Fatalf("Error parsing response: %v", err)
	}

	conversations := parseConversations(data)

	// Sort by creation date of first comment in each conversation
	sort.Slice(conversations, func(index1, index2 int) bool {
		var aTime, bTime time.Time
		if len(conversations[index1].Comments) > 0 {
			aTime, _ = time.Parse(time.RFC3339, conversations[index1].Comments[0].CreatedAt)
		}

		if len(conversations[index2].Comments) > 0 {
			bTime, _ = time.Parse(time.RFC3339, conversations[index2].Comments[0].CreatedAt)
		}

		return aTime.Before(bTime)
	})

	// Output JSON for programmatic use
	jsonOutput, err := json.MarshalIndent(conversations, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling conversations: %v", err)
	}

	fmt.Println(string(jsonOutput))
}

func main() {
	prNumber := getPRNumber()
	getPRComments(prNumber)
}
