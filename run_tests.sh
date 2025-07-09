#!/bin/bash

# run_tests.sh - Comprehensive test runner for the EHQ MCP Server

set -e

echo "🧪 Running EHQ MCP Server Test Suite"
echo "======================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

print_status "Go version: $(go version)"

# Download dependencies
print_status "Downloading dependencies..."
go mod tidy
if [ $? -eq 0 ]; then
    print_success "Dependencies downloaded successfully"
else
    print_error "Failed to download dependencies"
    exit 1
fi

# Run tests with coverage
print_status "Running unit tests with coverage..."

# Create coverage directory if it doesn't exist
mkdir -p coverage

# Run tests and generate coverage
go test -v -race -coverprofile=coverage/coverage.out ./...

if [ $? -eq 0 ]; then
    print_success "All tests passed!"

    # Generate HTML coverage report
    print_status "Generating coverage report..."
    go tool cover -html=coverage/coverage.out -o coverage/coverage.html

    # Show coverage summary
    print_status "Coverage summary:"
    go tool cover -func=coverage/coverage.out

    print_success "Coverage report generated: coverage/coverage.html"

    # Check coverage percentage
    COVERAGE=$(go tool cover -func=coverage/coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
    COVERAGE_INT=$(echo $COVERAGE | cut -d'.' -f1)

    if [ "$COVERAGE_INT" -ge 80 ]; then
        print_success "Great! Coverage is ${COVERAGE}% (≥80%)"
    elif [ "$COVERAGE_INT" -ge 60 ]; then
        print_warning "Coverage is ${COVERAGE}% (consider improving to ≥80%)"
    else
        print_warning "Coverage is ${COVERAGE}% (should be improved to ≥80%)"
    fi

else
    print_error "Some tests failed!"
    exit 1
fi

# Run specific test categories
echo ""
print_status "Running test categories..."

echo ""
print_status "Testing MCP Core..."
go test -v ./internal/mcp/...

echo ""
print_status "Testing Tools..."
go test -v ./pkg/tools/...

echo ""
print_status "Testing Resources..."
go test -v ./pkg/resources/...

echo ""
print_status "Testing API Client..."
go test -v ./pkg/api/...

echo ""
print_status "Testing Integration..."
go test -v ./cmd/server/...

# Run race condition tests
echo ""
print_status "Running race condition tests..."
go test -race ./...

if [ $? -eq 0 ]; then
    print_success "No race conditions detected!"
else
    print_error "Race conditions detected!"
    exit 1
fi

# Run benchmarks if available
echo ""
print_status "Running benchmarks..."
go test -bench=. -benchmem ./... || print_warning "No benchmarks found"

# Vet the code
echo ""
print_status "Running go vet..."
go vet ./...

if [ $? -eq 0 ]; then
    print_success "go vet passed!"
else
    print_error "go vet found issues!"
    exit 1
fi

# Check formatting
echo ""
print_status "Checking code formatting..."
UNFORMATTED=$(gofmt -l .)

if [ -z "$UNFORMATTED" ]; then
    print_success "All Go files are properly formatted!"
else
    print_warning "The following files need formatting:"
    echo "$UNFORMATTED"
    print_warning "Run 'gofmt -w .' to fix formatting"
fi

echo ""
print_success "🎉 All tests completed successfully!"
echo ""
print_status "Test Summary:"
echo "  ✅ Unit tests passed"
echo "  ✅ Race condition tests passed"
echo "  ✅ Code vetting passed"
echo "  📊 Coverage report: coverage/coverage.html"
echo ""
