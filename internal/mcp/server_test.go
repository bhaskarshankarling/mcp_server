package mcp

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestNewServer(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	if server.serverInfo.Name != "test-server" {
		t.Errorf("Expected server name 'test-server', got '%s'", server.serverInfo.Name)
	}

	if server.serverInfo.Version != "1.0.0" {
		t.Errorf("Expected server version '1.0.0', got '%s'", server.serverInfo.Version)
	}

	if server.tools == nil {
		t.Error("Expected tools map to be initialized")
	}

	if server.resources == nil {
		t.Error("Expected resources map to be initialized")
	}

	if server.handlers == nil {
		t.Error("Expected handlers map to be initialized")
	}

	// Check that default handlers are registered
	expectedHandlers := []string{
		"initialize",
		"tools/list",
		"tools/call",
		"resources/list",
		"resources/read",
	}

	for _, handler := range expectedHandlers {
		if _, exists := server.handlers[handler]; !exists {
			t.Errorf("Expected handler '%s' to be registered", handler)
		}
	}
}

func TestRegisterTool(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	tool := Tool{
		Name:        "test-tool",
		Description: "A test tool",
		InputSchema: map[string]interface{}{
			"type": "object",
		},
	}

	server.RegisterTool(tool)

	if registeredTool, exists := server.tools["test-tool"]; !exists {
		t.Error("Expected tool to be registered")
	} else if registeredTool.Name != "test-tool" {
		t.Errorf("Expected tool name 'test-tool', got '%s'", registeredTool.Name)
	}
}

func TestRegisterResource(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	resource := Resource{
		URI:         "test://resource",
		Name:        "Test Resource",
		Description: "A test resource",
		MimeType:    "text/plain",
	}

	server.RegisterResource(resource)

	if registeredResource, exists := server.resources["test://resource"]; !exists {
		t.Error("Expected resource to be registered")
	} else if registeredResource.URI != "test://resource" {
		t.Errorf("Expected resource URI 'test://resource', got '%s'", registeredResource.URI)
	}
}

func TestHandleInitialize(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	msg := &Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": MCPVersion,
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}

	response, err := server.handleInitialize(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.ID != 1 {
		t.Errorf("Expected response ID 1, got %v", response.ID)
	}

	// Check that the response contains the correct structure
	result, ok := response.Result.(InitializeResponse)
	if !ok {
		t.Error("Expected response result to be InitializeResponse")
	}

	if result.ProtocolVersion != MCPVersion {
		t.Errorf("Expected protocol version '%s', got '%s'", MCPVersion, result.ProtocolVersion)
	}

	if result.ServerInfo.Name != "test-server" {
		t.Errorf("Expected server name 'test-server', got '%s'", result.ServerInfo.Name)
	}
}

func TestHandleToolsList(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Register a test tool
	tool := Tool{
		Name:        "test-tool",
		Description: "A test tool",
		InputSchema: map[string]interface{}{
			"type": "object",
		},
	}
	server.RegisterTool(tool)

	msg := &Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}

	response, err := server.handleToolsList(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result, ok := response.Result.(ToolsListResponse)
	if !ok {
		t.Error("Expected response result to be ToolsListResponse")
	}

	if len(result.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(result.Tools))
	}

	if result.Tools[0].Name != "test-tool" {
		t.Errorf("Expected tool name 'test-tool', got '%s'", result.Tools[0].Name)
	}
}

func TestHandleToolsCall(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Register a test tool
	tool := Tool{
		Name:        "hello_world",
		Description: "Hello world tool",
		InputSchema: map[string]interface{}{
			"type": "object",
		},
	}
	server.RegisterTool(tool)

	msg := &Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      "hello_world",
			"arguments": map[string]interface{}{
				"name": "Test",
			},
		},
	}

	response, err := server.handleToolsCall(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result, ok := response.Result.(ToolsCallResponse)
	if !ok {
		t.Error("Expected response result to be ToolsCallResponse")
	}

	if len(result.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Content))
	}

	if result.Content[0].Type != "text" {
		t.Errorf("Expected content type 'text', got '%s'", result.Content[0].Type)
	}
}

func TestHandleToolsCallUnknownTool(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	msg := &Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      "unknown-tool",
			"arguments": map[string]interface{}{},
		},
	}

	response, err := server.handleToolsCall(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.Error == nil {
		t.Error("Expected error response for unknown tool")
	}

	if response.Error.Code != MethodNotFound {
		t.Errorf("Expected error code %d, got %d", MethodNotFound, response.Error.Code)
	}
}

