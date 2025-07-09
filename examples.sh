#!/bin/bash

# Simple one-liner examples for testing MCP server
# These demonstrate the exact command format for invoking your server

echo "🚀 Simple MCP Server Invocation Examples"
echo "========================================"
echo ""

cd "$(dirname "$0")"

# Ensure server is built
make build > /dev/null 2>&1

echo "💡 Here are simple one-liner commands to test your MCP server:"
echo ""

echo "1️⃣  Initialize connection:"
echo "   echo '{\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"initialize\", \"params\": {\"protocolVersion\": \"2024-11-05\", \"capabilities\": {}, \"clientInfo\": {\"name\": \"test\", \"version\": \"1.0\"}}}' | ./bin/ehq-mcp-server"
echo ""

echo "2️⃣  List all available tools:"
echo "   echo '{\"jsonrpc\": \"2.0\", \"id\": 2, \"method\": \"tools/list\", \"params\": {}}' | ./bin/ehq-mcp-server"
echo ""

echo "3️⃣  Call hello_world tool:"
echo "   echo '{\"jsonrpc\": \"2.0\", \"id\": 3, \"method\": \"tools/call\", \"params\": {\"name\": \"hello_world\", \"arguments\": {\"name\": \"World\"}}}' | ./bin/ehq-mcp-server"
echo ""

echo "4️⃣  Call echo tool:"
echo "   echo '{\"jsonrpc\": \"2.0\", \"id\": 4, \"method\": \"tools/call\", \"params\": {\"name\": \"echo\", \"arguments\": {\"message\": \"Hello MCP!\"}}}' | ./bin/ehq-mcp-server"
echo ""

echo "5️⃣  Get current time:"
echo "   echo '{\"jsonrpc\": \"2.0\", \"id\": 5, \"method\": \"tools/call\", \"params\": {\"name\": \"get_time\", \"arguments\": {\"format\": \"Kitchen\"}}}' | ./bin/ehq-mcp-server"
echo ""

echo "6️⃣  List resources:"
echo "   echo '{\"jsonrpc\": \"2.0\", \"id\": 6, \"method\": \"resources/list\", \"params\": {}}' | ./bin/ehq-mcp-server"
echo ""

echo "7️⃣  Read a resource:"
echo "   echo '{\"jsonrpc\": \"2.0\", \"id\": 7, \"method\": \"resources/read\", \"params\": {\"uri\": \"hello://world\"}}' | ./bin/ehq-mcp-server"
echo ""

echo "🧪 Let's try a few of these now:"
echo "================================"

echo ""
echo "🔹 Testing tool list:"
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}' | ./bin/ehq-mcp-server | jq .

echo ""
echo "🔹 Testing hello_world tool:"
echo '{"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {"name": "hello_world", "arguments": {"name": "CLI User"}}}' | ./bin/ehq-mcp-server | jq .

echo ""
echo "🔹 Testing get_time tool:"
echo '{"jsonrpc": "2.0", "id": 5, "method": "tools/call", "params": {"name": "get_time", "arguments": {"format": "Kitchen"}}}' | ./bin/ehq-mcp-server | jq .

echo ""
echo "✅ Your server is ready! Copy and paste any of the commands above to test manually."
