// Package resources provides implementations for MCP resources
package resources

import (
	"github.com/granicus/ehq-mcp-server/internal/mcp"
)

// HelloWorldResource returns a hello world resource
func HelloWorldResource() mcp.Resource {
	return mcp.Resource{
		URI:         "hello://world",
		Name:        "Hello World",
		Description: "A simple hello world resource",
		MimeType:    "text/plain",
	}
}

// ServerInfoResource returns server information resource
func ServerInfoResource() mcp.Resource {
	return mcp.Resource{
		URI:         "info://server",
		Name:        "Server Information",
		Description: "Information about this MCP server",
		MimeType:    "text/plain",
	}
}
