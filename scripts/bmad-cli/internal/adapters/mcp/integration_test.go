package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestMCPWebSocketIntegration tests the full WebSocket MCP flow
func TestMCPWebSocketIntegration(t *testing.T) {
	// Create and start the MCP server
	server := NewServer()
	httpServer := httptest.NewServer(http.HandlerFunc(server.HandleHTTP))
	defer httpServer.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := strings.Replace(httpServer.URL, "http", "ws", 1) + "/mcp"

	// Connect to the WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Test the MCP initialization flow
	t.Run("Initialize", func(t *testing.T) {
		// Send initialize request
		initRequest := MCPMessage{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "initialize",
			Params: InitializeParams{
				ProtocolVersion: "2024-11-05",
				ClientInfo: ClientInfo{
					Name:    "test-client",
					Version: "1.0.0",
				},
				Capabilities: MCPCapabilities{
					Tools: &ToolsCapability{},
				},
			},
		}

		if err := conn.WriteJSON(initRequest); err != nil {
			t.Fatalf("Failed to send initialize request: %v", err)
		}

		// Read initialize response
		var initResponse MCPMessage
		if err := conn.ReadJSON(&initResponse); err != nil {
			t.Fatalf("Failed to read initialize response: %v", err)
		}

		// Validate response
		if initResponse.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC '2.0', got '%s'", initResponse.JSONRPC)
		}

		if initResponse.ID != float64(1) { // JSON unmarshaling converts numbers to float64
			t.Errorf("Expected ID 1, got %v", initResponse.ID)
		}

		if initResponse.Error != nil {
			t.Errorf("Expected no error, got %v", initResponse.Error)
		}

		// Parse the result
		resultBytes, _ := json.Marshal(initResponse.Result)
		var result InitializeResult
		if err := json.Unmarshal(resultBytes, &result); err != nil {
			t.Fatalf("Failed to parse initialize result: %v", err)
		}

		if result.ProtocolVersion != "2024-11-05" {
			t.Errorf("Expected protocol version '2024-11-05', got '%s'", result.ProtocolVersion)
		}

		if result.ServerInfo.Name != "bmad-mcp-server" {
			t.Errorf("Expected server name 'bmad-mcp-server', got '%s'", result.ServerInfo.Name)
		}

		// Send initialized notification
		initializedNotification := MCPMessage{
			JSONRPC: "2.0",
			Method:  "initialized",
			Params:  map[string]interface{}{},
		}

		if err := conn.WriteJSON(initializedNotification); err != nil {
			t.Fatalf("Failed to send initialized notification: %v", err)
		}
	})

	t.Run("Ping", func(t *testing.T) {
		// Send ping request
		pingRequest := MCPMessage{
			JSONRPC: "2.0",
			ID:      2,
			Method:  "ping",
			Params:  PingParams{},
		}

		if err := conn.WriteJSON(pingRequest); err != nil {
			t.Fatalf("Failed to send ping request: %v", err)
		}

		// Read ping response
		var pingResponse MCPMessage
		if err := conn.ReadJSON(&pingResponse); err != nil {
			t.Fatalf("Failed to read ping response: %v", err)
		}

		// Validate response
		if pingResponse.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC '2.0', got '%s'", pingResponse.JSONRPC)
		}

		if pingResponse.ID != float64(2) {
			t.Errorf("Expected ID 2, got %v", pingResponse.ID)
		}

		if pingResponse.Error != nil {
			t.Errorf("Expected no error, got %v", pingResponse.Error)
		}
	})

	t.Run("ListTools", func(t *testing.T) {
		// Send tools/list request
		toolsRequest := MCPMessage{
			JSONRPC: "2.0",
			ID:      3,
			Method:  "tools/list",
			Params:  map[string]interface{}{},
		}

		if err := conn.WriteJSON(toolsRequest); err != nil {
			t.Fatalf("Failed to send tools/list request: %v", err)
		}

		// Read tools/list response
		var toolsResponse MCPMessage
		if err := conn.ReadJSON(&toolsResponse); err != nil {
			t.Fatalf("Failed to read tools/list response: %v", err)
		}

		// Validate response
		if toolsResponse.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC '2.0', got '%s'", toolsResponse.JSONRPC)
		}

		if toolsResponse.ID != float64(3) {
			t.Errorf("Expected ID 3, got %v", toolsResponse.ID)
		}

		if toolsResponse.Error != nil {
			t.Errorf("Expected no error, got %v", toolsResponse.Error)
		}

		// Parse the tools result
		resultBytes, _ := json.Marshal(toolsResponse.Result)
		var result map[string]interface{}
		if err := json.Unmarshal(resultBytes, &result); err != nil {
			t.Fatalf("Failed to parse tools result: %v", err)
		}

		tools, ok := result["tools"].([]interface{})
		if !ok {
			t.Fatal("Expected tools to be an array")
		}

		// Default tools should include 'echo' and 'status'
		if len(tools) < 2 {
			t.Errorf("Expected at least 2 default tools, got %d", len(tools))
		}
	})

	t.Run("CallTool", func(t *testing.T) {
		// Send tools/call request for echo tool
		toolCallRequest := MCPMessage{
			JSONRPC: "2.0",
			ID:      4,
			Method:  "tools/call",
			Params: ToolCallParams{
				Name:      "echo",
				Arguments: json.RawMessage(`{"message": "Hello, MCP!"}`),
			},
		}

		if err := conn.WriteJSON(toolCallRequest); err != nil {
			t.Fatalf("Failed to send tools/call request: %v", err)
		}

		// Read tools/call response
		var toolCallResponse MCPMessage
		if err := conn.ReadJSON(&toolCallResponse); err != nil {
			t.Fatalf("Failed to read tools/call response: %v", err)
		}

		// Validate response
		if toolCallResponse.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC '2.0', got '%s'", toolCallResponse.JSONRPC)
		}

		if toolCallResponse.ID != float64(4) {
			t.Errorf("Expected ID 4, got %v", toolCallResponse.ID)
		}

		if toolCallResponse.Error != nil {
			t.Errorf("Expected no error, got %v", toolCallResponse.Error)
		}

		// Parse the tool result
		resultBytes, _ := json.Marshal(toolCallResponse.Result)
		var result ToolResult
		if err := json.Unmarshal(resultBytes, &result); err != nil {
			t.Fatalf("Failed to parse tool result: %v", err)
		}

		if result.IsError {
			t.Error("Expected tool execution to succeed")
		}

		if len(result.Content) == 0 {
			t.Error("Expected tool result to have content")
		}
	})
}

