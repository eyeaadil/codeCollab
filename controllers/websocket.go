package controllers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader for upgrading HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins; customize for security
	},
}

// Connection pool to manage active WebSocket connections
var connections = struct {
	sync.RWMutex
	clients map[*websocket.Conn]bool
}{
	clients: make(map[*websocket.Conn]bool),
}

// Message represents a WebSocket message
type Message struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
}

// WebSocketHandler manages WebSocket connections and broadcasts messages
func WebSocketHandler(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Add connection to the pool
	connections.Lock()
	connections.clients[conn] = true
	connections.Unlock()

	// Remove connection from the pool on close
	defer func() {
		connections.Lock()
		delete(connections.clients, conn)
		connections.Unlock()
	}()

	// Listen for messages from the client
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Log the received message
		log.Printf("Received message: %+v\n", msg)

		// Broadcast the message to all connected clients
		broadcastMessage(msg)
	}
}

// broadcastMessage sends a message to all connected WebSocket clients
func broadcastMessage(msg Message) {
	connections.RLock()
	defer connections.RUnlock()

	for client := range connections.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Println("Error broadcasting message:", err)
			client.Close()
			delete(connections.clients, client)
		}
	}
}
