package router

import (
	"reflect"
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestRouter_MethodAndParamRouting(t *testing.T) {
	r := NewRouter()

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

	r.Get("/users/:id", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Body = []byte("User ID: " + req.Params["id"])
		return resp
	})

	r.Get("/static/*filepath", func(req *http.Request) *http.Response {
		resp := http.NewResponse()
		resp.StatusCode = http.StatusOK
		resp.Body = []byte("File: " + req.Params["filepath"])
		return resp
	})

	// Test 1: Hit existing GET route
	reqGet := &http.Request{Method: http.MethodGet, Path: "/hello", Params: make(map[string]string)}
	resp1 := r.ServeHTTP(reqGet)
	if resp1.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for GET /hello, got %d", resp1.StatusCode)
	}
	if string(resp1.Body) != "GET Hello" {
		t.Errorf("Expected body 'GET Hello', got %s", string(resp1.Body))
	}

	// Test 2: Hit existing POST route
	reqPost := &http.Request{Method: http.MethodPost, Path: "/hello", Params: make(map[string]string)}
	resp2 := r.ServeHTTP(reqPost)
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for POST /hello, got %d", resp2.StatusCode)
	}

	// Test 3: Hit existing path with unregistered method (should 405)
	reqPut := &http.Request{Method: http.MethodPut, Path: "/hello", Params: make(map[string]string)}
	resp3 := r.ServeHTTP(reqPut)
	if resp3.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 for PUT /hello, got %d", resp3.StatusCode)
	}

	// Test 4: Hit non-existent path (should 404)
	reqMissing := &http.Request{Method: http.MethodGet, Path: "/missing", Params: make(map[string]string)}
	resp4 := r.ServeHTTP(reqMissing)
	if resp4.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for /missing, got %d", resp4.StatusCode)
	}

	// Test 5: Parameterized route
	reqParam := &http.Request{
		Method: http.MethodGet,
		Path:   "/users/99",
		Params: make(map[string]string),
	}
	resp5 := r.ServeHTTP(reqParam)
	if resp5.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for /users/99, got %d", resp5.StatusCode)
	}
	if string(resp5.Body) != "User ID: 99" {
		t.Errorf("Expected body 'User ID: 99', got %s", string(resp5.Body))
	}
	expectedParams := map[string]string{"id": "99"}
	if !reflect.DeepEqual(reqParam.Params, expectedParams) {
		t.Errorf("Params not injected correctly. Got %v, want %v", reqParam.Params, expectedParams)
	}

	// Test 6: Wildcard route
	reqWild := &http.Request{
		Method: http.MethodGet,
		Path:   "/static/css/main.css",
		Params: make(map[string]string),
	}
	resp6 := r.ServeHTTP(reqWild)
	if resp6.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for /static/css/main.css, got %d", resp6.StatusCode)
	}
	if string(resp6.Body) != "File: css/main.css" {
		t.Errorf("Expected body 'File: css/main.css', got %s", string(resp6.Body))
	}
	expectedWildParams := map[string]string{"filepath": "css/main.css"}
	if !reflect.DeepEqual(reqWild.Params, expectedWildParams) {
		t.Errorf("Wildcard params not injected correctly. Got %v, want %v", reqWild.Params, expectedWildParams)
	}
}

// --- Benchmarks ---

func setupBenchmarkRouter() *Router {
	r := NewRouter()
	r.Get("/hello", func(req *http.Request) *http.Response { return nil })
	r.Get("/users/:id", func(req *http.Request) *http.Response { return nil })
	r.Get("/static/*filepath", func(req *http.Request) *http.Response { return nil })
	return r
}

func BenchmarkRouterStatic(b *testing.B) {
	r := setupBenchmarkRouter()
	req := &http.Request{Method: http.MethodGet, Path: "/hello", Params: make(map[string]string)}

	b.ReportAllocs()

	for b.Loop() {
		_ = r.ServeHTTP(req)
	}
}

func BenchmarkRouterParams(b *testing.B) {
	r := setupBenchmarkRouter()
	req := &http.Request{Method: http.MethodGet, Path: "/users/12345", Params: make(map[string]string)}

	b.ReportAllocs()

	for b.Loop() {
		_ = r.ServeHTTP(req)
	}
}

func BenchmarkRouterWildcard(b *testing.B) {
	r := setupBenchmarkRouter()
	req := &http.Request{Method: http.MethodGet, Path: "/static/css/main.css", Params: make(map[string]string)}

	b.ReportAllocs()

	for b.Loop() {
		_ = r.ServeHTTP(req)
	}
}

func BenchmarkRouter_Parallel(b *testing.B) {
	r := setupBenchmarkRouter()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// create a new req per goroutine to avoid map races
			localReq := &http.Request{Method: http.MethodGet, Path: "/users/999", Params: make(map[string]string)}
			_ = r.ServeHTTP(localReq)
		}
	})
}
