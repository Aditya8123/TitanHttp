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
| **50 Connections** | **TitanHTTP** | **117,763** | **0** | **2.72** | **+27.1%** |
| | net/http | 92,659 | 0.51 | 3.22 | |
| **100 Connections** | **TitanHTTP** | **113,354** | **0.78** | **4.1** | **+19.1%** |
| | net/http | 95,191 | 0.56 | 7.6 | |
| **200 Connections** | **TitanHTTP** | **122,774** | **1.5** | **4.75** | **+18.5%** |
| | net/http | 103,593 | 1.05 | 14.56 | |

#### Standard Workload Profiles
| Workload | Server | Throughput (Req/s) | P50 Latency (ms) | P99 Latency (ms) | RPS Difference |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Concurrent Connections** | **TitanHTTP** | **112,965** | **0.51** | **316.53** | **+106.9%** |
| | net/http | 54,606 | 0.66 | 387.01 | |
| **Throughput (Req/s)** | **TitanHTTP** | **151,316** | **0.65** | **3.56** | **+28.9%** |
| | net/http | 117,401 | 0.55 | 7.78 | |
| **Static File Serving** | **TitanHTTP** | **82,576** | **1.13** | **5.59** | **+34.3%** |
| | net/http | 61,493 | 1.06 | 15.92 | |
| **Keep-Alive Performance** | **TitanHTTP** | **121,369** | **1.00** | **4.03** | **+19.9%** |
| | net/http | 101,234 | 0.60 | 9.16 | |
| **Routing Performance** | **TitanHTTP** | **129,746** | **0.93** | **3.83** | **+32.1%** |
| | net/http | 98,216 | 0.63 | 8.9 | |
| **Large Payloads** | **TitanHTTP** | **60,352** | **1.67** | **6.00** | **+41.4%** |
| | net/http | 42,695 | 1.58 | 26.41 | |

#### High-Churn and TLS Infrastructure
| Profile | Server | Throughput (Req/s) | P50 Latency (ms) | P99 Latency (ms) | RPS Difference |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Connection Churn** | **TitanHTTP** | **3,568** | **26.26** | **56.23** | **+81.1%** |
| (Connection: close) | net/http | 1,969 | 38.65 | 228.17 | |
| **TLS Overhead** | **TitanHTTP** | **114,279** | **1.00** | **4.21** | **+23.4%** |
| | net/http | 92,624 | 0.70 | 9.00 | |

#### Memory Profiling
- **TitanHTTP Leak:** 0 MB (Tested with 500,000 sustained requests)
- **net/http Leak:** 0 MB

### 2. Framework Comparisons

Compared against popular Go frameworks (Gin, Fiber, Chi) handling a simple `/ping` route via `benchstat`:

```text
              │             sec/op                     │
NetHTTP-16                                      179.2n ± 16%
Gin-16                                          239.6n ± 16%
Chi-16                                          623.2n ±  8%
Fiber-16                                        12.28µ ±  7%
TitanHTTP-16                                    358.2n ±  8%
geomean                                         651.9n

              │                     B/op                      │
NetHTTP-16                                        4.000 ± 0%
Gin-16                                            48.00 ± 0%
Chi-16                                            372.0 ± 0%
Fiber-16                                        5.462Ki ± 0%
TitanHTTP-16                                      148.0 ± 0%
geomean                                           142.7

              │                   allocs/op                   │
NetHTTP-16                                        1.000 ± 0%
Gin-16                                            1.000 ± 0%
Chi-16                                            3.000 ± 0%
Fiber-16                                          21.00 ± 0%
TitanHTTP-16                                      3.000 ± 0%
geomean                                           2.853
```

### 3. Compliance and Security

- **HTTP Compliance Tests:** 20 Passed, 0 Failed
  - (Testing varied Methods, Chunked Transfer-Encoding, Content-Length, Keep-Alive, 100-Continue, Range Requests, Compression, and Invalid Syntaxes).
- **Unified Security Suite:** 8 Passed, 0 Failed
  - (Testing Slowloris, Request Smuggling, Body/Header Floods, Path Traversal, and Connection limits).
