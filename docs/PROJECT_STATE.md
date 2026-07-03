# PROJECT_STATE.md

> The current snapshot of TitanHTTP's progress. Read at every session start.

---

## Current Position

| Field | Value |
| --- | --- |
| **Active Phase** | Phase 7 — Advanced Backend Features |
| **Active Task** | Task 7.1 — Reverse Proxy |
| **Last Completed Subtask** | Incremental support (Task 6.5) |
| **Active Subtask** | Proxy requests |
| **Next Subtask** | Response forwarding |

> Note: This file is a living document tracking progress.
> Updated at the completion of Task 6.5 (HTTP/2).

---

## Phase 6 — Production Features

| Task | Status | Progress |
| --- | :---: | --- |
| 6.1 — Persistent Connections | ✅ Complete | 3 / 3 subtasks |
| 6.2 — Transfer Encoding | ✅ Complete | 3 / 3 subtasks |
| 6.3 — Compression | ✅ Complete | 3 / 3 subtasks |
| 6.4 — HTTPS | ✅ Complete | 3 / 3 subtasks |
| 6.5 — HTTP/2 | ✅ Complete | 3 / 3 subtasks |

### Task 7.1 — Reverse Proxy

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Proxy requests | 🚧 In Progress |
| 2 | Response forwarding | ⏳ Pending |
| 3 | Header management | ⏳ Pending |

### Task 6.1 — Persistent Connections

| # | Subtask | Status |
| --- | --- | :---: |
| 1 | Keep-Alive | ✅ |
| 2 | Connection reuse | ✅ |
| 3 | Idle timeout | ✅ |

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
