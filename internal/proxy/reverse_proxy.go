package proxy

import (
	"bufio"
	"net"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// NewReverseProxy returns a handler that proxies requests to the specified target.
// target should be in the format "host:port", e.g., "localhost:8081".
func NewReverseProxy(target string) router.Handler {
	return func(req *http.Request) *http.Response {
		return ForwardRequest(req, target)
	}
}

// ForwardRequest dials the target backend, forwards the given request,
// and returns the backend's response with a streaming body.
func ForwardRequest(req *http.Request, target string) *http.Response {
	// 1. Dial the backend
	conn, err := net.DialTimeout("tcp", target, 5*time.Second)
	if err != nil {
		return http.NewResponse500()
	}

	// Set deadlines to prevent hanging
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Add Header Management
	req.Headers.Set("connection", "close")
	
	// X-Forwarded-For: append client IP
	clientIP, _, _ := net.SplitHostPort(req.RemoteAddr)
	if existing := req.Headers.Get("x-forwarded-for"); existing != "" {
		req.Headers.Set("x-forwarded-for", existing + ", " + clientIP)
	} else {
		req.Headers.Set("x-forwarded-for", clientIP)
	}

	// X-Forwarded-Host and X-Forwarded-Proto
	if req.Headers.Get("x-forwarded-host") == "" {
		req.Headers.Set("x-forwarded-host", req.Headers.Get("host"))
	}
	if req.Headers.Get("x-forwarded-proto") == "" {
		req.Headers.Set("x-forwarded-proto", req.Scheme)
	}

	// 2. Write the request to the backend
	_, err = req.WriteTo(conn)
	if err != nil {
		conn.Close()
		return http.NewResponse500()
	}

	// 3. Read the response from the backend
	reader := bufio.NewReader(conn)
	resp, err := http.ParseResponse(reader)
	if err != nil {
		conn.Close()
		return http.NewResponse500()
	}

	// 4. Wrap the stream to ensure the connection is closed when the response
	// is fully streamed to the client.
	resp.Stream = &connReadCloser{
		reader: reader,
		conn:   conn,
	}

	return resp
}

// connReadCloser wraps a bufio.Reader and a net.Conn.
// It implements io.Reader by reading from the buffered reader (preserving any buffered bytes),
// and implements io.Closer by closing the underlying network connection.
type connReadCloser struct {
	reader *bufio.Reader
	conn   net.Conn
}

func (c *connReadCloser) Read(p []byte) (n int, err error) {
	return c.reader.Read(p)
}

func (c *connReadCloser) Close() error {
	return c.conn.Close()
}
