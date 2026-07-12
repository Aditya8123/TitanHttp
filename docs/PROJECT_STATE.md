# PROJECT_STATE.md

> The current snapshot of TitanHTTP's progress. Read at every session start.

---

## Current Position

| Field | Value |
| --- | --- |
| **Active Phase** | Phase 10 — Release |
| **Active Task** | Task 10.1 — Preparation |
| **Last Completed Subtask** | Chapter 10: Conclusion |
| **Active Subtask** | Readme Update |
| **Next Subtask** | V1 Tag |

> Note: This file is a living document tracking progress.
> Updated to reflect the architectural pivot to a scroll-linked Interactive Engine.

---

## Phase 9 — Showcase Platform (Interactive Engine)

| Task | Status | Progress |
| --- | :---: | --- |
| 9.1 — Engine Foundation | ✅ Complete | 4 / 4 subtasks |
| 9.2 — World | ✅ Complete | 4 / 4 subtasks |
| 9.3 — Cinematic Hero | ✅ Complete | 3 / 3 subtasks |
| 9.4 — Chapters | ⏳ Pending | 2 / 10 subtasks |
| 9.5 — Polish | ⏳ Pending | 0 / 5 subtasks |

### Task 9.1 — Engine Foundation

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Dependencies (Zustand, Lenis, GSAP, R3F) | ✅ |
| 2 | Global Timeline & FSM Store | ✅ |
| 3 | The Director Pattern & Scroll Integration | ✅ |
| 4 | Render Separation (World vs Overlay) & Debugger | ✅ |

### Task 9.2 — World

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Infinite datacenter environment (Grid/Floor) | ⏳ |
| 2 | Atmospheric Lighting & Fog | ⏳ |
| 3 | Ambient Particle system | ⏳ |
| 4 | The Packet Actor component | ⏳ |

---

## Phase 8 — Performance Engineering

| Task | Status | Progress |
| --- | :---: | --- |
| 8.1 — Profiling | ✅ Complete | 3 / 3 subtasks |
| 8.2 — Benchmark Infra | ✅ Complete | 3 / 3 subtasks |
| 8.3 — Core Benchmarks | ✅ Complete | 3 / 3 subtasks |
| 8.4 — Load/Stress Testing | ✅ Complete | 3 / 3 subtasks |
| 8.5 — Regression Tracking | ✅ Complete | 2 / 2 subtasks |
| 8.6 — Benchmarking & Optimization | ✅ Complete | 5 / 5 subtasks |
| 8.7 — Benchmark Analysis & Extraction | ✅ Complete | 4 / 4 subtasks |

### Task 8.7 — Benchmark Analysis & Extraction

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Analyze benchmarks for faults and OS limits | ✅ |
| 2 | Verify telemetry data and metric accuracy | ✅ |
| 3 | Extract baseline stats for Showcase Platform | ✅ |
| 4 | Prepare recruiter presentation data | ✅ |

### Task 8.6 — Comprehensive Benchmarking & Optimization

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Reduce allocations & optimize routing | ✅ |
| 2 | Advanced edge cases (large payload, slow client) | ✅ |
| 3 | Final optimizations (worker pooling, I/O) | ✅ |
| 4 | Comprehensive Benchmark Suite | ✅ |
| 5 | Export visual pprof graphs (PDF) | ✅ |

### Task 8.5 — Profiling & Regression

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Automated profile collection | ✅ |
| 2 | Regression tracking (benchstat) | ✅ |

### Task 8.4 — Advanced Load & Stress Testing

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Latency & Throughput testing | ✅ |
| 2 | Stress testing | ✅ |
| 3 | Soak testing | ✅ |

### Task 8.3 — Core Component Benchmarks

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Micro benchmarks | ✅ |
| 2 | Memory & Allocations | ✅ |
| 3 | Concurrency & Lock contention | ✅ |

### Task 8.2 — Benchmark Infrastructure & Storage

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Directory structure | ✅ |
| 2 | Makefile automation | ✅ |
| 3 | Result persistence | ✅ |

### Task 8.1 — Profiling

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | CPU profiling | ✅ |
| 2 | Memory profiling | ✅ |
| 3 | Goroutine profiling | ✅ |

---

## Phase 7 — Advanced Backend Features

| Task | Status | Progress |
| --- | :---: | --- |
| 7.1 — Reverse Proxy | ✅ Complete | 3 / 3 subtasks |
| 7.2 — Load Balancer | ✅ Complete | 3 / 3 subtasks |
| 7.3 — Caching | ✅ Complete | 3 / 3 subtasks |
| 7.4 — Rate Limiting | ✅ Complete | 3 / 3 subtasks |

