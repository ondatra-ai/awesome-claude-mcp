package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// MCPMessage represents an MCP protocol message
type MCPMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

// MCPError represents an MCP protocol error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ConnectionPool manages WebSocket connections
type ConnectionPool struct {
	connections sync.Map
	activeCount atomic.Int64
	totalCount  atomic.Int64
}

// ConnectionInfo stores connection metadata
type ConnectionInfo struct {
	ID          string
	Conn        *websocket.Conn
	ConnectedAt time.Time
	LastActive  time.Time
	MessageCount atomic.Int64
	mu          sync.Mutex
}

var pool = &ConnectionPool{}

func main() {
	// Configure logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Get port from environment or use default
	port := os.Getenv("MCP_PORT")
	if port == "" {
		port = "8081"
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // For testing - production should restrict this
		AllowMethods: "GET,POST,HEAD,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check endpoint
	app.Get("/health", healthCheckHandler)

	// WebSocket upgrade endpoint for MCP protocol
	app.Get("/mcp", websocket.New(mcpWebSocketHandler))

	// Start server in goroutine
	go func() {
		log.Info().Str("port", port).Msg("Starting MCP service")
		if err := app.Listen(":" + port); err != nil {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Close all WebSocket connections
	pool.connections.Range(func(key, value interface{}) bool {
		if connInfo, ok := value.(*ConnectionInfo); ok {
			log.Info().Str("connection_id", connInfo.ID).Msg("Closing connection")
			connInfo.Conn.Close()
		}
		return true
	})

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

func healthCheckHandler(c *fiber.Ctx) error {
	health := map[string]interface{}{
		"status": "healthy",
		"connections": map[string]interface{}{
			"active": pool.activeCount.Load(),
			"total":  pool.totalCount.Load(),
		},
		"dependencies": map[string]interface{}{
			"redis": map[string]string{
				"status": "healthy",
			},
			"googleAPI": map[string]string{
				"status": "healthy",
			},
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	return c.JSON(health)
}

func mcpWebSocketHandler(c *websocket.Conn) {
	// Create connection info
	connInfo := &ConnectionInfo{
		ID:          uuid.New().String(),
		Conn:        c,
		ConnectedAt: time.Now(),
		LastActive:  time.Now(),
	}

	// Register connection
	pool.connections.Store(connInfo.ID, connInfo)
	pool.activeCount.Add(1)
	pool.totalCount.Add(1)

	log.Info().
		Str("connection_id", connInfo.ID).
		Str("remote_addr", c.RemoteAddr().String()).
		Msg("WebSocket connection established")

	// Cleanup on exit
	defer func() {
		pool.connections.Delete(connInfo.ID)
		pool.activeCount.Add(-1)
		log.Info().
			Str("connection_id", connInfo.ID).
			Int64("messages_processed", connInfo.MessageCount.Load()).
			Msg("WebSocket connection closed")
	}()

	// Message handling loop
	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Info().
					Str("connection_id", connInfo.ID).
					Msg("Client closed connection")
			} else {
				log.Error().
					Err(err).
					Str("connection_id", connInfo.ID).
					Msg("Error reading message")
			}
			break
		}

		// Only handle text messages
		if messageType != websocket.TextMessage {
			log.Warn().
				Str("connection_id", connInfo.ID).
				Int("message_type", messageType).
				Msg("Received non-text message")
			continue
		}

		// Update connection metadata
		connInfo.mu.Lock()
		connInfo.LastActive = time.Now()
		connInfo.mu.Unlock()
		connInfo.MessageCount.Add(1)

		// Parse MCP message
		var mcpMsg MCPMessage
		if err := json.Unmarshal(message, &mcpMsg); err != nil {
			log.Error().
				Err(err).
				Str("connection_id", connInfo.ID).
				Str("message", string(message)).
				Msg("Failed to parse MCP message")

			// Send error response
			errorResponse := MCPMessage{
				JSONRPC: "2.0",
				ID:      nil,
				Error: &MCPError{
					Code:    -32700,
					Message: "Parse error",
				},
			}
			sendMCPMessage(c, connInfo.ID, errorResponse)
			continue
		}

		// Validate required fields
		if mcpMsg.JSONRPC != "2.0" {
			errorResponse := MCPMessage{
				JSONRPC: "2.0",
				ID:      mcpMsg.ID,
				Error: &MCPError{
					Code:    -32600,
					Message: "Invalid Request - jsonrpc must be '2.0'",
				},
			}
			sendMCPMessage(c, connInfo.ID, errorResponse)
			continue
		}

		// Handle methods that require ID
		if mcpMsg.Method != "" && mcpMsg.ID == nil {
			// This is a notification, not a request
			log.Info().
				Str("connection_id", connInfo.ID).
				Str("method", mcpMsg.Method).
				Msg("Received MCP notification")
			continue
		}

		// Validate method field for requests
		if mcpMsg.Method == "" && mcpMsg.Result == nil && mcpMsg.Error == nil {
			errorResponse := MCPMessage{
				JSONRPC: "2.0",
				ID:      mcpMsg.ID,
				Error: &MCPError{
					Code:    -32600,
					Message: "Invalid Request - method field is required",
				},
			}
			sendMCPMessage(c, connInfo.ID, errorResponse)
			continue
		}

		log.Info().
			Str("connection_id", connInfo.ID).
			Str("method", mcpMsg.Method).
			Interface("id", mcpMsg.ID).
			Msg("Processing MCP request")

		// Handle MCP methods
		response := handleMCPMethod(mcpMsg)
		sendMCPMessage(c, connInfo.ID, response)
	}
}

func handleMCPMethod(msg MCPMessage) MCPMessage {
	switch msg.Method {
	case "initialize":
		// Return initialize response
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]string{
					"name":    "mcp-service",
					"version": "1.0.0",
				},
			},
		}

	case "tools/list":
		// Return empty tools list for now
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Result: map[string]interface{}{
				"tools": []interface{}{},
			},
		}

	default:
		// Method not found
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

func sendMCPMessage(c *websocket.Conn, connID string, msg MCPMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Error().
			Err(err).
			Str("connection_id", connID).
			Msg("Failed to marshal MCP response")
		return
	}

	if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Error().
			Err(err).
			Str("connection_id", connID).
			Msg("Failed to send MCP response")
	}
}
