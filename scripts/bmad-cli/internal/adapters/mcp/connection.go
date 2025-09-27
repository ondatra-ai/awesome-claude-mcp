package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Connection represents a single MCP WebSocket connection
type Connection struct {
	id         string
	conn       *websocket.Conn
	send       chan []byte
	server     *Server
	mu         sync.RWMutex
	closed     bool
	lastPing   time.Time
	authorized bool
}

// ConnectionManager manages multiple MCP connections
type ConnectionManager struct {
	connections map[string]*Connection
	mu          sync.RWMutex
	upgrader    websocket.Upgrader
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from Claude and other MCP clients
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// NewConnection creates a new MCP connection
func NewConnection(id string, conn *websocket.Conn, server *Server) *Connection {
	return &Connection{
		id:       id,
		conn:     conn,
		send:     make(chan []byte, 256),
		server:   server,
		lastPing: time.Now(),
	}
}

// HandleWebSocket handles WebSocket connections for MCP
func (cm *ConnectionManager) HandleWebSocket(w http.ResponseWriter, r *http.Request, server *Server) {
	conn, err := cm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	connectionID := generateConnectionID()
	mcpConn := NewConnection(connectionID, conn, server)

	cm.mu.Lock()
	cm.connections[connectionID] = mcpConn
	cm.mu.Unlock()

	log.Printf("New MCP connection established: %s", connectionID)

	// Start connection handlers
	go mcpConn.writePump()
	go mcpConn.readPump(cm)
}

// readPump handles reading messages from the WebSocket
func (c *Connection) readPump(cm *ConnectionManager) {
	defer func() {
		c.close()
		cm.removeConnection(c.id)
	}()

	c.conn.SetReadLimit(512 * 1024) // 512KB max message size
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.lastPing = time.Now()
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		c.handleMessage(message)
	}
}

// writePump handles writing messages to the WebSocket
func (c *Connection) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming MCP messages
func (c *Connection) handleMessage(data []byte) {
	var msg MCPMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Failed to parse MCP message: %v", err)
		c.sendError(-32700, "Parse error", nil, msg.ID)
		return
	}

	// Validate JSON-RPC version
	if msg.JSONRPC != "2.0" {
		c.sendError(-32600, "Invalid Request", "JSON-RPC version must be 2.0", msg.ID)
		return
	}

	// Handle the message based on method
	c.server.HandleMessage(c, &msg)
}

// SendMessage sends a message through the WebSocket connection
func (c *Connection) SendMessage(msg *MCPMessage) error {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return fmt.Errorf("connection is closed")
	}
	c.mu.RUnlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("send channel is full")
	}
}

// sendError sends an error response
func (c *Connection) sendError(code int, message string, data interface{}, id interface{}) {
	errorMsg := &MCPMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	c.SendMessage(errorMsg)
}

// close closes the connection
func (c *Connection) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		c.closed = true
		close(c.send)
		c.conn.Close()
	}
}

// removeConnection removes a connection from the manager
func (cm *ConnectionManager) removeConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conn, exists := cm.connections[id]; exists {
		conn.close()
		delete(cm.connections, id)
		log.Printf("MCP connection removed: %s", id)
	}
}

// GetConnectionCount returns the number of active connections
func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.connections)
}

// BroadcastMessage sends a message to all connected clients
func (cm *ConnectionManager) BroadcastMessage(msg *MCPMessage) {
	cm.mu.RLock()
	connections := make([]*Connection, 0, len(cm.connections))
	for _, conn := range cm.connections {
		connections = append(connections, conn)
	}
	cm.mu.RUnlock()

	for _, conn := range connections {
		conn.SendMessage(msg)
	}
}

// generateConnectionID generates a unique connection ID
func generateConnectionID() string {
	return fmt.Sprintf("conn_%d", time.Now().UnixNano())
}
