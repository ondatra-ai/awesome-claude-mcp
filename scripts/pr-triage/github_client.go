package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type GitHubClient interface {
    GetCurrentPRNumber(ctx context.Context) (int, error)
    ListAllReviewThreads(ctx context.Context, prNumber int) ([]Thread, error)
    ResolveReply(ctx context.Context, threadID, body string, resolve bool) error
}

type ghCLIClient struct{}

func NewGitHubCLIClient() GitHubClient { return &ghCLIClient{} }

func (c *ghCLIClient) GetCurrentPRNumber(ctx context.Context) (int, error) {
	// Primary: use gh to view PR for current checkout
	out, err := runShell(ctx, "gh", "pr", "view", "--json", "number", "-q", ".number")
	if err == nil {
		n, convErr := strconv.Atoi(strings.TrimSpace(out))
		if convErr == nil {
			return n, nil
		}
	}
	// Fallback: detect current branch and list PRs for it
	branch, berr := getCurrentBranch(ctx)
	if berr != nil {
		return 0, berr
	}
	lout, lerr := runShell(ctx, "gh", "pr", "list", "--head", branch, "--json", "number,title,state,url", "--limit", "1")
	if lerr != nil {
		return 0, fmt.Errorf("gh pr list failed: %w, out=%s", lerr, lout)
	}
	var prs []pullRequest
	if err := json.Unmarshal([]byte(lout), &prs); err != nil {
		return 0, fmt.Errorf("parse pr list: %w", err)
	}
	if len(prs) == 0 {
		return 0, fmt.Errorf("no PR found for branch: %s", branch)
	}
	return prs[0].Number, nil
}

func (c *ghCLIClient) ListAllReviewThreads(ctx context.Context, prNumber int) ([]Thread, error) {
    owner, name, err := repoOwnerAndName(ctx)
    if err != nil {
        return nil, err
    }
    // Fetch ALL review threads (paginated)
    var all []Thread
    after := ""
    for {
        query := buildThreadsPageQuery(owner, name)
        args := []string{"gh", "api", "graphql", "-f", "query=" + query, "-F", fmt.Sprintf("prNumber=%d", prNumber)}
        if after != "" {
            args = append(args, "-F", "after="+after)
        }
        out, err := runShell(ctx, args[0], args[1:]...)
        if err != nil {
            return nil, fmt.Errorf("graphql list threads: %w, out=%s", err, out)
        }
        var page threadsPageResponse
        if err := json.Unmarshal([]byte(out), &page); err != nil {
            return nil, fmt.Errorf("parse threads page: %w", err)
        }
        for _, tn := range page.Data.Repository.PullRequest.ReviewThreads.Nodes {
            var cs []Comment
            for _, c := range tn.Comments.Nodes {
                if c.Outdated {
                    continue
                }
                file := ""
                line := 0
                if c.Path != nil { file = *c.Path }
                if c.Line != nil { line = *c.Line }
                cs = append(cs, Comment{File: file, Line: line, URL: c.URL, Body: c.Body})
            }
            all = append(all, Thread{ID: tn.ID, IsResolved: tn.IsResolved, Comments: cs})
        }
        if page.Data.Repository.PullRequest.ReviewThreads.PageInfo.HasNextPage && page.Data.Repository.PullRequest.ReviewThreads.PageInfo.EndCursor != nil {
            after = *page.Data.Repository.PullRequest.ReviewThreads.PageInfo.EndCursor
            continue
        }
        break
    }
    // Filter after complete pagination: collect all unresolved threads with at least one non-outdated comment
    var eligible []Thread
    for _, th := range all {
        if !th.IsResolved && len(th.Comments) > 0 {
            eligible = append(eligible, th)
        }
    }
    if len(eligible) == 0 {
        return nil, fmt.Errorf("no eligible review threads found for PR %d", prNumber)
    }
    return eligible, nil
}

func (c *ghCLIClient) ResolveReply(ctx context.Context, threadID, body string, resolve bool) error {
	if resolve {
		// Resolve the review thread via GraphQL
		query := "mutation($tid:ID!){ resolveReviewThread(input:{threadId:$tid}){ thread { isResolved } } }"
		out, err := runShell(ctx, "gh", "api", "graphql", "-f", "query="+query, "-F", "tid="+threadID)
		if err != nil {
			return fmt.Errorf("resolve thread: %w, out=%s", err, out)
		}
		_ = out
		return nil
	}
	query := "mutation($tid:ID!, $body:String!){ addPullRequestReviewThreadReply(input:{pullRequestReviewThreadId:$tid, body:$body}){ comment { id } } }"
	out, err := runShell(ctx, "gh", "api", "graphql", "-f", "query="+query, "-F", "tid="+threadID, "-F", "body="+body)
	if err != nil {
		return fmt.Errorf("reply thread: %w, out=%s", err, out)
	}
	_ = out
	return nil
}

// helpers for repo and mapping GraphQL to local types

func repoOwnerAndName(ctx context.Context) (string, string, error) {
	out, err := runShell(ctx, "gh", "repo", "view", "--json", "owner,name", "-q", ".owner.login + \" \" + .name")
	if err != nil {
		return "", "", fmt.Errorf("repo view: %w", err)
	}
	parts := strings.Split(strings.TrimSpace(out), " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected repo view output: %s", out)
	}
	return parts[0], parts[1], nil
}

type graphQLResponse struct {
	Data struct {
		Repository struct {
			PullRequest struct {
				ReviewThreads struct {
					Nodes []threadNode `json:"nodes"`
				} `json:"reviewThreads"`
			} `json:"pullRequest"`
		} `json:"repository"`
	} `json:"data"`
}

