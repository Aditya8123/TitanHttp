package http

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
}

// NewRequest creates a new Request with initialized maps.
func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
		Params:  make(map[string]string),
	}
}

// Validate checks if the parsed request conforms to protocol requirements.
func (r *Request) Validate() error {
	// HTTP/1.1 requires a Host header (RFC 2616, Section 14.23)
	if r.Version == "HTTP/1.1" {
		if _, ok := r.Headers["host"]; !ok {
			return ErrMissingHostHeader
		}
	}

	return nil
}
