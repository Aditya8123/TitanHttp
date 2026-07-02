# Testing Philosophy

> Testing in TitanHTTP is not an afterthought; it is the structural integrity of the project.

Because we are building core infrastructure and parsing raw byte streams, the cost of a bug is extremely high. A malformed HTTP header parser can lead to infinite loops, crashes, or security vulnerabilities (e.g., HTTP Request Smuggling). 

Therefore, testing is mandatory and strictly enforced before any Subtask can be marked as complete.

## Test Pyramid

### 1. Unit Tests (The Foundation)
The vast majority of our tests will be unit tests targeting specific components in isolation.
- **Parser:** We will heavily test the HTTP parser with edge cases: missing CRLFs, malformed headers, massive payloads, and invalid characters.
- **Router:** Ensuring routes match correctly, parameters are extracted perfectly, and wildcards do not bleed.
- **Worker Pool:** Validating that jobs are dispatched and completed without race conditions or memory leaks.

*Tooling:* Go's standard `testing` package. Table-driven tests are strongly preferred for covering a wide array of input/output scenarios cleanly.

### 2. Integration Tests (The Pipeline)
These tests ensure that components work together. They will simulate the request lifecycle without opening an actual network port.
- Passing raw bytes into a mocked connection interface.
- Verifying the parser correctly hands the Request to the Router.
- Verifying the Router invokes the correct Handler.
- Verifying the Response Writer outputs the correct raw HTTP bytes.

### 3. End-to-End (E2E) Tests (The Reality Check)
These tests spin up the actual HTTP server on a local port and hit it with a real HTTP client (like `curl` or Go's standard `net/http` client).
- **Concurrency Testing:** Firing thousands of requests simultaneously to ensure the worker pool scales and handles connection churn correctly.
- **Keep-Alive:** Ensuring multiple requests over a single TCP connection are processed correctly.

## The Rule of Regressions

If a bug is found in TitanHTTP, the fix must always be accompanied by a test that explicitly reproduces the bug. This guarantees the bug will never return. 

Tests must be deterministic. Flaky tests (tests that fail intermittently, often due to timing issues in concurrency) are treated as failing tests and must be fixed immediately.

---
> *Design Note: Code without tests is a theory. Tests are the proof.*
