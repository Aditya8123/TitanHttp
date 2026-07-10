package middleware

import (
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/cache"
	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestCacheMiddleware(t *testing.T) {
	memCache := cache.NewMemoryCache()
	middleware := CacheMiddleware(memCache)

	callCount := 0
	handler := middleware(func(req *http.Request) *http.Response {
		callCount++
		resp := http.NewResponse200()
		resp.Body = []byte("Hello")
		return resp
	})

	req := http.NewRequest()
	req.Method = http.MethodGet
	req.Path = "/test"

	// First call - Cache Miss
	resp1 := handler(req)
	if callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", callCount)
	}
	if string(resp1.Body) != "Hello" {
		t.Errorf("Expected body 'Hello', got '%s'", string(resp1.Body))
	}

	// Second call - Cache Hit
	resp2 := handler(req)
	if callCount != 1 {
		t.Errorf("Expected callCount to remain 1 on hit, got %d", callCount)
	}
	if string(resp2.Body) != "Hello" {
		t.Errorf("Expected body 'Hello', got '%s'", string(resp2.Body))
	}

	// Third call - non-GET request bypasses cache
	reqPost := http.NewRequest()
	reqPost.Method = http.MethodPost
	reqPost.Path = "/test"
	
	handler(reqPost)
	if callCount != 2 {
		t.Errorf("Expected callCount to be 2 after POST, got %d", callCount)
	}
}

func TestCacheMiddleware_Validation(t *testing.T) {
	memCache := cache.NewMemoryCache()
	middleware := CacheMiddleware(memCache)

	callCount := 0
	handler := middleware(func(req *http.Request) *http.Response {
		callCount++
		resp := http.NewResponse200()
		resp.Body = []byte("Validated")
		
		switch req.Path {
		case "/no-store":
			resp.Headers["cache-control"] = "no-store"
		case "/max-age":
			resp.Headers["cache-control"] = "max-age=60"
		}
		
		return resp
	})

	// 1. Response no-store
	req1 := http.NewRequest()
	req1.Method = http.MethodGet
	req1.Path = "/no-store"

	handler(req1)
	handler(req1) // Should miss again
	if callCount != 2 {
		t.Errorf("Expected callCount to be 2 for no-store response, got %d", callCount)
	}

	// 2. Request no-cache
	req2 := http.NewRequest()
	req2.Method = http.MethodGet
	req2.Path = "/max-age"
	
	handler(req2) // Miss 1 (callCount=3)
	handler(req2) // Hit 1 (callCount=3)
	if callCount != 3 {
		t.Errorf("Expected callCount to be 3 for max-age, got %d", callCount)
	}

	req2.Headers["cache-control"] = "no-cache"
	handler(req2) // Miss (bypasses cache) (callCount=4)
	if callCount != 4 {
		t.Errorf("Expected callCount to be 4 after bypass, got %d", callCount)
	}
}
