package http

import (
	"bytes"
	"fmt"
	"strconv"
)

// StatusCode represents an HTTP response status code
type StatusCode int

const (
	// 2xx Success
	StatusOK      StatusCode = 200
	StatusCreated StatusCode = 201

	// 4xx Client Errors
	StatusBadRequest       StatusCode = 400
	StatusUnauthorized     StatusCode = 401
	StatusForbidden        StatusCode = 403
	StatusNotFound         StatusCode = 404
	StatusMethodNotAllowed StatusCode = 405

	// 5xx Server Errors
	StatusInternalServerError StatusCode = 500
)

// Response represents an HTTP response to be sent to a client.
// In HTTP/1.1, a response consists of a Status-Line, Headers, and an optional Body.
type Response struct {
	// Version is the HTTP protocol version (e.g., HTTP/1.1).
	Version string

	// StatusCode is the 3-digit integer result code of the attempt to understand and satisfy the request.
	StatusCode StatusCode

	// StatusText is the short textual description of the Status-Code.
	StatusText string

	// Headers stores the key-value pairs of the HTTP headers.
	Headers map[string]string

	// Body contains the payload of the response, if any.
	Body []byte
}

// NewResponse creates a new Response with initialized maps and a default HTTP/1.1 version.
func NewResponse() *Response {
	return &Response{
		Version: "HTTP/1.1",
		Headers: make(map[string]string),
	}
}

var statusText = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusCreated:             "Created",
	StatusBadRequest:          "Bad Request",
	StatusUnauthorized:        "Unauthorized",
	StatusForbidden:           "Forbidden",
	StatusNotFound:            "Not Found",
	StatusMethodNotAllowed:    "Method Not Allowed",
	StatusInternalServerError: "Internal Server Error",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code StatusCode) string {
	return statusText[code]
}

// Bytes serializes the Response object into a raw HTTP byte stream.
func (r *Response) Bytes() []byte {
	var b bytes.Buffer

	text := r.StatusText
	if text == "" {
		text = statusText[r.StatusCode]
	}

	// Status-Line
	b.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Version, r.StatusCode, text))

	// Ensure Content-Length is set correctly based on the body payload
	if len(r.Body) > 0 {
		r.Headers["Content-Length"] = strconv.Itoa(len(r.Body))
	} else {
		r.Headers["Content-Length"] = "0"
	}

	// Headers
	for k, v := range r.Headers {
		b.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	// Empty line signifying end of headers
	b.WriteString("\r\n")

	// Body payload
	if len(r.Body) > 0 {
		b.Write(r.Body)
	}

	return b.Bytes()
}

// NewResponse200 generates a standard 200 OK response.
func NewResponse200() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusOK
	return resp
}

// NewResponse400 generates a standard 400 Bad Request response.
func NewResponse400() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusBadRequest
	resp.Body = []byte("400 Bad Request\n")
	return resp
}

// NewResponse401 generates a standard 401 Unauthorized response.
func NewResponse401() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusUnauthorized
	resp.Body = []byte("401 Unauthorized\n")
	return resp
}

// NewResponse403 generates a standard 403 Forbidden response.
func NewResponse403() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusForbidden
	resp.Body = []byte("403 Forbidden\n")
	return resp
}

// NewResponse404 generates a standard 404 Not Found response.
func NewResponse404() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusNotFound
	resp.Body = []byte("404 Not Found\n")
	return resp
}

// NewResponse405 generates a standard 405 Method Not Allowed response.
func NewResponse405() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusMethodNotAllowed
	resp.Body = []byte("405 Method Not Allowed\n")
	return resp
}

// NewResponse500 generates a standard 500 Internal Server Error response.
func NewResponse500() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusInternalServerError
	resp.Body = []byte("500 Internal Server Error\n")
	return resp
}
