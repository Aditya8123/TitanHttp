# Architectural Decisions

> A log of the major engineering choices made while building TitanHTTP.
>
> We use Architectural Decision Records (ADRs) to document the context, trade-offs, and reasoning behind significant design choices. This ensures that the *why* is preserved alongside the *how*.

---

## ADR 001: Building From Scratch in Pure Go

**Status:** Accepted

### Context
Go provides a robust, production-ready `net/http` package in its standard library. For almost any real-world application, utilizing this package is the correct engineering choice. However, the goal of TitanHTTP is not just to build a backend, but to deeply understand web infrastructure. 

### Decision
We will **not** use the `net/http` package to implement core server functionality. Instead, we will use the `net` package to interact directly with TCP sockets. We will manually implement HTTP request parsing, routing, response formatting, and connection management.

### Trade-offs & Consequences
- **Pro:** Deep educational value. It forces a strong understanding of the HTTP protocol specifications (RFC 7230, etc.) and raw socket programming.
- **Pro:** Complete control over memory allocations and concurrency models.
- **Con:** Increased development time and complexity.
- **Con:** We take on the responsibility of handling edge cases (e.g., malformed headers, slowloris attacks) that the standard library typically handles for free.

---

## ADR 002: Strict Development Hierarchy

**Status:** Accepted

### Context
Building a complete HTTP server from scratch is a complex, multi-layered task. Without a rigid structure, it is easy to get distracted by premature optimization or feature creep.

### Decision
Development will strictly follow the hierarchy defined in `phases.md` (Phase → Task → Subtask).
- Only **one** subtask may be actively implemented at any given time.
- A phase cannot be marked complete until all code, tests, and `.academy/` learning materials for its tasks are reviewed and finished.

### Trade-offs & Consequences
- **Pro:** Guaranteed forward momentum without technical debt.
- **Pro:** Easier for recruiters and reviewers to follow the commit history and understand the progression of complexity.
- **Con:** Feels artificially constrained if an engineer wants to jump ahead and implement a "fun" feature early (e.g., compression before basic routing).

---
---

## ADR 003: Radix Tree for Dynamic Routing

**Status:** Accepted

### Context
Our initial routing implementation used an O(1) hash map (`map[http.Method]map[string]Handler`). While this is extremely fast for exact path matches, it completely breaks down when introducing dynamic path parameters (e.g., `/users/:id/posts/:post_id`). We need a data structure capable of parameter extraction and wildcard matching without sacrificing performance by falling back to slow regular expressions.

### Decision
We will replace the Hash Map with a **Radix Tree** (a space-optimized Trie). The router will maintain one Radix Tree per HTTP Method. The `Request` struct will be extended with a `Params map[string]string` field to hold the extracted values, bypassing standard `context` injection for raw performance and simplicity.

### Trade-offs & Consequences
- **Pro:** Sub-microsecond routing lookups (O(k) where k is path depth).
- **Pro:** Built-in parameter extraction and prioritization (Exact > Parameter > Wildcard).
- **Con:** The `internal/router` package becomes significantly more complex to maintain and debug compared to a map.
- **Con:** Edge cases with conflicting parameter names at the same tree depth require strict validation during route registration.

---

## ADR 004: Middleware Pipeline Architecture

**Status:** Accepted

### Context
As the server complexity grows, we need a way to execute cross-cutting concerns (e.g., logging, panic recovery, authentication) across many routes without duplicating code inside every handler.

### Decision
We adopted the **Decorator Pattern** for middleware. A middleware is a function that takes a `router.Handler` and returns a new `router.Handler`. We built a global `router.Use()` chain, and also allow composing middlewares around specific routes (e.g., `middleware.AuthPlaceholder(myHandler)`). We avoided `net/http`'s `HandlerFunc` to maintain strict compatibility with our custom `Request` and `Response` structs.

### Trade-offs & Consequences
- **Pro:** Highly composable and idiopathic to Go web engineering.
- **Pro:** Allows route-specific protections (e.g., Auth only on `/protected`).
- **Con:** Middlewares wrap handlers in closures, which slightly increases the call stack depth and introduces a tiny amount of allocation overhead compared to inline execution.

---

## ADR 005: Bounded Worker Pool for Connection Handling

