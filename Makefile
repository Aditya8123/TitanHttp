.PHONY: all build run test lint fmt clean dev

# Variables
BINARY_NAME=titanhttp
MAIN_PACKAGE=./cmd/titanhttp

all: lint test build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)

run: build
	@echo "Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME)

dev:
	@echo "Running in development mode (hot-reload)..."
	air

test:
	@echo "Running tests..."
	go test -v -race ./...

lint:
	@echo "Running linter..."
	golangci-lint run

fmt:
	@echo "Formatting code..."
	gofumpt -w .

clean:
	@echo "Cleaning build cache..."
	go clean
	rm -rf bin/
