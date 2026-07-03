package middleware

import (
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/middleware/rate"
)

func TestRateLimitMiddleware(t *testing.T) {
	// 1 token capacity, 0 refill (never refills)
	tb := rate.NewTokenBucket(1, 0)
	middleware := RateLimitMiddleware(tb)

	callCount := 0
	handler := middleware(func(req *http.Request) *http.Response {
		callCount++
		return http.NewResponse200()
	})

	// Request 1: Should be allowed
	req1 := http.NewRequest()
	req1.RemoteAddr = "192.168.1.1:12345"

	resp1 := handler(req1)
	if resp1.StatusCode != 200 {
		t.Errorf("Expected 200 OK, got %d", resp1.StatusCode)
	}
	if callCount != 1 {
		t.Errorf("Expected handler to be called once")
	}

	// Request 2: Should be rejected (429) from same IP (different port)
	req2 := http.NewRequest()
	req2.RemoteAddr = "192.168.1.1:54321"

	resp2 := handler(req2)
	if resp2.StatusCode != 429 {
		t.Errorf("Expected 429 Too Many Requests, got %d", resp2.StatusCode)
	}
	if callCount != 1 {
		t.Errorf("Expected handler to NOT be called, callCount is %d", callCount)
	}

	// Request 3: Should be allowed (different IP)
	req3 := http.NewRequest()
	req3.RemoteAddr = "10.0.0.1:12345"

	resp3 := handler(req3)
	if resp3.StatusCode != 200 {
		t.Errorf("Expected 200 OK, got %d", resp3.StatusCode)
	}
	if callCount != 2 {
		t.Errorf("Expected handler to be called twice")
	}
}
