package router

import (
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestRouter_AddRoute_And_ServeHTTP(t *testing.T) {
	r := NewRouter()

	// Register a route
	r.AddRoute("/hello", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Body = []byte("Hello, World!")
		return resp
	})

	// Test 1: Hit existing route
	reqHello := &http.Request{
		Method: http.MethodGet,
		Path:   "/hello",
	}
	resp1 := r.ServeHTTP(reqHello)
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for /hello, got %d", resp1.StatusCode)
	}
	if string(resp1.Body) != "Hello, World!" {
		t.Errorf("Expected body 'Hello, World!', got %s", string(resp1.Body))
	}

	// Test 2: Hit non-existent route (should 404)
	reqMissing := &http.Request{
		Method: http.MethodGet,
		Path:   "/missing",
	}
	resp2 := r.ServeHTTP(reqMissing)
	if resp2.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for /missing, got %d", resp2.StatusCode)
	}
}
