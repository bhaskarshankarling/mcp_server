package tools

import (
	"testing"

	"github.com/granicus/ehq-mcp-server/internal/mcp"
)

func TestHelloWorldTool(t *testing.T) {
	tool := HelloWorldTool()

	if tool.Name != "hello_world" {
		t.Errorf("Expected tool name 'hello_world', got '%s'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("Expected non-empty description")
	}

	if tool.InputSchema == nil {
		t.Error("Expected input schema to be defined")
	}

	// Check input schema structure
	schema, ok := tool.InputSchema.(map[string]interface{})
	if !ok {
		t.Error("Expected input schema to be a map")
	}

	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Error("Expected properties to be defined")
	}

	nameProperty, ok := properties["name"].(map[string]interface{})
	if !ok {
		t.Error("Expected 'name' property to be defined")
	}

	if nameProperty["type"] != "string" {
		t.Errorf("Expected name property type 'string', got '%v'", nameProperty["type"])
	}
}

func TestEchoTool(t *testing.T) {
	tool := EchoTool()

	if tool.Name != "echo" {
		t.Errorf("Expected tool name 'echo', got '%s'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("Expected non-empty description")
	}

	if tool.InputSchema == nil {
		t.Error("Expected input schema to be defined")
	}

	// Check input schema structure
	schema, ok := tool.InputSchema.(map[string]interface{})
	if !ok {
		t.Error("Expected input schema to be a map")
	}

	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Error("Expected properties to be defined")
	}

	messageProperty, ok := properties["message"].(map[string]interface{})
	if !ok {
		t.Error("Expected 'message' property to be defined")
	}

	if messageProperty["type"] != "string" {
		t.Errorf("Expected message property type 'string', got '%v'", messageProperty["type"])
	}

	// Check required fields
	required, ok := schema["required"].([]string)
	if !ok {
		t.Error("Expected required fields to be defined")
	}

	if len(required) != 1 || required[0] != "message" {
		t.Errorf("Expected required field 'message', got %v", required)
	}
}

func TestGetTimeTool(t *testing.T) {
	tool := GetTimeTool()

	if tool.Name != "get_time" {
		t.Errorf("Expected tool name 'get_time', got '%s'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("Expected non-empty description")
	}

	if tool.InputSchema == nil {
		t.Error("Expected input schema to be defined")
	}

	// Check input schema structure
	schema, ok := tool.InputSchema.(map[string]interface{})
	if !ok {
		t.Error("Expected input schema to be a map")
	}

	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Error("Expected properties to be defined")
	}

	formatProperty, ok := properties["format"].(map[string]interface{})
	if !ok {
		t.Error("Expected 'format' property to be defined")
	}

	if formatProperty["type"] != "string" {
		t.Errorf("Expected format property type 'string', got '%v'", formatProperty["type"])
	}

	// Check enum values
	enum, ok := formatProperty["enum"].([]string)
	if !ok {
		t.Error("Expected enum values to be defined")
	}

	expectedEnums := []string{"RFC3339", "Unix", "Kitchen"}
	if len(enum) != len(expectedEnums) {
		t.Errorf("Expected %d enum values, got %d", len(expectedEnums), len(enum))
	}

	for i, expected := range expectedEnums {
		if i >= len(enum) || enum[i] != expected {
			t.Errorf("Expected enum value '%s', got '%s'", expected, enum[i])
		}
	}
}

func TestGetProjectsTool(t *testing.T) {
	tool := GetProjectsTool()

	if tool.Name != "get_projects" {
		t.Errorf("Expected tool name 'get_projects', got '%s'", tool.Name)
	}

	if tool.Description == "" {
		t.Error("Expected non-empty description")
	}

	if tool.InputSchema == nil {
		t.Error("Expected input schema to be defined")
	}

	// Check input schema structure
	schema, ok := tool.InputSchema.(map[string]interface{})
	if !ok {
		t.Error("Expected input schema to be a map")
	}

	if schema["type"] != "object" {
		t.Errorf("Expected schema type 'object', got '%v'", schema["type"])
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Error("Expected properties to be defined")
	}

	searchProperty, ok := properties["search"].(map[string]interface{})
	if !ok {
		t.Error("Expected 'search' property to be defined")
	}

	if searchProperty["type"] != "string" {
		t.Errorf("Expected search property type 'string', got '%v'", searchProperty["type"])
	}
}

func TestAllToolsReturnValidMCPTool(t *testing.T) {
	tools := []mcp.Tool{
		HelloWorldTool(),
		EchoTool(),
		GetTimeTool(),
		GetProjectsTool(),
	}

	for _, tool := range tools {
		if tool.Name == "" {
			t.Error("Tool name should not be empty")
		}

		if tool.Description == "" {
			t.Error("Tool description should not be empty")
		}

		if tool.InputSchema == nil {
			t.Error("Tool input schema should not be nil")
		}

		// Validate that input schema is a proper object
		schema, ok := tool.InputSchema.(map[string]interface{})
		if !ok {
			t.Errorf("Tool '%s' input schema should be a map", tool.Name)
			continue
		}

		if schema["type"] != "object" {
			t.Errorf("Tool '%s' schema type should be 'object'", tool.Name)
		}

		// Properties should exist
		if _, ok := schema["properties"]; !ok {
			t.Errorf("Tool '%s' schema should have properties", tool.Name)
		}
	}
}

func TestToolNamesAreUnique(t *testing.T) {
	tools := []mcp.Tool{
		HelloWorldTool(),
		EchoTool(),
		GetTimeTool(),
		GetProjectsTool(),
	}

	names := make(map[string]bool)
	for _, tool := range tools {
		if names[tool.Name] {
			t.Errorf("Duplicate tool name found: %s", tool.Name)
		}
		names[tool.Name] = true
	}
}
