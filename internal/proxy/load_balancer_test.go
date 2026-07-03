package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	titanhttp "github.com/Aditya8123/TitanHttp/internal/http"
)

func TestLoadBalancer_BackendPool(t *testing.T) {
	// Create two dummy backends
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-ID", "1")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-ID", "2")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend2.Close()

	target1 := strings.TrimPrefix(backend1.URL, "http://")
	target2 := strings.TrimPrefix(backend2.URL, "http://")

	// Create the Load Balancer
	lbHandler := NewLoadBalancer([]string{target1, target2})

	// Make a request
	req := titanhttp.NewRequest()
	req.Method = titanhttp.MethodGet
	req.Path = "/"
	req.Version = "HTTP/1.1"
	req.Headers["host"] = target1

	// With Round Robin, current starts at 0. First call adds 1 -> idx 1 (backend2)
	// Second call adds 1 -> idx 2 % 2 = 0 (backend1).
	// Let's fire 4 requests and verify they alternate: 2, 1, 2, 1.
	expectedBackends := []string{"2", "1", "2", "1"}
	
	for i, expected := range expectedBackends {
		resp := lbHandler(req)
		if resp.StatusCode != 200 {
			t.Fatalf("Req %d: Expected 200, got %d", i, resp.StatusCode)
		}
		if got := resp.Headers["x-backend-id"]; got != expected {
			t.Errorf("Req %d: Expected to hit backend %s, got %v", i, expected, got)
		}

		// Read stream to prevent leaking
		if resp.Stream != nil {
			io.ReadAll(resp.Stream)
			if closer, ok := resp.Stream.(io.Closer); ok {
				closer.Close()
			}
		}
	}
}

func TestLoadBalancer_HealthCheckFailover(t *testing.T) {
	// Create two dummy backends
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-ID", "1")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-ID", "2")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend2.Close()

	target1 := strings.TrimPrefix(backend1.URL, "http://")
	target2 := strings.TrimPrefix(backend2.URL, "http://")

	// We can test failover by making the handler skip dead backends.
	lb := &LoadBalancer{
		backends: make([]*Backend, 0, 2),
	}
	
	b1 := &Backend{URL: target1, Alive: false}
	b2 := &Backend{URL: target2, Alive: true}
	lb.backends = append(lb.backends, b1, b2)
	
	req := titanhttp.NewRequest()
	req.Method = titanhttp.MethodGet
	req.Path = "/"
	req.Version = "HTTP/1.1"
	req.Headers["host"] = target2

	// Should hit backend 2 repeatedly
	for i := 0; i < 3; i++ {
		resp := lb.ServeHTTP(req)
		if resp.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", resp.StatusCode)
		}
		if resp.Headers["x-backend-id"] != "2" {
			t.Errorf("Expected to hit backend 2, got %v", resp.Headers["x-backend-id"])
		}
	}
}

func TestLoadBalancer_NoHealthyBackends(t *testing.T) {
	lbHandler := NewLoadBalancer([]string{})

	req := titanhttp.NewRequest()
	resp := lbHandler(req)

	if resp.StatusCode != 503 {
		t.Fatalf("Expected 503 when no backends are available, got %d", resp.StatusCode)
	}
}
