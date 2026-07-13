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
- **`bombardier`:** To blast the server with thousands of concurrent connections and measure the resulting throughput and tail latency.

## The Optimization Loop

Optimization will follow a strict, scientific loop:
1. **Measure:** Run a benchmark to establish a baseline.
2. **Profile:** Identify the bottleneck using `pprof`.
3. **Hypothesize & Change:** Implement a targeted code change (e.g., using a `sync.Pool` for byte buffers).
4. **Verify:** Re-run the benchmark to prove the change improved performance without introducing regressions in the tests.

---
> *Design Note: Optimization without measurement is just guessing.*

## Final Benchmark Results (Cinematic Suite)

The following metrics were generated automatically by our benchmarking suite running across Micro Benchmarks, Framework Comparisons, Unified Functional Tests, Compliance, and Security tests.

### 1. Unified Cinematic Suite (TitanHTTP vs net/http)

#### Concurrency Scaling Performance
| Concurrency (Connections) | Server | Throughput (Req/s) | P50 Latency (ms) | P99 Latency (ms) | TitanHTTP Speedup |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **50 Connections** | **TitanHTTP** | **99,390** | **0.50** | **3.02** | **+41.1%** |
| | net/http | 70,403 | 0.54 | 5.04 | |
| **100 Connections** | **TitanHTTP** | **112,567** | **0.84** | **3.22** | **+28.1%** |
| | net/http | 87,843 | 0.63 | 6.74 | |
| **200 Connections** | **TitanHTTP** | **101,621** | **1.59** | **7.55** | **+15.5%** |
| | net/http | 87,969 | 1.16 | 18.17 | |

#### Standard Workload Profiles
| Workload | Server | Throughput (Req/s) | P50 Latency (ms) | P99 Latency (ms) | RPS Difference |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Concurrent Connections** | **TitanHTTP** | **68,172** | **0.52** | **537.83** | **+79.1%** |
| | net/http | 38,043 | 0.66 | 556.85 | |
| **Throughput (Ping)** | **TitanHTTP** | **123,483** | **1.00** | **4.76** | **+30.3%** |
| | net/http | 94,769 | 0.66 | 11.06 | |
| **Static File Serving** | **TitanHTTP** | **62,584** | **1.57** | **8.89** | **+28.1%** |
| | net/http | 48,844 | 1.18 | 20.97 | |
| **Keep-Alive Performance** | **TitanHTTP** | **115,524** | **1.01** | **5.07** | **+17.9%** |
| | net/http | 97,984 | 0.93 | 7.06 | |
| **Routing Performance** | **TitanHTTP** | **95,290** | **1.09** | **6.06** | **+25.7%** |
| | net/http | 75,753 | 0.70 | 13.57 | |
| **Large Payloads (10KB)** | **74,266** | **1.09** | **11.55** | **+112.2%** |
| | net/http | 34,990 | 1.19 | 27.83 | |

#### High-Churn and TLS Infrastructure
| Profile | Server | Throughput (Req/s) | P50 Latency (ms) | P99 Latency (ms) | RPS Difference |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Connection Churn** | TitanHTTP | 3,365 | 29.71 | 55.43 | -5.7% |
| (Connection: close) | **net/http** | **3,568** | **28.02** | **51.81** | **Winner** |
| **TLS Overhead (HTTPS)** | **TitanHTTP** | **102,450** | **0.97** | **4.10** | **+8.3%** |
| | net/http | 94,592 | 0.75 | 10.35 | |

#### Memory Profiling
- **TitanHTTP Leak:** 0 MB (Tested with 500,000 sustained requests)
- **net/http Leak:** 0 MB

### 2. Internal Micro-Benchmarks (Zero-Allocation Architecture)

Through strict memory management and custom implementations (like our Radix Tree Router and Sharded LRU Cache), we have driven allocations down to zero in critical hot-paths.

| Component / Benchmark | Time (ns/op) | Memory (B/op) | Allocations (op) |
| :--- | :--- | :--- | :--- |
| **LRU Cache (Parallel Hits)** | **14.37 ns** | **0 B** | **0 allocs** |
| **Router (Static Match)** | **80.37 ns** | **0 B** | **0 allocs** |
| **Router (Param Match)** | **127.85 ns** | **0 B** | **0 allocs** |
| **Rate Limiter (Sliding Window)** | **190.55 ns** | **21.5 B** | **0 allocs** |
| **Request Parser (Full)** | **179.05 ns** | **128 B** | **7 allocs** |

### 3. Framework Comparisons

Compared against popular Go frameworks (Gin, Fiber, Chi) handling a simple `/ping` route via `benchstat`:

| Rank | Benchmark Name | Median Time | Memory | Allocs |
| :--- | :--- | :--- | :--- | :--- |
| **#1** | **BenchmarkTitanHTTP-16** | **118.5 ns** | **4 B** | **1 allocs** |
| #2 | BenchmarkNetHTTP-16 | 198.15 ns | 4 B | 1 allocs |
| #3 | BenchmarkGin-16 | 244.6 ns | 48 B | 1 allocs |
| #4 | BenchmarkChi-16 | 626 ns | 372 B | 3 allocs |
| #5 | BenchmarkFiber-16 | 11,234.5 ns | 5,595.2 B | 21 allocs |

### 4. Compliance and Security

- **HTTP Compliance Tests:** 20 Passed, 0 Failed
  - (Testing varied Methods, Chunked Transfer-Encoding, Content-Length, Keep-Alive, 100-Continue, Range Requests, Compression, and Invalid Syntaxes).

- **Unified Security Suite Matrix:**

| Attack Scenario | TitanHTTP | net/http |
| :--- | :--- | :--- |
| **Header Exhaustion (Bomb)** | **PASS** (Blocked at 775) | FAIL (Parsed 2000 junk) |
| **Malformed HTTP Request** | **PASS** (400 Bad Request) | **PASS** (400 Bad Request) |
| **Slow POST Body Trickle** | FAIL (Vulnerable) | FAIL (Vulnerable) |
| **Idle Connection Flood** | FAIL (Vulnerable) | FAIL (Vulnerable) |

TitanHTTP successfully limits request headers (mitigating header bombs) more strictly than the default `net/http` configuration, though both require explicit Read/Write timeouts to mitigate Slowloris attacks.
