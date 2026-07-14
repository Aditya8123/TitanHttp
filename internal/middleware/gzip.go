package middleware

import (
	"compress/gzip"
	"io"
	"strings"
	"sync"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// CompressibleTypes maps MIME types that benefit from Gzip compression.
var CompressibleTypes = map[string]bool{
	"text/plain":                true,
	"text/html":                 true,
	"text/css":                  true,
	"application/json":          true,
	"application/javascript":    true,
	"application/xml":           true,
	"text/xml":                  true,
	"image/svg+xml":             true,
	"text/plain; charset=utf-8": true,
	"text/html; charset=utf-8":  true,
	"text/css; charset=utf-8":   true,
}

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

// Gzip is a middleware that dynamically compresses the response payload
// using gzip if the client supports it and the content type is compressible.
func Gzip(next router.Handler) router.Handler {
	return func(req *http.Request) *http.Response {
		resp := next(req)

		// 1. Negotiation: Check if client accepts gzip
		acceptEncoding := req.Headers.Get("accept-encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			return resp
		}

		// 2. Filter by Content-Type
		contentType := resp.Headers.Get("Content-Type")
		if !CompressibleTypes[contentType] {
			return resp
		}

		// Prevent double-compression
		if resp.Headers.Get("Content-Encoding") != "" {
			return resp
		}

		// We can't compress an empty response
		if len(resp.Body) == 0 && resp.Stream == nil {
			return resp
		}

		// 3. Prepare headers for chunked streaming
		resp.Headers.Set("Content-Encoding", "gzip")
		resp.Headers.Del("Content-Length")

		// 4. Wrap the response in an io.Pipe and pooled gzip.Writer
		pr, pw := io.Pipe()

		gw := gzipWriterPool.Get().(*gzip.Writer)
		gw.Reset(pw)

		// Capture original body and stream to prevent data races
		origStream := resp.Stream
		origBody := resp.Body

		// Replace the response payload with the compressed stream
		resp.Stream = pr
		resp.Body = nil // Clear in-memory body to save space

		// Spawn a background goroutine to stream compressed data to the pipe
		go func() {
			defer pw.Close()
			defer func() {
				gw.Close()
				gzipWriterPool.Put(gw)
			}()

			if origStream != nil {
				io.Copy(gw, origStream)
				if closer, ok := origStream.(io.Closer); ok {
					closer.Close()
				}
			} else if len(origBody) > 0 {
				gw.Write(origBody)
			}
		}()

		return resp
	}
}
