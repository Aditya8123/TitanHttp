package router

import (
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestRouter_MethodRouting(t *testing.T) {
	r := NewRouter()

	// Register a GET and POST route
	r.Get("/hello", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Body = []byte("GET Hello")
		return resp
	})

	r.Post("/hello", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Body = []byte("POST Hello")
		return resp
	})

	// Test 1: Hit existing GET route
	reqGet := &http.Request{
		Method: http.MethodGet,
		Path:   "/hello",
	}
	resp1 := r.ServeHTTP(reqGet)
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for GET /hello, got %d", resp1.StatusCode)
	}
	if string(resp1.Body) != "GET Hello" {
		t.Errorf("Expected body 'GET Hello', got %s", string(resp1.Body))
	}

	// Test 2: Hit existing POST route
	reqPost := &http.Request{
		Method: http.MethodPost,
		Path:   "/hello",
	}
	resp2 := r.ServeHTTP(reqPost)
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for POST /hello, got %d", resp2.StatusCode)
	}
	if string(resp2.Body) != "POST Hello" {
		t.Errorf("Expected body 'POST Hello', got %s", string(resp2.Body))
	}

	// Test 3: Hit existing path with unregistered method (should 405)
	reqPut := &http.Request{
		Method: http.MethodPut,
		Path:   "/hello",
	}
	resp3 := r.ServeHTTP(reqPut)
	if resp3.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for PUT /hello, got %d", resp3.StatusCode)
	}

	// Test 4: Hit non-existent path (should 404)
	reqMissing := &http.Request{
		Method: http.MethodGet,
		Path:   "/missing",
	}
	resp4 := r.ServeHTTP(reqMissing)
	if resp4.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for /missing, got %d", resp4.StatusCode)
	}
}
