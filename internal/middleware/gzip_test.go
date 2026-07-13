package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestGzipMiddleware(t *testing.T) {
	tests := []struct {
		name                 string
		acceptEncoding       string
		contentType          string
		initialBody          []byte
		expectCompression    bool
	}{
		{
			name:              "Compressible content with gzip accepted",
			acceptEncoding:    "gzip, deflate",
			contentType:       "text/plain",
			initialBody:       []byte("Hello, TitanHTTP! This is a test of the Gzip compression middleware. It should compress this text."),
			expectCompression: true,
		},
		{
			name:              "Incompressible content (image) with gzip accepted",
			acceptEncoding:    "gzip",
			contentType:       "image/jpeg",
			initialBody:       []byte("fake image data"),
			expectCompression: false,
		},
		{
			name:              "Compressible content without gzip accepted",
			acceptEncoding:    "",
			contentType:       "text/html",
			initialBody:       []byte("<h1>Hello</h1>"),
			expectCompression: false,
		},
		{
			name:              "Empty body is not compressed",
			acceptEncoding:    "gzip",
			contentType:       "text/plain",
			initialBody:       []byte(""),
			expectCompression: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock handler that returns the initial body
			mockHandler := func(req *http.Request) *http.Response {
				resp := http.NewResponse()
				resp.StatusCode = http.StatusOK
				resp.Headers.Set("Content-Type", tt.contentType)
				resp.Body = tt.initialBody
				return resp
			}

			// Wrap handler with Gzip middleware
			handler := Gzip(mockHandler)

			// Create request
			req := &http.Request{
				Method:  http.MethodGet,
				Path:    "/",
			}
			if tt.acceptEncoding != "" {
				req.Headers.Set("accept-encoding", tt.acceptEncoding)
			}

			// Execute
			resp := handler(req)

			// Validate
			if tt.expectCompression {
				if resp.Headers.Get("Content-Encoding") != "gzip" {
					t.Errorf("Expected Content-Encoding: gzip, got %s", resp.Headers.Get("Content-Encoding"))
				}
				if resp.Headers.Get("Content-Length") != "" {
					t.Errorf("Expected Content-Length to be removed")
				}
				if resp.Body != nil {
					t.Errorf("Expected resp.Body to be nil")
				}
				if resp.Stream == nil {
					t.Fatalf("Expected resp.Stream to be present")
				}

				// Read compressed data from stream
				compressedData, err := io.ReadAll(resp.Stream)
				if err != nil {
					t.Fatalf("Failed to read compressed stream: %v", err)
				}

				// Decompress to verify
				gr, err := gzip.NewReader(bytes.NewReader(compressedData))
				if err != nil {
					t.Fatalf("Failed to create gzip reader: %v", err)
				}
				defer gr.Close()

				decompressedData, err := io.ReadAll(gr)
				if err != nil {
					t.Fatalf("Failed to decompress data: %v", err)
				}

				if !bytes.Equal(decompressedData, tt.initialBody) {
					t.Errorf("Decompressed data mismatch. Expected %q, got %q", string(tt.initialBody), string(decompressedData))
				}
			} else {
				if resp.Headers.Get("Content-Encoding") == "gzip" {
					t.Errorf("Did not expect Content-Encoding: gzip")
				}
				if resp.Stream != nil {
					t.Errorf("Did not expect resp.Stream for uncompressed response")
				}
				if !bytes.Equal(resp.Body, tt.initialBody) {
					t.Errorf("Body mismatch. Expected %q, got %q", string(tt.initialBody), string(resp.Body))
				}
			}
		})
	}
}