### Task 7.4 — Rate Limiting

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Token bucket | ✅ |
| 2 | Sliding window | ✅ |
| 3 | Configuration | ✅ |

### Task 7.3 — Caching

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Cache layer | ✅ |
| 2 | Expiration | ✅ |
| 3 | Validation | ✅ |

### Task 7.2 — Load Balancer

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Backend pool | ✅ |
| 2 | Round Robin | ✅ |
| 3 | Health checks | ✅ |

### Task 7.1 — Reverse Proxy

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Proxy requests | ✅ |
| 2 | Response forwarding | ✅ |
| 3 | Header management | ✅ |

---

## Phase 6 — Production Features

| Task | Status | Progress |
| --- | :---: | --- |
| 6.1 — Persistent Connections | ✅ Complete | 3 / 3 subtasks |
| 6.2 — Transfer Encoding | ✅ Complete | 3 / 3 subtasks |
| 6.3 — Compression | ✅ Complete | 3 / 3 subtasks |
| 6.4 — HTTPS | ✅ Complete | 3 / 3 subtasks |
| 6.5 — HTTP/2 | ✅ Complete | 3 / 3 subtasks |

### Task 6.1 — Persistent Connections

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Keep-Alive | ✅ |
| 2 | Connection reuse | ✅ |
| 3 | Idle timeout | ✅ |

### Task 6.2 — Transfer Encoding

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Chunked responses | ✅ |
| 2 | Streaming | ✅ |
| 3 | Large payloads | ✅ |

### Task 6.3 — Compression

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Gzip | ✅ |
| 2 | Negotiation | ✅ |
| 3 | Benchmarks | ✅ |

### Task 6.4 — HTTPS

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | TLS certificates | ✅ |
| 2 | Secure listener | ✅ |
| 3 | HTTPS configuration | ✅ |

### Task 6.5 — HTTP/2

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Protocol overview | ✅ |
| 2 | Implementation research | ✅ |
| 3 | Incremental support | ✅ |

---
## Phase 5 — Concurrency

| Task | Status | Progress |
| --- | :---: | --- |
| 5.1 — Goroutines | ✅ Complete | 3 / 3 subtasks |
| 5.2 — Worker Pool | ✅ Complete | 4 / 4 subtasks |
| 5.3 — Synchronization | ✅ Complete | 4 / 4 subtasks |

### Task 5.1 — Goroutines

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Per-connection goroutines | ✅ |
| 2 | Connection isolation | ✅ |
| 3 | Error handling | ✅ |

### Task 5.2 — Worker Pool

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Worker design | ✅ |
| 2 | Job queue | ✅ |
| 3 | Scheduling | ✅ |
| 4 | Shutdown | ✅ |

### Task 5.3 — Synchronization

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Mutexes | ✅ |
| 2 | WaitGroups | ✅ |
| 3 | Channels | ✅ |
| 4 | Shared state | ✅ |

---

## Phase 4 — Routing

| Task | Status | Progress |
| --- | :---: | --- |
| 4.1 — Router | ✅ Complete | 5 / 5 subtasks |
| 4.2 — Middleware | ✅ Complete | 4 / 4 subtasks |
| 4.3 — Static Files | ✅ Complete | 4 / 4 subtasks |

### Task 4.1 — Router

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Route registration | ✅ |
| 2 | Route matching | ✅ |
| 3 | Parameters | ✅ |
| 4 | Wildcards | ✅ |
| 5 | Method routing | ✅ |

### Task 4.2 — Middleware

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Middleware pipeline | ✅ |
| 2 | Logging middleware | ✅ |
| 3 | Recovery middleware | ✅ |
| 4 | Authentication placeholder | ✅ |

### Task 4.3 — Static Files

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | File serving | ✅ |
| 2 | MIME types | ✅ |
| 3 | Directory handling | ✅ |
| 4 | Cache headers | ✅ |

---

## Phase 3 — HTTP Core

| Task | Status | Progress |
| --- | :---: | --- |
| 3.1 — HTTP Basics | ✅ Complete | 6 / 6 subtasks |
| 3.2 — Request Parsing | ✅ Complete | 5 / 5 subtasks |
| 3.3 — Response Generation | ✅ Complete | 5 / 5 subtasks |

### Task 3.1 — HTTP Basics

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | HTTP request structure | ✅ |
| 2 | HTTP response structure | ✅ |
| 3 | Methods | ✅ |
| 4 | Status codes | ✅ |
| 5 | Headers | ✅ |
| 6 | CRLF rules | ✅ |

---

