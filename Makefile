
# Variables
BINARY_NAME=tetris-optimizer
MAIN_PATH=./cmd
GO_FILES=$(shell find . -name "*.go")

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Run with arguments
.PHONY: run-with-args
run-with-args: build
	@echo "Running $(BINARY_NAME) with arguments..."
	@./$(BINARY_NAME) $(ARGS)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...

# Generate detailed coverage report
.PHONY: coverage-report
coverage-report:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out
	@rm -f coverage.html

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@go vet ./...

# Run all quality checks
.PHONY: check
check: fmt lint test

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Initialize module (run once)
.PHONY: init
init:
	@echo "Initializing Go module..."
	@go mod init github.com/stkisengese/tetris-optimizer

# Development target - format, lint, test, and build
.PHONY: dev
dev: fmt lint test build

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  run-with-args - Run with arguments (use ARGS=...)"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  coverage-report - Generate HTML coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  check         - Run all quality checks"
	@echo "  deps          - Install dependencies"
	@echo "  init          - Initialize Go module"
	@echo "  dev           - Development workflow (fmt, lint, test, build)"
	@echo "  help          - Show this help message"
