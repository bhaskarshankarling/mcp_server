# Makefile for EHQ MCP Server

.PHONY: build run test clean deps fmt vet lint install dev help

# Variables
BINARY_NAME=ehq-mcp-server
BUILD_DIR=bin
CMD_DIR=cmd/server
MAIN_FILE=$(CMD_DIR)/main.go

# Default target
all: build

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "🔍 Vetting code..."
	go vet ./...

# Build the application
build: deps fmt vet
	@echo "🏗️  Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for development (without optimizations)
build-dev: deps fmt
	@echo "🏗️  Building $(BINARY_NAME) for development..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Development build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
run: build-dev
	@echo "🚀 Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with debug logging
run-debug: build-dev
	@echo "🚀 Running $(BINARY_NAME) with debug logging..."
	./$(BUILD_DIR)/$(BINARY_NAME) -debug

# Run directly with go run
dev:
	@echo "🔧 Running in development mode..."
	go run ./$(CMD_DIR) -debug

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Test MCP server functionality
test-mcp: build-dev
	@echo "🧪 Testing MCP server functionality..."
	./quick_test.sh

# Interactive MCP server testing
test-interactive: build-dev
	@echo "🧪 Running interactive MCP tests..."
	./interactive_test.sh

# Show MCP usage examples
examples: build-dev
	@echo "💡 Showing MCP server examples..."
	./examples.sh

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean

# Install the binary globally
install: build
	@echo "📥 Installing $(BINARY_NAME) globally..."
	go install ./$(CMD_DIR)

# Cross-platform builds
build-all: deps fmt vet
	@echo "🏗️  Building for all platforms..."
	@mkdir -p $(BUILD_DIR)

	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)

	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)

	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)

	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)

	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)

	@echo "✅ Cross-platform builds complete!"
	@ls -la $(BUILD_DIR)/

# Show version
version:
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		./$(BUILD_DIR)/$(BINARY_NAME) -version; \
	else \
		go run ./$(CMD_DIR) -version; \
	fi

# Lint code (requires golangci-lint)
lint:
	@echo "🔎 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2"; \
	fi

# Generate documentation
docs:
	@echo "📚 Generating documentation..."
	go doc -all ./... > docs.txt
	@echo "📖 Documentation generated: docs.txt"

# Quick development setup
setup: deps
	@echo "🛠️  Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "📥 Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2; \
	fi
	@echo "✅ Development environment ready!"

# Help target
help:
	@echo "🚀 EHQ MCP Server - Available Commands:"
	@echo ""
	@echo "🏗️  Building:"
	@echo "  build         Build optimized binary"
	@echo "  build-dev     Build development binary"
	@echo "  build-all     Build for all platforms"
	@echo "  install       Install binary globally"
	@echo ""
	@echo "🚀 Running:"
	@echo "  run           Build and run the server"
	@echo "  run-debug     Build and run with debug logging"
	@echo "  dev           Run directly with go run (development mode)"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test          Run all tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  test-mcp      Test MCP server functionality"
	@echo "  test-interactive  Run interactive MCP tests"
	@echo "  examples      Show MCP usage examples"
	@echo "  test-coverage Run tests with coverage report"
	@echo ""
	@echo "🔧 Development:"
	@echo "  deps          Install/update dependencies"
	@echo "  fmt           Format code"
	@echo "  vet           Vet code for issues"
	@echo "  lint          Lint code (requires golangci-lint)"
	@echo "  setup         Setup development environment"
	@echo ""
	@echo "📖 Utilities:"
	@echo "  version       Show version information"
	@echo "  docs          Generate documentation"
	@echo "  clean         Clean build artifacts"
	@echo "  help          Show this help message"
