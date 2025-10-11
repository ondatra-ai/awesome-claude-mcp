package app

import (
	"context"
	"log/slog"
)

// MultiHandler writes to both file and console with different levels.
type MultiHandler struct {
	fileHandler    slog.Handler
	consoleHandler slog.Handler
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.fileHandler.Enabled(ctx, level) || h.consoleHandler.Enabled(ctx, level)
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	// Always write to file
	err := h.fileHandler.Handle(ctx, record)
	if err != nil {
		// Don't fail if file write fails, continue to console
	}

	// Write to console only for info and above
	if record.Level >= slog.LevelInfo {
		err := h.consoleHandler.Handle(ctx, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MultiHandler{
		fileHandler:    h.fileHandler.WithAttrs(attrs),
		consoleHandler: h.consoleHandler.WithAttrs(attrs),
	}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	return &MultiHandler{
		fileHandler:    h.fileHandler.WithGroup(name),
		consoleHandler: h.consoleHandler.WithGroup(name),
	}
}
