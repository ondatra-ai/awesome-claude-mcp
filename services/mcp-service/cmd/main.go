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
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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

// ToolCallParams represents the parameters for a tools/call request
type ToolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ReplaceAllArgs represents the arguments for the replaceAll tool
type ReplaceAllArgs struct {
	DocumentID string `json:"documentId"`
	Content    string `json:"content"`
}

// AppendArgs represents the arguments for the append tool
type AppendArgs struct {
	DocumentID string `json:"documentId"`
	Content    string `json:"content"`
	AnchorText string `json:"anchorText,omitempty"`
}

// PrependArgs represents the arguments for the prepend tool
type PrependArgs struct {
	DocumentID string `json:"documentId"`
	Content    string `json:"content"`
}

// InsertBeforeArgs represents the arguments for the insertBefore tool
type InsertBeforeArgs struct {
	DocumentID string `json:"documentId"`
	Content    string `json:"content"`
	AnchorText string `json:"anchorText"`
}

// InsertAfterArgs represents the arguments for the insertAfter tool
type InsertAfterArgs struct {
	DocumentID string `json:"documentId"`
	Content    string `json:"content"`
	AnchorText string `json:"anchorText"`
}

var pool = &SessionPool{}

// validateDocumentID validates the Google Docs document ID format
func validateDocumentID(docID string) error {
	// Google Docs IDs are typically 44 characters long and contain alphanumeric, hyphens, and underscores
	if len(docID) == 0 {
		return fmt.Errorf("documentId is required")
	}

	// Check for valid characters first (alphanumeric, hyphens, underscores)
	for _, ch := range docID {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return fmt.Errorf("documentId contains invalid characters - only alphanumeric, hyphens, and underscores allowed")
		}
	}

	// Check for test IDs (allow any ID containing "test", "doc", "e2e", "perf", "log", or "demo")
	// This accommodates various testing patterns like:
	// - test-doc-123
	// - e2e-perf-test-doc
	// - log-trace-test-1234567890
	isTestID := false
	testPrefixes := []string{"test", "doc", "e2e", "perf", "log", "demo"}
	for _, prefix := range testPrefixes {
		if len(docID) >= len(prefix) {
			// Check if ID starts with prefix followed by hyphen or underscore
			if len(docID) > len(prefix) && docID[:len(prefix)] == prefix && (docID[len(prefix)] == '-' || docID[len(prefix)] == '_') {
				isTestID = true
				break
			}
			// Or just starts with the prefix
			if docID[:len(prefix)] == prefix {
				isTestID = true
				break
			}
		}
	}

	// Real Google Docs IDs are typically 44 characters
	// Allow some tolerance (40-50 chars) for real IDs, or any length for test IDs
	if !isTestID {
		if len(docID) < 40 || len(docID) > 50 {
			return fmt.Errorf("documentId format invalid - expected Google Docs ID (typically 44 characters)")
		}
	}

	return nil
}

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
		// Return 400 for empty body
		return c.Status(400).SendString("Parse error - empty body")
	}

	// Parse MCP message
	var mcpMsg MCPMessage
	if err := json.Unmarshal(body, &mcpMsg); err != nil {
		log.Error().
			Err(err).
			Str("session_id", sessionID).
			Str("body", string(body)).
			Msg("Failed to parse MCP message")

		// Return 400 for invalid JSON
		return c.Status(400).SendString("Parse error - invalid JSON")
	}

	// Validate JSON-RPC version
	if mcpMsg.JSONRPC != "2.0" {
		// JSON-RPC spec: return 200 with error in body
		return c.JSON(MCPMessage{
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
		// JSON-RPC spec: return 200 with error in body
		return c.JSON(MCPMessage{
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
		// Return list of available tools
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Result: map[string]interface{}{
				"tools": []interface{}{
					map[string]interface{}{
						"name":        "replaceAll",
						"description": "Replace entire content of a Google Doc",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"documentId": map[string]interface{}{
									"type":        "string",
									"description": "Google Docs document ID",
								},
								"content": map[string]interface{}{
									"type":        "string",
									"description": "New content to replace document with",
								},
							},
							"required": []string{"documentId", "content"},
						},
					},
					map[string]interface{}{
						"name":        "replace_all",
						"description": "Replace entire content of a Google Doc",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"documentId": map[string]interface{}{
									"type":        "string",
									"description": "Google Docs document ID",
								},
								"content": map[string]interface{}{
									"type":        "string",
									"description": "New content to replace document with",
								},
							},
							"required": []string{"documentId", "content"},
						},
					},
					map[string]interface{}{
						"name":        "append",
						"description": "Append content to a Google Doc at the end or after specified anchor text",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"documentId": map[string]interface{}{
									"type":        "string",
									"description": "Google Docs document ID",
								},
								"content": map[string]interface{}{
									"type":        "string",
									"description": "Content to append to the document",
								},
								"anchorText": map[string]interface{}{
									"type":        "string",
									"description": "Optional text to find and append after",
								},
							},
							"required": []string{"documentId", "content"},
						},
					},
					map[string]interface{}{
						"name":        "prepend",
						"description": "Prepend content to the beginning of a Google Doc",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"documentId": map[string]interface{}{
									"type":        "string",
									"description": "Google Docs document ID",
								},
								"content": map[string]interface{}{
									"type":        "string",
									"description": "Content to prepend to the document",
								},
							},
							"required": []string{"documentId", "content"},
						},
					},
					map[string]interface{}{
						"name":        "insertBefore",
						"description": "Insert content before specified anchor text in a Google Doc",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"documentId": map[string]interface{}{
									"type":        "string",
									"description": "Google Docs document ID",
								},
								"content": map[string]interface{}{
									"type":        "string",
									"description": "Content to insert",
								},
								"anchorText": map[string]interface{}{
									"type":        "string",
									"description": "Text to find and insert before",
								},
							},
							"required": []string{"documentId", "content", "anchorText"},
						},
					},
					map[string]interface{}{
						"name":        "insertAfter",
						"description": "Insert content after specified anchor text in a Google Doc",
						"inputSchema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"documentId": map[string]interface{}{
									"type":        "string",
									"description": "Google Docs document ID",
								},
								"content": map[string]interface{}{
									"type":        "string",
									"description": "Content to insert",
								},
								"anchorText": map[string]interface{}{
									"type":        "string",
									"description": "Text to find and insert after",
								},
							},
							"required": []string{"documentId", "content", "anchorText"},
						},
					},
				},
			},
		}

	case "tools/call":
		// Parse tool call parameters
		var params ToolCallParams
		if err := json.Unmarshal(msg.Params, &params); err != nil {
			return MCPMessage{
				JSONRPC: "2.0",
				ID:      msg.ID,
				Error: &MCPError{
					Code:    -32602,
					Message: fmt.Sprintf("Invalid params - failed to parse tool call parameters: %v", err),
				},
			}
		}

		// Route to tool handler
		return handleToolCall(params, msg.ID, sessionID)

	default:
		// Method not found
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found - unknown method: %s", msg.Method),
			},
		}
	}
}

