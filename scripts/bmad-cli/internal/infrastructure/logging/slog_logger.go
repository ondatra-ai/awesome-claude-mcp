package logging

import (
	"log/slog"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
	return &SlogLogger{
		logger: slog.Default(),
	}
}

func (l *SlogLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Error(msg string, err error, args ...interface{}) {
	allArgs := append([]interface{}{"error", err}, args...)
	l.logger.Error(msg, allArgs...)
}

func (l *SlogLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}
