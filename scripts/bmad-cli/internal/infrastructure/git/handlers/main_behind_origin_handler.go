package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"bmad-cli/internal/domain/ports"
)

// MainBehindOriginHandler checks if main is behind origin/main
type MainBehindOriginHandler struct {
	BaseBranchHandler
	gitService ports.GitPort
}

// NewMainBehindOriginHandler creates a new handler
func NewMainBehindOriginHandler(gitService ports.GitPort) *MainBehindOriginHandler {
	return &MainBehindOriginHandler{
		gitService: gitService,
	}
}

// Handle checks if main branch needs to be updated
func (h *MainBehindOriginHandler) Handle(ctx context.Context, branchCtx *BranchContext) error {
	// Only check if currently on main
	if branchCtx.CurrentBranch != "main" {
		slog.Debug("Not on main branch, skipping origin check")
		return h.callNext(ctx, branchCtx)
	}

	slog.Debug("Checking if main is behind origin/main")

	isBehind, err := h.gitService.IsMainBehindOrigin(ctx)
	if err != nil {
		slog.Error("Failed to check if main is behind origin", "error", err)
		return fmt.Errorf("failed to check if main is behind origin: %w", err)
	}

	if isBehind {
		slog.Error("Main branch is behind origin/main")
		return fmt.Errorf("main branch is behind origin/main - please pull the latest changes first")
	}

	slog.Debug("Main branch is up to date with origin")
	return h.callNext(ctx, branchCtx)
}