type threadNode struct {
	ID         string `json:"id"`
	IsResolved bool   `json:"isResolved"`
	Comments   struct {
		Nodes []commentNode `json:"nodes"`
	} `json:"comments"`
}

type commentNode struct {
    Path     *string `json:"path"`
    Line     *int    `json:"line"`
    Body     string  `json:"body"`
    Outdated bool    `json:"outdated"`
    URL      string  `json:"url"`
}

// pageInfo supports pagination cursors
type pageInfo struct {
    HasNextPage bool    `json:"hasNextPage"`
    EndCursor   *string `json:"endCursor"`
}

// threadsPageResponse is a minimal GraphQL response shape for paginating threads
type threadsPageResponse struct {
    Data struct {
        Repository struct {
            PullRequest struct {
                ReviewThreads struct {
                    PageInfo pageInfo `json:"pageInfo"`
                    Nodes    []struct {
                        ID         string `json:"id"`
                        IsResolved bool   `json:"isResolved"`
                        Comments   struct {
                            Nodes []commentNode `json:"nodes"`
                        } `json:"comments"`
                    } `json:"nodes"`
                } `json:"reviewThreads"`
            } `json:"pullRequest"`
        } `json:"repository"`
    } `json:"data"`
}

// commentsPage models the response for fetching a page of comments for a thread by node id
type commentsPage struct {
    Data struct {
        Node struct {
            ID        string `json:"id"`
            IsResolved bool  `json:"isResolved"`
            Comments  struct {
                PageInfo pageInfo      `json:"pageInfo"`
                Nodes    []commentNode `json:"nodes"`
            } `json:"comments"`
        } `json:"node"`
    } `json:"data"`
}

// buildThreadsPageQuery returns a query that supports pagination via $after
func buildThreadsPageQuery(owner, name string) string {
    return fmt.Sprintf(`
query($prNumber: Int!, $after: String) {
  repository(owner: "%s", name: "%s") {
    pullRequest(number: $prNumber) {
      reviewThreads(first: 100, after: $after) {
        pageInfo { hasNextPage endCursor }
        nodes {
          id
          isResolved
          comments(first: 50) {
            nodes { path line body outdated url }
          }
        }
      }
    }
  }
}`, owner, name)
}

// buildCommentsPageQuery fetches comments for a thread by node id with pagination
func buildCommentsPageQuery() string {
    return `
query($tid: ID!, $after: String) {
  node(id: $tid) {
    ... on PullRequestReviewThread {
      id
      isResolved
      comments(first: 100, after: $after) {
        pageInfo { hasNextPage endCursor }
        nodes { path line body outdated url }
      }
    }
  }
}`
}

// fetchAllComments retrieves all non-outdated comments for a given thread id
func fetchAllComments(ctx context.Context, threadID string) ([]Comment, error) {
    var all []Comment
    after := ""
    query := buildCommentsPageQuery()
    for {
        args := []string{"gh", "api", "graphql", "-f", "query=" + query, "-F", "tid=" + threadID}
        if after != "" {
            args = append(args, "-F", "after="+after)
        }
        out, err := runShell(ctx, args[0], args[1:]...)
        if err != nil {
            return nil, fmt.Errorf("comments page: %w, out=%s", err, out)
        }
        var page commentsPage
        if err := json.Unmarshal([]byte(out), &page); err != nil {
            return nil, fmt.Errorf("parse comments page: %w", err)
        }
        for _, c := range page.Data.Node.Comments.Nodes {
            if c.Outdated {
                continue
            }
            file := ""
            line := 0
            if c.Path != nil {
                file = *c.Path
            }
            if c.Line != nil {
                line = *c.Line
            }
            all = append(all, Comment{File: file, Line: line, URL: c.URL, Body: c.Body})
        }
        if page.Data.Node.Comments.PageInfo.HasNextPage && page.Data.Node.Comments.PageInfo.EndCursor != nil {
            after = *page.Data.Node.Comments.PageInfo.EndCursor
            continue
        }
        break
    }
    return all, nil
}

func buildQuery(owner, name string) string {
    return fmt.Sprintf(`
query($prNumber: Int!) {
  repository(owner: "%s", name: "%s") {
    pullRequest(number: $prNumber) {
      # Fetch a page of threads; we'll filter unresolved client-side and return one
      reviewThreads(first: 100) {
        nodes {
          id
          isResolved
          # Fetch a reasonable slice of comments for context
          comments(first: 50) {
            nodes {
              path
              line
              body
              outdated
              url
            }
          }
        }
      }
    }
  }
}`, owner, name)
}

func toThreads(resp graphQLResponse) []Thread {
    var out []Thread
    for _, t := range resp.Data.Repository.PullRequest.ReviewThreads.Nodes {
        // Skip resolved threads; we only want active conversations
        if t.IsResolved {
            continue
        }
        var cs []Comment
        for _, c := range t.Comments.Nodes {
            if c.Outdated {
                continue
            }
			file := ""
			line := 0
			if c.Path != nil {
				file = *c.Path
			}
			if c.Line != nil {
				line = *c.Line
			}
			cs = append(cs, Comment{File: file, Line: line, URL: c.URL, Body: c.Body})
		}
		if len(cs) == 0 {
			continue
		}
		out = append(out, Thread{ID: t.ID, IsResolved: t.IsResolved, Comments: cs})
	}
	return out
}

// Local helpers copied from get-pr-number logic (simplified)
type pullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  string `json:"state"`
	URL    string `json:"url"`
}

func getCurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git branch: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
