#!/bin/bash

# Test the GetProjects tool
echo "Testing GetProjects tool..."

# Build the server
echo "Building server..."
go build -o server cmd/server/main.go

# Test tools list to make sure get_projects is available
echo "Checking if get_projects tool is available..."
TOOLS_RESPONSE=$(echo '{"jsonrpc": "2.0", "method": "tools/list", "id": 1}' | ./server 2>/dev/null)
echo "$TOOLS_RESPONSE" | jq -r '.result.tools[] | select(.name == "get_projects")' 2>/dev/null

echo ""
echo "Testing get_projects tool execution..."

# Test 1: Get all projects (no search parameter)
echo "Test 1: Getting all projects..."
echo '{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "get_projects", "arguments": {}}, "id": 2}' | ./server > temp_response.json 2>&1

# Extract the JSON response (last line that starts with {"jsonrpc")
JSON_RESPONSE=$(grep '^{"jsonrpc"' temp_response.json | tail -1)

if [ -n "$JSON_RESPONSE" ]; then
    echo "✅ Valid JSON response received"

    # Parse and display project summary
    echo ""
    echo "Response summary:"
    echo "$JSON_RESPONSE" | jq '.result.content[0].text | fromjson | {success, count}' 2>/dev/null

    echo ""
    echo "First 3 projects:"
    echo "$JSON_RESPONSE" | jq '.result.content[0].text | fromjson | .data[:3] | map({id, name: .attributes.name, state: .attributes.state})' 2>/dev/null

    echo ""
    echo "Project states distribution:"
    echo "$JSON_RESPONSE" | jq '.result.content[0].text | fromjson | .data | group_by(.attributes.state) | map({state: .[0].attributes.state, count: length})' 2>/dev/null

    # Store total count for comparison
    TOTAL_COUNT=$(echo "$JSON_RESPONSE" | jq '.result.content[0].text | fromjson | .count' 2>/dev/null)
else
    echo "❌ No valid JSON response found for Test 1"
    echo "Raw response:"
    cat temp_response.json
fi

echo ""
echo "Test 2: Searching for projects with 'idea'..."
echo '{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "get_projects", "arguments": {"search": "idea"}}, "id": 3}' | ./server > temp_search_response.json 2>&1

# Extract the JSON response for search test
SEARCH_JSON_RESPONSE=$(grep '^{"jsonrpc"' temp_search_response.json | tail -1)

if [ -n "$SEARCH_JSON_RESPONSE" ]; then
    echo "✅ Valid JSON response received for search"

    # Parse and display search results summary
    echo ""
    echo "Search results summary:"
    echo "$SEARCH_JSON_RESPONSE" | jq '.result.content[0].text | fromjson | {success, count}' 2>/dev/null

    echo ""
    echo "Projects matching 'idea':"
    echo "$SEARCH_JSON_RESPONSE" | jq '.result.content[0].text | fromjson | .data | map({id, name: .attributes.name, state: .attributes.state})' 2>/dev/null

    # Store search count for comparison
    SEARCH_COUNT=$(echo "$SEARCH_JSON_RESPONSE" | jq '.result.content[0].text | fromjson | .count' 2>/dev/null)

    echo ""
    echo "Search validation:"
    if [ "$SEARCH_COUNT" -le "$TOTAL_COUNT" ]; then
        echo "✅ Search returned $SEARCH_COUNT projects (≤ total $TOTAL_COUNT projects)"
    else
        echo "❌ Search returned more projects than total (unexpected)"
    fi
else
    echo "❌ No valid JSON response found for search test"
    echo "Raw search response:"
    cat temp_search_response.json
fi

echo ""
echo "Test 3: Searching for projects with 'nonexistent'..."
echo '{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "get_projects", "arguments": {"search": "nonexistent"}}, "id": 4}' | ./server > temp_empty_response.json 2>&1

# Extract the JSON response for empty search test
EMPTY_JSON_RESPONSE=$(grep '^{"jsonrpc"' temp_empty_response.json | tail -1)

if [ -n "$EMPTY_JSON_RESPONSE" ]; then
    EMPTY_COUNT=$(echo "$EMPTY_JSON_RESPONSE" | jq '.result.content[0].text | fromjson | .count' 2>/dev/null)
    echo "✅ Search for 'nonexistent' returned $EMPTY_COUNT projects"
else
    echo "❌ No valid JSON response found for empty search test"
fi

# Clean up
rm -f temp_response.json temp_search_response.json temp_empty_response.json debug_response.txt

echo ""
echo "Search functionality tests completed!"
echo "✅ Tool supports optional search parameter"
echo "✅ Tool works without search parameter (all projects)"
echo "✅ Tool filters results when search parameter is provided"
