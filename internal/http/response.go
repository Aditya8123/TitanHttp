package http

import (
	"bytes"
	"io"
	"strconv"
	"sync"
)

// StatusCode represents an HTTP response status code
type StatusCode int

const (
	// 2xx Success
	StatusOK             StatusCode = 200
	StatusCreated        StatusCode = 201
	StatusPartialContent StatusCode = 206

	// 4xx Client Errors
	StatusBadRequest       StatusCode = 400
	StatusUnauthorized     StatusCode = 401
	StatusForbidden        StatusCode = 403
	StatusNotFound         StatusCode = 404
	StatusMethodNotAllowed StatusCode = 405
	StatusURITooLong       StatusCode = 414
	StatusTooManyRequests  StatusCode = 429
	StatusRequestHeaderFieldsTooLarge StatusCode = 431

	// 5xx Server Errors
	StatusInternalServerError StatusCode = 500
	StatusNotImplemented      StatusCode = 501
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

	// Headers stores the parsed HTTP headers using a zero-allocation slice-backed structure.
	Headers Header

	// Body contains the payload of the response, if any (loaded in memory).
	Body []byte

	// Stream contains a reader for the response payload, enabling chunked or streaming transfer.
	Stream io.Reader
}

var responsePool = sync.Pool{
	New: func() interface{} {
		return &Response{
			Version: "HTTP/1.1",
		}
	},
}

// AcquireResponse fetches a clean Response object from the pool.
func AcquireResponse() *Response {
	resp := responsePool.Get().(*Response)
	resp.Version = "HTTP/1.1"
	return resp
}

// ReleaseResponse cleans up the Response object and returns it to the pool.
func ReleaseResponse(resp *Response) {
	resp.StatusCode = 0
	resp.StatusText = ""
	resp.Body = nil
	resp.Stream = nil
	resp.Headers.Reset()
	responsePool.Put(resp)
}

// NewResponse creates a new Response (deprecated for internal routing, use AcquireResponse).
func NewResponse() *Response {
	return AcquireResponse()
}

var statusText = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusCreated:             "Created",
	StatusPartialContent:      "Partial Content",
	StatusBadRequest:          "Bad Request",
	StatusUnauthorized:        "Unauthorized",
	StatusForbidden:           "Forbidden",
	StatusNotFound:            "Not Found",
	StatusMethodNotAllowed:    "Method Not Allowed",
	StatusURITooLong:          "URI Too Long",
	StatusTooManyRequests:     "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge: "Request Header Fields Too Large",
	StatusInternalServerError: "Internal Server Error",
	StatusNotImplemented:      "Not Implemented",
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
	hasContentLength := r.Headers.Get("Content-Length") != ""
	isChunked := r.Stream != nil && !hasContentLength

	if len(r.Body) > 0 && !hasContentLength && r.Stream == nil {
		r.Headers.Set("Content-Length", strconv.Itoa(len(r.Body)))
	} else if len(r.Body) == 0 && r.Stream == nil && !hasContentLength {
		r.Headers.Set("Content-Length", "0")
	}

	if isChunked {
		r.Headers.Set("Transfer-Encoding", "chunked")
	}

	var headerBuf bytes.Buffer
	// Status-Line: HTTP/1.1 200 OK\r\n
	headerBuf.WriteString(r.Version)
	headerBuf.WriteByte(' ')
	headerBuf.WriteString(strconv.Itoa(int(r.StatusCode)))
	headerBuf.WriteByte(' ')
	headerBuf.WriteString(text)
	headerBuf.WriteString("\r\n")

	// Headers
	for _, entry := range r.Headers.Entries() {
		headerBuf.WriteString(entry.Key())
		headerBuf.WriteString(": ")
		headerBuf.WriteString(entry.Value())
		headerBuf.WriteString("\r\n")
	}

	// Empty line signifying end of headers
	headerBuf.WriteString("\r\n")

	// Write headers to the wire first
	n, err := w.Write(headerBuf.Bytes())
	totalWritten += int64(n)
	if err != nil {
		return totalWritten, err
	}

	// Write the body directly to the wire, avoiding bytes.Buffer copy allocation
	if r.Stream == nil && len(r.Body) > 0 {
		nBody, errBody := w.Write(r.Body)
		totalWritten += int64(nBody)
		if errBody != nil {
			return totalWritten, errBody
		}
	}

	// Stream payload
	if r.Stream != nil {
		if isChunked {
			buf := make([]byte, 8192)
			for {
				readBytes, err := r.Stream.Read(buf)
				if readBytes > 0 {
					// chunkHeader
					headerBuf.Reset()
					headerBuf.WriteString(strconv.FormatInt(int64(readBytes), 16))
					headerBuf.WriteString("\r\n")
					n, wErr := w.Write(headerBuf.Bytes())
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
			n2, err2 := io.Copy(w, r.Stream)
			totalWritten += n2
			if err2 != nil {
				return totalWritten, err2
			}
		}

		// If Stream is an io.Closer (like os.File), close it
		if closer, ok := r.Stream.(io.Closer); ok {
			closer.Close()
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
	resp.StatusCode = StatusTooManyRequests
	resp.Body = []byte("429 Too Many Requests\n")
	return resp
}

// NewResponse414 generates a standard 414 URI Too Long response.
func NewResponse414() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusURITooLong
	resp.Body = []byte("414 URI Too Long\n")
	return resp
}

// NewResponse431 generates a standard 431 Request Header Fields Too Large response.
func NewResponse431() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusRequestHeaderFieldsTooLarge
	resp.Body = []byte("431 Request Header Fields Too Large\n")
	return resp
}

// NewResponse501 generates a standard 501 Not Implemented response.
func NewResponse501() *Response {
	resp := NewResponse()
	resp.StatusCode = StatusNotImplemented
	resp.Body = []byte("501 Not Implemented\n")
	return resp
}
