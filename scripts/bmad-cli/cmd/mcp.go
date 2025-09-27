package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bmad-cli/internal/adapters/mcp"
	"bmad-cli/internal/app"
	"github.com/spf13/cobra"
)

// NewMCPCommand creates a new MCP command
func NewMCPCommand(container *app.Container) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start MCP server for Claude communication",
		Long: `Start the Model Context Protocol (MCP) server that allows Claude AI
to communicate with this service. The server provides WebSocket and HTTP endpoints
for MCP protocol communication.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer(port)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the MCP server on")

	return cmd
}

// runMCPServer starts and runs the MCP server
func runMCPServer(port int) error {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create and configure MCP server
	server := mcp.NewServer()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	addr := fmt.Sprintf(":%d", port)
	errChan := make(chan error, 1)

	go func() {
		log.Printf("Starting MCP server on port %d", port)
		log.Printf("WebSocket endpoint: ws://localhost:%d/mcp", port)
		log.Printf("HTTP status endpoint: http://localhost:%d/mcp/status", port)
		log.Printf("HTTP info endpoint: http://localhost:%d/mcp/info", port)

		if err := server.Start(ctx, addr); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-sigChan:
		log.Println("Received shutdown signal, gracefully shutting down...")
		cancel()
		return nil
	case err := <-errChan:
		log.Printf("Server error: %v", err)
		return err
	case <-ctx.Done():
		return nil
	}
}
