package http

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// Method represents an HTTP request method
type Method string

const (
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodPatch   Method = "PATCH"
	MethodDelete  Method = "DELETE"
	MethodConnect Method = "CONNECT"
	MethodOptions Method = "OPTIONS"
	MethodTrace   Method = "TRACE"
)

// Request represents a parsed HTTP request.
// In HTTP/1.1, a request consists of a Request-Line, Headers, and an optional Body.
type Request struct {
	// Method is the HTTP method (e.g., GET, POST).
	Method Method

	// Path is the requested URI (e.g., /index.html).
	Path string

	// Version is the HTTP protocol version (e.g., HTTP/1.1).
	Version string

	// Headers stores the key-value pairs of the HTTP headers.
	// HTTP headers are case-insensitive, but we typically store them in a canonical format.
	Headers map[string]string

	// Body contains the payload of the request, if any.
	Body []byte

	// Params stores dynamic path parameters extracted by the router (e.g., /users/:id).
	Params map[string]string

	// RemoteAddr is the network address of the client that sent the request.
	RemoteAddr string

	// Scheme is the protocol scheme (e.g., "http" or "https").
	Scheme string
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &Request{
			Headers: make(map[string]string),
			Params:  make(map[string]string),
		}
	},
}

// AcquireRequest fetches a clean Request object from the pool.
func AcquireRequest() *Request {
	req := requestPool.Get().(*Request)
	// Maps are already allocated, just make sure they are empty.
	// (Usually done on release, but safe to do here).
	return req
}

// ReleaseRequest cleans up the Request object and returns it to the pool.
func ReleaseRequest(req *Request) {
	req.Method = ""
	req.Path = ""
	req.Version = ""
	req.Body = nil
	req.RemoteAddr = ""
	req.Scheme = ""

	for k := range req.Headers {
		delete(req.Headers, k)
	}
	for k := range req.Params {
		delete(req.Params, k)
	}

	requestPool.Put(req)
}

// NewRequest creates a new Request (deprecated for internal routing, use AcquireRequest).
func NewRequest() *Request {
	return AcquireRequest()
}

// Validate checks if the parsed request conforms to protocol requirements.
func (r *Request) Validate() error {
	// HTTP/1.1 requires a Host header (RFC 2616, Section 14.23)
	if r.Version == "HTTP/1.1" {
		if _, ok := r.Headers["host"]; !ok {
			return ErrMissingHostHeader
		}
	}

	// Prevent HTTP Request Smuggling (CL-TE)
	if _, hasCL := r.Headers["content-length"]; hasCL {
		if _, hasTE := r.Headers["transfer-encoding"]; hasTE {
			return ErrConflictingHeaders
		}
	}

	return nil
}

// WantsKeepAlive determines if the client wants to maintain a persistent connection.
func (r *Request) WantsKeepAlive() bool {
	connHeader := r.Headers["connection"]
	
	if r.Version == "HTTP/1.1" {
		// HTTP/1.1 is keep-alive by default, unless "close" is specified.
		return connHeader != "close"
	}
	
	// HTTP/1.0 is close by default, unless "keep-alive" is explicitly specified.
	return connHeader == "keep-alive"
}

// WriteTo serializes the Request object into a raw HTTP byte stream and writes it to w.
// This is primarily used by the reverse proxy to forward requests to backend servers.
func (r *Request) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64
	var buf bytes.Buffer

	// Request-Line
	buf.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, r.Path, r.Version))

	// Headers
	for k, v := range r.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	// Empty line signifying end of headers
	buf.WriteString("\r\n")

	n, err := w.Write(buf.Bytes())
	totalWritten += int64(n)
	if err != nil {
		return totalWritten, err
	}

	// Body
	if len(r.Body) > 0 {
		n, err = w.Write(r.Body)
		totalWritten += int64(n)
		if err != nil {
			return totalWritten, err
		}
	}

	return totalWritten, nil
}

