package server

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	titanhttp "github.com/Aditya8123/TitanHttp/internal/http"
)

func TestServerTLS(t *testing.T) {
	// 1. Generate test certificates
	tmpDir := t.TempDir()
	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	if err := generateTestCertificate(certFile, keyFile); err != nil {
		t.Fatalf("Failed to generate test certificates: %v", err)
	}

	// 2. Start TitanHTTP in TLS mode
	// Use port 0 to automatically allocate a free port
	srv := NewServer("127.0.0.1:0")

	// Add a simple route for testing
	srv.Router().Get("/secure", func(req *titanhttp.Request) *titanhttp.Response {
		resp := titanhttp.NewResponse()
		resp.StatusCode = titanhttp.StatusOK
		resp.Headers.Set("Content-Type", "text/plain")
		resp.Body = []byte("Secure Payload")
		return resp
	})

	// Start the server in a goroutine using the generated cert and key
	go func() {
		if err := srv.StartTLS(certFile, keyFile); err != nil {
			t.Errorf("Server failed to start in TLS mode: %v", err)
		}
	}()

	// Give the server a moment to bind the port
	time.Sleep(100 * time.Millisecond)
	
	addr := srv.Addr()
	if addr == "127.0.0.1:0" || addr == "" {
		t.Fatalf("Server failed to allocate a port")
	}

	// 3. Make an HTTPS request to the server
	// We must configure the client to skip verification since we're using a self-signed cert
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get("https://" + addr + "/secure")
	if err != nil {
		t.Fatalf("Failed to make HTTPS request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if !strings.Contains(string(bodyBytes), "Secure Payload") {
		t.Errorf("Expected response to contain 'Secure Payload', got: %s", string(bodyBytes))
	}

	// 4. Clean shutdown
	client.CloseIdleConnections()
	srv.Shutdown(context.Background())
}
