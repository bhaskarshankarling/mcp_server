// Package transport provides WebSocket transport for MCP
package transport

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// WebSocketTransport implements MCP over WebSocket
type WebSocketTransport struct {
	port     int
	logger   *logrus.Logger
	upgrader websocket.Upgrader
}

// NewWebSocketTransport creates a new WebSocket transport
func NewWebSocketTransport(port int) *WebSocketTransport {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &WebSocketTransport{
		port:   port,
		logger: logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

// Start starts the WebSocket server
func (t *WebSocketTransport) Start(handler func([]byte) ([]byte, error)) error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := t.upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.logger.Errorf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		t.logger.Info("🔌 WebSocket client connected")

		for {
			// Read message
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					t.logger.Errorf("WebSocket error: %v", err)
				}
				break
			}

			t.logger.Debugf("📨 Received: %s", string(message))

			// Process message
			response, err := handler(message)
			if err != nil {
				t.logger.Errorf("Handler error: %v", err)
				continue
			}

			// Send response
			if err := conn.WriteMessage(websocket.TextMessage, response); err != nil {
				t.logger.Errorf("Write error: %v", err)
				break
			}

			t.logger.Debugf("📤 Sent: %s", string(response))
		}

		t.logger.Info("🔌 WebSocket client disconnected")
	})

	t.logger.Infof("🕸️  WebSocket transport starting on port %d", t.port)
	t.logger.Infof("📡 WebSocket endpoint: ws://localhost:%d/ws", t.port)

	return http.ListenAndServe(fmt.Sprintf(":%d", t.port), nil)
}

// SetLogLevel sets the logging level
func (t *WebSocketTransport) SetLogLevel(level logrus.Level) {
	t.logger.SetLevel(level)
}