// TestMCPConcurrentConnections tests multiple concurrent connections
func TestMCPConcurrentConnections(t *testing.T) {
	// Create and start the MCP server
	server := NewServer()
	httpServer := httptest.NewServer(http.HandlerFunc(server.HandleHTTP))
	defer httpServer.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := strings.Replace(httpServer.URL, "http", "ws", 1) + "/mcp"

	// Number of concurrent connections to test
	numConnections := 5
	connections := make([]*websocket.Conn, numConnections)

	// Create multiple connections
	for i := 0; i < numConnections; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("Failed to connect WebSocket %d: %v", i, err)
		}
		connections[i] = conn
		defer conn.Close()
	}

	// Test that all connections can initialize
	for i, conn := range connections {
		t.Run(fmt.Sprintf("Connection_%d", i), func(t *testing.T) {
			// Send initialize request
			initRequest := MCPMessage{
				JSONRPC: "2.0",
				ID:      i + 1,
				Method:  "initialize",
				Params: InitializeParams{
					ProtocolVersion: "2024-11-05",
					ClientInfo: ClientInfo{
						Name:    fmt.Sprintf("test-client-%d", i),
						Version: "1.0.0",
					},
					Capabilities: MCPCapabilities{
						Tools: &ToolsCapability{},
					},
				},
			}

			if err := conn.WriteJSON(initRequest); err != nil {
				t.Fatalf("Failed to send initialize request for connection %d: %v", i, err)
			}

			// Read initialize response
			var initResponse MCPMessage
			if err := conn.ReadJSON(&initResponse); err != nil {
				t.Fatalf("Failed to read initialize response for connection %d: %v", i, err)
			}

			// Validate response
			if initResponse.Error != nil {
				t.Errorf("Connection %d initialization failed: %v", i, initResponse.Error)
			}

			// Send initialized notification
			initializedNotification := MCPMessage{
				JSONRPC: "2.0",
				Method:  "initialized",
				Params:  map[string]interface{}{},
			}

			if err := conn.WriteJSON(initializedNotification); err != nil {
				t.Fatalf("Failed to send initialized notification for connection %d: %v", i, err)
			}
		})
	}

	// Verify connection count
	expectedCount := numConnections
	actualCount := server.connManager.GetConnectionCount()
	if actualCount != expectedCount {
		t.Errorf("Expected %d active connections, got %d", expectedCount, actualCount)
	}
}

// TestMCPErrorHandling tests various error scenarios
func TestMCPErrorHandling(t *testing.T) {
	// Create and start the MCP server
	server := NewServer()
	httpServer := httptest.NewServer(http.HandlerFunc(server.HandleHTTP))
	defer httpServer.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := strings.Replace(httpServer.URL, "http", "ws", 1) + "/mcp"

	// Connect to the WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	tests := []struct {
		name          string
		message       MCPMessage
		expectedError int
	}{
		{
			name: "Invalid method",
			message: MCPMessage{
				JSONRPC: "2.0",
				ID:      1,
				Method:  "invalid/method",
			},
			expectedError: -32601,
		},
		{
			name: "Tools/call without initialization",
			message: MCPMessage{
				JSONRPC: "2.0",
				ID:      2,
				Method:  "tools/call",
				Params: ToolCallParams{
					Name:      "echo",
					Arguments: json.RawMessage(`{}`),
				},
			},
			expectedError: -32002,
		},
		{
			name: "Tools/list without initialization",
			message: MCPMessage{
				JSONRPC: "2.0",
				ID:      3,
				Method:  "tools/list",
			},
			expectedError: -32002,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Send the message
			if err := conn.WriteJSON(tt.message); err != nil {
				t.Fatalf("Failed to send message: %v", err)
			}

			// Read the response
			var response MCPMessage
			if err := conn.ReadJSON(&response); err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Validate error response
			if response.Error == nil {
				t.Error("Expected error response, got none")
			} else if response.Error.Code != tt.expectedError {
				t.Errorf("Expected error code %d, got %d", tt.expectedError, response.Error.Code)
			}

			if response.ID != tt.message.ID {
				t.Errorf("Expected ID %v, got %v", tt.message.ID, response.ID)
			}
		})
	}
}

// TestMCPMessageTimeout tests connection timeout handling
func TestMCPMessageTimeout(t *testing.T) {
	// Create and start the MCP server
	server := NewServer()
	httpServer := httptest.NewServer(http.HandlerFunc(server.HandleHTTP))
	defer httpServer.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := strings.Replace(httpServer.URL, "http", "ws", 1) + "/mcp"

	// Connect to the WebSocket
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Set a short read deadline for testing
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

	// Try to read with timeout
	var response MCPMessage
	err = conn.ReadJSON(&response)
	if err == nil {
		t.Error("Expected timeout error, got successful read")
	}

	// Verify it's a timeout error
	if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "deadline") {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}
