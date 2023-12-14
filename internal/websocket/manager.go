package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

type ConnectionManager struct {
	connections   map[*websocket.Conn]struct{}             // Map to store connections
	messageRoutes map[string]func(*websocket.Conn, []byte) // Map to store message routes based on message types
	mu            sync.RWMutex                             // Mutex to ensure safe concurrent access
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections:   make(map[*websocket.Conn]struct{}),
		messageRoutes: make(map[string]func(*websocket.Conn, []byte)),
	}
}

func (cm *ConnectionManager) AddConnection(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[conn] = struct{}{}
}

func (cm *ConnectionManager) RemoveConnection(conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, conn)
}

func (manager *ConnectionManager) SendMessage(conn *websocket.Conn, message []byte) error {
	// Send the message to the WebSocket connection
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}

func (manager *ConnectionManager) HandleWebSocket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Check the request method for CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Store WebSocket connection in the ConnectionManager
	manager.AddConnection(conn)

	// Handle incoming messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			// Handle errors or connection closure
			log.Println("WebSocket connection closed:", err)
			break
		}

		// Call the corresponding message route function if registered
		if routeFunc, ok := manager.messageRoutes[string(message)]; ok {
			routeFunc(conn, message) // Execute the registered function
		}

		// You may handle other message types here if needed
		_ = messageType
	}

	// Remove the connection when it's closed
	manager.RemoveConnection(conn)
}

func (cm *ConnectionManager) On(messageType string, callback func(*websocket.Conn, []byte)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Register callback function for a specific message type
	cm.messageRoutes[messageType] = callback
}
