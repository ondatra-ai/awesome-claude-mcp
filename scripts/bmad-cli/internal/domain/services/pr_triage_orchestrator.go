package services

import (
	"context"
	"fmt"
	"strings"

	"bmad-cli/internal/domain/models"
	"bmad-cli/internal/domain/ports"
	"bmad-cli/internal/infrastructure/logging"
)

const lowRiskThreshold = 5

type PRTriageOrchestrator struct {
	github ports.GitHubService
	ai     ports.AIService
	logger logging.Logger
}

func NewPRTriageOrchestrator(
	github ports.GitHubService,
	ai ports.AIService,
	logger logging.Logger,
) *PRTriageOrchestrator {
	return &PRTriageOrchestrator{
		github: github,
		ai:     ai,
		logger: logger,
	}
}

func (o *PRTriageOrchestrator) Run(ctx context.Context) error {
	prNum, err := o.github.GetPRNumber(ctx)
	if err != nil {
		return err
	}

	threads, err := o.github.FetchThreads(ctx, prNum)
	if err != nil {
		return err
	}

	for _, thread := range threads {
		comment := o.firstRelevantComment(thread.Comments)

		if comment.Outdated {
			_ = o.github.ResolveThread(ctx, thread.ID, "This thread resolved as outdated.")
			continue
		}

		threadCtx := models.ThreadContext{
			PRNumber: prNum,
			Thread:   thread,
			Comment:  comment,
		}

		result, err := o.ai.AnalyzeThread(ctx, threadCtx)
		if err != nil {
			return err
		}

		o.printHeuristic(result)

		if result.Score < lowRiskThreshold {
			o.logger.Info("Applying code changes")

			summary, err := o.ai.ImplementChanges(ctx, threadCtx)
			if err != nil {
				return fmt.Errorf("apply failed for thread %s: %w", thread.ID, err)
			}

			o.logger.Info("Applied code changes; resolving.")

			_ = o.github.ResolveThread(ctx, thread.ID, "Applied low-risk default strategy; resolving.")
			o.printActionBlock(thread.ID, comment.URL, comment.File, comment.Line, summary)
		}
	}

	return nil
}

func (o *PRTriageOrchestrator) firstRelevantComment(comments []models.Comment) models.Comment {
	if len(comments) > 0 {
		return comments[0]
	}
	return models.Comment{}
}

func (o *PRTriageOrchestrator) printHeuristic(result models.HeuristicAnalysisResult) {
	o.logger.Info("Analysis complete", "summary", strings.TrimSpace(result.Summary), "score", result.Score)

	if len(result.ProposedActions) > 0 {
		o.logger.Info("Recommended action", "action", result.ProposedActions[0])
	}
}

func (o *PRTriageOrchestrator) printActionBlock(id, _ string, file string, line int, summary string) {
	if summary == "" {
		summary = "Applied reviewer's suggestion in minimal, scoped change"
	}

	o.logger.Info("Action completed", "thread", id, "location", fmt.Sprintf("%s:%d", file, line), "summary", summary)
}
