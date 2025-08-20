.PHONY: help build test lint fmt vet clean coverage deps check-deps install-lint

# Default target
help: ## Display this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the module
	@echo "Building..."
	@go build ./...

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality targets
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: check-deps ## Run linter
	@echo "Running linter..."
	@golangci-lint run

# Development targets
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

clean: ## Clean build artifacts and coverage files
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html

# CI/CD targets
ci: fmt vet lint test ## Run all CI checks (format, vet, lint, test)
	@echo "All CI checks passed!"

# Install development tools
install-lint: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	@which golangci-lint > /dev/null 2>&1 || \
		(echo "golangci-lint not found, installing..." && \
		 curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin)

# Check if required tools are installed
check-deps: ## Check if required dependencies are installed
	@which golangci-lint > /dev/null 2>&1 || \
		(echo "❌ golangci-lint not found. Run 'make install-lint' to install it." && exit 1)
	@echo "✅ All dependencies are installed"