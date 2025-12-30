.PHONY: help build build-all test test-verbose coverage clean install lint run

# Variables
BINARY_NAME=journal
VERSION?=0.1.0
BUILD_DIR=bin
MAIN_PATH=cmd/journal/main.go
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)"

# Default target
.DEFAULT_GOAL := help

help: ## Display this help message
	@echo "Journal CLI - Build System"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary for the current platform
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build binaries for all platforms (Linux, macOS, Windows)
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	
	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "All builds complete!"
	@ls -lh $(BUILD_DIR)

test: ## Run all tests
	@echo "Running tests..."
	go test -v -race -coverprofile=$(COVERAGE_FILE) ./...
	@echo ""
	@echo "Coverage summary:"
	@go tool cover -func=$(COVERAGE_FILE) | grep total

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	go test -v -race -coverprofile=$(COVERAGE_FILE) ./...

coverage: test ## Generate HTML coverage report
	@echo "Generating coverage report..."
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@echo "Opening in browser..."
	@which xdg-open > /dev/null && xdg-open $(COVERAGE_HTML) || \
	 which open > /dev/null && open $(COVERAGE_HTML) || \
	 echo "Please open $(COVERAGE_HTML) manually"

clean: ## Remove build artifacts and test files
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "Clean complete!"

install: build ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) $(MAIN_PATH)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	@echo "Running linter..."
	golangci-lint run ./...

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated!"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted!"

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "Vet complete!"
