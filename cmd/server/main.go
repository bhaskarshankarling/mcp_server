// EHQ MCP Server - A Model Context Protocol server for EngagementHQ
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
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
		debug   = flag.Bool("debug", false, "Enable debug logging")
		version = flag.Bool("version", false, "Show version information")
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

	// Create transport
	transport := transport.NewStdioTransport()
	if *debug {
		transport.SetLogLevel(logrus.DebugLevel)
	}

	fmt.Fprintf(os.Stderr, "🚀 %s v%s starting...\n", ServerName, ServerVersion)
	fmt.Fprintf(os.Stderr, "📡 Protocol: %s\n", mcp.MCPVersion)
	fmt.Fprintf(os.Stderr, "🔧 Tools: %d registered\n", 4)
	fmt.Fprintf(os.Stderr, "📄 Resources: %d registered\n", 2)
	fmt.Fprintf(os.Stderr, "📨 Listening on stdio...\n")

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server loop
	go func() {
		for {
			// Read message from transport
			data, err := transport.ReadMessage()
			if err != nil {
				if err == io.EOF {
					fmt.Fprintf(os.Stderr, "📪 Client disconnected\n")
					os.Exit(0)
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
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Fprintf(os.Stderr, "\n🛑 Shutting down gracefully...\n")
}
