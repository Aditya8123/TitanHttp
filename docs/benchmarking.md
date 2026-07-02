# Benchmarking and Performance Engineering

> "Backend engineering is deterministic, not bouncy." — TitanHTTP Design Principles

In Phase 8, TitanHTTP shifts focus from feature completeness to raw performance. Building a server from scratch allows us unprecedented control over memory allocations and CPU usage. This document outlines the metrics we care about and the tools we will use to measure them.

## Key Metrics

1. **Throughput (Requests per Second - RPS):** How many HTTP requests can the server process in one second under varying concurrent connection loads.
2. **Latency (P50, P90, P99):** The time it takes for a request to be fully processed. We care deeply about the tail latency (P99)—consistency under load is more important than absolute maximum speed.
3. **Memory Allocations (B/op, allocs/op):** In Go, reducing heap allocations directly reduces Garbage Collection (GC) pauses. Our parser and router must be heavily optimized to allocate as little memory as possible per request.

## Tooling

### Go's Built-in Benchmarks
We will use Go's standard `testing` package with `Benchmark` functions for micro-benchmarks.
- Example: Benchmarking the speed of parsing a 1KB HTTP header vs a 10KB HTTP header.
- We will track memory usage aggressively using `go test -bench . -benchmem`.

### Profiling (`pprof`)
When optimizing, we will rely on data, not intuition. Go's `pprof` tool will be used to visualize:
- **CPU Profiles:** Identifying which functions take the most execution time (e.g., string concatenation vs byte slice manipulation).
- **Heap Profiles:** Pinpointing exactly where memory is being allocated.
- **Goroutine Profiles:** Ensuring we are not leaking goroutines or getting blocked on mutex contention.

### Load Testing
To simulate real-world traffic and determine max RPS, we will use industry-standard CLI tools:
- **`wrk` or `vegeta`:** To blast the server with thousands of concurrent connections and measure the resulting throughput and tail latency.

## The Optimization Loop

Optimization will follow a strict, scientific loop:
1. **Measure:** Run a benchmark to establish a baseline.
2. **Profile:** Identify the bottleneck using `pprof`.
3. **Hypothesize & Change:** Implement a targeted code change (e.g., using a `sync.Pool` for byte buffers).
4. **Verify:** Re-run the benchmark to prove the change improved performance without introducing regressions in the tests.

---
> *Design Note: Optimization without measurement is just guessing.*
