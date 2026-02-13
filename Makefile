.PHONY: help build run test clean docker-build docker-run install dev

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
install: ## Install dependencies
	go mod tidy

dev: ## Run the application in development mode
	go run ./cmd/api

run: ## Run the application
	go run ./cmd/api

# Building
build: ## Build the application
	@mkdir -p internal/apispec
	@cp api/openapi.yaml internal/apispec/openapi.yaml
	go build -o bin/learn-go ./cmd/api

# Testing
test: ## Run all tests
	go test ./...

test-coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Docker
docker-build: ## Build Docker image
	docker build -t learn-go .

docker-run: ## Run Docker container
	docker run -p 8080:8080 learn-go

# Cleanup
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Linting
lint: ## Run linter
	golangci-lint run

# Formatting
fmt: ## Format code
	go fmt ./...
