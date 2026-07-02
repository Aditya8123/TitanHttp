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
> *"Code tells you how; comments tell you why."*
