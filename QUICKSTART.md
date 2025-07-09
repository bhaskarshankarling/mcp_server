# 🎯 **How to Invoke and Send Messages to Your EHQ MCP Server**

## 📋 **Quick Reference**

Your MCP server accepts JSON-RPC 2.0 messages via **stdin** and responds via **stdout**.

### **Basic Invocation Pattern:**
```bash
echo 'JSON_MESSAGE' | ./bin/ehq-mcp-server
```

---

## 🚀 **Available Methods**

### **1. Testing Commands (Built-in)**
```bash
make test-mcp          # Automated functional test
make test-interactive  # Detailed interactive test
make examples         # Show usage examples
./quick_test.sh       # Quick automated test
./interactive_test.sh # Comprehensive test suite
./examples.sh         # One-liner examples
./manual_test.sh      # Interactive manual testing
```

### **2. Direct Command Examples**

#### **📋 List Available Tools:**
```bash
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./bin/ehq-mcp-server
```

#### **🔧 Call a Tool:**
```bash
# Hello World Tool
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/call", "params": {"name": "hello_world", "arguments": {"name": "Your Name"}}}' | ./bin/ehq-mcp-server

# Echo Tool
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "echo", "arguments": {"message": "Hello MCP!"}}}' | ./bin/ehq-mcp-server

# Get Time Tool
echo '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "get_time", "arguments": {"format": "Kitchen"}}}' | ./bin/ehq-mcp-server
```

#### **📄 List Resources:**
```bash
echo '{"jsonrpc": "2.0", "id": 5, "method": "resources/list", "params": {}}' | ./bin/ehq-mcp-server
```

#### **📖 Read a Resource:**
```bash
echo '{"jsonrpc": "2.0", "id": 6, "method": "resources/read", "params": {"uri": "hello://world"}}' | ./bin/ehq-mcp-server
```

---

## 🤖 **AI Client Integration**

### **Claude Desktop Configuration:**
1. Open Claude Desktop settings
2. Add this to your MCP configuration:
```json
{
  "mcpServers": {
    "ehq-mcp-server": {
      "command": "/Users/bhaskar.shankarling/Projects/Live/Granicus/EngagementHQ/ehq_mcp_server/bin/ehq-mcp-server",
      "args": []
    }
  }
}
```
3. Restart Claude Desktop
4. Your tools will be available automatically! ✨

### **Other MCP Clients:**
Most follow similar patterns - just point to your binary:
```json
{
  "command": "/path/to/ehq-mcp-server",
  "transport": "stdio"
}
```

---

## 💻 **Programmatic Integration**

### **Node.js Example:**
```javascript
const { spawn } = require('child_process');

const server = spawn('./bin/ehq-mcp-server');

// Send message
server.stdin.write(JSON.stringify({
  jsonrpc: "2.0",
  id: 1,
  method: "tools/call",
  params: { name: "hello_world", arguments: { name: "Node.js" } }
}) + '\n');

// Read response
server.stdout.on('data', (data) => {
  console.log('Response:', JSON.parse(data.toString()));
});
```

### **Python Example:**
```python
import subprocess
import json

server = subprocess.Popen(['./bin/ehq-mcp-server'],
                         stdin=subprocess.PIPE,
                         stdout=subprocess.PIPE,
                         text=True)

# Send message
message = {
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {"name": "hello_world", "arguments": {"name": "Python"}}
}
server.stdin.write(json.dumps(message) + '\n')
server.stdin.flush()

# Read response
response = server.stdout.readline()
print("Response:", json.loads(response))
```

### **cURL Alternative (via Named Pipe):**
```bash
# Create named pipe
mkfifo mcp_pipe

# Start server with pipe
./bin/ehq-mcp-server < mcp_pipe &

# Send messages
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' > mcp_pipe
```

---

## 🛠️ **Available Tools & Resources**

### **🔧 Tools:**
| Tool | Parameters | Description |
|------|------------|-------------|
| `hello_world` | `name` (optional) | Friendly greeting |
| `echo` | `message` (required) | Echo back message |
| `get_time` | `format` (optional) | Current time in various formats |

### **📄 Resources:**
| URI | Description |
|-----|-------------|
| `hello://world` | Simple hello world text |
| `info://server` | Server information |

---

## 🔍 **Message Format**

### **Request Structure:**
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

### **Response Structure:**
```json
{
  "jsonrpc": "2.0",
  "id": 123,
  "result": {
    // method-specific result
  }
}
```

### **Error Structure:**
```json
{
  "jsonrpc": "2.0",
  "id": 123,
  "error": {
    "code": -32601,
    "message": "Method not found"
  }
}
```

---

## 🚀 **Next Steps**

1. **Test your server:** `make test-mcp`
2. **Try examples:** `make examples`
3. **Integrate with Claude:** Use the configuration above
4. **Extend functionality:** Add more tools in `pkg/tools/`
5. **Add resources:** Create new resources in `pkg/resources/`

---

## 📚 **More Information**

- **INTEGRATION.md**: Detailed integration guide
- **EXAMPLES.md**: Development examples and patterns
- **README.md**: Complete project documentation

🎉 **Your EHQ MCP Server is ready to power AI interactions!**