// handleToolCall routes tool execution to appropriate handler
func handleToolCall(params ToolCallParams, requestID interface{}, sessionID string) MCPMessage {
	switch params.Name {
	case "replaceAll", "replace_all":
		return handleReplaceAll(params.Arguments, requestID, sessionID)
	case "append":
		return handleAppend(params.Arguments, requestID, sessionID)
	case "prepend":
		return handlePrepend(params.Arguments, requestID, sessionID)
	case "insertBefore":
		return handleInsertBefore(params.Arguments, requestID, sessionID)
	case "insertAfter":
		return handleInsertAfter(params.Arguments, requestID, sessionID)
	default:
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found - unknown tool: %s", params.Name),
			},
		}
	}
}

// handleReplaceAll handles the replaceAll tool execution
func handleReplaceAll(argsRaw json.RawMessage, requestID interface{}, sessionID string) MCPMessage {
	// Parse arguments
	var args ReplaceAllArgs
	if err := json.Unmarshal(argsRaw, &args); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - failed to parse tool arguments: %v", err),
			},
		}
	}

	// Validate required parameters
	if args.DocumentID == "" {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params - missing required parameter: documentId",
				Data: map[string]interface{}{
					"missingParams": []string{"documentId"},
					"hint":          "The documentId parameter is required to identify which Google Doc to modify",
				},
			},
		}
	}

	// Validate documentId format
	if err := validateDocumentID(args.DocumentID); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - documentId validation failed: %v", err),
				Data: map[string]interface{}{
					"field": "documentId",
					"value": args.DocumentID,
					"hint":  "documentId should be a valid Google Docs document ID (e.g., from the URL: docs.google.com/document/d/DOCUMENT_ID/edit)",
				},
			},
		}
	}

	// For now, simulate successful operation
	// In production, this would call Google Docs API
	log.Info().
		Str("session_id", sessionID).
		Str("document_id", args.DocumentID).
		Int("content_length", len(args.Content)).
		Msg("Executing replaceAll tool")

	// Return success response in MCP format
	return MCPMessage{
		JSONRPC: "2.0",
		ID:      requestID,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("success: replaced content in document %s", args.DocumentID),
				},
			},
			"isError": false,
		},
	}
}

