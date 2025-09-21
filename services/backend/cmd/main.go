package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// VersionResponse represents the API version information
type VersionResponse struct {
	Version string `json:"version"`
}

// HealthResponse represents the health check status of the service
type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}

func setupRoutes(app *fiber.App) {
	// Version endpoint - returns application version
	app.Get("/version", func(c *fiber.Ctx) error {
		log.Info().
			Str("endpoint", "/version").
			Str("method", "GET").
			Str("user_agent", c.Get("User-Agent")).
			Msg("Version endpoint accessed")

		return c.JSON(VersionResponse{
			Version: "1.0.0",
		})
	})

	// Health check endpoint - returns service health status
	app.Get("/health", func(c *fiber.Ctx) error {
		log.Info().
			Str("endpoint", "/health").
			Str("method", "GET").
			Msg("Health check endpoint accessed")

		return c.JSON(HealthResponse{
			Status:    "healthy",
			Service:   "MCP Google Docs Editor - Backend",
			Timestamp: time.Now(),
		})
	})
}

func setupLogger() {
	// Configure zerolog for structured JSON logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	// Set log level from environment variable, default to info
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func createFiberApp() *fiber.App {
	// Production-optimized Fiber setup per tech-stack.md
	app := fiber.New(fiber.Config{
		Prefork:       false, // Disable for development
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "MCP Google Docs Editor - Backend",
		AppName:       "MCP Google Docs Editor - Backend",
	})

	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		log.Fatal().Msg("CORS_ALLOWED_ORIGINS environment variable is required")
	}

	// Middleware setup
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins, // Comma separated list of allowed origins
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))

	app.Use(recover.New())

	return app
}

func gracefulShutdown(app *fiber.App) {
	// Create channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		log.Info().
			Str("port", port).
			Str("service", "backend").
			Msg("Backend server starting")

		if err := app.Listen(":" + port); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for signal
	<-sigChan

	log.Info().Msg("Graceful shutdown initiated")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		log.Info().Msg("Server shutdown complete")
	}
}

func main() {
	// Setup structured logging
	setupLogger()

	log.Info().
		Str("service", "backend").
		Str("version", "1.0.0").
		Msg("MCP Google Docs Editor Backend starting")

	// Create and configure Fiber app
	app := createFiberApp()

	// Setup routes
	setupRoutes(app)

	// Start server with graceful shutdown
	gracefulShutdown(app)

	log.Info().Msg("Application terminated")
}
