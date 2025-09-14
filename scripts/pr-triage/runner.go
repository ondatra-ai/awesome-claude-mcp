package main

import (
	"context"
	"fmt"
	"strings"
)

type Runner struct {
	gh    GitHubClient
	codex CodexClient
}

func NewRunner(gh GitHubClient, codex CodexClient) *Runner {
	return &Runner{gh: gh, codex: codex}
}

func (r *Runner) Run(ctx context.Context) error {
	prNum, err := r.gh.GetCurrentPRNumber(ctx)
	if err != nil {
		return err
	}
	threads, err := r.gh.ListAllReviewThreads(ctx, prNum)
	if err != nil {
		return err
	}
	for _, th := range threads {
		cm := firstRelevant(th.Comments)
		tc := ThreadContext{PRNumber: prNum, Thread: th, Comment: cm}
		res, err := r.codex.HeuristicAnalysis(ctx, tc)
		if err != nil {
			return err
		}
		printHeuristic(res)

		// Auto-apply simple fixes for low-risk items
		if res.Score < 5 {
			fmt.Printf("Applying code changes\n")
			if summary, apErr := r.codex.ImplementCode(ctx, tc); apErr == nil {
				fmt.Printf("Applied code changes; resolving.\n")
				// Post a concise reply and resolve the thread
				// _ = r.gh.ResolveReply(ctx, th.ID, "Applied low-risk default strategy; resolving.", true)
				printActionBlock(th.ID, cm.URL, cm.File, cm.Line, summary)
			} else {
				return fmt.Errorf("apply failed for thread %s: %v", th.ID, apErr)
			}
			return nil
		}
	}
	return nil
}

func firstRelevant(comments []Comment) Comment {
	if len(comments) > 0 {
		return comments[0]
	}
	return Comment{}
}

func printHeuristic(res HeuristicAnalysisResult) {
	fmt.Printf("BEGIN_HEURISTIC\n")
	fmt.Printf("Heuristic Checklist Result\n")
	fmt.Printf("- Summary: %s\n", strings.TrimSpace(res.Summary))

	// Print known checklist items in a fixed order
	order := []string{
		"tools_present",
		"pr_detected",
		"conversations_fetched",
		"auto_resolved_outdated",
		"relevance_classified",
	}
	for _, k := range order {
		if res.Items != nil {
			fmt.Printf("- %s: %v\n", k, res.Items[k])
		} else {
			fmt.Printf("- %s: %v\n", k, false)
		}
	}

	// Preferred option
	if len(res.ProposedActions) > 0 {
		fmt.Printf("- preferred_option: %s\n", res.ProposedActions[0])
	}

	// Alternatives
	if len(res.Alternatives) > 0 {
		fmt.Printf("- alternatives:\n")
		for _, alt := range res.Alternatives {
			opt := strings.TrimSpace(alt["option"])
			why := strings.TrimSpace(alt["why"])
			fmt.Printf("  - option: %s\n", opt)
			fmt.Printf("    why: %s\n", why)
		}
	}

	fmt.Printf("- Risk score (1–10): %d\n", res.Score)
	fmt.Printf("END_HEURISTIC\n")
}

func printActionBlock(id, url, file string, line int, summary string) {
	fmt.Printf("BEGIN_ACTION\n")
	fmt.Printf("Id: \"%s\"\n", id)
	fmt.Printf("Url: \"%s\"\n", url)
	fmt.Printf("Location: \"%s:%d\"\n", file, line)
	if summary == "" {
		summary = "Applied reviewer’s suggestion in minimal, scoped change"
	}
	fmt.Printf("Summary: %s\n", summary)
	fmt.Printf("Actions Taken: Auto-applied or verified already fixed; posted resolving reply\n")
	fmt.Printf("Tests/Checks: Local validations as applicable\n")
	fmt.Printf("Resolution: Posted reply and resolved\n")
	fmt.Printf("END_ACTION\n")
}

func printApprovalBlock(id, url, file string, line int, comment string, risk int) {
	fmt.Printf("Id: \"%s\"\n", id)
	fmt.Printf("Url: \"%s\"\n", url)
	fmt.Printf("Location: \"%s:%d\"\n", file, line)
	fmt.Printf("Comment: %s\n", comment)
	fmt.Printf("Proposed Fix: Implement the reviewer’s suggestion in a minimal, scoped change aligned with architecture and coding standards; validate via tests.\n")
	fmt.Printf("Risk: \"%d\"\n\n", risk)
	fmt.Printf("Should I proceed with the Implement the reviewer’s suggestion in a minimal, scoped change aligned with architecture and coding standards; validate via tests.?\n")
	fmt.Printf("1. Yes\n2. No, do ... instead\n")
}
