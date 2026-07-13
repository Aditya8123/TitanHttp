package http

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestParseRequestLine(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError error
		expectedReq   *Request
	}{
		{
			name:          "Valid GET request",
			input:         "GET /index.html HTTP/1.1\r\n",
			expectedError: nil,
			expectedReq: &Request{
				Method:  MethodGet,
				Path:    "/index.html",
				Version: "HTTP/1.1",
			},
		},
		{
			name:          "Valid POST request",
			input:         "POST /api/data HTTP/1.1\r\n",
			expectedError: nil,
			expectedReq: &Request{
				Method:  MethodPost,
				Path:    "/api/data",
				Version: "HTTP/1.1",
			},
		},
		{
			name:          "Missing CRLF",
			input:         "GET /index.html HTTP/1.1\n", // only LF
			expectedError: ErrMalformedRequest,
		},
		{
			name:          "Missing parts",
			input:         "GET /index.html\r\n",
			expectedError: ErrMalformedRequest,
		},
		{
			name:          "Invalid URI",
			input:         "GET index.html HTTP/1.1\r\n", // missing leading slash
			expectedError: ErrInvalidURI,
		},
		{
			name:          "Invalid Version",
			input:         "GET /index.html HTP/1.1\r\n", // typo in version
			expectedError: ErrInvalidVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			req := NewRequest()

			err := parseRequestLine(reader, req)

			if err != tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil {
				if req.Method != tt.expectedReq.Method {
					t.Errorf("expected Method %q, got %q", tt.expectedReq.Method, req.Method)
				}
				if req.Path != tt.expectedReq.Path {
					t.Errorf("expected Path %q, got %q", tt.expectedReq.Path, req.Path)
				}
				if req.Version != tt.expectedReq.Version {
					t.Errorf("expected Version %q, got %q", tt.expectedReq.Version, req.Version)
				}
			}
		})
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedError   error
		expectedHeaders map[string]string
	}{
		{
			name:          "Valid headers",
			input:         "Host: localhost:8080\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
			expectedError: nil,
			expectedHeaders: map[string]string{
				"host":       "localhost:8080",
				"user-agent": "curl/7.81.0",
				"accept":     "*/*",
			},
		},
		{
			name:          "Headers with extra whitespace",
			input:         "   Content-Type   :   application/json   \r\n\r\n",
			expectedError: nil,
			expectedHeaders: map[string]string{
				"content-type": "application/json",
			},
		},
		{
			name:          "No headers",
			input:         "\r\n",
			expectedError: nil,
			expectedHeaders: map[string]string{},
		},
		{
			name:          "Missing CRLF",
			input:         "Host: localhost:8080\n\r\n",
			expectedError: ErrMalformedHeader,
		},
		{
			name:          "Missing colon",
			input:         "Host localhost 8080\r\n\r\n",
			expectedError: ErrMalformedHeader,
		},
		{
			name:          "Missing key",
			input:         ": localhost:8080\r\n\r\n",
			expectedError: ErrMalformedHeader,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			req := NewRequest()

			err := parseHeaders(reader, &req.Headers)

			if err != tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil {
				if len(req.Headers.Entries()) != len(tt.expectedHeaders) {
					t.Fatalf("expected %d headers, got %d", len(tt.expectedHeaders), len(req.Headers.Entries()))
				}

				for k, expectedVal := range tt.expectedHeaders {
					val := req.Headers.Get(k)
					if val == "" {
						t.Errorf("missing expected header %q", k)
					}
					if val != expectedVal {
						t.Errorf("for header %q, expected %q, got %q", k, expectedVal, val)
					}
				}
			}
		})
	}
}

func TestParseBody(t *testing.T) {
	tests := []struct {
		name          string
		headers       map[string]string
		input         string
		expectedError error
		expectedBody  string
	}{
		{
			name:          "No body",
			headers:       map[string]string{},
			input:         "",
			expectedError: nil,
			expectedBody:  "",
		},
		{
			name: "Valid body",
			headers: map[string]string{
				"content-length": "11",
			},
			input:         "hello world",
			expectedError: nil,
			expectedBody:  "hello world",
		},
		{
			name: "Invalid content length string",
			headers: map[string]string{
				"content-length": "abc",
			},
			input:         "hello world",
			expectedError: ErrInvalidContentLength,
		},
		{
			name: "Negative content length",
			headers: map[string]string{
				"content-length": "-1",
			},
			input:         "hello world",
			expectedError: ErrInvalidContentLength,
		},
		{
			name: "Too large content length",
			headers: map[string]string{
				"content-length": "20000000", // > 10MB
			},
			input:         "hello world",
			expectedError: ErrBodyTooLarge,
		},
		{
			name: "Premature EOF",
			headers: map[string]string{
				"content-length": "20",
			},
			input:         "hello world", // only 11 bytes
			expectedError: io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			req := NewRequest()
			for k, v := range tt.headers {
				req.Headers.Set(k, v)
			}

			err := ParseBody(reader, req)

			if err != tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}

			if tt.expectedError == nil {
				if string(req.RawBody) != tt.expectedBody {
					t.Errorf("expected body %q, got %q", tt.expectedBody, string(req.RawBody))
				}
			}
		})
	}
}

