.PHONY: all build run test lint fmt clean dev bench-micro bench-component bench-load bench-stress bench-soak profile-cpu profile-heap bench-all install-tools

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

install-tools:
	@echo "Installing bombardier load testing tool..."
	go install github.com/codesenberg/bombardier@latest
	@echo "Installing benchstat for regression analysis..."
	go install golang.org/x/perf/cmd/benchstat@latest
	@echo "Tools installed successfully to your GOPATH/bin."

bench-all: bench-micro bench-component

bench-micro:
	@echo "Running Micro Benchmarks..."
	go test -bench="." -benchmem ./internal/... > benchmarks/results/micro/micro_latest.bench
	@echo "Generating CPU/Memory profiles for internal/http..."
	go test -bench="." -benchmem -cpuprofile=benchmarks/profiles/micro_cpu.prof -memprofile=benchmarks/profiles/micro_mem.prof ./internal/http
	@echo "Results saved to benchmarks/results/micro/micro_latest.bench and profiles to benchmarks/profiles/"

bench-component:
	@echo "Running Component Benchmarks..."
	go test -bench="." -benchmem -run=^$$ ./internal/server/... > benchmarks/results/component/component_latest.bench

bench-load:
	@echo "Running Load Tests..."
	powershell -File ./scripts/bench/load_test.ps1

bench-stress:
	@echo "Running Stress Tests..."
	powershell -File ./scripts/bench/stress_test.ps1

bench-soak:
	@echo "Running Soak Tests..."
	powershell -File ./scripts/bench/soak_test.ps1

profile-cpu:
	@echo "Capturing CPU Profile (30s)..."
	go tool pprof -raw http://localhost:6060/debug/pprof/profile?seconds=30 > benchmarks/profiles/cpu_latest.pprof
	@echo "Generating CPU SVG graph..."
	go tool pprof -svg benchmarks/profiles/cpu_latest.pprof > benchmarks/profiles/cpu_latest.svg
	@echo "Saved CPU profile and SVG to benchmarks/profiles/"

profile-heap:
	@echo "Capturing Heap Profile..."
	go tool pprof -raw http://localhost:6060/debug/pprof/heap > benchmarks/profiles/heap_latest.pprof
	@echo "Generating Heap SVG graph..."
	go tool pprof -svg benchmarks/profiles/heap_latest.pprof > benchmarks/profiles/heap_latest.svg
	@echo "Saved Heap profile and SVG to benchmarks/profiles/"

profile-hotspots:
	@echo "Analyzing top 20 CPU hotspots..."
	go tool pprof -top -cum -nodecount=20 benchmarks/profiles/cpu_latest.pprof
