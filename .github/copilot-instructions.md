# GitHub Copilot Instructions for EHQ MCP Server

## Project Overview

This is the **EHQ MCP Server** - a Model Context Protocol (MCP) server implementation in Go for EngagementHQ. The server provides AI-powered tools and resources that can be consumed by AI assistants like Claude, ChatGPT, and other MCP-compatible clients.

## 🏗️ Architecture & Structure

### Core Components

- **MCP Protocol**: Implements MCP 2024-11-05 specification with JSON-RPC 2.0
- **HTTP Transport**: REST API endpoints for MCP communication
- **Tools Framework**: Executable functions for AI assistants
- **EHQ API Integration**: EngagementHQ platform integration

### Project Structure

```
ehq_mcp_server/
├── cmd/
│   └── server/            # HTTP MCP server
├── internal/              # Private packages
│   ├── mcp/              # Core MCP implementation
│   └── transport/        # HTTP transport implementation
├── pkg/                   # Public packages
│   ├── api/              # EHQ API client
│   └── tools/            # MCP tools
└── bin/                   # Built binaries
```

## 🎯 Coding Guidelines

### Go Best Practices

1. **Follow Go conventions**: Use gofmt, golint, and go vet
2. **Package naming**: Use lowercase, single words when possible
3. **Error handling**: Always handle errors explicitly, never ignore them
4. **Context usage**: Use context.Context for cancellation and timeouts
5. **Interfaces**: Keep interfaces small and focused
6. **Documentation**: Use godoc comments for all exported functions

### MCP-Specific Guidelines

1. **Tool Definitions**: All tools must implement `mcp.Tool` interface
2. **JSON-RPC**: Follow JSON-RPC 2.0 specification exactly
3. **Schema Validation**: Use JSON Schema for input/output validation
4. **Error Responses**: Return proper MCP error codes and messages
5. **HTTP Transport**: Keep business logic separate from HTTP transport layer

### Code Style

```go
// Preferred: Clear, descriptive function names
func CreateEHQProjectTool() mcp.Tool {
    return mcp.Tool{
        Name:        "create_project",
        Description: "Creates a new project in EngagementHQ",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "name": map[string]interface{}{
                    "type":        "string",
                    "description": "Project name",
                },
            },
            "required": []string{"name"},
        },
    }
}

// Tool handler with proper error handling
func HandleCreateProject(params map[string]interface{}) (*mcp.ToolResult, error) {
    name, ok := params["name"].(string)
    if !ok {
        return nil, fmt.Errorf("name parameter is required and must be a string")
    }

    // Implementation here
    return &mcp.ToolResult{
        Content: []mcp.Content{
            {
                Type: "text",
                Text: fmt.Sprintf("Created project: %s", name),
            },
        },
    }, nil
}
```

## 🛠️ Common Patterns

### Adding New Tools

1. **Define the tool** in `pkg/tools/`:
```go
func NewExampleTool() mcp.Tool {
    return mcp.Tool{
        Name:        "example_tool",
        Description: "Description of what this tool does",
        InputSchema: // JSON schema definition
    }
}
```

2. **Implement the handler**:
```go
func HandleExampleTool(params map[string]interface{}) (*mcp.ToolResult, error) {
    // Validate inputs
    // Perform operation
    // Return result
}
```

3. **Register in server** (`cmd/server/main.go`):
```go
server.RegisterTool(tools.NewExampleTool(), tools.HandleExampleTool)
```

### EHQ API Integration

- Use the `pkg/api/EHQClient` for all EngagementHQ API calls
- Always handle authentication properly
- Implement proper error handling and retries
- Use structured logging for API calls

## 📝 Testing Guidelines

### Test Structure

```go
func TestExampleTool(t *testing.T) {
    tests := []struct {
        name        string
        params      map[string]interface{}
        expectError bool
        expectText  string
    }{
        {
            name:       "valid input",
            params:     map[string]interface{}{"name": "test"},
            expectText: "Expected output",
        },
        {
            name:        "missing parameter",
            params:      map[string]interface{}{},
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := HandleExampleTool(tt.params)

            if tt.expectError {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.Contains(t, result.Content[0].Text, tt.expectText)
        })
    }
}
```

### Test Files

- Create `*_test.go` files alongside implementation files
- Use table-driven tests for multiple scenarios
- Test both success and error cases
- Mock external dependencies (EHQ API calls)

## 🔧 Development Workflow

### Building & Running

```bash
# Build HTTP server
go build -o bin/ehq-mcp-server ./cmd/server

# Run with debug logging
./bin/ehq-mcp-server -debug

# Run HTTP server on custom port
./bin/ehq-mcp-server -port 8080
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage/coverage.html

# Run specific test files
go test ./pkg/tools/
```

## 🚨 Important Considerations

### Security

- **Never log sensitive data** (API keys, tokens, user data)
- **Validate all inputs** before processing
- **Use secure HTTP clients** with proper timeouts
- **Implement rate limiting** for API calls

### Performance

- **Use connection pooling** for HTTP clients
- **Implement caching** where appropriate
- **Handle large payloads** efficiently
- **Use context for timeouts** and cancellation

### MCP Protocol Compliance

- **Follow JSON-RPC 2.0** specification exactly
- **Return proper error codes** (see MCP specification)
- **Handle protocol negotiation** correctly
- **Support required MCP methods** (initialize, tools/list, tools/call)

## 📚 Key Dependencies

- **Go 1.22+**: Language version
- **github.com/sirupsen/logrus**: Structured logging
- **Standard library**: JSON-RPC, HTTP, etc.

## 🔍 Debugging Tips

1. **Use debug logging**: Run with `-debug` flag
2. **Test with curl**: For HTTP transport testing
3. **Use MCP test tools**: For protocol validation
4. **Check logs**: Structured logging provides detailed information
5. **Validate JSON**: Ensure proper JSON-RPC format

## 💡 When Suggesting Code

1. **Always consider MCP protocol requirements**
2. **Follow existing patterns** in the codebase
3. **Include proper error handling**
4. **Add appropriate tests**
5. **Update documentation** if needed
6. **Consider HTTP transport best practices**
7. **Follow Go idioms and best practices**

## 🎯 Common Tasks

- **Adding new tools**: Follow the tool pattern in `pkg/tools/`
- **EHQ API integration**: Use `pkg/api/EHQClient`
- **HTTP transport modifications**: Edit `internal/transport/`
- **Protocol changes**: Modify `internal/mcp/`
- **Configuration**: Environment variables and command-line flags

Remember: This is a production HTTP MCP server that needs to be reliable, secure, and performant. Always prioritize code quality, testing, and MCP protocol compliance.
