package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	server := NewServer()

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.serverInfo.Name != "bmad-mcp-server" {
		t.Errorf("Expected server name 'bmad-mcp-server', got '%s'", server.serverInfo.Name)
	}

	if server.serverInfo.Version != "1.0.0" {
		t.Errorf("Expected server version '1.0.0', got '%s'", server.serverInfo.Version)
	}

	if server.capabilities.Tools == nil {
		t.Error("Expected tools capability to be initialized")
	}
}

func TestRegisterTool(t *testing.T) {
	server := NewServer()

	tool := &Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: ToolSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}

	server.RegisterTool(tool)

	server.mu.RLock()
	registeredTool, exists := server.tools["test_tool"]
	server.mu.RUnlock()

	if !exists {
		t.Error("Expected tool to be registered")
	}

	if registeredTool.Name != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got '%s'", registeredTool.Name)
	}
}

func TestHandleStatus(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/mcp/status", nil)
	w := httptest.NewRecorder()

	server.handleStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got '%s'", contentType)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", response["status"])
	}

	if response["protocol_version"] != "2024-11-05" {
		t.Errorf("Expected protocol_version '2024-11-05', got '%v'", response["protocol_version"])
	}
}

func TestHandleInfo(t *testing.T) {
	server := NewServer()

	req := httptest.NewRequest("GET", "/mcp/info", nil)
	w := httptest.NewRecorder()

	server.handleInfo(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var response ServerInfo
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Name != "bmad-mcp-server" {
		t.Errorf("Expected name 'bmad-mcp-server', got '%s'", response.Name)
	}

	if response.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", response.Version)
	}
}

func TestMCPMessageValidation(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		expectError bool
		errorCode   int
	}{
		{
			name:        "valid message",
			message:     `{"jsonrpc": "2.0", "method": "ping", "id": 1}`,
			expectError: false,
		},
		{
			name:        "invalid JSON",
			message:     `{"jsonrpc": "2.0", "method": "ping", "id": 1`,
			expectError: true,
			errorCode:   -32700,
		},
		{
			name:        "invalid JSON-RPC version",
			message:     `{"jsonrpc": "1.0", "method": "ping", "id": 1}`,
			expectError: true,
			errorCode:   -32600,
		},
	}

	server := NewServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock connection for testing
			conn := &Connection{
				id:     "test",
				send:   make(chan []byte, 1),
				server: server,
			}

			// Handle the message
			conn.handleMessage([]byte(tt.message))

			if tt.expectError {
				select {
				case response := <-conn.send:
					var msg MCPMessage
					if err := json.Unmarshal(response, &msg); err != nil {
						t.Fatalf("Failed to unmarshal error response: %v", err)
					}

					if msg.Error == nil {
						t.Error("Expected error response, got none")
					} else if msg.Error.Code != tt.errorCode {
						t.Errorf("Expected error code %d, got %d", tt.errorCode, msg.Error.Code)
					}
				case <-time.After(100 * time.Millisecond):
					t.Error("Expected error response, got timeout")
				}
			} else {
				// For valid messages, we should not get an immediate error
				select {
				case <-conn.send:
					// This is fine, we might get a response
				case <-time.After(100 * time.Millisecond):
					// This is also fine, some methods might not send immediate responses
				}
			}
		})
	}
}

func TestDefaultToolsRegistration(t *testing.T) {
	server := NewServer()
	server.registerDefaultTools()

	server.mu.RLock()
	defer server.mu.RUnlock()

	expectedTools := []string{"echo", "status"}

	for _, toolName := range expectedTools {
		if _, exists := server.tools[toolName]; !exists {
			t.Errorf("Expected default tool '%s' to be registered", toolName)
		}
	}

	if len(server.tools) != len(expectedTools) {
		t.Errorf("Expected %d default tools, got %d", len(expectedTools), len(server.tools))
	}
}

func TestHTTPRouting(t *testing.T) {
	server := NewServer()

	tests := []struct {
		path           string
		expectedStatus int
	}{
		{"/mcp/status", http.StatusOK},
		{"/mcp/info", http.StatusOK},
		{"/nonexistent", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			server.HandleHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d for path %s, got %d", tt.expectedStatus, tt.path, w.Code)
			}
		})
	}
}

func TestConnectionManager(t *testing.T) {
	cm := NewConnectionManager()

	if cm == nil {
		t.Fatal("Expected connection manager to be created, got nil")
	}

	if cm.GetConnectionCount() != 0 {
		t.Errorf("Expected 0 connections initially, got %d", cm.GetConnectionCount())
	}

	if cm.connections == nil {
		t.Error("Expected connections map to be initialized")
	}

	if cm.upgrader.CheckOrigin == nil {
		t.Error("Expected CheckOrigin function to be set")
	}
}

func TestMCPMessageStructures(t *testing.T) {
	// Test MCPMessage JSON marshaling/unmarshaling
	msg := MCPMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "test",
		Params:  map[string]string{"key": "value"},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal MCPMessage: %v", err)
	}

	var unmarshaled MCPMessage
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal MCPMessage: %v", err)
	}

	if unmarshaled.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC '2.0', got '%s'", unmarshaled.JSONRPC)
	}

	if unmarshaled.Method != "test" {
		t.Errorf("Expected method 'test', got '%s'", unmarshaled.Method)
	}
}

func TestToolSchema(t *testing.T) {
	schema := ToolSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"message": map[string]interface{}{
				"type":        "string",
				"description": "A test message",
			},
		},
		Required: []string{"message"},
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal ToolSchema: %v", err)
	}

	var unmarshaled ToolSchema
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ToolSchema: %v", err)
	}

	if unmarshaled.Type != "object" {
		t.Errorf("Expected type 'object', got '%s'", unmarshaled.Type)
	}

	if len(unmarshaled.Required) != 1 || unmarshaled.Required[0] != "message" {
		t.Errorf("Expected required fields ['message'], got %v", unmarshaled.Required)
	}
}
