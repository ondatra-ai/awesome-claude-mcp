package commands

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"bmad-cli/internal/adapters/ai"
	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/config"
)

const lowRiskThreshold = 5

type PRTriageCommand struct {
	github          ports.GitHubService
	threadProcessor *ai.ThreadProcessor
	config          *config.ViperConfig
}

func NewPRTriageCommand(
	github ports.GitHubService,
	threadProcessor *ai.ThreadProcessor,
	config *config.ViperConfig,
) *PRTriageCommand {
	return &PRTriageCommand{
		github:          github,
		threadProcessor: threadProcessor,
		config:          config,
	}
}

func (c *PRTriageCommand) Execute(ctx context.Context) error {
	prNum, err := c.github.GetPRNumber(ctx)
	if err != nil {
		return err
	}

	threads, err := c.github.FetchThreads(ctx, prNum)
	if err != nil {
		return err
	}

	for _, thread := range threads {
		comment := c.firstRelevantComment(thread.Comments)

		if comment.Outdated {
			_ = c.github.ResolveThread(ctx, thread.ID, "This thread resolved as outdated.")
			continue
		}

		threadCtx := models.ThreadContext{
			PRNumber: prNum,
			Thread:   thread,
			Comment:  comment,
		}

		result, err := c.analyzeThread(ctx, threadCtx)
		if err != nil {
			return err
		}

		c.printHeuristic(result)

		if result.Score < lowRiskThreshold {
			slog.Info("Applying code changes")

			summary, err := c.implementChanges(ctx, threadCtx)
			if err != nil {
				return fmt.Errorf("apply failed for thread %s: %w", thread.ID, err)
			}

			slog.Info("Applied code changes; resolving.")

			_ = c.github.ResolveThread(ctx, thread.ID, "Applied low-risk default strategy; resolving.")
			c.printActionBlock(thread.ID, comment.URL, comment.File, comment.Line, summary)
		}
	}

	return nil
}

func (c *PRTriageCommand) analyzeThread(ctx context.Context, threadContext models.ThreadContext) (models.HeuristicAnalysisResult, error) {
	return c.threadProcessor.AnalyzeThread(ctx, threadContext)
}

func (c *PRTriageCommand) implementChanges(ctx context.Context, threadContext models.ThreadContext) (string, error) {
	return c.threadProcessor.ImplementChanges(ctx, threadContext)
}

func (c *PRTriageCommand) firstRelevantComment(comments []models.Comment) models.Comment {
	if len(comments) > 0 {
		return comments[0]
	}
	return models.Comment{}
}

func (c *PRTriageCommand) printHeuristic(result models.HeuristicAnalysisResult) {
	slog.Info("Analysis complete", "summary", strings.TrimSpace(result.Summary), "score", result.Score)

	if len(result.ProposedActions) > 0 {
		slog.Info("Recommended action", "action", result.ProposedActions[0])
	}
}

func (c *PRTriageCommand) printActionBlock(id, _ string, file string, line int, summary string) {
	if summary == "" {
		summary = "Applied reviewer's suggestion in minimal, scoped change"
	}

	slog.Info("Action completed", "thread", id, "location", fmt.Sprintf("%s:%d", file, line), "summary", summary)
}
