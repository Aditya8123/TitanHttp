package middleware

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestAuthPlaceholder(t *testing.T) {
	// Success case
	handler := AuthPlaceholder(func(req *http.Request) *http.Response {
		return http.NewResponse200()
	})

	req := http.NewRequest()
	req.Headers.Set("Authorization", "Bearer titan-secret-token")
	resp := handler(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	// Failure case
	req2 := http.NewRequest()
	req2.Headers.Set("Authorization", "Bearer invalid-token")
	resp2 := handler(req2)
	if resp2.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", resp2.StatusCode)
	}
}

func TestRecovery(t *testing.T) {
	handler := Recovery(func(req *http.Request) *http.Response {
		panic("test panic")
	})

	req := http.NewRequest()
	resp := handler(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected 500 Internal Server Error, got %d", resp.StatusCode)
	}
}

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	LogOutput = &buf
	defer func() {
		LogOutput = nil
	}()

	handler := Logger(func(req *http.Request) *http.Response {
		return http.NewResponse200()
	})

	req := http.NewRequest()
	req.Method = http.MethodGet
	req.Path = "/test-log"
	req.Version = "HTTP/1.1"

	resp := handler(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	output := buf.String()
	if !strings.Contains(output, "[TitanHTTP] → GET /test-log HTTP/1.1") {
		t.Errorf("Log output missing request details, got:\n%s", output)
	}
	if !strings.Contains(output, "[TitanHTTP] ← 200 OK") {
		t.Errorf("Log output missing response details, got:\n%s", output)
	}
}

func TestRange(t *testing.T) {
	handler := Range(func(req *http.Request) *http.Response {
		resp := http.NewResponse200()
		resp.Body = []byte("0123456789")
		return resp
	})

	// Match Range
	req := http.NewRequest()
	req.Headers.Set("Range", "bytes=2-6")
	resp := handler(req)
	if resp.StatusCode != http.StatusPartialContent {
		t.Errorf("Expected 206 Partial Content, got %d", resp.StatusCode)
	}
	if string(resp.Body) != "23456" {
		t.Errorf("Expected body '23456', got %q", string(resp.Body))
	}
	if val := resp.Headers.Get("Content-Range"); val != "bytes 2-6/10" {
		t.Errorf("Expected Content-Range 'bytes 2-6/10', got %q", val)
	}

	// Bypass Range when header is empty
	req2 := http.NewRequest()
	resp2 := handler(req2)
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp2.StatusCode)
	}
	if string(resp2.Body) != "0123456789" {
		t.Errorf("Expected full body, got %q", string(resp2.Body))
	}
}
