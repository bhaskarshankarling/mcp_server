package main

import (
	"encoding/json"
	"testing"

	"github.com/granicus/ehq-mcp-server/internal/mcp"
	"github.com/granicus/ehq-mcp-server/pkg/resources"
	"github.com/granicus/ehq-mcp-server/pkg/tools"
)

func TestServerIntegration(t *testing.T) {
	// Create a new server
	server := mcp.NewServer("ehq-mcp-server", "1.0.0")

	// Register tools
	server.RegisterTool(tools.HelloWorldTool())
	server.RegisterTool(tools.EchoTool())
	server.RegisterTool(tools.GetTimeTool())
	server.RegisterTool(tools.GetProjectsTool())

	// Register resources
	server.RegisterResource(resources.HelloWorldResource())
	server.RegisterResource(resources.ServerInfoResource())

	// Test initialize flow
	t.Run("Initialize", func(t *testing.T) {
		initMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "initialize",
			Params: map[string]interface{}{
				"protocolVersion": mcp.MCPVersion,
				"capabilities":    map[string]interface{}{},
				"clientInfo": map[string]interface{}{
					"name":    "test-client",
					"version": "1.0.0",
				},
			},
		}

		msgBytes, err := json.Marshal(initMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("Expected no error in response, got %v", response.Error)
		}

		// Check response structure
		result, ok := response.Result.(map[string]interface{})
		if !ok {
			t.Error("Expected response result to be a map")
		}

		protocolVersion, ok := result["protocolVersion"].(string)
		if !ok || protocolVersion != mcp.MCPVersion {
			t.Errorf("Expected protocol version '%s', got '%s'", mcp.MCPVersion, protocolVersion)
		}
	})

	// Test tools list
	t.Run("ToolsList", func(t *testing.T) {
		listMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      2,
			Method:  "tools/list",
		}

		msgBytes, err := json.Marshal(listMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("Expected no error in response, got %v", response.Error)
		}

		// Check that we have tools
		result, ok := response.Result.(map[string]interface{})
		if !ok {
			t.Error("Expected response result to be a map")
		}

		toolsInterface, ok := result["tools"]
		if !ok {
			t.Error("Expected tools in response")
		}

		toolsSlice, ok := toolsInterface.([]interface{})
		if !ok {
			t.Error("Expected tools to be a slice")
		}

		if len(toolsSlice) != 4 {
			t.Errorf("Expected 4 tools, got %d", len(toolsSlice))
		}
	})

	// Test tool call - hello_world
	t.Run("ToolCallHelloWorld", func(t *testing.T) {
		callMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      3,
			Method:  "tools/call",
			Params: map[string]interface{}{
				"name": "hello_world",
				"arguments": map[string]interface{}{
					"name": "Integration Test",
				},
			},
		}

		msgBytes, err := json.Marshal(callMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("Expected no error in response, got %v", response.Error)
		}
	})

	// Test tool call - echo
	t.Run("ToolCallEcho", func(t *testing.T) {
		callMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      4,
			Method:  "tools/call",
			Params: map[string]interface{}{
				"name": "echo",
				"arguments": map[string]interface{}{
					"message": "Test message",
				},
			},
		}

		msgBytes, err := json.Marshal(callMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("Expected no error in response, got %v", response.Error)
		}
	})

	// Test resources list
	t.Run("ResourcesList", func(t *testing.T) {
		listMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      5,
			Method:  "resources/list",
		}

		msgBytes, err := json.Marshal(listMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("Expected no error in response, got %v", response.Error)
		}

		// Check that we have resources
		result, ok := response.Result.(map[string]interface{})
		if !ok {
			t.Error("Expected response result to be a map")
		}

		resourcesInterface, ok := result["resources"]
		if !ok {
			t.Error("Expected resources in response")
		}

		resourcesSlice, ok := resourcesInterface.([]interface{})
		if !ok {
			t.Error("Expected resources to be a slice")
		}

		if len(resourcesSlice) != 2 {
			t.Errorf("Expected 2 resources, got %d", len(resourcesSlice))
		}
	})

	// Test resource read
	t.Run("ResourceRead", func(t *testing.T) {
		readMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      6,
			Method:  "resources/read",
			Params: map[string]interface{}{
				"uri": "hello://world",
			},
		}

		msgBytes, err := json.Marshal(readMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error != nil {
			t.Errorf("Expected no error in response, got %v", response.Error)
		}
	})
}

func TestServerCapabilities(t *testing.T) {
	server := mcp.NewServer("test-server", "1.0.0")

	// Register tools and resources
	server.RegisterTool(tools.HelloWorldTool())
	server.RegisterResource(resources.HelloWorldResource())

	// Send initialize request
	initMsg := mcp.Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": mcp.MCPVersion,
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	msgBytes, err := json.Marshal(initMsg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	responseBytes, err := server.HandleMessage(msgBytes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var response mcp.Message
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check capabilities
	result, ok := response.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected response result to be a map")
	}

	capabilities, ok := result["capabilities"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected capabilities in response")
	}

	// Check tools capability
	if _, ok := capabilities["tools"]; !ok {
		t.Error("Expected tools capability")
	}

	// Check resources capability
	if _, ok := capabilities["resources"]; !ok {
		t.Error("Expected resources capability")
	}

	// Check logging capability
	if _, ok := capabilities["logging"]; !ok {
		t.Error("Expected logging capability")
	}
}

func TestErrorHandling(t *testing.T) {
	server := mcp.NewServer("test-server", "1.0.0")

	// Test unknown method
	t.Run("UnknownMethod", func(t *testing.T) {
		unknownMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "unknown/method",
		}

		msgBytes, err := json.Marshal(unknownMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error == nil {
			t.Error("Expected error response")
		}

		if response.Error.Code != mcp.MethodNotFound {
			t.Errorf("Expected error code %d, got %d", mcp.MethodNotFound, response.Error.Code)
		}
	})

	// Test invalid tool call
	t.Run("InvalidToolCall", func(t *testing.T) {
		callMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      2,
			Method:  "tools/call",
			Params: map[string]interface{}{
				"name": "nonexistent_tool",
				"arguments": map[string]interface{}{},
			},
		}

		msgBytes, err := json.Marshal(callMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error == nil {
			t.Error("Expected error response")
		}

		if response.Error.Code != mcp.MethodNotFound {
			t.Errorf("Expected error code %d, got %d", mcp.MethodNotFound, response.Error.Code)
		}
	})

	// Test invalid resource read
	t.Run("InvalidResourceRead", func(t *testing.T) {
		readMsg := mcp.Message{
			JSONRPC: "2.0",
			ID:      3,
			Method:  "resources/read",
			Params: map[string]interface{}{
				"uri": "nonexistent://resource",
			},
		}

		msgBytes, err := json.Marshal(readMsg)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		responseBytes, err := server.HandleMessage(msgBytes)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		var response mcp.Message
		err = json.Unmarshal(responseBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Error == nil {
			t.Error("Expected error response")
		}

		if response.Error.Code != mcp.MethodNotFound {
			t.Errorf("Expected error code %d, got %d", mcp.MethodNotFound, response.Error.Code)
		}
	})
}
