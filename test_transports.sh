#!/bin/bash

# Test script for multi-transport MCP server
# Demonstrates different ways to invoke the MCP server

echo "🚀 Multi-Transport MCP Server Testing"
echo "====================================="
echo ""

cd "$(dirname "$0")"

# Build both servers
echo "📦 Building servers..."
make build > /dev/null 2>&1
go build -o bin/ehq-mcp-multi-server ./cmd/multi-server > /dev/null 2>&1

echo "🧪 Available Transport Methods:"
echo ""

echo "1️⃣  STDIO Transport (Default - works with Claude Desktop)"
echo "   ./bin/ehq-mcp-server"
echo "   echo '{\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"tools/list\", \"params\": {}}' | ./bin/ehq-mcp-server"
echo ""

echo "2️⃣  HTTP Transport (Web/API integration)"
echo "   ./bin/ehq-mcp-multi-server -http 8080"
echo "   curl -X POST http://localhost:8080/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"tools/list\", \"params\": {}}'"
echo ""

echo "3️⃣  WebSocket Transport (Real-time web apps)"
echo "   ./bin/ehq-mcp-multi-server -ws 8081"
echo "   # Connect via WebSocket to ws://localhost:8081/ws"
echo ""

echo "4️⃣  Multiple Transports (All at once)"
echo "   ./bin/ehq-mcp-multi-server -http 8080 -ws 8081"
echo ""

echo "🧪 Let's test STDIO (current method):"
echo "=====================================

echo ""
echo "🔹 Testing tools/list via STDIO:"
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list", "params": {}}' | ./bin/ehq-mcp-server 2>/dev/null | jq .

echo ""
echo "💡 Want to test HTTP/WebSocket? Run these commands:"
echo ""
echo "   # Terminal 1: Start HTTP server"
echo "   ./bin/ehq-mcp-multi-server -http 8080"
echo ""
echo "   # Terminal 2: Test HTTP"
echo "   curl -X POST http://localhost:8080/mcp \\"
echo "        -H 'Content-Type: application/json' \\"
echo "        -d '{\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"tools/list\", \"params\": {}}'"
echo ""
echo "   # Health check"
echo "   curl http://localhost:8080/health"
echo ""

echo "🎯 Summary:"
echo "==========="
echo "✅ STDIO: Ready for Claude Desktop and AI clients"
echo "🌐 HTTP: Ready for web applications and APIs"
echo "🕸️  WebSocket: Ready for real-time applications"
echo "🔧 Multi: Can run all transports simultaneously"
echo ""
echo "📋 For Claude Desktop, use the STDIO version (default behavior)"
