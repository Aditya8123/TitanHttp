package http

import (
	"bytes"
	"context"
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

	// Headers stores the parsed HTTP headers using a zero-allocation slice-backed structure.
	Headers Header

	// RawBody contains the fast-path byte slice of the request payload, if it was less than 64KB.
	// This enables zero-allocation access for small JSON and text payloads.
	RawBody []byte

	// Body provides standard io.ReadCloser streaming access to the payload.
	// This supports large payloads and chunked transfer encoding.
	Body io.ReadCloser

	// Params stores dynamic path parameters extracted by the router (e.g., /users/:id).
	Params map[string]string

	// RemoteAddr is the network address of the client that sent the request.
	RemoteAddr string

	// Scheme is the protocol scheme (e.g., "http" or "https").
	Scheme string

	// ctx is the context for this request, used for cancellation signaling.
	ctx context.Context
}

// Context returns the request's context.
func (r *Request) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}
	return context.Background()
}

// WithContext returns a shallow copy of r with its context changed to ctx.
func (r *Request) WithContext(ctx context.Context) *Request {
	if ctx == nil {
		panic("nil context")
	}
	r.ctx = ctx
	return r
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &Request{
			Params: make(map[string]string),
		}
	},
}

// AcquireRequest fetches a clean Request object from the pool.
func AcquireRequest() *Request {
	req := requestPool.Get().(*Request)
	// Reset maps that we re-use
	return req
}

// ReleaseRequest cleans up the Request object and returns it to the pool.
func ReleaseRequest(req *Request) {
	req.Method = ""
	req.Path = ""
	req.Version = ""
	req.RawBody = nil
	req.Body = nil
	req.RemoteAddr = ""
	req.Scheme = ""
	req.ctx = nil

	req.Headers.Reset()

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
		if r.Headers.Get("host") == "" {
			return ErrMissingHostHeader
		}
	}

	// Prevent HTTP Request Smuggling (CL-TE)
	if r.Headers.Get("content-length") != "" {
		if r.Headers.Get("transfer-encoding") != "" {
			return ErrConflictingHeaders
		}
	}

	return nil
}

// WantsKeepAlive determines if the client wants to maintain a persistent connection.
func (r *Request) WantsKeepAlive() bool {
	connHeader := r.Headers.Get("connection")

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
	for _, entry := range r.Headers.Entries() {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", entry.Key(), entry.Value()))
	}

	// Empty line signifying end of headers
	buf.WriteString("\r\n")

	n, err := w.Write(buf.Bytes())
	totalWritten += int64(n)
	if err != nil {
		return totalWritten, err
	}

	// Body
	if len(r.RawBody) > 0 {
		n, err = w.Write(r.RawBody)
		totalWritten += int64(n)
		if err != nil {
			return totalWritten, err
		}
	} else if r.Body != nil {
		written, err := io.Copy(w, r.Body)
		totalWritten += written
		if err != nil {
			return totalWritten, err
		}
	}

	return totalWritten, nil
}
