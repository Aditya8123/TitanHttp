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
	req.Headers["Host"] = "localhost:8080"
	req.Headers["Content-Type"] = "application/json"
	req.Headers["Content-Length"] = "17"
	req.Body = []byte(`{"message":"hi"}`)

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

	if !strings.Contains(output, "Host: localhost:8080\r\n") {
		t.Errorf("Missing or incorrect Host header")
	}

	if !strings.Contains(output, "Content-Type: application/json\r\n") {
		t.Errorf("Missing or incorrect Content-Type header")
	}

	if !strings.Contains(output, "Content-Length: 17\r\n") {
		t.Errorf("Missing or incorrect Content-Length header")
	}

	if !strings.HasSuffix(output, "\r\n\r\n{\"message\":\"hi\"}") {
		t.Errorf("Body is incorrect, got: %q", output)
	}
}
