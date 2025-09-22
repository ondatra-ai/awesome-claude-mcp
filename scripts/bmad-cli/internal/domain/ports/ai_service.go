package ports

import (
	"context"

	"bmad-cli/internal/domain/models"
)

type AIService interface {
	AnalyzeThread(ctx context.Context, threadContext models.ThreadContext) (models.HeuristicAnalysisResult, error)
	ImplementChanges(ctx context.Context, threadContext models.ThreadContext) (string, error)
}
