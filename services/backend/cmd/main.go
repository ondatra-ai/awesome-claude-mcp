package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Create Fiber app with production config
	app := fiber.New(fiber.Config{
		Prefork:       false, // Disable for development
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "MCP-Google-Docs-Backend",
		AppName:       "MCP Google Docs Editor - Backend",
	})

	// Add middleware
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	// Routes
	app.Get("/version", getVersion)
	app.Get("/health", getHealth)

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Backend server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getVersion returns the current API version
func getVersion(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"version": "1.0.0",
	})
}

// getHealth returns service health status
func getHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "backend",
	})
}
