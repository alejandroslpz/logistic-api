.PHONY: build run test clean help

# Variables
BINARY_NAME=logistics-api
BUILD_DIR=bin

# Build the application
build:
	@echo "📦 Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

# Run the application
run: build
	@echo "🚀 Starting $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run with live reload (if you have air installed)
dev:
	@echo "🔥 Starting development server..."
	@air

# Test the application
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Install development dependencies
install-dev:
	@echo "📥 Installing development dependencies..."
	@go install github.com/cosmtrek/air@latest

# Generate JWT secret
gen-jwt:
	@echo "🔐 Generating JWT secret..."
	@openssl rand -base64 32

# Database status
db-status:
	@echo "📊 Checking database connection..."
	@curl -s http://localhost:8080/health | jq .

# Show help
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Build and run the application"
	@echo "  dev        - Run with live reload"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  install-dev- Install development dependencies"
	@echo "  gen-jwt    - Generate JWT secret"
	@echo "  db-status  - Check database connection"
	@echo "  help       - Show this help"