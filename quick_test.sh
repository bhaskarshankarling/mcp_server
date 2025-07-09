#!/bin/bash

# Simple MCP Server Test
# Tests basic functionality by sending JSON-RPC messages via stdin

echo "🧪 Testing EHQ MCP Server via stdin/stdout"
echo "==========================================="

cd "$(dirname "$0")"

# Build the server
echo "📦 Building server..."
make build > /dev/null 2>&1

echo "🚀 Testing server responses..."

# Create a temporary file with test messages
cat > test_messages.txt << 'EOF'
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}
{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}
{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "hello_world", "arguments": {"name": "Test User"}}}
{"jsonrpc": "2.0", "id": 4, "method": "resources/list", "params": {}}
{"jsonrpc": "2.0", "id": 5, "method": "resources/read", "params": {"uri": "hello://world"}}
EOF

echo "📨 Sending test messages to server..."
echo ""

# Run the server with our test input
timeout 10 ./bin/ehq-mcp-server < test_messages.txt 2>/dev/null

echo ""
echo "🎉 Test completed!"
echo "✅ If you see JSON responses above, your server is working!"

# Clean up
rm -f test_messages.txt

echo ""
echo "💡 To use this server with an AI client:"
echo "   Add this to your MCP client configuration:"
echo "   {"
echo "     \"command\": \"$(pwd)/bin/ehq-mcp-server\","
echo "     \"args\": []"
echo "   }"
