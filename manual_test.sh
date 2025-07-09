#!/bin/bash

# Manual MCP Server Testing Script
# Allows you to send custom messages interactively

cd "$(dirname "$0")"

echo "🔧 Manual MCP Server Testing"
echo "============================="
echo ""
echo "This script will start the MCP server and allow you to send custom messages."
echo "Type your JSON-RPC messages and press Enter."
echo "Type 'quit' or 'exit' to stop."
echo ""

# Build if needed
if [ ! -f "./bin/ehq-mcp-server" ]; then
    echo "📦 Building server..."
    make build > /dev/null 2>&1
fi

echo "🚀 Starting MCP server..."
echo "📝 You can now type JSON-RPC messages:"
echo ""

# Examples to show user
echo "💡 Example messages you can try:"
echo "   Initialize: {\"jsonrpc\": \"2.0\", \"id\": 1, \"method\": \"initialize\", \"params\": {\"protocolVersion\": \"2024-11-05\", \"capabilities\": {}, \"clientInfo\": {\"name\": \"manual-test\", \"version\": \"1.0.0\"}}}"
echo "   List tools: {\"jsonrpc\": \"2.0\", \"id\": 2, \"method\": \"tools/list\", \"params\": {}}"
echo "   Call tool:  {\"jsonrpc\": \"2.0\", \"id\": 3, \"method\": \"tools/call\", \"params\": {\"name\": \"hello_world\", \"arguments\": {\"name\": \"YourName\"}}}"
echo ""
echo "📨 Enter your messages (one per line):"
echo "----------------------------------------"

# Create a named pipe for bidirectional communication
PIPE=$(mktemp -u)
mkfifo "$PIPE"

# Start the server with the pipe
./bin/ehq-mcp-server < "$PIPE" &
SERVER_PID=$!

# Function to cleanup
cleanup() {
    echo ""
    echo "🛑 Shutting down..."
    kill $SERVER_PID 2>/dev/null || true
    rm -f "$PIPE"
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Read user input and send to server
exec 3>"$PIPE"
while true; do
    read -r input

    if [ "$input" = "quit" ] || [ "$input" = "exit" ]; then
        break
    fi

    if [ -n "$input" ]; then
        echo "$input" >&3
        sleep 0.1  # Give server time to respond
    fi
done

cleanup