func TestHandleResourcesList(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Register a test resource
	resource := Resource{
		URI:         "test://resource",
		Name:        "Test Resource",
		Description: "A test resource",
		MimeType:    "text/plain",
	}
	server.RegisterResource(resource)

	msg := &Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "resources/list",
	}

	response, err := server.handleResourcesList(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result, ok := response.Result.(ResourcesListResponse)
	if !ok {
		t.Error("Expected response result to be ResourcesListResponse")
	}

	if len(result.Resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(result.Resources))
	}

	if result.Resources[0].URI != "test://resource" {
		t.Errorf("Expected resource URI 'test://resource', got '%s'", result.Resources[0].URI)
	}
}

func TestHandleResourcesRead(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Register a test resource
	resource := Resource{
		URI:         "hello://world",
		Name:        "Hello World Resource",
		Description: "A hello world resource",
		MimeType:    "text/plain",
	}
	server.RegisterResource(resource)

	msg := &Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "resources/read",
		Params: map[string]interface{}{
			"uri": "hello://world",
		},
	}

	response, err := server.handleResourcesRead(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	result, ok := response.Result.(ResourcesReadResponse)
	if !ok {
		t.Error("Expected response result to be ResourcesReadResponse")
	}

	if len(result.Contents) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(result.Contents))
	}

	if result.Contents[0].URI != "hello://world" {
		t.Errorf("Expected content URI 'hello://world', got '%s'", result.Contents[0].URI)
	}
}

func TestHandleMessage(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Test valid initialize message
	initMsg := Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": MCPVersion,
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

	var response Message
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error != nil {
		t.Errorf("Expected no error in response, got %v", response.Error)
	}

	// Check response ID (could be int or float64 from JSON)
	responseID := response.ID
	if responseID == nil {
		t.Error("Expected response ID to be set")
	} else {
		// Convert to string for comparison to handle JSON number parsing
		expectedID := "1"
		actualID := ""
		switch v := responseID.(type) {
		case int:
			actualID = fmt.Sprintf("%d", v)
		case float64:
			actualID = fmt.Sprintf("%.0f", v)
		case string:
			actualID = v
		default:
			actualID = fmt.Sprintf("%v", v)
		}

		if actualID != expectedID {
			t.Errorf("Expected response ID %s, got %s (type: %T)", expectedID, actualID, responseID)
		}
	}
}

func TestHandleMessageInvalidJSON(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	invalidJSON := []byte(`{invalid json}`)

	responseBytes, err := server.HandleMessage(invalidJSON)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var response Message
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error == nil {
		t.Error("Expected error response for invalid JSON")
	}

	if response.Error.Code != ParseError {
		t.Errorf("Expected error code %d, got %d", ParseError, response.Error.Code)
	}
}

func TestHandleMessageUnknownMethod(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	unknownMsg := Message{
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

	var response Message
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error == nil {
		t.Error("Expected error response for unknown method")
	}

	if response.Error.Code != MethodNotFound {
		t.Errorf("Expected error code %d, got %d", MethodNotFound, response.Error.Code)
	}
}

func TestExecuteHelloWorld(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Test with name argument
	args := map[string]interface{}{
		"name": "World",
	}
	result := server.ExecuteHelloWorld(args)
	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Test without name argument
	result = server.ExecuteHelloWorld(map[string]interface{}{})
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestExecuteEcho(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Test with message argument
	args := map[string]interface{}{
		"message": "Hello, World!",
	}
	result := server.ExecuteEcho(args)
	expected := "Echo: Hello, World!"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test without message argument
	result = server.ExecuteEcho(map[string]interface{}{})
	if result != "Error: message parameter is required" {
		t.Errorf("Expected error message, got '%s'", result)
	}
}

func TestExecuteGetTime(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Test default format
	result := server.ExecuteGetTime(map[string]interface{}{})
	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Test Unix format
	args := map[string]interface{}{
		"format": "Unix",
	}
	result = server.ExecuteGetTime(args)
	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Test Kitchen format
	args = map[string]interface{}{
		"format": "Kitchen",
	}
	result = server.ExecuteGetTime(args)
	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestSetLogLevel(t *testing.T) {
	server := NewServer("test-server", "1.0.0")

	// Create a test hook to capture log entries
	hook := test.NewGlobal()
	server.logger.AddHook(hook)

	// Set log level to Debug
	server.SetLogLevel(logrus.DebugLevel)

	// Log a debug message
	server.logger.Debug("test debug message")

	// Check that the debug message was captured
	if len(hook.Entries) == 0 {
		t.Error("Expected debug message to be logged")
	}

	// Clear entries
	hook.Reset()

	// Set log level to Error
	server.SetLogLevel(logrus.ErrorLevel)

	// Log a debug message (should not be captured)
	server.logger.Debug("test debug message")

	// Check that the debug message was not captured
	if len(hook.Entries) != 0 {
		t.Error("Expected debug message to not be logged at Error level")
	}
}
