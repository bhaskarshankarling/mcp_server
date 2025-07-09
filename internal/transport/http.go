// Package transport provides HTTP transport for MCP
package transport

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// HTTPTransport implements MCP over HTTP
type HTTPTransport struct {
	port   int
	logger *logrus.Logger
	server *http.Server
}

// NewHTTPTransport creates a new HTTP transport
func NewHTTPTransport(port int) *HTTPTransport {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &HTTPTransport{
		port:   port,
		logger: logger,
	}
}

// Start starts the HTTP server
func (t *HTTPTransport) Start(handler func([]byte) ([]byte, error)) error {
	mux := http.NewServeMux()

	// Handle MCP requests
	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		response, err := handler(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})

	t.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", t.port),
		Handler: mux,
	}

	t.logger.Infof("🌐 HTTP transport starting on port %d", t.port)
	t.logger.Infof("📡 MCP endpoint: http://localhost:%d/mcp", t.port)
	t.logger.Infof("❤️  Health check: http://localhost:%d/health", t.port)

	return t.server.ListenAndServe()
}

// Stop stops the HTTP server
func (t *HTTPTransport) Stop() error {
	if t.server != nil {
		return t.server.Close()
	}
	return nil
}

// SetLogLevel sets the logging level
func (t *HTTPTransport) SetLogLevel(level logrus.Level) {
	t.logger.SetLevel(level)
}
