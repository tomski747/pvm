# Build variables
BINARY_NAME=pvm
BUILD_DIR=bin
MAIN_PACKAGE=./cmd/pvm

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/$(BUILD_DIR)
GOFILES=$(wildcard *.go)

# Go build flags
LDFLAGS=-ldflags "-s -w"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean test coverage lint help

all: test build ## Run tests and build

build: ## Build the binary
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

clean: ## Remove build directory
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean -testcache
	@echo "Cleaned!"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@mkdir -p $(BUILD_DIR)
	@go test -coverprofile=$(BUILD_DIR)/coverage.out ./...
	@go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report generated: $(BUILD_DIR)/coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint is not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

install: build ## Install binary to GOPATH/bin
	@echo "Installing..."
	@go install $(MAIN_PACKAGE)
	@echo "Installed $(BINARY_NAME) to $(GOPATH)/bin/$(BINARY_NAME)"

help: ## Display this help screen
	@echo "Available commands:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help 