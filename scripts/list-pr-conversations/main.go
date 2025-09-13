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
    "strings"
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
func buildGraphQLQuery(owner, name string) string {
    // owner and name are derived from the current repo via gh
    return fmt.Sprintf(`
query($prNumber: Int!) {
  repository(owner: "%s", name: "%s") {
    pullRequest(number: $prNumber) {
      reviewThreads(first: 100) {
        nodes {
          id
          isResolved
          comments(first: 50) {
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
}` , owner, name)
}

// parseConversations parses the GraphQL response into conversations
func parseConversations(data GraphQLResponse) []Conversation {
    var conversations []Conversation

    for _, thread := range data.Data.Repository.PullRequest.ReviewThreads.Nodes {
        // Include ALL threads (resolved and unresolved) to enable full reporting.
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
    // Detect current repo owner/name using gh
    repoInfoCmd := exec.Command("gh", "repo", "view", "--json", "owner,name", "-q", ".owner.login + \" \" + .name")
    repoOut, err := repoInfoCmd.Output()
    if err != nil {
        log.Fatalf("Error getting repo info: %v", err)
    }
    parts := strings.Split(strings.TrimSpace(string(repoOut)), " ")
    if len(parts) != 2 {
        log.Fatalf("Unexpected repo info format: %s", string(repoOut))
    }
    owner := parts[0]
    name := parts[1]

    query := buildGraphQLQuery(owner, name)

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
