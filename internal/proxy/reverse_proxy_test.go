package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	titanhttp "github.com/Aditya8123/TitanHttp/internal/http"
)

func TestReverseProxy(t *testing.T) {
	var target string

	// 1. Create a dummy backend server using standard library
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/data" {
			t.Errorf("expected /api/data, got %s", r.URL.Path)
		}
		
		body, _ := io.ReadAll(r.Body)
		if string(body) != "client-data" {
			t.Errorf("expected 'client-data', got %q", string(body))
		}

		if r.Header.Get("X-Forwarded-For") != "192.168.1.100" {
			t.Errorf("expected X-Forwarded-For '192.168.1.100', got %q", r.Header.Get("X-Forwarded-For"))
		}
		if r.Header.Get("X-Forwarded-Host") != target {
			t.Errorf("expected X-Forwarded-Host %q, got %q", target, r.Header.Get("X-Forwarded-Host"))
		}
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			t.Errorf("expected X-Forwarded-Proto 'https', got %q", r.Header.Get("X-Forwarded-Proto"))
		}

		w.Header().Set("X-Backend-Header", "Hello from backend")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("backend-response"))
	}))
	defer backend.Close()

	// Extract host:port from backend URL (e.g., http://127.0.0.1:51234)
	target = strings.TrimPrefix(backend.URL, "http://")

	// 2. Create our reverse proxy handler
	proxyHandler := NewReverseProxy(target)

	// 3. Create a TitanHTTP request
	req := titanhttp.NewRequest()
	req.Method = titanhttp.MethodPost
	req.Path = "/api/data"
	req.Version = "HTTP/1.1"
	req.Headers.Set("host", target)
	req.Headers.Set("content-length", "11") // length of "client-data"
	req.RemoteAddr = "192.168.1.100:54321"
	req.Scheme = "https"
	req.Body = []byte("client-data")

	// 4. Send request through the proxy
	resp := proxyHandler(req)

	// 5. Verify the response
	if resp.StatusCode != titanhttp.StatusCreated {
		t.Fatalf("expected status 201, got %d", resp.StatusCode)
	}

	if resp.Headers.Get("x-backend-header") != "Hello from backend" {
		t.Errorf("expected header 'Hello from backend', got %v", resp.Headers.Get("x-backend-header"))
	}

	if resp.Stream == nil {
		t.Fatalf("expected response stream to be set")
	}

	// Read the proxied response body
	respBody, err := io.ReadAll(resp.Stream)
	if err != nil {
		t.Fatalf("failed to read response stream: %v", err)
	}

	if string(respBody) != "backend-response" {
		t.Errorf("expected body 'backend-response', got %q", string(respBody))
	}

	// Verify that the stream can be closed without error
	if closer, ok := resp.Stream.(io.Closer); ok {
		err := closer.Close()
		if err != nil {
			t.Errorf("failed to close stream: %v", err)
		}
	} else {
		t.Error("expected response stream to implement io.Closer")
	}
}
