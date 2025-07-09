// Package tools provides implementations for MCP tools
package tools

import (
	"github.com/granicus/ehq-mcp-server/internal/mcp"
)

// HelloWorldTool implements a simple hello world tool
func HelloWorldTool() mcp.Tool {
	return mcp.Tool{
		Name:        "hello_world",
		Description: "Returns a friendly hello world message",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name to greet (optional)",
				},
			},
		},
	}
}

// EchoTool implements an echo tool that returns the input
func EchoTool() mcp.Tool {
	return mcp.Tool{
		Name:        "echo",
		Description: "Echoes back the provided message",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"message": map[string]interface{}{
					"type":        "string",
					"description": "Message to echo back",
				},
			},
			"required": []string{"message"},
		},
	}
}

// GetTimeTool implements a tool that returns the current time
func GetTimeTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_time",
		Description: "Returns the current date and time",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"format": map[string]interface{}{
					"type":        "string",
					"description": "Time format (optional, defaults to RFC3339)",
					"enum":        []string{"RFC3339", "Unix", "Kitchen"},
				},
			},
		},
	}
}

// GetProjectsTool implements a tool that fetches projects from EHQ API
func GetProjectsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "get_projects",
		Description: "Fetches projects from the EHQ API using authentication",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"search": map[string]interface{}{
					"type":        "string",
					"description": "Search term to filter projects by name or description (optional)",
				},
			},
		},
	}
}
