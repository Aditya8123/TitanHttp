# PROJECT_STATE.md

> The current snapshot of TitanHTTP's progress. Read at every session start.

---

## Current Position

| Field | Value |
| --- | --- |
| **Active Phase** | Phase 2 — Networking Fundamentals |
| **Active Task** | Task 2.4 — Connection Lifecycle ⬜ |
| **Last Completed Subtask** | Handle errors (Task 2.3) |
| **Active Subtask** | — (Task 2.3 complete; next: Task 2.4 — Connection Lifecycle) |
| **Next Subtask** | Blocking I/O |

> Note: This file is a living document tracking progress.
> Updated at the completion of Task 2.3 — Socket Programming.

---

## Phase 2 — Networking Fundamentals

| Task | Status | Progress |
| --- | :---: | --- |
| 2.1 — Learn TCP | ✅ Complete | 6 / 6 subtasks |
| 2.2 — Learn Go | ✅ Complete | 5 / 5 subtasks |
| 2.3 — Socket Programming | ✅ Complete | 6 / 6 subtasks |
| 2.4 — Connection Lifecycle | ⬜ Not Started | — |

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
| `lessons/01_tcp_fundamentals.md` | TCP vs UDP, IP, Ports, Handshakes |
| `lessons/02_go_syntax_and_types.md` | Go fundamentals, types, conversions |
| `lessons/03_go_pointers_and_memory.md` | Passing by value vs pointer, memory allocation |
| `lessons/04_go_structs_and_interfaces.md` | State (structs) and behavior (interfaces) |
| `lessons/05_go_error_handling.md` | The error interface, returning and wrapping errors |
| `lessons/06_go_concurrency.md` | Goroutines and Channels |
| `lessons/07_socket_programming.md` | Socket programming in Go, net.Listen |
| `lessons/08_accepting_connections.md` | The Accept syscall, net.Conn, and connection loops |
| `lessons/09_reading_bytes.md` | Reading bytes from sockets, net.Conn |
| `lessons/10_writing_bytes.md` | Writing HTTP responses to sockets, net.Conn |
| `lessons/11_defer_and_closure.md` | Robust socket cleanup using defer |
| `glossary.md` | Core networking terminology definitions |

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
- Established the learning system in `.academy/`: created `lessons/01_tcp_fundamentals.md` and `glossary.md` covering IPs, ports, and connection lifecycles. Task 2.1 — Learn TCP complete (6/6 subtasks).
- Added comprehensive Go tutorials to `.academy/lessons/` (02 through 06) covering syntax, pointers, interfaces, error handling, and concurrency. Task 2.2 — Learn Go complete (5/5 subtasks).
- Created `internal/server` package and implemented `Server` struct with `Start()` method using `net.Listen`. Added `07_socket_programming.md` lesson. Task 2.3 — Create listener subtask complete.
- Implemented infinite `for` loop in `Start()` to `Accept()` incoming TCP connections and log their remote address. Added `08_accepting_connections.md` lesson. Task 2.3 — Accept connections subtask complete.
- Added buffer allocation and `conn.Read()` calls inside the accept loop to display raw incoming client bytes. Added `09_reading_bytes.md` lesson. Task 2.3 — Read bytes subtask complete.
- Sent raw text-based HTTP response to client using `conn.Write()` before connection closure. Added `10_writing_bytes.md` lesson. Task 2.3 — Write bytes subtask complete.
- Extracted connection logic to `handleConnection` and implemented robust cleanup using `defer`. Added `11_defer_and_closure.md`. Task 2.3 complete (6/6 subtasks).