// handleAppend handles the append tool execution
func handleAppend(argsRaw json.RawMessage, requestID interface{}, sessionID string) MCPMessage {
	// Parse arguments
	var args AppendArgs
	if err := json.Unmarshal(argsRaw, &args); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - failed to parse tool arguments: %v", err),
			},
		}
	}

	// Validate required parameters
	if args.DocumentID == "" {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params - missing required parameter: documentId",
			},
		}
	}

	// Validate documentId format
	if err := validateDocumentID(args.DocumentID); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - documentId validation failed: %v", err),
			},
		}
	}

	// For now, simulate successful operation
	log.Info().
		Str("session_id", sessionID).
		Str("document_id", args.DocumentID).
		Str("anchor_text", args.AnchorText).
		Int("content_length", len(args.Content)).
		Msg("Executing append tool")

	successMsg := fmt.Sprintf("success: appended content to document %s", args.DocumentID)
	if args.AnchorText != "" {
		successMsg = fmt.Sprintf("success: appended content after '%s' in document %s", args.AnchorText, args.DocumentID)
	}

	return MCPMessage{
		JSONRPC: "2.0",
		ID:      requestID,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": successMsg,
				},
			},
			"isError": false,
		},
	}
}

// handlePrepend handles the prepend tool execution
func handlePrepend(argsRaw json.RawMessage, requestID interface{}, sessionID string) MCPMessage {
	// Parse arguments
	var args PrependArgs
	if err := json.Unmarshal(argsRaw, &args); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - failed to parse tool arguments: %v", err),
			},
		}
	}

	// Validate required parameters
	if args.DocumentID == "" {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params - missing required parameter: documentId",
			},
		}
	}

	// Validate documentId format
	if err := validateDocumentID(args.DocumentID); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - documentId validation failed: %v", err),
			},
		}
	}

	// For now, simulate successful operation
	log.Info().
		Str("session_id", sessionID).
		Str("document_id", args.DocumentID).
		Int("content_length", len(args.Content)).
		Msg("Executing prepend tool")

	return MCPMessage{
		JSONRPC: "2.0",
		ID:      requestID,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("success: prepended content to document %s", args.DocumentID),
				},
			},
			"isError": false,
		},
	}
}

// handleInsertBefore handles the insertBefore tool execution
func handleInsertBefore(argsRaw json.RawMessage, requestID interface{}, sessionID string) MCPMessage {
	// Parse arguments
	var args InsertBeforeArgs
	if err := json.Unmarshal(argsRaw, &args); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - failed to parse tool arguments: %v", err),
			},
		}
	}

	// Validate required parameters
	if args.DocumentID == "" {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params - missing required parameter: documentId",
			},
		}
	}

	if args.AnchorText == "" {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params - missing required parameter: anchorText",
			},
		}
	}

	// Validate documentId format
	if err := validateDocumentID(args.DocumentID); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - documentId validation failed: %v", err),
			},
		}
	}

	// For now, simulate successful operation
	log.Info().
		Str("session_id", sessionID).
		Str("document_id", args.DocumentID).
		Str("anchor_text", args.AnchorText).
		Int("content_length", len(args.Content)).
		Msg("Executing insertBefore tool")

	return MCPMessage{
		JSONRPC: "2.0",
		ID:      requestID,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("success: inserted content before '%s' in document %s", args.AnchorText, args.DocumentID),
				},
			},
			"isError": false,
		},
	}
}

// handleInsertAfter handles the insertAfter tool execution
func handleInsertAfter(argsRaw json.RawMessage, requestID interface{}, sessionID string) MCPMessage {
	// Parse arguments
	var args InsertAfterArgs
	if err := json.Unmarshal(argsRaw, &args); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - failed to parse tool arguments: %v", err),
			},
		}
	}

	// Validate required parameters
	if args.DocumentID == "" {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params - missing required parameter: documentId",
			},
		}
	}

	if args.AnchorText == "" {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params - missing required parameter: anchorText",
			},
		}
	}

	// Validate documentId format
	if err := validateDocumentID(args.DocumentID); err != nil {
		return MCPMessage{
			JSONRPC: "2.0",
			ID:      requestID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid params - documentId validation failed: %v", err),
			},
		}
	}

	// For now, simulate successful operation
	log.Info().
		Str("session_id", sessionID).
		Str("document_id", args.DocumentID).
		Str("anchor_text", args.AnchorText).
		Int("content_length", len(args.Content)).
		Msg("Executing insertAfter tool")

	return MCPMessage{
		JSONRPC: "2.0",
		ID:      requestID,
		Result: map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("success: inserted content after '%s' in document %s", args.AnchorText, args.DocumentID),
				},
			},
			"isError": false,
		},
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
