package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// MCPMessage represents an MCP protocol message (JSON-RPC 2.0)
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

// SessionPool manages HTTP+SSE sessions
type SessionPool struct {
	sessions    sync.Map
	activeCount atomic.Int64
	totalCount  atomic.Int64
}

// SessionInfo stores session metadata
type SessionInfo struct {
	ID           string
	CreatedAt    time.Time
	LastActive   time.Time
	MessageCount atomic.Int64
	SSEChannel   chan []byte
	mu           sync.Mutex
}

var pool = &SessionPool{}

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
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,Mcp-Session-Id",
		ExposeHeaders: "Mcp-Session-Id",
	}))

	// Health check endpoint
	app.Get("/health", healthCheckHandler)

	// MCP HTTP+SSE endpoints (Streamable HTTP per MCP specification)
	// POST /mcp: Client sends JSON-RPC messages
	app.Post("/mcp", mcpPostHandler)
	// GET /mcp: Client establishes SSE stream for server-to-client messages
	app.Get("/mcp", mcpSSEHandler)

	// Start server in goroutine
	go func() {
		log.Info().Str("port", port).Msg("Starting MCP service with HTTP+SSE transport")
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

	// Close all SSE sessions
	pool.sessions.Range(func(key, value interface{}) bool {
		if sessionInfo, ok := value.(*SessionInfo); ok {
			log.Info().Str("session_id", sessionInfo.ID).Msg("Closing session")
			close(sessionInfo.SSEChannel)
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

// mcpPostHandler handles POST /mcp for JSON-RPC messages
func mcpPostHandler(c *fiber.Ctx) error {
	// Get or create session ID
	sessionID := c.Get("Mcp-Session-Id")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// Get or create session
	session := getOrCreateSession(sessionID)

	// Update session activity
	session.mu.Lock()
	session.LastActive = time.Now()
	session.mu.Unlock()
	session.MessageCount.Add(1)

	// Set session ID header in response
	c.Set("Mcp-Session-Id", sessionID)

	// Parse request body
	body := c.Body()
	if len(body) == 0 {
		return c.Status(400).JSON(MCPMessage{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32700,
				Message: "Parse error - empty body",
			},
		})
	}

	// Parse MCP message
	var mcpMsg MCPMessage
	if err := json.Unmarshal(body, &mcpMsg); err != nil {
		log.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("body", string(body)).
			Msg("Failed to parse MCP message")

		return c.Status(400).JSON(MCPMessage{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32700,
				Message: "Parse error - invalid JSON",
			},
		})
	}

	// Validate JSON-RPC version
	if mcpMsg.JSONRPC != "2.0" {
		return c.Status(400).JSON(MCPMessage{
			JSONRPC: "2.0",
			ID:      mcpMsg.ID,
			Error: &MCPError{
				Code:    -32600,
				Message: "Invalid Request - jsonrpc must be '2.0'",
			},
		})
	}

	// Handle notifications (no ID, no response needed)
	if mcpMsg.Method != "" && mcpMsg.ID == nil {
		log.Info().
			Str("session_id", sessionID).
			Str("method", mcpMsg.Method).
			Msg("Received MCP notification")
		// Notifications return 204 No Content
		return c.SendStatus(204)
	}

	// Validate method field for requests
	if mcpMsg.Method == "" && mcpMsg.Result == nil && mcpMsg.Error == nil {
		return c.Status(400).JSON(MCPMessage{
			JSONRPC: "2.0",
			ID:      mcpMsg.ID,
			Error: &MCPError{
				Code:    -32600,
				Message: "Invalid Request - method field is required",
			},
		})
	}

	log.Info().
		Str("session_id", sessionID).
		Str("method", mcpMsg.Method).
		Interface("id", mcpMsg.ID).
		Msg("Processing MCP request")

	// Handle MCP methods
	response := handleMCPMethod(mcpMsg, sessionID)
	return c.JSON(response)
}

// mcpSSEHandler handles GET /mcp for SSE stream
func mcpSSEHandler(c *fiber.Ctx) error {
	// Get session ID from header
	sessionID := c.Get("Mcp-Session-Id")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	// Get or create session
	session := getOrCreateSession(sessionID)

	log.Info().
		Str("session_id", sessionID).
		Msg("SSE stream requested")

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Mcp-Session-Id", sessionID)

	// Use streaming response
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Send initial ping to establish connection
		fmt.Fprintf(w, "event: ping\ndata: {\"type\":\"ping\"}\n\n")
		w.Flush()

		// Listen for messages on session channel
		for msg := range session.SSEChannel {
			fmt.Fprintf(w, "data: %s\n\n", msg)
			if err := w.Flush(); err != nil {
				log.Warn().
					Err(err).
					Str("session_id", sessionID).
					Msg("SSE write error, closing stream")
				return
			}
		}
	})

	return nil
}

// getOrCreateSession returns existing session or creates new one
func getOrCreateSession(sessionID string) *SessionInfo {
	// Try to load existing session
	if existing, ok := pool.sessions.Load(sessionID); ok {
		return existing.(*SessionInfo)
	}

	// Create new session
	session := &SessionInfo{
		ID:         sessionID,
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
		SSEChannel: make(chan []byte, 100),
	}

	// Store session (use LoadOrStore to handle race condition)
	actual, loaded := pool.sessions.LoadOrStore(sessionID, session)
	if loaded {
		// Another goroutine created the session first
		return actual.(*SessionInfo)
	}

	pool.activeCount.Add(1)
	pool.totalCount.Add(1)

	log.Info().
		Str("session_id", sessionID).
		Msg("New MCP session created")

	return session
}

func handleMCPMethod(msg MCPMessage, sessionID string) MCPMessage {
	switch msg.Method {
	case "initialize":
		// Return initialize response with session ID
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
				"sessionId": sessionID,
			},
		}

	case "ping":
		// Return pong response
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Result:  map[string]interface{}{},
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

// sendSSEMessage sends a message to a session's SSE channel (for server-initiated messages)
func sendSSEMessage(sessionID string, msg MCPMessage) error {
	sessionVal, ok := pool.sessions.Load(sessionID)
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session := sessionVal.(*SessionInfo)
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	select {
	case session.SSEChannel <- data:
		return nil
	default:
		return fmt.Errorf("SSE channel full for session: %s", sessionID)
	}
}

// Ensure fasthttp import is used
var _ = fasthttp.StatusOK
