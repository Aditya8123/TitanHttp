package router

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

func TestStaticFiles(t *testing.T) {
	// Create a temporary directory for static files
	tmpDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tmpDir, "hello.txt")
	testContent := []byte("hello world from static file")
	if err := os.WriteFile(testFile, testContent, 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a test subdirectory
	subDir := filepath.Join(tmpDir, "assets")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatalf("Failed to create test subdirectory: %v", err)
	}

	// Create a css file in the subdirectory
	cssFile := filepath.Join(subDir, "style.css")
	cssContent := []byte("body { color: red; }")
	if err := os.WriteFile(cssFile, cssContent, 0o600); err != nil {
		t.Fatalf("Failed to create css file: %v", err)
	}

	// Create another subdirectory without index.html
	emptySubDir := filepath.Join(tmpDir, "empty")
	if err := os.Mkdir(emptySubDir, 0o755); err != nil {
		t.Fatalf("Failed to create empty subdirectory: %v", err)
	}

	// Create a subdirectory with index.html
	indexSubDir := filepath.Join(tmpDir, "docs")
	if err := os.Mkdir(indexSubDir, 0o755); err != nil {
		t.Fatalf("Failed to create docs subdirectory: %v", err)
	}
	indexFile := filepath.Join(indexSubDir, "index.html")
	indexContent := []byte("<h1>Docs</h1>")
	if err := os.WriteFile(indexFile, indexContent, 0o600); err != nil {
		t.Fatalf("Failed to create index file: %v", err)
	}

	r := NewRouter()
	r.Static("/static/", tmpDir)

	tests := []struct {
		name            string
		path            string
		expectedStatus  http.StatusCode
		expectedBody    []byte
		expectedHeaders map[string]string
	}{
		{
			name:           "Valid file",
			path:           "/static/hello.txt",
			expectedStatus: http.StatusOK,
			expectedBody:   testContent,
			expectedHeaders: map[string]string{
				"Content-Type":  "text/plain; charset=utf-8",
				"Cache-Control": "public, max-age=3600",
			},
		},
		{
			name:           "CSS file",
			path:           "/static/assets/style.css",
			expectedStatus: http.StatusOK,
			expectedBody:   cssContent,
			expectedHeaders: map[string]string{
				"Content-Type":  "text/css; charset=utf-8",
				"Cache-Control": "public, max-age=3600",
			},
		},
		{
			name:           "Missing file",
			path:           "/static/missing.txt",
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:           "Path traversal attempt",
			path:           "/static/../etc/passwd",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:           "Directory access without index",
			path:           "/static/empty",
			expectedStatus: http.StatusForbidden,
			expectedBody:   nil,
		},
		{
			name:           "Directory access with index",
			path:           "/static/docs",
			expectedStatus: http.StatusOK,
			expectedBody:   indexContent,
			expectedHeaders: map[string]string{
				"Content-Type":  "text/html; charset=utf-8",
				"Cache-Control": "public, max-age=3600",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Method: http.MethodGet,
				Path:   tt.path,
			}

			resp := r.ServeHTTP(req)

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectedBody != nil {
				if !bytes.Equal(resp.Body, tt.expectedBody) {
					t.Errorf("Expected body %q, got %q", string(tt.expectedBody), string(resp.Body))
				}
			}

			if tt.expectedHeaders != nil {
				for k, v := range tt.expectedHeaders {
					if resp.Headers[k] != v {
						t.Errorf("Expected header %s: %s, got %s", k, v, resp.Headers[k])
					}
				}
			}
		})
	}
}
