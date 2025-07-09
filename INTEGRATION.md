# MCP Server Integration Guide

## How to Invoke and Send Messages to Your EHQ MCP Server

### 1. 📋 **Command Line Testing Methods**

#### A. Quick Test (Automated)
```bash
./quick_test.sh
```

#### B. Interactive Test (Detailed)
```bash
./interactive_test.sh
```

#### C. Manual Test (Custom messages)
```bash
./manual_test.sh
```

#### D. Direct Command Line
```bash
# Single message
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./bin/ehq-mcp-server

# Multiple messages from file
./bin/ehq-mcp-server < messages.txt
```

### 2. 🤖 **Integration with AI Clients**

#### Claude Desktop Integration

1. **Locate Claude Desktop config file:**
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

2. **Add your server to the config:**
```json
{
  "mcpServers": {
    "ehq-mcp-server": {
      "command": "/Users/bhaskar.shankarling/Projects/Live/Granicus/EngagementHQ/ehq_mcp_server/bin/ehq-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```

3. **Restart Claude Desktop** - Your server will be available automatically!

#### Other MCP Clients

Most MCP clients follow a similar pattern:
```json
{
  "servers": {
    "ehq-server": {
      "command": "/path/to/ehq-mcp-server",
      "args": [],
      "transport": "stdio"
    }
  }
}
```

### 3. 🔗 **Programmatic Integration**

#### Node.js Client Example
```javascript
const { spawn } = require('child_process');

class MCPClient {
    constructor(serverPath) {
        this.server = spawn(serverPath);
        this.messageId = 1;
    }

    sendMessage(method, params = {}) {
        const message = {
            jsonrpc: "2.0",
            id: this.messageId++,
            method: method,
            params: params
        };

        this.server.stdin.write(JSON.stringify(message) + '\n');

        return new Promise((resolve) => {
            this.server.stdout.once('data', (data) => {
                resolve(JSON.parse(data.toString()));
            });
        });
    }

    async initialize() {
        return this.sendMessage('initialize', {
            protocolVersion: "2024-11-05",
            capabilities: {},
            clientInfo: { name: "node-client", version: "1.0.0" }
        });
    }

    async listTools() {
        return this.sendMessage('tools/list');
    }

    async callTool(name, arguments) {
        return this.sendMessage('tools/call', { name, arguments });
    }
}

// Usage
async function example() {
    const client = new MCPClient('./bin/ehq-mcp-server');

    await client.initialize();
    const tools = await client.listTools();
    console.log('Available tools:', tools);

    const result = await client.callTool('hello_world', { name: 'Node.js' });
    console.log('Tool result:', result);
}
```

#### Python Client Example
```python
import subprocess
import json
import threading

class MCPClient:
    def __init__(self, server_path):
        self.server = subprocess.Popen(
            [server_path],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        self.message_id = 1

    def send_message(self, method, params=None):
        message = {
            "jsonrpc": "2.0",
            "id": self.message_id,
            "method": method,
            "params": params or {}
        }
        self.message_id += 1

        self.server.stdin.write(json.dumps(message) + '\n')
        self.server.stdin.flush()

        response = self.server.stdout.readline()
        return json.loads(response)

    def initialize(self):
        return self.send_message('initialize', {
            "protocolVersion": "2024-11-05",
            "capabilities": {},
            "clientInfo": {"name": "python-client", "version": "1.0.0"}
        })

    def list_tools(self):
        return self.send_message('tools/list')

    def call_tool(self, name, arguments):
        return self.send_message('tools/call', {
            "name": name,
            "arguments": arguments
        })

# Usage
client = MCPClient('./bin/ehq-mcp-server')
client.initialize()

tools = client.list_tools()
print("Available tools:", tools)

result = client.call_tool('hello_world', {'name': 'Python'})
print("Tool result:", result)
```

### 4. 🌐 **HTTP/WebSocket Integration** (Future Enhancement)

For web applications, you might want to add HTTP or WebSocket transports:

#### HTTP Transport (Example Structure)
```go
// Add to internal/transport/http.go
func NewHTTPTransport(port int) *HTTPTransport {
    // HTTP server implementation
}
```

#### WebSocket Transport (Example Structure)
```go
// Add to internal/transport/websocket.go
func NewWebSocketTransport(port int) *WebSocketTransport {
    // WebSocket server implementation
}
```

### 5. 📋 **Message Format Reference**

#### Standard MCP Message Structure
```json
{
  "jsonrpc": "2.0",
  "id": 123,
  "method": "method_name",
  "params": {
    // method-specific parameters
  }
}
```

#### Available Methods in Your Server

| Method | Description | Parameters |
|--------|-------------|------------|
| `initialize` | Initialize connection | `protocolVersion`, `capabilities`, `clientInfo` |
| `tools/list` | List available tools | none |
| `tools/call` | Execute a tool | `name`, `arguments` |
| `resources/list` | List available resources | none |
| `resources/read` | Read a resource | `uri` |

#### Your Available Tools

| Tool | Parameters | Description |
|------|------------|-------------|
| `hello_world` | `name` (optional) | Returns greeting message |
| `echo` | `message` (required) | Echoes back the message |
| `get_time` | `format` (optional) | Returns current time |

#### Your Available Resources

| Resource URI | Description |
|-------------|-------------|
| `hello://world` | Simple hello world text |
| `info://server` | Server information |

### 6. 🔧 **Development & Debugging**

#### Run with Debug Logging
```bash
./bin/ehq-mcp-server -debug
```

#### Monitor Server Logs
```bash
# Run in one terminal
./bin/ehq-mcp-server -debug 2> server.log

# Monitor logs in another terminal
tail -f server.log
```

#### Validate JSON Messages
```bash
# Use jq to validate JSON format
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | jq .
```

### 7. 🚀 **Production Deployment**

#### Systemd Service (Linux)
```ini
[Unit]
Description=EHQ MCP Server
After=network.target

[Service]
Type=simple
User=mcpuser
ExecStart=/path/to/ehq-mcp-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Docker Container
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ehq-mcp-server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ehq-mcp-server .
CMD ["./ehq-mcp-server"]
```

---

🎯 **Choose the method that best fits your use case!**

- **Testing/Development**: Use the shell scripts
- **AI Integration**: Use Claude Desktop or other MCP clients
- **Custom Applications**: Use the programmatic examples
- **Production**: Use systemd/Docker deployment
