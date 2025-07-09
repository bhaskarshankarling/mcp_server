# EHQ MCP Server Examples

This document provides examples of how to interact with the EHQ MCP Server.

## Running the Server

### Basic Usage
```bash
# Build and run
make run

# Run with debug logging
make run-debug

# Run directly with go
go run ./cmd/server -debug
```

### Testing the Server

```bash
# Run our automated test
./quick_test.sh

# Manual testing with individual commands
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' | ./bin/ehq-mcp-server
```

## Example MCP Messages

### 1. Initialize Connection
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": {},
    "clientInfo": {
      "name": "test-client",
      "version": "1.0.0"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "logging": {},
      "resources": {},
      "tools": {}
    },
    "serverInfo": {
      "name": "EHQ MCP Server",
      "version": "1.0.0"
    }
  }
}
```

### 2. List Available Tools
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/list",
  "params": {}
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "tools": [
      {
        "name": "hello_world",
        "description": "Returns a friendly hello world message",
        "inputSchema": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string",
              "description": "Name to greet (optional)"
            }
          }
        }
      },
      {
        "name": "echo",
        "description": "Echoes back the provided message",
        "inputSchema": {
          "type": "object",
          "properties": {
            "message": {
              "type": "string",
              "description": "Message to echo back"
            }
          },
          "required": ["message"]
        }
      },
      {
        "name": "get_time",
        "description": "Returns the current date and time",
        "inputSchema": {
          "type": "object",
          "properties": {
            "format": {
              "type": "string",
              "description": "Time format (optional, defaults to RFC3339)",
              "enum": ["RFC3339", "Unix", "Kitchen"]
            }
          }
        }
      }
    ]
  }
}
```

### 3. Call a Tool
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "hello_world",
    "arguments": {
      "name": "Developer"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Hello, Developer! 🌍\nWelcome to the EHQ MCP Server!\nCurrent time: 2025-07-04 16:53:35"
      }
    ]
  }
}
```

### 4. List Available Resources
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "resources/list",
  "params": {}
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "result": {
    "resources": [
      {
        "uri": "hello://world",
        "name": "Hello World",
        "description": "A simple hello world resource",
        "mimeType": "text/plain"
      },
      {
        "uri": "info://server",
        "name": "Server Information",
        "description": "Information about this MCP server",
        "mimeType": "text/plain"
      }
    ]
  }
}
```

### 5. Read a Resource
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "resources/read",
  "params": {
    "uri": "hello://world"
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "result": {
    "contents": [
      {
        "uri": "hello://world",
        "mimeType": "text/plain",
        "text": "Hello, World! This is a sample resource from the MCP server."
      }
    ]
  }
}
```

## Get Projects Tool Example

The `get_projects` tool fetches project data from the EHQ API.

### Tool Usage

```bash
# Test the get_projects tool
echo '{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "get_projects", "arguments": {}}, "id": 1}' | ./server

# Expected response format (success):
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [{
      "type": "text",
      "text": "{\"success\":true,\"data\":[{\"type\":\"projects\",\"attributes\":{\"name\":\"Project 1\"}},{\"type\":\"projects\",\"attributes\":{\"name\":\"Project 2\"}}],\"count\":2}"
    }]
  }
}

# Expected response format (authentication error):
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [{
      "type": "text",
      "text": "{\"error\":\"Authentication failed: authentication failed with status 403\"}"
    }]
  }
}
```

### Configuration

To use with real credentials, update the hardcoded values in `internal/mcp/server.go`:

```go
// In ExecuteGetProjects method:
client := api.NewEHQClient("https://your-ehq-domain.com")
err := client.Authenticate("your_username", "your_password")
```

### API Flow

1. **Authentication**: `POST /api/v2/tokens` with login/password
2. **Fetch Projects**: `GET /api/v2/projects?filterable=true` with JWT token
3. **Response**: JSON-API compliant format with project data
````