**Status:** Accepted

### Context
Handling thousands of concurrent connections by spawning unbounded goroutines (`go handleConnection(conn)`) can lead to resource exhaustion, memory out-of-bounds, and application crashes under sudden traffic spikes (e.g., DDOS). 

### Decision
We will employ a **Bounded Worker Pool** pattern. When the server starts, a fixed number of worker goroutines are spawned, blocking on a shared job channel (`chan net.Conn`). The main listener pushes accepted connections to this channel. If the queue reaches capacity, new connections receive an immediate `503 Service Unavailable` response and are closed, enforcing aggressive load shedding.

### Trade-offs & Consequences
- **Pro:** Hard boundary on CPU and memory usage, ensuring predictable performance under load.
- **Pro:** Load shedding protects the server from catastrophic cascading failure.
- **Con:** Connections might be dropped during massive bursts unless the queue size is tuned appropriately.
- **Con:** More complex connection lifecycle and shutdown sequences compared to naive one-goroutine-per-connection.

---

## ADR 006: Synchronization and Lock-Free Metrics

**Status:** Accepted

### Context
With a concurrent worker pool, shared state (like the Router table, active connections, request totals, and connection closures) needs protection from Data Races, which could corrupt memory or crash the application.

### Decision
We will enforce thread-safety using fine-grained synchronization primitives:
1. `sync.RWMutex` for the Router, allowing unlimited concurrent reads (pattern matching) while locking only for route registration.
2. `sync/atomic` for all server metrics (`totalRequests`, `activeConns`), circumventing mutex overhead for highly contended counters.
3. `sync.Mutex` purely to protect the `net.Listener` variable during graceful shutdown.

### Trade-offs & Consequences
- **Pro:** `sync/atomic` provides lock-free, high-throughput metric tracking, keeping latency negligible.
- **Pro:** `RWMutex` guarantees safe dynamic route injection without slowing down traffic parsing.
- **Con:** Atomic primitives restrict complex state interactions (e.g., you cannot easily transactionally update two metrics at once without CAS loops or a Mutex).

---

## ADR 007: Autonomous Background Sweepers for Resource Cleanup

**Status:** Accepted

### Context
In Phase 7, we introduced stateful infrastructure layers: a Memory Cache and Rate Limiters (Token Bucket and Sliding Window). In a high-traffic environment (or under a DDoS attack), tracking thousands of unique IPs or caching thousands of responses can silently exhaust server RAM (Out of Memory).

### Decision
We will employ **Autonomous Background Sweepers** in these components. When initializing a Cache or Rate Limiter, a background goroutine is spawned with a `time.Ticker`. This routine periodically locks the state map, scans for expired TTLs or stale IPs, and `delete()`s them to reclaim memory.

### Trade-offs & Consequences
- **Pro:** Completely eliminates memory leaks caused by unbounded state accumulation.
- **Pro:** The cleanup overhead is decoupled from the critical path of handling a client request (handlers don't have to pause to clean the whole map).
- **Con:** Introduces hidden goroutines that must be accounted for during server shutdown to prevent leaks of the sweepers themselves (though acceptable for the global lifespan of limiters).

---

## ADR 008: Lazy Evaluation in Token Bucket Rate Limiting

**Status:** Accepted

### Context
A naive implementation of a Token Bucket rate limiter spawns a background `time.Ticker` loop that iterates over every single IP in the bucket map and adds tokens every second. With 100,000 active IPs, this creates massive CPU overhead and locks the map constantly, freezing actual traffic.

### Decision
We will use **Lazy Evaluation** for token refills. Tokens are not actually refilled in the background. Instead, when a request arrives, the limiter calculates `(time.Now() - lastRefill) * rate`, adds the accumulated tokens to the bucket *at that exact moment*, and then processes the request.

### Trade-offs & Consequences
- **Pro:** O(1) CPU usage. Refilling costs practically zero CPU cycles because it only happens exactly when needed, mathematically.
- **Pro:** Eliminates map-wide lock contention, allowing the rate limiter to scale to millions of IPs.
- **Con:** The logic for time-delta math and token clamping is slightly more complex to test than a naive background adder.

---
> *"Code tells you how; comments tell you why."*
