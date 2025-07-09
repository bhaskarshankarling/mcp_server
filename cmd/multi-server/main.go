// Multi-transport MCP Server - supports stdio, HTTP, and WebSocket
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/granicus/ehq-mcp-server/internal/mcp"
	"github.com/granicus/ehq-mcp-server/internal/transport"
	"github.com/granicus/ehq-mcp-server/pkg/resources"
	"github.com/granicus/ehq-mcp-server/pkg/tools"
	"github.com/sirupsen/logrus"
)

const (
	ServerName    = "EHQ MCP Server"
	ServerVersion = "1.0.0"
)

func main() {
	var (
		debug     = flag.Bool("debug", false, "Enable debug logging")
		version   = flag.Bool("version", false, "Show version information")
		httpPort  = flag.Int("http", 0, "Enable HTTP transport on specified port (e.g., -http 8080)")
		wsPort    = flag.Int("ws", 0, "Enable WebSocket transport on specified port (e.g., -ws 8081)")
		stdioOnly = flag.Bool("stdio-only", false, "Use only stdio transport (default behavior)")
	)
	flag.Parse()

	if *version {
		fmt.Printf("%s v%s\n", ServerName, ServerVersion)
		fmt.Printf("Protocol Version: %s\n", mcp.MCPVersion)
		os.Exit(0)
	}

	// Create the MCP server
	server := mcp.NewServer(ServerName, ServerVersion)

	// Set log level
	if *debug {
		server.SetLogLevel(logrus.DebugLevel)
	}

	// Register tools
	server.RegisterTool(tools.HelloWorldTool())
	server.RegisterTool(tools.EchoTool())
	server.RegisterTool(tools.GetTimeTool())
	server.RegisterTool(tools.GetProjectsTool())


	// Register resources
	server.RegisterResource(resources.HelloWorldResource())
	server.RegisterResource(resources.ServerInfoResource())

	fmt.Fprintf(os.Stderr, "🚀 %s v%s starting...\n", ServerName, ServerVersion)
	fmt.Fprintf(os.Stderr, "📡 Protocol: %s\n", mcp.MCPVersion)
	fmt.Fprintf(os.Stderr, "🔧 Tools: %d registered\n", 4)
	fmt.Fprintf(os.Stderr, "📄 Resources: %d registered\n", 2)

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Start stdio transport (default)
	if *stdioOnly || (*httpPort == 0 && *wsPort == 0) {
		fmt.Fprintf(os.Stderr, "📨 Starting stdio transport...\n")
		wg.Add(1)
		go func() {
			defer wg.Done()
			runStdioTransport(server, *debug)
		}()
	}

	// Start HTTP transport if requested
	if *httpPort > 0 {
		fmt.Fprintf(os.Stderr, "🌐 Starting HTTP transport on port %d...\n", *httpPort)
		wg.Add(1)
		go func() {
			defer wg.Done()
			runHTTPTransport(server, *httpPort, *debug, ctx)
		}()
	}

	// Start WebSocket transport if requested
	if *wsPort > 0 {
		fmt.Fprintf(os.Stderr, "🕸️  Starting WebSocket transport on port %d...\n", *wsPort)
		wg.Add(1)
		go func() {
			defer wg.Done()
			runWebSocketTransport(server, *wsPort, *debug, ctx)
		}()
	}

	// Wait for shutdown signal
	<-sigChan
	fmt.Fprintf(os.Stderr, "\n🛑 Shutting down gracefully...\n")
	cancel()
	wg.Wait()
}

func runStdioTransport(server *mcp.Server, debug bool) {
	transport := transport.NewStdioTransport()
	if debug {
		transport.SetLogLevel(logrus.DebugLevel)
	}

	for {
		// Read message from transport
		data, err := transport.ReadMessage()
		if err != nil {
			if err == io.EOF {
				fmt.Fprintf(os.Stderr, "📪 Client disconnected\n")
				return
			}
			fmt.Fprintf(os.Stderr, "❌ Error reading message: %v\n", err)
			continue
		}

		// Process message through MCP server
		response, err := server.HandleMessage(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error handling message: %v\n", err)
			continue
		}

		// Send response back through transport
		if err := transport.WriteMessage(response); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error writing response: %v\n", err)
			continue
		}
	}
}

func runHTTPTransport(server *mcp.Server, port int, debug bool, ctx context.Context) {
	httpTransport := transport.NewHTTPTransport(port)
	if debug {
		httpTransport.SetLogLevel(logrus.DebugLevel)
	}

	go func() {
		<-ctx.Done()
		httpTransport.Stop()
	}()

	if err := httpTransport.Start(server.HandleMessage); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "❌ HTTP transport error: %v\n", err)
	}
}

func runWebSocketTransport(server *mcp.Server, port int, debug bool, ctx context.Context) {
	wsTransport := transport.NewWebSocketTransport(port)
	if debug {
		wsTransport.SetLogLevel(logrus.DebugLevel)
	}

	if err := wsTransport.Start(server.HandleMessage); err != nil {
		fmt.Fprintf(os.Stderr, "❌ WebSocket transport error: %v\n", err)
	}
}
