package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrade upgrades the HTTP connection to a WebSocket connection.
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// ReadMessage reads a message from the WebSocket connection.
func ReadMessage(conn *websocket.Conn) (string, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return "", err
	}

	return string(message), nil
}

// WriteMessage writes a message to the WebSocket connection.
func WriteMessage(conn *websocket.Conn, message string) error {
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return err
	}

	return nil
}

// Close closes the WebSocket connection.
func Close(conn *websocket.Conn) error {
	err := conn.Close()
	if err != nil {
		return err
	}

	return nil
}
