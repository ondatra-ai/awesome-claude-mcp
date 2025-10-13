package app

import (
	"log"
	"log/slog"
	"os"
)

const (
	fileModeReadWrite = 0644 // Standard file permission for read/write files
	fileModeDirectory = 0755 // Standard directory permission
)

func configureLogging() {
	log.SetFlags(0)

	// Ensure tmp directory exists
	err := os.MkdirAll("./tmp", fileModeDirectory)
	if err != nil {
		log.Println("Warning: failed to create tmp directory:", err)
	}

	// Open log file for JSON output (all levels)
	logFile, err := os.OpenFile("./tmp/bmad-cli.log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, fileModeReadWrite)
	if err != nil {
		log.Println("Warning: failed to open log file:", err)
		// Fallback to console only
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
			ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
				if attr.Key == slog.TimeKey {
					return slog.Attr{}
				}

				if attr.Key == slog.LevelKey {
					level := attr.Value.String()
					switch level {
					case "INFO":
						return slog.String("", "ℹ️")
					case "WARN":
						return slog.String("", "⚠️")
					case "ERROR":
						return slog.String("", "❌")
					case "DEBUG":
						return slog.String("", "🐛")
					default:
						return slog.String("", level)
					}
				}

				return attr
			},
		}
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))

		return
	}

	// Create JSON handler for file (all levels including debug)
	fileOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // Log everything to file
	}
	fileHandler := slog.NewJSONHandler(logFile, fileOpts)

	// Create text handler for console (info and above only)
	consoleOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Only info, warn, error to console
	}
	consoleHandler := slog.NewTextHandler(os.Stdout, consoleOpts)

	// Create multi-handler that writes to both file and console
	multiHandler := &MultiHandler{
		fileHandler:    fileHandler,
		consoleHandler: consoleHandler,
	}

	slog.SetDefault(slog.New(multiHandler))
}
