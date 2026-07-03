package http

import (
	"bytes"
	"fmt"
	"io"
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
	StatusTooManyRequests  StatusCode = 429

	// 5xx Server Errors
	StatusInternalServerError StatusCode = 500
)

// Response represents an HTTP response to be sent to a client.
// In HTTP/1.1, a response consists of a Status-Line, Headers, and an optional Body or Stream.
type Response struct {
	// Version is the HTTP protocol version (e.g., HTTP/1.1).
	Version string

	// StatusCode is the 3-digit integer result code of the attempt to understand and satisfy the request.
	StatusCode StatusCode

	// StatusText is the short textual description of the Status-Code.
	StatusText string

	// Headers stores the key-value pairs of the HTTP headers.
	Headers map[string]string

	// Body contains the payload of the response, if any (loaded in memory).
	Body []byte

	// Stream contains a reader for the response payload, enabling chunked or streaming transfer.
	Stream io.Reader
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
	StatusTooManyRequests:     "Too Many Requests",
	StatusInternalServerError: "Internal Server Error",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code StatusCode) string {
	return statusText[code]
}

// WriteTo serializes the Response object directly into an io.Writer.
// It supports streaming and automatic chunked transfer encoding for unknown lengths.
func (r *Response) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	text := r.StatusText
	if text == "" {
		text = statusText[r.StatusCode]
	}

	// Calculate body size or chunked mode
	_, hasContentLength := r.Headers["Content-Length"]
	isChunked := r.Stream != nil && !hasContentLength

	if len(r.Body) > 0 && !hasContentLength && r.Stream == nil {
		r.Headers["Content-Length"] = strconv.Itoa(len(r.Body))
	} else if len(r.Body) == 0 && r.Stream == nil && !hasContentLength {
		r.Headers["Content-Length"] = "0"
	}

	if isChunked {
		r.Headers["Transfer-Encoding"] = "chunked"
	}

	var headerBuf bytes.Buffer
	// Status-Line
	headerBuf.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Version, r.StatusCode, text))

	// Headers
	for k, v := range r.Headers {
		headerBuf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	// Empty line signifying end of headers
	headerBuf.WriteString("\r\n")

	n, err := w.Write(headerBuf.Bytes())
	totalWritten += int64(n)
	if err != nil {
		return totalWritten, err
	}

	// Body payload
	if r.Stream != nil {
		if isChunked {
			buf := make([]byte, 8192)
			for {
				readBytes, err := r.Stream.Read(buf)
				if readBytes > 0 {
					chunkHeader := fmt.Sprintf("%x\r\n", readBytes)
					n, wErr := w.Write([]byte(chunkHeader))
					totalWritten += int64(n)
					if wErr != nil {
						return totalWritten, wErr
					}

					n, wErr = w.Write(buf[:readBytes])
					totalWritten += int64(n)
					if wErr != nil {
						return totalWritten, wErr
					}

					n, wErr = w.Write([]byte("\r\n"))
					totalWritten += int64(n)
					if wErr != nil {
						return totalWritten, wErr
					}
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					return totalWritten, err
				}
			}
			// Write the final zero-length chunk
			n, wErr := w.Write([]byte("0\r\n\r\n"))
			totalWritten += int64(n)
			if wErr != nil {
				return totalWritten, wErr
			}
		} else {
			written, err := io.Copy(w, r.Stream)
			totalWritten += written
			if err != nil {
				return totalWritten, err
			}
		}
		
		// If Stream is an io.Closer (like os.File), close it
		if closer, ok := r.Stream.(io.Closer); ok {
			closer.Close()
		}
	} else if len(r.Body) > 0 {
		n, err := w.Write(r.Body)
		totalWritten += int64(n)
		if err != nil {
			return totalWritten, err
		}
	}

	return totalWritten, nil
}

// Bytes serializes the Response object into a raw HTTP byte stream.
// Warning: This buffers the entire response in memory. Use WriteTo for large payloads.
func (r *Response) Bytes() []byte {
	var b bytes.Buffer
	r.WriteTo(&b)
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

// NewResponse503 returns a pre-configured 503 Service Unavailable response.
func NewResponse503() *Response {
	resp := NewResponse()
	resp.StatusCode = 503
	resp.StatusText = "Service Unavailable"
	return resp
}

// NewResponse429 returns a pre-configured 429 Too Many Requests response.
func NewResponse429() *Response {
	resp := NewResponse()
	resp.StatusCode = 429
	resp.StatusText = "Too Many Requests"
	return resp
}
