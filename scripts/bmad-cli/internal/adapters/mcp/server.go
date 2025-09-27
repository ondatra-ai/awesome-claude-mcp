package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Server represents the MCP server
type Server struct {
	connManager  *ConnectionManager
	capabilities MCPCapabilities
	serverInfo   ServerInfo
	tools        map[string]*Tool
	mu           sync.RWMutex
	initialized  bool
}

// NewServer creates a new MCP server
func NewServer() *Server {
	s := &Server{
		connManager: NewConnectionManager(),
		capabilities: MCPCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
		},
		serverInfo: ServerInfo{
			Name:    "bmad-mcp-server",
			Version: "1.0.0",
		},
		tools: make(map[string]*Tool),
	}

	// Register default tools
	s.registerDefaultTools()

	return s
}

// RegisterTool registers a new tool with the server
func (s *Server) RegisterTool(tool *Tool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools[tool.Name] = tool
	log.Printf("Registered MCP tool: %s", tool.Name)
}

// HandleHTTP handles HTTP requests for MCP endpoints
func (s *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/mcp":
		s.connManager.HandleWebSocket(w, r, s)
	case "/mcp/status":
		s.handleStatus(w, r)
	case "/mcp/info":
		s.handleInfo(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleStatus returns server status
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := map[string]interface{}{
		"status":           "healthy",
		"connections":      s.connManager.GetConnectionCount(),
		"protocol_version": "2024-11-05",
		"capabilities":     s.capabilities,
		"tools_count":      len(s.tools),
	}

	json.NewEncoder(w).Encode(status)
}

// handleInfo returns server information
func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.serverInfo)
}

// HandleMessage handles incoming MCP messages
func (s *Server) HandleMessage(conn *Connection, msg *MCPMessage) {
	switch msg.Method {
	case "initialize":
		s.handleInitialize(conn, msg)
	case "initialized":
		s.handleInitialized(conn, msg)
	case "ping":
		s.handlePing(conn, msg)
	case "tools/list":
		s.handleToolsList(conn, msg)
	case "tools/call":
		s.handleToolsCall(conn, msg)
	default:
		conn.sendError(-32601, "Method not found", fmt.Sprintf("Unknown method: %s", msg.Method), msg.ID)
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(conn *Connection, msg *MCPMessage) {
	var params InitializeParams
	if msg.Params != nil {
		paramsBytes, err := json.Marshal(msg.Params)
		if err != nil {
			conn.sendError(-32602, "Invalid params", err.Error(), msg.ID)
			return
		}
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			conn.sendError(-32602, "Invalid params", err.Error(), msg.ID)
			return
		}
	}

	log.Printf("MCP Initialize request from client: %s %s", params.ClientInfo.Name, params.ClientInfo.Version)

	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities:    s.capabilities,
		ServerInfo:      s.serverInfo,
	}

	response := &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  result,
	}

	conn.SendMessage(response)
}

// handleInitialized handles the initialized notification
func (s *Server) handleInitialized(conn *Connection, msg *MCPMessage) {
	s.mu.Lock()
	s.initialized = true
	conn.authorized = true
	s.mu.Unlock()

	log.Printf("MCP connection %s is now initialized", conn.id)
}

// handlePing handles ping requests
func (s *Server) handlePing(conn *Connection, msg *MCPMessage) {
	response := &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  PongResult{},
	}
	conn.SendMessage(response)
}

// handleToolsList handles tools/list requests
func (s *Server) handleToolsList(conn *Connection, msg *MCPMessage) {
	if !conn.authorized {
		conn.sendError(-32002, "Server not initialized", "Must call initialize first", msg.ID)
		return
	}

	s.mu.RLock()
	tools := make([]*Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}
	s.mu.RUnlock()

	result := map[string]interface{}{
		"tools": tools,
	}

	response := &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  result,
	}

	conn.SendMessage(response)
}

// handleToolsCall handles tools/call requests
func (s *Server) handleToolsCall(conn *Connection, msg *MCPMessage) {
	if !conn.authorized {
		conn.sendError(-32002, "Server not initialized", "Must call initialize first", msg.ID)
		return
	}

	var params ToolCallParams
	if msg.Params != nil {
		paramsBytes, err := json.Marshal(msg.Params)
		if err != nil {
			conn.sendError(-32602, "Invalid params", err.Error(), msg.ID)
			return
		}
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			conn.sendError(-32602, "Invalid params", err.Error(), msg.ID)
			return
		}
	}

	s.mu.RLock()
	tool, exists := s.tools[params.Name]
	s.mu.RUnlock()

	if !exists {
		conn.sendError(-32602, "Invalid params", fmt.Sprintf("Unknown tool: %s", params.Name), msg.ID)
		return
	}

	// Execute the tool (for now, just return a success message)
	result := s.executeTool(tool, params.Arguments)

	response := &MCPMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
		Result:  result,
	}

	conn.SendMessage(response)
}

// executeTool executes a tool and returns the result
func (s *Server) executeTool(tool *Tool, arguments json.RawMessage) ToolResult {
	// For now, return a simple success message
	// In a real implementation, this would call the actual tool logic
	content := []interface{}{
		TextContent{
			Type: "text",
			Text: fmt.Sprintf("Tool '%s' executed successfully with arguments: %s", tool.Name, string(arguments)),
		},
	}

	return ToolResult{
		Content: content,
		IsError: false,
	}
}

// Start starts the MCP server on the specified address
func (s *Server) Start(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HandleHTTP)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("Starting MCP server on %s", addr)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("MCP server error: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	return server.Shutdown(context.Background())
}

// registerDefaultTools registers default tools for demonstration
func (s *Server) registerDefaultTools() {
	// Example tool: echo
	echoTool := &Tool{
		Name:        "echo",
		Description: "Echo back the input message",
		InputSchema: ToolSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "The message to echo back",
				},
			},
			Required: []string{"message"},
		},
	}
	s.RegisterTool(echoTool)

	// Example tool: status
	statusTool := &Tool{
		Name:        "status",
		Description: "Get server status information",
		InputSchema: ToolSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
	s.RegisterTool(statusTool)
}
