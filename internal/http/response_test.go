package http

import (
	"strings"
	"testing"
)

func TestResponseBytes(t *testing.T) {
	resp := NewResponse()
	resp.StatusCode = StatusOK
	resp.Headers["Content-Type"] = "text/plain"
	resp.Body = []byte("Hello TitanHTTP")

	bytes := resp.Bytes()
	str := string(bytes)

	if !strings.HasPrefix(str, "HTTP/1.1 200 OK\r\n") {
		t.Errorf("Expected status line to be HTTP/1.1 200 OK\\r\\n, got %q", str)
	}

	if !strings.Contains(str, "Content-Length: 15\r\n") {
		t.Errorf("Expected Content-Length: 15, got %q", str)
	}

	if !strings.Contains(str, "Content-Type: text/plain\r\n") {
		t.Errorf("Expected Content-Type: text/plain, got %q", str)
	}

	if !strings.HasSuffix(str, "\r\nHello TitanHTTP") {
		t.Errorf("Expected body to be Hello TitanHTTP at the end, got %q", str)
	}
}

func TestResponseHelpers(t *testing.T) {
	resp400 := NewResponse400()
	if resp400.StatusCode != StatusBadRequest {
		t.Errorf("Expected 400 status, got %d", resp400.StatusCode)
	}

	resp404 := NewResponse404()
	if resp404.StatusCode != StatusNotFound {
		t.Errorf("Expected 404 status, got %d", resp404.StatusCode)
	}

	resp500 := NewResponse500()
	if resp500.StatusCode != StatusInternalServerError {
		t.Errorf("Expected 500 status, got %d", resp500.StatusCode)
	}
}

// --- Benchmarks ---

type devNullWriter struct{}

func (w devNullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkResponseWriteTo_LargePayload(b *testing.B) {
	payload := make([]byte, 1024*1024) // 1MB payload
	for i := range payload {
		payload[i] = 'a'
	}

	resp := NewResponse()
	resp.StatusCode = StatusOK
	resp.Headers["Content-Type"] = "text/plain"
	resp.Body = payload

	writer := devNullWriter{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = resp.WriteTo(writer)
	}
}
