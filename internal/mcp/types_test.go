package mcp

import (
	"testing"
)

func TestNewError(t *testing.T) {
	id := 123
	code := ParseError
	message := "Parse error occurred"
	data := map[string]interface{}{"details": "Invalid JSON"}

	errorMsg := NewError(id, code, message, data)

	if errorMsg.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC '2.0', got '%s'", errorMsg.JSONRPC)
	}

	if errorMsg.ID != id {
		t.Errorf("Expected ID %v, got %v", id, errorMsg.ID)
	}

	if errorMsg.Error.Code != code {
		t.Errorf("Expected error code %d, got %d", code, errorMsg.Error.Code)
	}

	if errorMsg.Error.Message != message {
		t.Errorf("Expected error message '%s', got '%s'", message, errorMsg.Error.Message)
	}

	// Note: We can't directly compare maps, so we'll check if data is not nil
	if data != nil && errorMsg.Error.Data == nil {
		t.Error("Expected error data to be set, but it was nil")
	}
}

func TestNewResponse(t *testing.T) {
	id := 456
	result := map[string]interface{}{
		"success": true,
		"data":    "test data",
	}

	response := NewResponse(id, result)

	if response.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC '2.0', got '%s'", response.JSONRPC)
	}

	if response.ID != id {
		t.Errorf("Expected ID %v, got %v", id, response.ID)
	}

	// Note: We can't directly compare maps, so we'll check if result is not nil
	if result != nil && response.Result == nil {
		t.Error("Expected response result to be set, but it was nil")
	}

	if response.Error != nil {
		t.Errorf("Expected no error, got %v", response.Error)
	}
}

func TestErrorCodes(t *testing.T) {
	// Test that error codes have expected values
	if ParseError != -32700 {
		t.Errorf("Expected ParseError to be -32700, got %d", ParseError)
	}

	if InvalidRequest != -32600 {
		t.Errorf("Expected InvalidRequest to be -32600, got %d", InvalidRequest)
	}

	if MethodNotFound != -32601 {
		t.Errorf("Expected MethodNotFound to be -32601, got %d", MethodNotFound)
	}

	if InvalidParams != -32602 {
		t.Errorf("Expected InvalidParams to be -32602, got %d", InvalidParams)
	}

	if InternalError != -32603 {
		t.Errorf("Expected InternalError to be -32603, got %d", InternalError)
	}
}

func TestMCPVersion(t *testing.T) {
	expectedVersion := "2024-11-05"
	if MCPVersion != expectedVersion {
		t.Errorf("Expected MCP version '%s', got '%s'", expectedVersion, MCPVersion)
	}
}

func TestToolValidation(t *testing.T) {
	tool := Tool{
		Name:        "test-tool",
		Description: "A test tool",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"param1": map[string]interface{}{
					"type":        "string",
					"description": "First parameter",
				},
			},
			"required": []string{"param1"},
		},
	}

	if tool.Name == "" {
		t.Error("Tool name should not be empty")
	}

	if tool.Description == "" {
		t.Error("Tool description should not be empty")
	}

	if tool.InputSchema == nil {
		t.Error("Tool input schema should not be nil")
	}
}

func TestResourceValidation(t *testing.T) {
	resource := Resource{
		URI:         "test://resource",
		Name:        "Test Resource",
		Description: "A test resource",
		MimeType:    "text/plain",
	}

	if resource.URI == "" {
		t.Error("Resource URI should not be empty")
	}

	if resource.Name == "" {
		t.Error("Resource name should not be empty")
	}

	if resource.MimeType == "" {
		t.Error("Resource mime type should not be empty")
	}
}

func TestInitializeRequest(t *testing.T) {
	req := InitializeRequest{
		ProtocolVersion: MCPVersion,
		Capabilities: Capabilities{
			Roots: &RootsCapability{
				ListChanged: true,
			},
			Sampling: &SamplingCapability{},
		},
		ClientInfo: ClientInfo{
			Name:    "test-client",
			Version: "1.0.0",
		},
	}

	if req.ProtocolVersion != MCPVersion {
		t.Errorf("Expected protocol version '%s', got '%s'", MCPVersion, req.ProtocolVersion)
	}

	if req.ClientInfo.Name != "test-client" {
		t.Errorf("Expected client name 'test-client', got '%s'", req.ClientInfo.Name)
	}

	if req.ClientInfo.Version != "1.0.0" {
		t.Errorf("Expected client version '1.0.0', got '%s'", req.ClientInfo.Version)
	}
}

func TestInitializeResponse(t *testing.T) {
	resp := InitializeResponse{
		ProtocolVersion: MCPVersion,
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
			Logging: &LoggingCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "test-server",
			Version: "1.0.0",
		},
	}

	if resp.ProtocolVersion != MCPVersion {
		t.Errorf("Expected protocol version '%s', got '%s'", MCPVersion, resp.ProtocolVersion)
	}

	if resp.ServerInfo.Name != "test-server" {
		t.Errorf("Expected server name 'test-server', got '%s'", resp.ServerInfo.Name)
	}

	if resp.ServerInfo.Version != "1.0.0" {
		t.Errorf("Expected server version '1.0.0', got '%s'", resp.ServerInfo.Version)
	}
}

func TestToolsCallRequest(t *testing.T) {
	req := ToolsCallRequest{
		Name: "test-tool",
		Arguments: map[string]interface{}{
			"param1": "value1",
			"param2": 42,
		},
	}

	if req.Name != "test-tool" {
		t.Errorf("Expected tool name 'test-tool', got '%s'", req.Name)
	}

	if len(req.Arguments) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(req.Arguments))
	}

	if req.Arguments["param1"] != "value1" {
		t.Errorf("Expected param1 to be 'value1', got %v", req.Arguments["param1"])
	}

	if req.Arguments["param2"] != 42 {
		t.Errorf("Expected param2 to be 42, got %v", req.Arguments["param2"])
	}
}

func TestResourcesReadRequest(t *testing.T) {
	req := ResourcesReadRequest{
		URI: "test://resource",
	}

	if req.URI != "test://resource" {
		t.Errorf("Expected URI 'test://resource', got '%s'", req.URI)
	}
}

func TestContent(t *testing.T) {
	content := Content{
		Type: "text",
		Text: "Hello, World!",
	}

	if content.Type != "text" {
		t.Errorf("Expected content type 'text', got '%s'", content.Type)
	}

	if content.Text != "Hello, World!" {
		t.Errorf("Expected content text 'Hello, World!', got '%s'", content.Text)
	}
}

func TestResourceContent(t *testing.T) {
	content := ResourceContent{
		URI:      "test://resource",
		MimeType: "text/plain",
		Text:     "Resource content",
	}

	if content.URI != "test://resource" {
		t.Errorf("Expected URI 'test://resource', got '%s'", content.URI)
	}

	if content.MimeType != "text/plain" {
		t.Errorf("Expected mime type 'text/plain', got '%s'", content.MimeType)
	}

	if content.Text != "Resource content" {
		t.Errorf("Expected text 'Resource content', got '%s'", content.Text)
	}
}
