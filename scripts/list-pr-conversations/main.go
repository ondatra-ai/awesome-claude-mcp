// usr/bin/env go run "$0" "$@"; exit
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

type Conversation struct {
	ID         string    `json:"id"`
	IsResolved bool      `json:"isResolved"`
	Comments   []Comment `json:"comments"`
}

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

type ThreadNode struct {
	ID         string `json:"id"`
	IsResolved bool   `json:"isResolved"`
	Comments   struct {
		Nodes []CommentNode `json:"nodes"`
	} `json:"comments"`
}

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

func repoOwnerAndName() (string, string) {
	cmd := exec.Command("gh", "repo", "view", "--json", "owner,name", "-q", ".owner.login + \" \" + .name")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to get repo info: %v", err)
	}
	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(parts) != 2 {
		log.Fatalf("unexpected repo info: %s", string(out))
	}
	return parts[0], parts[1]
}

func currentPRNumber(arg string) int {
	if arg != "" {
		if n, err := strconv.Atoi(arg); err == nil {
			return n
		}
		log.Fatalf("invalid PR number: %q", arg)
	}
	cmd := exec.Command("gh", "pr", "view", "--json", "number", "-q", ".number")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to detect PR number: %v", err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		log.Fatalf("invalid PR number: %v", err)
	}
	return n
}

func buildQuery(owner, name string) string {
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
              author { login }
            }
          }
        }
      }
    }
  }
}`, owner, name)
}

func toConversations(resp GraphQLResponse) []Conversation {
	var convs []Conversation
	for _, t := range resp.Data.Repository.PullRequest.ReviewThreads.Nodes {
		// Skip resolved threads entirely
		if t.IsResolved {
			continue
		}
		var comments []Comment
		for _, c := range t.Comments.Nodes {
			// Skip outdated comments
			if c.Outdated {
				continue
			}
			comments = append(comments, Comment{
				File:      c.Path,
				Line:      c.Line,
				Author:    c.Author.Login,
				Body:      c.Body,
				CreatedAt: c.CreatedAt,
				Outdated:  c.Outdated,
				Resolved:  t.IsResolved,
				DiffHunk:  c.DiffHunk,
				URL:       c.URL,
			})
		}
		// If no comments remain after filtering, skip the thread
		if len(comments) == 0 {
			continue
		}
		convs = append(convs, Conversation{
			ID:         t.ID,
			IsResolved: t.IsResolved,
			Comments:   comments,
		})
	}
	// sort for deterministic order
	sort.Slice(convs, func(i, j int) bool {
		var ai, aj time.Time
		if len(convs[i].Comments) > 0 {
			ai, _ = time.Parse(time.RFC3339, convs[i].Comments[0].CreatedAt)
		}
		if len(convs[j].Comments) > 0 {
			aj, _ = time.Parse(time.RFC3339, convs[j].Comments[0].CreatedAt)
		}
		return ai.Before(aj)
	})
	return convs
}

func main() {
	// Usage: list-pr-conversations <output-file> [PR_NUMBER]
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: list-pr-conversations <output-file> [PR_NUMBER]")
		os.Exit(1)
	}
	outPath := strings.TrimSpace(os.Args[1])
	var prArg string
	if len(os.Args) >= 3 {
		prArg = strings.TrimSpace(os.Args[2])
	}

	owner, name := repoOwnerAndName()
	prNum := currentPRNumber(prArg)
	query := buildQuery(owner, name)

	cmd := exec.Command("gh", "api", "graphql", "-f", "query="+query, "-F", fmt.Sprintf("prNumber=%d", prNum))
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to query GraphQL API: %v", err)
	}
	var resp GraphQLResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		log.Fatalf("failed to parse GraphQL response: %v", err)
	}
	convs := toConversations(resp)

	// ensure output dir
	if dir := filepath.Dir(outPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Fatalf("failed to create output dir: %v", err)
		}
	}
	data, err := json.MarshalIndent(convs, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal conversations: %v", err)
	}
	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		log.Fatalf("failed to write %s: %v", outPath, err)
	}
	// do not print JSON to stdout; writing to file only
}