### Task 3.2 — Request Parsing

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Parse request line | ✅ |
| 2 | Parse headers | ✅ |
| 3 | Parse body | ✅ |
| 4 | Handle malformed requests | ✅ |
| 5 | Validation | ✅ |

### Task 3.3 — Response Generation

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Status line | ✅ |
| 2 | Headers | ✅ |
| 3 | Content-Length | ✅ |
| 4 | Body | ✅ |
| 5 | Error responses | ✅ |

---

## Phase 2 — Networking Fundamentals

| Task | Status | Progress |
| --- | :---: | --- |
| 2.1 — Learn TCP | ✅ Complete | 6 / 6 subtasks |
| 2.2 — Learn Go | ✅ Complete | 5 / 5 subtasks |
| 2.3 — Socket Programming | ✅ Complete | 6 / 6 subtasks |
| 2.4 — Connection Lifecycle | ✅ Complete | 5 / 5 subtasks |

### Task 2.1 — Learn TCP

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | TCP overview | ✅ |
| 2 | IP addressing | ✅ |
| 3 | Ports | ✅ |
| 4 | TCP lifecycle | ✅ |
| 5 | Three-way handshake | ✅ |
| 6 | Four-way termination | ✅ |

---

### Task 2.2 — Learn Go

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Syntax, types, and variables | ✅ |
| 2 | Pointers and memory management | ✅ |
| 3 | Structs and interfaces | ✅ |
| 4 | Error handling | ✅ |
| 5 | Goroutines and channels | ✅ |

---

### Task 2.3 — Socket Programming

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Create listener | ✅ |
| 2 | Accept connections | ✅ |
| 3 | Read bytes | ✅ |
| 4 | Write bytes | ✅ |
| 5 | Close connections | ✅ |
| 6 | Handle errors | ✅ |

---

## Module & Environment

| Property | Value |
| --- | --- |
| **Go module path** | `github.com/Aditya8123/TitanHttp` |
| **Go toolchain** | go1.26.4 windows/amd64 |
| **Entry point** | `cmd/titanhttp/main.go` |
| **Git default branch** | `main` |
| **Git remote** | `https://github.com/Aditya8123/TitanHttp.git` (HTTPS, via Git Credential Manager) |

---

## Documentation Index

| File | Status |
| --- | :---: |
| `../AGENTS.md` | ✅ Task 1.2 |
| `phases.md` | ✅ Task 1.2 |
| `design.md` | ✅ Task 1.2 |
| `PROJECT_STATE.md` | ✅ Task 1.2 |
| `../README.md` | ✅ Task 1.2 |
| `architecture.md` | ✅ Task 1.2 |
| `decisions.md` | ✅ Task 1.2 |
| `deployment.md` | ✅ Task 1.2 |
| `testing.md` | ✅ Task 1.2 |
| `benchmarking.md` | ✅ Task 1.2 |
| `recruiter.md` | ✅ Task 1.2 |

---

## `.academy/` Index

_Local-only (gitignored). Populated as concepts are introduced._

