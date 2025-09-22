package prtriage

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

const (
	lowRiskThreshold = 5
)

type Runner struct {
	prFetcher      *PRNumberFetcher
	threadsFetcher *ThreadsFetcher
	resolver       *ThreadResolver
	analyzer       *ThreadAnalyzer
	implementer    *ThreadImplementer
}

func NewRunner(
	prFetcher *PRNumberFetcher,
	threadsFetcher *ThreadsFetcher,
	resolver *ThreadResolver,
	analyzer *ThreadAnalyzer,
	implementer *ThreadImplementer,
) *Runner {
	return &Runner{
		prFetcher:      prFetcher,
		threadsFetcher: threadsFetcher,
		resolver:       resolver,
		analyzer:       analyzer,
		implementer:    implementer,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	// 1. Get PR number
	prNum, err := r.prFetcher.Fetch(ctx)
	if err != nil {
		return err
	}

	// 2. Fetch all threads for this PR
	threads, err := r.threadsFetcher.FetchAll(ctx, prNum)
	if err != nil {
		return err
	}

	// 3. Process each thread
	for _, thread := range threads {
		comment := firstRelevant(thread.Comments)

		// Skip outdated threads and resolve them
		if comment.Outdated {
			_ = r.resolver.Resolve(ctx, thread.ID, "This thread resolved as outdated.")

			continue
		}

		// 4. Analyze the thread
		threadCtx := ThreadContext{PRNumber: prNum, Thread: thread, Comment: comment}

		res, err := r.analyzer.Analyze(ctx, threadCtx)
		if err != nil {
			return err
		}

		printHeuristic(res)

		// 5. Auto-apply simple fixes for low-risk items
		if res.Score < lowRiskThreshold {
			slog.Info("Applying code changes")

			// 6. Implement the changes
			summary, err := r.implementer.Implement(ctx, threadCtx)
			if err != nil {
				return fmt.Errorf("apply failed for thread %s: %w", thread.ID, err)
			}

			slog.Info("Applied code changes; resolving.")

			// 7. Resolve the thread with summary
			_ = r.resolver.Resolve(ctx, thread.ID, "Applied low-risk default strategy; resolving.")
			printActionBlock(thread.ID, comment.URL, comment.File, comment.Line, summary)
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
	slog.Info("Analysis complete", "summary", strings.TrimSpace(res.Summary), "score", res.Score)

	if len(res.ProposedActions) > 0 {
		slog.Info("Recommended action", "action", res.ProposedActions[0])
	}
}

func printActionBlock(id, _ string, file string, line int, summary string) {
	if summary == "" {
		summary = "Applied reviewer's suggestion in minimal, scoped change"
	}

	slog.Info("Action completed", "thread", id, "location", fmt.Sprintf("%s:%d", file, line), "summary", summary)
}

// RunPRTriage creates and runs the complete PR triage process
func RunPRTriage(ctx context.Context, engineType string) error {
	// Create AI client
	aiClient, err := CreateAIClient(engineType)
	if err != nil {
		return fmt.Errorf("failed to create AI client: %w", err)
	}

	// Create GitHub client
	ghClient := NewGitHubCLIClient()

	// Create GitHub operation components
	prFetcher := NewPRNumberFetcher(ghClient)
	threadsFetcher := NewThreadsFetcher(ghClient)
	resolver := NewThreadResolver(ghClient)

	// Create AI operation components
	analyzer := NewThreadAnalyzer(aiClient)
	implementer := NewThreadImplementer(aiClient)

	// Create runner with all components
	runner := NewRunner(prFetcher, threadsFetcher, resolver, analyzer, implementer)

	// Run pr-triage logic
	return runner.Run(ctx)
}
