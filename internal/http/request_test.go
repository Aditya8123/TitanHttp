package http

import (
	"bytes"
	"strings"
	"testing"
)

func TestRequest_WriteTo(t *testing.T) {
	req := NewRequest()
	req.Method = MethodPost
	req.Path = "/api/test"
	req.Version = "HTTP/1.1"
	req.Headers.Set("Host", "localhost:8080")
	req.Headers.Set("Content-Type", "application/json")
	req.Headers.Set("Content-Length", "17")
	req.RawBody = []byte(`{"message":"hi"}`)

	var buf bytes.Buffer
	n, err := req.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	output := buf.String()
	if n != int64(len(output)) {
		t.Errorf("Expected written bytes %d, got %d", len(output), n)
	}

	// We check for exact string, but headers order in map iteration is random.
	// So we should check for presence of all parts.
	
	if !strings.HasPrefix(output, "POST /api/test HTTP/1.1\r\n") {
		t.Errorf("Request line is incorrect, got: %s", output)
	}

	if !strings.Contains(output, "host: localhost:8080") {
		t.Errorf("Missing or incorrect Host header: \n%s", output)
	}

	if !strings.Contains(output, "content-type: application/json") {
		t.Errorf("Missing or incorrect Content-Type header: \n%s", output)
	}

	if !strings.Contains(output, "content-length: 17") {
		t.Errorf("Missing or incorrect Content-Length header: \n%s", output)
	}

	if !strings.HasSuffix(output, "\r\n\r\n{\"message\":\"hi\"}") {
		t.Errorf("Body is incorrect, got: %q", output)
	}
}