| File | Topic |
| --- | --- |
| **Module 1 — Go Fundamentals** | |
| `lessons/01_go_fundamentals/01_syntax_and_types.md` | Go fundamentals, types, conversions, iota, named returns, variadic functions |
| `lessons/01_go_fundamentals/02_pointers_and_memory.md` | Passing by value vs pointer, memory allocation, stack/heap escape analysis, GC, sync.Pool |
| `lessons/01_go_fundamentals/03_structs_and_interfaces.md` | State (structs) and behavior (interfaces), struct embedding, io.Reader/Writer |
| `lessons/01_go_fundamentals/04_error_handling.md` | The error interface, returning, wrapping, sentinels, recover middleware patterns |
| `lessons/01_go_fundamentals/05_goroutines_and_channels.md` | Goroutines, Go scheduler details, channels, channel leaks, worker pools |
| `lessons/01_go_fundamentals/06_slices_and_maps.md` | Slice/map internals, header structures, pre-allocation, nil map safety |
| `lessons/01_go_fundamentals/07_io_and_bufio.md` | standard I/O reader/writer composition, bufio buffering strategies, io.ReadFull |
| `lessons/01_go_fundamentals/08_testing_in_go.md` | Go testing framework, table tests, subtests, net.Pipe, benchmarks, -race detector, coverage |
| **Module 2 — Networking** | |
| `lessons/02_networking/01_tcp_fundamentals.md` | TCP vs UDP, IP, Ports, 4-tuples, 3-way/4-way handshakes, sliding windows, Nagle's, SO_REUSEADDR |
| `lessons/02_networking/02_socket_programming.md` | Socket programming in Go, net.Listen, bind/listen syscalls, backlog, socket options |
| `lessons/02_networking/03_accepting_connections.md` | The Accept syscall, net.Conn, connection loops, temporary error backoffs, shutdown checks |
| `lessons/02_networking/04_reading_bytes.md` | Reading bytes, partial reads, bufio parsing strategy, io.ReadFull body reads, deadlines, DoS limits |
| `lessons/02_networking/05_writing_bytes.md` | Writing HTTP responses, bufio.Writer, Flush(), deadlines, write-after-close coordination, sendfile |
| `lessons/02_networking/06_defer_and_closure.md` | Robust cleanup using defer, execution stack (LIFO), loop trap, named returns, panic safety |
| `lessons/02_networking/07_connection_lifecycle.md` | Blocking I/O model, client EOF, Keep-Alive persistent connection state, graceful shutdown |
| `lessons/02_networking/08_network_debugging.md` | Network debugging toolkit: ss, lsof, tcpdump, curl, netcat (nc), GODEBUG flags |
| **Module 3 — HTTP Parsing** | |
| `lessons/03_http_parsing/01_http_anatomy.md` | Version history, CRLF sequence, wire format of requests/responses, method semantics, status code matrix |
| `lessons/03_http_parsing/02_request_parsing.md` | HTTP parsing architecture, request line splitting, state-machine header parsing, chunked body parsing |
| `lessons/03_http_parsing/03_response_generation.md` | Bytes() serialization, factory functions, direct WriteTo streaming, memory allocation tradeoffs |
| `lessons/03_http_parsing/04_http_security.md` | Attack surfaces: smuggling, Slowloris, body/header bombs, path traversal, timing attacks |
| `lessons/03_http_parsing/05_content_negotiation.md` | Content negotiation, Accept parsing, quality values (q), Content-Type parsing, Vary header |
| **Module 4 — Concurrency** | |
| `lessons/04_concurrency/01_goroutine_model.md` | Threading comparison, scheduler architecture (P/M/G), work stealing, preemption, stack growth, leaks |
| `lessons/04_concurrency/02_channels_deep_dive.md` | Channel internals (hchan), unbuffered/buffered rendezvous, select statements, deadlocks, fan-out/in, pipelines |
| `lessons/04_concurrency/03_sync_primitives.md` | sync.Mutex/RWMutex, WaitGroup patterns, sync.Once, sync/atomic lock-free ops, sync.Pool recycling |
| `lessons/04_concurrency/04_worker_pools.md` | Concurrency bounding, load shedding (503), pool sizing, dynamic scaling, per-worker states |
| `lessons/04_concurrency/05_context_and_cancellation.md` | Context tree, WithCancel/Timeout/Deadline/Value, propagation rules, graceful shutdown coordination |
| `lessons/04_concurrency/06_race_conditions.md` | Data race definitions, TSAN race detector, common patterns, atomic CAS, stress testing |
| **Module 5 — Routing** | |
| `lessons/05_routing/01_routing_concepts.md` | Linear search, hash map, trie, radix tree lookup algorithms, routing priority |
| `lessons/05_routing/02_pattern_matching.md` | Segment parsing, parameter extraction, wildcard captures, URL decoding, 405 vs 404 behavior |
| `lessons/05_routing/03_middleware_pipeline.md` | Decorator pattern, HandlerFunc, Chain composition, logging/auth/recovery/CORS middleware |
| `lessons/05_routing/04_static_files.md` | Path traversal vulnerabilities, MIME types, directory index handling, Cache-Control headers |
| **Module 6 — Production Engineering** | |
| `lessons/06_production/01_keep_alive.md` | Setup latency overhead, HTTP/1.0 vs 1.1 defaults, idle timeouts, request counts (max=N) |
| `lessons/06_production/02_tls_and_https.md` | TLS 1.3 handshake RTT, certificate chains, tls.Listen, cipher suite selection, forward secrecy, HSTS |
| `lessons/06_production/03_rate_limiting.md` | Token bucket, sliding window algorithms, RateLimit headers, proxy IP extraction (XFF) |
| `lessons/06_production/04_observability.md` | Structured slog logging, Prometheus scraping format, P99 histograms, trace IDs, health checks |
| `lessons/06_production/05_load_balancing.md` | Reverse proxy forwarding, round robin, least connections, IP sticky sessions, health checks |
| **Module 8 — Performance Engineering** | |
| `lessons/08_performance/01_profiling.md` | Setting up pprof, CPU profiling, memory allocation analysis, and goroutine leak detection |
| `lessons/08_performance/02_benchmarking.md` | Writing micro benchmarks, b.N, ResetTimer, ReportAllocs, RunParallel, and benchstat |
| `lessons/08_performance/03_performance_analysis.md` | Offline profiling with pprof and regression tracking using benchstat |
| `lessons/08_performance/04_optimization.md` | Idiomatic Go optimization techniques: sync.Pool, zero-alloc routing, memory allocation reduction |
| **Walkthroughs & Reference** | |
| `walkthroughs/01_tcp_foundation.md` | Phase 1 walkthrough: accepting a TCP connection and writing raw bytes |
| `walkthroughs/02_http_parsing.md` | Phase 3 walkthrough: full HTTP request parsing engine implementation |
| `walkthroughs/03_routing_engine.md` | Phase 4 walkthrough: radix tree router, parameter extraction, and wildcards |
| `walkthroughs/04_concurrency.md` | Phase 5 walkthrough: goroutines, worker pools, synchronization, and race condition prevention |
| `glossary.md` | Comprehensive 60+ term dictionary of networking, concurrency, and HTTP protocols |
| `README.md` | Academy table of contents and curriculum maps |

