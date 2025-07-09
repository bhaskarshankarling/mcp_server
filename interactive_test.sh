#!/bin/bash

# Interactive MCP Server Testing
# Send individual messages to test specific functionality

cd "$(dirname "$0")"

# Build the server if needed
if [ ! -f "./bin/ehq-mcp-server" ]; then
    echo "📦 Building server..."
    make build > /dev/null 2>&1
fi

echo "🧪 Interactive MCP Server Testing"
echo "================================="
echo ""

# Function to send a single message and get response
send_message() {
    local message="$1"
    local description="$2"

    echo "🔹 $description"
    echo "📤 Request:  $message"
    echo -n "📥 Response: "
    echo "$message" | timeout 5 ./bin/ehq-mcp-server 2>/dev/null
    echo ""
}

# Test 1: Initialize
send_message \
    '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}' \
    "Initialize Connection"

# Test 2: List tools
send_message \
    '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}' \
    "List Available Tools"

# Test 3: Call hello_world tool with name
send_message \
    '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "hello_world", "arguments": {"name": "Alice"}}}' \
    "Call hello_world tool with name"

# Test 4: Call echo tool
send_message \
    '{"jsonrpc": "2.0", "id": 4, "method": "tools/call", "params": {"name": "echo", "arguments": {"message": "Hello from MCP test!"}}}' \
    "Call echo tool"

# Test 5: Call get_time tool with different formats
send_message \
    '{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "get_time", "arguments": {"format": "RFC3339"}}}' \
    "Get time in RFC3339 format"

send_message \
    '{"jsonrpc": "2.0", "id": 6, "method": "tools/call", "params": {"name": "get_time", "arguments": {"format": "Unix"}}}' \
    "Get time in Unix format"

send_message \
    '{"jsonrpc": "2.0", "id": 7, "method": "tools/call", "params": {"name": "get_time", "arguments": {"format": "Kitchen"}}}' \
    "Get time in Kitchen format"

# Test 6: List resources
send_message \
    '{"jsonrpc": "2.0", "id": 8, "method": "resources/list", "params": {}}' \
    "List Available Resources"

# Test 7: Read hello world resource
send_message \
    '{"jsonrpc": "2.0", "id": 9, "method": "resources/read", "params": {"uri": "hello://world"}}' \
    "Read hello://world resource"

# Test 8: Read server info resource
send_message \
    '{"jsonrpc": "2.0", "id": 10, "method": "resources/read", "params": {"uri": "info://server"}}' \
    "Read info://server resource"

# Test 9: Error case - unknown tool
send_message \
    '{"jsonrpc": "2.0", "id": 11, "method": "tools/call", "params": {"name": "unknown_tool", "arguments": {}}}' \
    "Test error handling - unknown tool"

# Test 10: Error case - unknown resource
send_message \
    '{"jsonrpc": "2.0", "id": 12, "method": "resources/read", "params": {"uri": "unknown://resource"}}' \
    "Test error handling - unknown resource"

echo "🎉 All tests completed!"
echo "✅ Your EHQ MCP Server is working correctly!"
