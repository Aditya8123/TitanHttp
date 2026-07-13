package server

import (
	"context"
	"crypto/tls"
	"io"
	"path/filepath"
	"testing"
	"time"
)

func TestServerHTTP2_ALPN(t *testing.T) {
	// 1. Generate test certificates
	tmpDir := t.TempDir()
	certFile := filepath.Join(tmpDir, "cert.pem")
	keyFile := filepath.Join(tmpDir, "key.pem")

	if err := generateTestCertificate(certFile, keyFile); err != nil {
		t.Fatalf("Failed to generate test certificates: %v", err)
	}

	// 2. Start TitanHTTP in TLS mode
	srv := NewServer("127.0.0.1:0")
	defer srv.Shutdown(context.Background())

	go func() {
		if err := srv.StartTLS(certFile, keyFile); err != nil {
			t.Errorf("Server failed to start in TLS mode: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	
	addr := srv.Addr()
	if addr == "127.0.0.1:0" || addr == "" {
		t.Fatalf("Server failed to allocate a port")
	}

	// 3. Connect via TLS forcing ALPN to "h2"
	conf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h2"},
	}

	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		t.Fatalf("Failed to dial TLS: %v", err)
	}
	defer conn.Close()

	// Verify that the server agreed to "h2"
	if conn.ConnectionState().NegotiatedProtocol != "h2" {
		t.Errorf("Expected ALPN to negotiate 'h2', got '%s'", conn.ConnectionState().NegotiatedProtocol)
	}

	// 4. Send the HTTP/2 preface
	preface := []byte("PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n")
	_, err = conn.Write(preface)
	if err != nil {
		t.Fatalf("Failed to write HTTP/2 preface: %v", err)
	}

	// Read from the connection. We expect EOF because the server stub closes it.
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		// Sometimes a connection reset can occur depending on OS, but typically we want EOF
		// We just don't want a panic or infinite hang
	}
	if n > 0 {
		t.Errorf("Expected connection to be gracefully closed without data, got %d bytes: %s", n, string(buf[:n]))
	}
}