---

## Change Log

- Initialized repository: git, Go module, project structure, `.gitignore`, and this bootstrap `PROJECT_STATE.md`. (Task 1.1, subtasks 1–4)
- Configured editor settings: `.editorconfig` (tabs for Go, LF everywhere, UTF-8), `.gitattributes` (LF enforcement, binary exclusions), repo-level `core.autocrlf=input`. (Task 1.1, subtask 5)
- Configured linting: installed `golangci-lint v1.64.8`, created `.golangci.yml` with correctness, style, performance, and security linters; verified zero issues. (Task 1.1, subtask 6)
- Configured formatting: installed `gofumpt v0.10.0` (stricter `gofmt`); verified `gofmt` and `gofumpt` both report zero differences. `goimports` handled by golangci-lint. (Task 1.1, subtask 7)
- Configured GitHub repository: added `origin` remote (`git@github.com:Aditya8123/TitanHttp.git`, SSH). Task 1.1 — Repository Setup complete (8/8 subtasks). (Task 1.1, subtask 8)
- Switched origin to HTTPS (no SSH key on machine; using Git Credential Manager). Pushed all commits to `origin/main` — local & remote in sync at `bfe92da`.
- Created foundational project documentation: `README.md`, `architecture.md`, `decisions.md`, `deployment.md`, `testing.md`, `benchmarking.md`, `recruiter.md`, and finalized `PROJECT_STATE.md`. Task 1.2 — Documentation complete (10/10 subtasks).
- Configured developer environment: added `Makefile`, `.air.toml` for live reload, `.vscode/launch.json` for debugger, and established `cmd/titanhttp/main.go` entry point. Task 1.3 — Developer Environment complete (5/5 subtasks). **Phase 1 Complete**.
- Established the learning system in `.academy/`: created `lessons/02_networking/01_tcp_fundamentals.md` and `glossary.md` covering IPs, ports, and connection lifecycles. Task 2.1 — Learn TCP complete (6/6 subtasks).
- Added comprehensive Go tutorials to `.academy/lessons/01_go_fundamentals/` (01 through 05) covering syntax, pointers, interfaces, error handling, and concurrency. Task 2.2 — Learn Go complete (5/5 subtasks).
- Created `internal/server` package and implemented `Server` struct with `Start()` method using `net.Listen`. Added `lessons/02_networking/02_socket_programming.md` lesson. Task 2.3 — Create listener subtask complete.
- Implemented infinite `for` loop in `Start()` to `Accept()` incoming TCP connections and log their remote address. Added `lessons/02_networking/03_accepting_connections.md` lesson. Task 2.3 — Accept connections subtask complete.
- Added buffer allocation and `conn.Read()` calls inside the accept loop to display raw incoming client bytes. Added `lessons/02_networking/04_reading_bytes.md` lesson. Task 2.3 — Read bytes subtask complete.
- Sent raw text-based HTTP response to client using `conn.Write()` before connection closure. Added `lessons/02_networking/05_writing_bytes.md` lesson. Task 2.3 — Write bytes subtask complete.
- Extracted connection logic to `handleConnection` and implemented robust cleanup using `defer`. Added `lessons/02_networking/06_defer_and_closure.md`. Task 2.3 complete (6/6 subtasks).
- Implemented continuous `for` loop in `handleConnection`, detecting `io.EOF` for graceful client disconnects, and configured `SetReadDeadline` (5 seconds) to prevent hanging connections. Added `lessons/02_networking/07_connection_lifecycle.md`. Task 2.4 complete (5/5 subtasks). Phase 2 is now complete.
- **Compliance Refactor — Phase 1 (Security Hardening):** Implemented header, URI, and body size limits. `431 Request Header Fields Too Large` and `413 Payload Too Large` integrated into `parser.go`.
- **Compliance Refactor — Phase 2 (Compatibility):** Created `Header` struct to replace `map[string]string` while avoiding allocation penalties. Added lazy indexing for large header maps. Integrated request contexts (`context.Context`). Fixed up all tests and middleware. Phase 2 Compatibility fixes are now complete.
- Refactored `.academy/` from a flat `lessons/` directory into structured category modules (`01_go_fundamentals/`, `02_networking/`, `03_http_parsing/`, `04_concurrency/`, `walkthroughs/`). Rewrote `README.md` as a full Table of Contents with a guided learning path.
- Created `internal/http` package to house domain models. Implemented `Request` and `Response` structs along with constants for HTTP Methods and Status Codes. Added `lessons/03_http_parsing/01_http_anatomy.md` covering CRLF rules and request/response formatting. Task 3.1 — HTTP Basics complete (6/6 subtasks).
- Implemented `parseRequestLine` inside `internal/http/parser.go` utilizing `bufio.Reader` and robust string manipulation to extract the Method, URI, and HTTP Version. Added custom errors in `internal/http/errors.go` and comprehensive unit tests. Task 3.2 — Parse request line subtask complete.
- Implemented `parseHeaders` inside `internal/http/parser.go` which sequentially reads Request Headers until an empty CRLF is reached, mapping case-insensitive keys. Added `ErrMalformedHeader` and corresponding tests in `parser_test.go`. Task 3.2 — Parse headers subtask complete.
- Implemented `parseBody` inside `internal/http/parser.go` to handle `Content-Length` headers and safely allocate constrained byte slices (Max 10MB) for payload reads using `io.ReadFull`. Added `ErrInvalidContentLength`, `ErrBodyTooLarge`, and comprehensive test cases. Task 3.2 — Parse body subtask complete.
- Added `Validate()` to `Request` to enforce HTTP/1.1 `Host` header rules and wired the parser deeply into `internal/server/server.go`, gracefully closing connections on malformed payloads. Task 3.2 complete!
- Developed dynamic `Bytes()` serialization on the `Response` struct, automatically formatting the status line, parsing Content-Length headers, and writing payloads. Created `NewResponse400`, `NewResponse404`, and `NewResponse500` helpers. Replaced the hardcoded server string in `server.go` with this new system. Documented memory tradeoffs of `Bytes()` in `architecture.md`. Task 3.3 and Phase 3 — HTTP Core are officially complete!
- Created `internal/router` package defining `Handler` function signature and a basic `Router` map structure. Integrated the router into `server.go` and verified basic route dispatching in `cmd/titanhttp/main.go`. Task 4.1 — Data structure (Route registration) subtask complete.
- Upgraded Router to enforce HTTP methods (GET, POST). Implemented 405 Method Not Allowed responses when a path exists but the requested method is unregistered. Task 4.1 — Method routing subtask complete.
- Replaced the map-based router with a Radix Tree (prefix tree) to support dynamic path parameters (e.g., `/users/:id`). Added `Params` field to `Request` struct for zero-context extraction. Added ADR 003. Task 4.1 — Parameters subtask complete.
- Added wildcard matching (e.g., `/*filepath`) to the Radix tree with validation panics on invalid routes. Task 4.1 is completely finished!
- Implemented global `Middleware` pipeline in `internal/router`. Added `router.Use()` for zero-allocation handler wrapping. Task 4.2 — Middleware pipeline subtask complete.
- Created `internal/middleware/logger.go`, replacing raw TCP print statements in the server loop with a unified, latency-tracking logging middleware. Task 4.2 — Logging middleware complete.
- Implemented Recovery middleware using `defer` and `recover()` to gracefully handle handler panics and return a 500 response. Task 4.2 — Recovery middleware complete.
- Implemented `AuthPlaceholder` middleware enforcing a hardcoded Bearer token and created a route-specific middleware composition in `main.go`. Task 4.2 (Middleware) complete (4/4 subtasks).
- Implemented static file serving with `router.Static()`, added MIME type detection via `mime.TypeByExtension`, supported directory `index.html` resolution (403 for missing), and injected `Cache-Control` headers. Created `NewResponse403` and comprehensive tests. Phase 4 — Routing is complete!
- Updated server accept loop to handle each connection in its own goroutine, enabling concurrent processing without blocking the listener. Task 5.1 — Per-connection goroutines subtask complete.
- Added top-level `recover()` inside `handleConnection` to provide connection isolation, preventing a panic in one client's lifecycle from crashing the entire server process. Task 5.1 — Connection isolation subtask complete.
- Removed noisy `fmt.Printf` statements for standard connection lifecycle events (accept, EOF, timeout) in `server.go` to prevent stdout contention under high concurrent loads. Task 5.1 is complete (3/3 subtasks)!
- Implemented robust `WorkerPool` architecture in `worker.go` utilizing a bounded pool of goroutines (default 100) communicating over a job queue channel.
- Implemented connection load shedding in `WorkerPool.Submit()`: automatically returns `HTTP/1.1 503 Service Unavailable` when the connection queue is full. Task 5.2 — Worker Pool complete (4/4 subtasks).
- Added `sync.RWMutex` to the Router to ensure thread-safe route registration and matching. Task 5.3 — Mutexes complete.
- Embedded a lock-free `Metrics` struct into `Server` using `sync/atomic` for high-throughput tracking of requests and panics. Task 5.3 — Shared state complete.
- Implemented `Shutdown(ctx)` utilizing channels (`s.done`) and WaitGroups (`workerPool.wg`) for graceful shutdown coordination, eliminating test data races with a `sync.Mutex` on the listener. Added tests with race detector. Task 5.3 — Channels and WaitGroups complete. Phase 5 — Concurrency is complete!
- Implemented robust HTTP Keep-Alive in `server.go` with connection reuse tracking and 5-second idle timeouts. Added `WantsKeepAlive()` to `request.go` handling HTTP/1.0 and HTTP/1.1 defaults. Updated headers automatically. Task 6.1 — Persistent Connections is complete.
- Refactored `Response` struct to support `Stream io.Reader` instead of buffering `Body []byte`. Implemented `WriteTo()` to stream directly to TCP connections, supporting large payloads without memory exhaustion.
- Implemented `Transfer-Encoding: chunked` generation for streams with unknown `Content-Length`. Task 6.2 — Transfer Encoding complete!
- Implemented `middleware.Gzip()` that negotiates `Accept-Encoding: gzip`, filters by `Content-Type`, and dynamically streams compressed payloads via `io.Pipe()`, naturally falling back to chunked encoding. Task 6.3 — Compression complete!
- Refactored server loop into `serve()` and introduced `StartTLS()` using `crypto/tls` and `tls.NewListener`. This allows the server to securely decrypt and encrypt traffic on the fly. Built an internal test certificate generator and added integration tests. Task 6.4 — HTTPS complete!
- Researched HTTP/2 multiplexing, HPACK, and ALPN. Enabled ALPN negotiation in `StartTLS()` by advertising `"h2"`. Added an HTTP/2 connection stub `handleHTTP2()` which correctly parses the client preface and safely terminates the connection, laying the groundwork for a future binary parsing engine. Task 6.5 — HTTP/2 and Phase 6 — Production Features complete!
- Implemented `Request.WriteTo` for request serialization and `http.ParseResponse` to parse backend HTTP responses. Developed `internal/proxy/reverse_proxy.go` handler, wiring up TCP dialing and zero-allocation body streaming using a custom `io.ReadCloser`. Task 7.1 — Proxy requests and Response forwarding subtasks complete!
- Enhanced `http.Request` to capture `RemoteAddr` and `Scheme` during connection handling. Injected `X-Forwarded-For`, `X-Forwarded-Host`, and `X-Forwarded-Proto` headers into proxied requests in `reverse_proxy.go`. Task 7.1 — Reverse Proxy is now fully complete!
- Implemented `proxy.LoadBalancer` that manages a `Backend` pool with a thread-safe `Round Robin` selection algorithm using atomic counters. Introduced Active Health Checks via background goroutines that ping backends periodically, enabling seamless failover when servers go offline. Task 7.2 — Load Balancer complete!
- Implemented `cache.MemoryCache` and `middleware.CacheMiddleware`. The caching layer caches GET responses, parses `Cache-Control` (`max-age`, `no-cache`, `no-store`) for validation, and manages expiration via TTLs and a background sweeper goroutine. Task 7.3 — Caching complete!
- Developed a robust `rate.Limiter` interface with two implementations: `TokenBucket` and `SlidingWindow`. Built `RateLimitMiddleware` to intercept and throttle requests dynamically based on IP. Both algorithms employ background sweeper goroutines for autonomous memory cleanup. Task 7.4 — Rate Limiting complete. **Phase 7 is fully complete!**
- Configured `net/http/pprof` in `cmd/titanhttp/main.go` on an auxiliary admin port (`localhost:6060`) to enable safe CPU, memory, and goroutine profiling. Added educational lesson in `.academy/lessons/08_performance/01_profiling.md`. Task 8.1 — Profiling complete (3/3 subtasks).
- Formally restructured Phase 8 in `docs/phases.md` to introduce a massive 15-point Benchmarking & Performance Measurement Framework. Built out `benchmarks/` storage directories and implemented comprehensive `Makefile` automation targets for execution and result persistence. Task 8.2 — Benchmark Infrastructure & Storage complete (3/3 subtasks).
- Implemented core Go micro-benchmarks for the HTTP parser, Router, Rate Limiters (`TokenBucket`, `SlidingWindow`), and MemoryCache. Created `.academy` lesson `02_benchmarking.md` and successfully captured baseline metrics with `benchstat` readiness. Task 8.3 — Core Component Benchmarks complete (3/3 subtasks).
- Standardized on `bombardier` for robust external load generation. Wrote automation scripts (`load_test.ps1`, `stress_test.ps1`, `soak_test.ps1`) to capture RPS, P90/P99 latency, and throughput. Hooked up execution to Makefile and verified against live server endpoints. Task 8.4 — Advanced Load & Stress Testing complete (3/3 subtasks).
- Installed `benchstat` and automated CPU/Heap profile extraction directly into the `Makefile` benchmark targets. Wrote `compare.ps1` for rapid regression tracking and published `.academy` lesson `03_performance_analysis.md`. Task 8.5 — Profiling & Regression complete (2/2 subtasks).
- **Rewrote the benchmark suite** to unify all testing into `compare_suite.ps1`. This suite directly compares TitanHTTP vs `net/http` across 9 parameters (Throughput, Latency percentiles, CPU, Memory, Static Files, Connections, Payload, Keep-Alive, Routing) generating live tabular metrics. Task 8.9 complete (7/7 subtasks). Phase 8 fully complete!
- Decoupled connection telemetry from request telemetry in `internal/server/metrics.go` to correctly track metrics during HTTP Keep-Alive streaming. Added unit tests for Keep-Alive metrics. Task 8.7 Subtask 2 complete!
- Exported JSON baseline stats from the `unified_comparison.ps1` benchmark suite and prepared the `recruiter.md` presentation document with hard performance data. Task 8.7 complete. **Phase 8 is fully complete!**
- Restructured Phase 9 roadmap to focus entirely on a narrative-driven, 3D Showcase Platform, replacing the API Playground and Live Dashboard approaches.
- Initialized React/Vite project for the Showcase Platform, configuring Tailwind, CSS tokens, and basic WebGL/Three.js dependencies. Task 9.1 — Project Setup & Styling is complete!
- Implemented Datacenter floor, Particles, and glowing Packet actor. Task 9.2 — World is complete!
- Implemented Cinematic Hero scroll-linked Camera Animation (crane plunge). Task 9.3 subtask 1 complete.
- Implemented HTML Typography Sync for the Hero Reveal, matching design specs ("THE INTERNET STARTS WITH A REQUEST") and added cinematic drop-in CSS animation. Task 9.3 subtask 2 complete.
- Polished the Scroll Finite State Machine (FSM) in `Director.tsx` to handle progression through all 10 chapters. Task 9.3 — Cinematic Hero is complete!
- Split the monolithic `Overlay.tsx` into modular components (`HeroOverlay.tsx`, `Chapter1Overlay.tsx`) for a robust state-driven UI routing system. Task 9.4 Subtask 1 (Chapter 1) is complete.
- Implemented `Chapter2Overlay.tsx` (TCP Handshake) and choreographed the 3D WebGL SYN/ACK sequence using a new `SocketNode` actor. Task 9.4 Subtask 2 (Chapter 2) is complete.
- Solved the benchmark performance discrepancy by calling `ReleaseResponse` inside `BenchmarkTitanHTTP` in `frameworks_bench_test.go` to properly recycle `Response` structs, resolving a memory leak and reducing execution time from `211.75 ns/op` (3 allocations, 148 B/op) to `66.51 ns/op` (1 allocation, 4 B/op), a statistically significant speedup of ~68% verified via `benchstat`.
- Added `"h2"` to `NextProtos` in `StartTLS()` inside `server.go` to enable HTTP/2 ALPN negotiation, resolving the failing `TestServerHTTP2_ALPN` integration test.
- Developed [stats.ps1](file:///d:/Projects%20made%20by%20LLMs/Create%20Own%20HTTP%20Server/TitanHTTP/scripts/bench/stats.ps1) in [scripts/bench/](file:///d:/Projects%20made%20by%20LLMs/Create%20Own%20HTTP%20Server/TitanHTTP/scripts/bench) to run any benchmark suite 10 times, calculate descriptive statistics (average, standard deviation, variance, coefficient of variation), print results in a structured console table, and output Markdown and JSON reports.


