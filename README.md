# TitanHTTP

> A cinematic orchestration of raw infrastructure. Visualizing the invisible.

TitanHTTP is a production-inspired HTTP server built completely from scratch in Go. It is not just a backend system; it is an immersive, educational, and premium engineering experience designed to demystify the core of web infrastructure.

Here, data packets don't just transfer—they move like light through a datacenter. Invisible systems become visible. Complexity is revealed progressively, and motion always has meaning.

## ⚡ The Mission

The objective of TitanHTTP is to deeply master networking, operating systems, concurrency, and software architecture, and to showcase this mastery through high-end storytelling and data visualization. 

We deliberately avoid using Go's standard `net/http` for core server functionality. Instead, we build from the ground up—from raw TCP sockets to HTTP parsing, routing, and high-concurrency worker pools.

## 📊 Performance Benchmarks

TitanHTTP was rigorously profiled against Go's standard `net/http` library under simulated high-throughput and C10K scenarios. By utilizing zero-allocation parsing, a custom radix-tree router, and a bounded thread-safe worker pool, the architecture yields significant performance gains.

| Profile | TitanHTTP | `net/http` | Improvement |
| :--- | :--- | :--- | :--- |
| **Max Throughput** | **151,316 req/s** | 117,401 req/s | **+28.9%** |
| **Routing Speed** | **129,746 req/s** | 98,216 req/s | **+32.1%** |
| **Large Payloads** | **60,352 req/s** | 42,695 req/s | **+41.4%** |
| **Connection Churn** | **3,568 req/s** | 1,969 req/s | **+81.1%** |
| **P99 Latency (C10K)**| **4.75 ms** | 14.56 ms | **-67.4%** |
| **Memory Leaks** | **0 MB** | 0 MB | **Stable** |

> *Tests executed natively via `bombardier` handling massive concurrent loads. Full results, including framework comparisons (Gin, Fiber, Chi) and security tests, are available in [benchmarking.md](./docs/benchmarking.md).*

## 📖 Experience Chapters

The project is structured as a guided narrative. A visitor follows a single packet through the entire lifecycle:

* **Chapter 1:** Why HTTP Exists
* **Chapter 2:** TCP
* **Chapter 3:** Building a Socket
* **Chapter 4:** Reading Bytes
* **Chapter 5:** Parsing Requests
* **Chapter 6:** Routing
* **Chapter 7:** Concurrency
* **Chapter 8:** Production Features
* **Chapter 9:** Benchmarks
* **Chapter 10:** Source Code

*(See [phases.md](./docs/phases.md) for the strict engineering roadmap.)*

## 📚 Documentation Index

Our documentation is treated as a first-class, premium editorial experience. Read them to understand the architecture, decisions, and design language driving this project.

* [Architecture Blueprint](./docs/architecture.md) — How the HTTP lifecycle and components are structured.
* [Design Language](./docs/design.md) — The visual tokens, aesthetics, and emotion timeline.
* [Architectural Decisions (ADR)](./docs/decisions.md) — Context and trade-offs for major engineering choices.
* [Deployment Strategy](./docs/deployment.md) — How TitanHTTP is shipped.
* [Testing Philosophy](./docs/testing.md) — Ensuring rock-solid infrastructure.
* [Benchmarking](./docs/benchmarking.md) — Performance engineering metrics.
* [Recruiter Summary](./docs/recruiter.md) — A 10-minute executive tour of the technical challenges.
* [Project State](./docs/PROJECT_STATE.md) — The current snapshot of our development progress.
* [AI Orchestrator Rules](./AGENTS.md) — The strict guidelines that govern the AI Engineering Team building this project.

## 🛠️ Quick Start

```bash
# Clone the repository
git clone https://github.com/Aditya8123/TitanHttp.git

# Enter the datacenter
cd TitanHttp

# Boot the server
go run cmd/titanhttp/main.go
```

---
> *"Would a senior backend engineer enjoy reviewing this code, and would a recruiter trust the developer who built it?"* — The TitanHTTP standard.
