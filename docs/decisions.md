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
> *"Code tells you how; comments tell you why."*
