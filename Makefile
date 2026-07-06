.PHONY: all build run test lint fmt clean dev bench-micro bench-component bench-load bench-stress bench-soak profile-cpu profile-heap bench-all

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

# --- Benchmarking Framework ---

bench-all: bench-micro bench-component

bench-micro:
	@echo "Running Micro Benchmarks..."
	go test -bench=. -benchmem ./internal/... > benchmarks/micro/micro_latest.bench
	@echo "Results saved to benchmarks/micro/micro_latest.bench"

bench-component:
	@echo "Running Component Benchmarks..."
	@echo "Component benchmarks will be added here" > benchmarks/component/component_latest.bench

bench-load:
	@echo "Running Load Tests..."
	@echo "Load tests require an external tool (e.g. wrk) and will be implemented in a future task"

bench-stress:
	@echo "Running Stress Tests..."
	@echo "Stress tests will be implemented in a future task"

bench-soak:
	@echo "Running Soak Tests..."
	@echo "Soak tests will be implemented in a future task"

profile-cpu:
	@echo "Capturing CPU Profile (30s)..."
	go tool pprof -raw http://localhost:6060/debug/pprof/profile?seconds=30 > benchmarks/profiles/cpu_latest.pprof
	@echo "Saved CPU profile to benchmarks/profiles/cpu_latest.pprof"

profile-heap:
	@echo "Capturing Heap Profile..."
	go tool pprof -raw http://localhost:6060/debug/pprof/heap > benchmarks/profiles/heap_latest.pprof
	@echo "Saved Heap profile to benchmarks/profiles/heap_latest.pprof"
