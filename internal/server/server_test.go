package server

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestServer_GracefulShutdown(t *testing.T) {
	// Start server on a random port
	srv := NewServer("127.0.0.1:0")
	
	// Add a slow route that blocks for 200ms
	srv.Router().Get("/slow", func(req *http.Request) *http.Response {
		time.Sleep(200 * time.Millisecond)
		resp := http.NewResponse200()
		resp.Body = []byte("done")
		return resp
	})

	// Start server in background
	go srv.Start()
	time.Sleep(50 * time.Millisecond) // Wait for server to bind

	// Determine actual bound address safely
	addr := srv.Addr()

	// Launch a request that will take 200ms
	reqCompleted := make(chan struct{})
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Write([]byte("GET /slow HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n"))
			io.ReadAll(conn)
			conn.Close()
		}
		close(reqCompleted)
	}()

	// Wait briefly to ensure request is being handled
	time.Sleep(50 * time.Millisecond)

	// Call Shutdown with a 500ms timeout
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Expected graceful shutdown, got error: %v", err)
	}

	// Wait for the request to complete to verify it wasn't forcefully killed before shutdown finished
	select {
	case <-reqCompleted:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Request did not complete during graceful shutdown")
	}
}

func TestServer_Metrics(t *testing.T) {
	srv := NewServer("127.0.0.1:0")
	
	srv.Router().Get("/ping", func(req *http.Request) *http.Response {
		return http.NewResponse200()
	})

	go srv.Start()
	time.Sleep(50 * time.Millisecond)

	addr := srv.Addr()

	// Send 5 concurrent requests
	for i := 0; i < 5; i++ {
		conn, _ := net.Dial("tcp", addr)
		conn.Write([]byte("GET /ping HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n"))
		io.ReadAll(conn)
		conn.Close()
	}

	// Give the server a moment to finish metric updates
	time.Sleep(50 * time.Millisecond)
	
	// Shut down server
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	// Verify metrics
	if reqs := srv.metrics.totalRequests.Load(); reqs != 5 {
		t.Errorf("Expected 5 total requests, got %d", reqs)
	}
	
	if active := srv.metrics.activeConns.Load(); active != 0 {
		t.Errorf("Expected 0 active connections after shutdown, got %d", active)
	}
}

func TestServer_MetricsKeepAlive(t *testing.T) {
	srv := NewServer("127.0.0.1:0")
	
	srv.Router().Get("/ping", func(req *http.Request) *http.Response {
		return http.NewResponse200()
	})

	go srv.Start()
	time.Sleep(50 * time.Millisecond)

	addr := srv.Addr()

	// Establish a single Keep-Alive connection
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}

	// Send 3 requests over the same connection
	reqStr := "GET /ping HTTP/1.1\r\nHost: localhost\r\nConnection: keep-alive\r\n\r\n"
	for i := 0; i < 3; i++ {
		_, err := conn.Write([]byte(reqStr))
		if err != nil {
			t.Fatalf("Failed to write request %d: %v", i, err)
		}
		
		// Read response
		buf := make([]byte, 1024)
		conn.Read(buf)
	}

	// Wait for metrics to update
	time.Sleep(50 * time.Millisecond)

	// activeConns should be 1, totalRequests should be 3
	if active := srv.metrics.activeConns.Load(); active != 1 {
		t.Errorf("Expected 1 active connection, got %d", active)
	}
	if reqs := srv.metrics.totalRequests.Load(); reqs != 3 {
		t.Errorf("Expected 3 total requests, got %d", reqs)
	}

	conn.Close()

	// Shut down server
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	// activeConns should be 0
	if active := srv.metrics.activeConns.Load(); active != 0 {
		t.Errorf("Expected 0 active connections after shutdown, got %d", active)
	}
}