func TestRequestValidation(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		headers       map[string]string
		expectedError error
	}{
		{
			name:    "HTTP/1.1 with Host header",
			version: "HTTP/1.1",
			headers: map[string]string{
				"host": "localhost:8080",
			},
			expectedError: nil,
		},
		{
			name:          "HTTP/1.1 without Host header",
			version:       "HTTP/1.1",
			headers:       map[string]string{},
			expectedError: ErrMissingHostHeader,
		},
		{
			name:          "HTTP/1.0 without Host header",
			version:       "HTTP/1.0",
			headers:       map[string]string{},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewRequest()
			req.Version = tt.version
			for k, v := range tt.headers {
				req.Headers.Set(k, v)
			}

			err := req.Validate()
			if err != tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	rawResponse := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 5\r\n" +
		"\r\n" +
		"hello"

	reader := bufio.NewReader(strings.NewReader(rawResponse))
	resp, err := ParseResponse(reader)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Version != "HTTP/1.1" {
		t.Errorf("expected version HTTP/1.1, got %v", resp.Version)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, got %v", resp.StatusCode)
	}

	if resp.StatusText != "OK" {
		t.Errorf("expected status text 'OK', got %q", resp.StatusText)
	}

	if resp.Headers.Get("content-type") != "text/plain" {
		t.Errorf("expected content-type 'text/plain', got %v", resp.Headers.Get("content-type"))
	}

	if resp.Stream == nil {
		t.Errorf("expected resp.Stream to be set")
	} else {
		// Read the body from the stream
		body, err := io.ReadAll(resp.Stream)
		if err != nil {
			t.Fatalf("failed to read from stream: %v", err)
		}
		if string(body) != "hello" {
			t.Errorf("expected body 'hello', got %q", string(body))
		}
	}
}

// --- Benchmarks ---

func BenchmarkParseRequestLine(b *testing.B) {
	reqLine := "GET /api/v1/users/12345 HTTP/1.1\r\n"
	req := NewRequest()
	
	sr := strings.NewReader(reqLine)
	reader := bufio.NewReader(sr)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sr.Reset(reqLine)
		reader.Reset(sr)
		_ = parseRequestLine(reader, req)
	}
}

func BenchmarkParseHeaders(b *testing.B) {
	headersRaw := "Host: localhost:8080\r\nUser-Agent: curl/7.81.0\r\nAccept: application/json\r\nConnection: keep-alive\r\n\r\n"
	
	sr := strings.NewReader(headersRaw)
	reader := bufio.NewReader(sr)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sr.Reset(headersRaw)
		reader.Reset(sr)
		var headers Header
		_ = parseHeaders(reader, &headers)
	}
}

func BenchmarkParseBody_Large(b *testing.B) {
	bodyContent := strings.Repeat("a", 1024*1024) // 1MB body
	sr := strings.NewReader(bodyContent)
	reader := bufio.NewReader(sr)

	req := requestPool.Get().(*Request)
	defer requestPool.Put(req)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sr.Reset(bodyContent)
		reader.Reset(sr)
		req.Headers.Set("Content-Length", "1048576") // 1MB
		_ = ParseBody(reader, req)
	}
}

func BenchmarkParseRequestFull_Parallel(b *testing.B) {
	rawRequest := "POST /api/upload HTTP/1.1\r\n" +
		"Host: localhost:8080\r\n" +
		"Content-Length: 15\r\n" +
		"Content-Type: application/json\r\n" +
		"\r\n" +
		`{"key":"value"}`

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		sr := strings.NewReader(rawRequest)
		reader := bufio.NewReader(sr)
		for pb.Next() {
			sr.Reset(rawRequest)
			reader.Reset(sr)
			req, _ := ParseRequest(reader)
			if req != nil {
				ReleaseRequest(req)
			}
		}
	})
}
