package mcp

import (
	"encoding/json"
	"testing"
)

func BenchmarkServerHandleMessage(b *testing.B) {
	server := NewServer("benchmark-server", "1.0.0")

	// Prepare a test message
	initMsg := Message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": MCPVersion,
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "benchmark-client",
				"version": "1.0.0",
			},
		},
	}

	msgBytes, err := json.Marshal(initMsg)
	if err != nil {
		b.Fatalf("Failed to marshal message: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := server.HandleMessage(msgBytes)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkServerRegisterTool(b *testing.B) {
	server := NewServer("benchmark-server", "1.0.0")

	tool := Tool{
		Name:        "benchmark-tool",
		Description: "A benchmark tool",
		InputSchema: map[string]interface{}{
			"type": "object",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use a unique name for each iteration to avoid conflicts
		tool.Name = "benchmark-tool-" + string(rune(i))
		server.RegisterTool(tool)
	}
}

func BenchmarkServerRegisterResource(b *testing.B) {
	server := NewServer("benchmark-server", "1.0.0")

	resource := Resource{
		URI:         "benchmark://resource",
		Name:        "Benchmark Resource",
		Description: "A benchmark resource",
		MimeType:    "text/plain",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use a unique URI for each iteration to avoid conflicts
		resource.URI = "benchmark://resource-" + string(rune(i))
		server.RegisterResource(resource)
	}
}

func BenchmarkExecuteHelloWorld(b *testing.B) {
	server := NewServer("benchmark-server", "1.0.0")
	args := map[string]interface{}{
		"name": "Benchmark",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.ExecuteHelloWorld(args)
	}
}

func BenchmarkExecuteEcho(b *testing.B) {
	server := NewServer("benchmark-server", "1.0.0")
	args := map[string]interface{}{
		"message": "Benchmark message",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.ExecuteEcho(args)
	}
}

func BenchmarkExecuteGetTime(b *testing.B) {
	server := NewServer("benchmark-server", "1.0.0")
	args := map[string]interface{}{
		"format": "RFC3339",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server.ExecuteGetTime(args)
	}
}

func BenchmarkJSONMarshalMessage(b *testing.B) {
	msg := Message{
		JSONRPC: "2.0",
		ID:      123,
		Method:  "test/method",
		Params: map[string]interface{}{
			"param1": "value1",
			"param2": 42,
			"param3": true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(msg)
		if err != nil {
			b.Fatalf("Marshal error: %v", err)
		}
	}
}

func BenchmarkJSONUnmarshalMessage(b *testing.B) {
	msgBytes := []byte(`{
		"jsonrpc": "2.0",
		"id": 123,
		"method": "test/method",
		"params": {
			"param1": "value1",
			"param2": 42,
			"param3": true
		}
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg Message
		err := json.Unmarshal(msgBytes, &msg)
		if err != nil {
			b.Fatalf("Unmarshal error: %v", err)
		}
	}
}

func BenchmarkNewError(b *testing.B) {
	id := 123
	code := InternalError
	message := "Benchmark error"
	data := map[string]interface{}{
		"details": "error details",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewError(id, code, message, data)
	}
}

func BenchmarkNewResponse(b *testing.B) {
	id := 123
	result := map[string]interface{}{
		"success": true,
		"data":    "benchmark data",
		"count":   42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewResponse(id, result)
	}
}
