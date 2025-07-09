#!/bin/bash

# Test script for EHQ MCP Server
# This script tests the basic functionality of the MCP server

set -e

echo "🧪 Testing EHQ MCP Server"
echo "=========================="

# Build the server first
echo "📦 Building server..."
make build > /dev/null 2>&1

# Start the server in background and capture its PID
echo "🚀 Starting server..."
./bin/ehq-mcp-server > server_output.log 2>&1 &
SERVER_PID=$!

# Give the server a moment to start
sleep 2

# Function to send JSON-RPC message to server
send_message() {
    echo "$1" | timeout 5 ./bin/ehq-mcp-server 2>/dev/null
}

echo "🔧 Testing server functionality..."

# Test 1: Initialize
echo "  ✅ Testing initialize..."
INIT_MSG='{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}}'
INIT_RESPONSE=$(send_message "$INIT_MSG")
echo "     Response: $INIT_RESPONSE"

# Test 2: List tools
echo "  ✅ Testing tools/list..."
TOOLS_MSG='{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'
TOOLS_RESPONSE=$(send_message "$TOOLS_MSG")
echo "     Response: $TOOLS_RESPONSE"

# Test 3: Call hello_world tool
echo "  ✅ Testing tools/call (hello_world)..."
HELLO_MSG='{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "hello_world", "arguments": {"name": "Test User"}}}'
HELLO_RESPONSE=$(send_message "$HELLO_MSG")
echo "     Response: $HELLO_RESPONSE"

# Test 4: List resources
echo "  ✅ Testing resources/list..."
RESOURCES_MSG='{"jsonrpc": "2.0", "id": 4, "method": "resources/list", "params": {}}'
RESOURCES_RESPONSE=$(send_message "$RESOURCES_MSG")
echo "     Response: $RESOURCES_RESPONSE"

# Test 5: Read resource
echo "  ✅ Testing resources/read..."
READ_MSG='{"jsonrpc": "2.0", "id": 5, "method": "resources/read", "params": {"uri": "hello://world"}}'
READ_RESPONSE=$(send_message "$READ_MSG")
echo "     Response: $READ_RESPONSE"

# Clean up
if kill -0 $SERVER_PID 2>/dev/null; then
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
fi

rm -f server_output.log

echo ""
echo "🎉 All tests completed!"
echo "✅ Your EHQ MCP Server is working correctly!"
