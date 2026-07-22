package router

import (
	"strings"
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestMiddlewareChainOrder(t *testing.T) {
	var executionOrder []string

	// Middleware A
	middlewareA := func(next Handler) Handler {
		return func(req *http.Request) *http.Response {
			executionOrder = append(executionOrder, "A_Pre")
			resp := next(req)
			executionOrder = append(executionOrder, "A_Post")
			return resp
		}
	}

	// Middleware B
	middlewareB := func(next Handler) Handler {
		return func(req *http.Request) *http.Response {
			executionOrder = append(executionOrder, "B_Pre")
			resp := next(req)
			executionOrder = append(executionOrder, "B_Post")
			return resp
		}
	}

	// Base Handler
	baseHandler := func(req *http.Request) *http.Response {
		executionOrder = append(executionOrder, "Handler")
		return http.NewResponse()
	}

	// Chain: A -> B -> Handler
	chain := Chain(middlewareA, middlewareB)
	finalHandler := chain(baseHandler)

	// Execute
	req := http.NewRequest()
	finalHandler(req)

	// Verify order
	expectedOrder := []string{"A_Pre", "B_Pre", "Handler", "B_Post", "A_Post"}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d steps, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, step := range executionOrder {
		if step != expectedOrder[i] {
			t.Errorf("Step %d: expected %s, got %s", i, expectedOrder[i], step)
		}
	}
}

func TestMiddlewareShortCircuit(t *testing.T) {
	var executedHandler bool

	// Middleware that short-circuits (doesn't call next)
	shortCircuitMiddleware := func(next Handler) Handler {
		return func(req *http.Request) *http.Response {
			resp := http.NewResponse()
			resp.StatusCode = http.StatusMethodNotAllowed
			resp.Body = []byte("Short Circuit")
			return resp
		}
	}

	// Base Handler
	baseHandler := func(req *http.Request) *http.Response {
		executedHandler = true
		return http.NewResponse()
	}

	chain := Chain(shortCircuitMiddleware)
	finalHandler := chain(baseHandler)

	req := http.NewRequest()
	resp := finalHandler(req)

	if executedHandler {
		t.Errorf("Expected handler to NOT execute due to short circuit")
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	if strings.TrimSpace(string(resp.Body)) != "Short Circuit" {
		t.Errorf("Expected body 'Short Circuit', got '%s'", string(resp.Body))
	}
}

// --- Benchmarks ---

func dummyMiddleware(next Handler) Handler {
	return func(req *http.Request) *http.Response {
		// simulate some very light work
		_ = req.Method
		return next(req)
	}
}

func BenchmarkMiddlewareChain_10Layers(b *testing.B) {
	baseHandler := func(req *http.Request) *http.Response {
		return http.NewResponse()
	}

	middlewares := make([]Middleware, 10)
	for i := 0; i < 10; i++ {
		middlewares[i] = dummyMiddleware
	}

	chain := Chain(middlewares...)
	finalHandler := chain(baseHandler)

	req := http.NewRequest()
	req.Method = http.MethodGet

	b.ReportAllocs()

	for b.Loop() {
		_ = finalHandler(req)
	}
}
