package http

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

const MaxBodySize = 10 * 1024 * 1024 // 10 MB limit

// ParseRequest reads from a bufio.Reader and constructs an HTTP Request.
func ParseRequest(reader *bufio.Reader) (*Request, error) {
	req := NewRequest()

	// 1. Parse the request line
	err := parseRequestLine(reader, req)
	if err != nil {
		return nil, err
	}

	// 2. Parse headers
	err = parseHeaders(reader, req.Headers)
	if err != nil {
		return nil, err
	}

	// 3. Parse body
	err = parseBody(reader, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// parseRequestLine reads the first line and extracts Method, Path, and Version.
func parseRequestLine(reader *bufio.Reader, req *Request) error {
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// HTTP lines must end with CRLF (\r\n)
	if !strings.HasSuffix(line, "\r\n") {
		return ErrMalformedRequest
	}

	// Strip the CRLF
	line = line[:len(line)-2]

	// Split by space. Request-Line = Method SP Request-URI SP HTTP-Version CRLF
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return ErrMalformedRequest
	}

	req.Method = Method(parts[0])
	req.Path = parts[1]
	req.Version = parts[2]

	// Basic validation
	if req.Method == "" {
		return ErrInvalidMethod
	}

	if req.Path == "" || !strings.HasPrefix(req.Path, "/") {
		return ErrInvalidURI
	}

	if !strings.HasPrefix(req.Version, "HTTP/") {
		return ErrInvalidVersion
	}

	return nil
}

// parseHeaders reads the headers and adds them to the provided headers map.
func parseHeaders(reader *bufio.Reader, headers map[string]string) error {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if !strings.HasSuffix(line, "\r\n") {
			return ErrMalformedHeader
		}

		// Strip CRLF
		line = line[:len(line)-2]

		// An empty line signifies the end of headers
		if line == "" {
			break
		}

		// Split by the first colon
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return ErrMalformedHeader
		}

		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return ErrMalformedHeader
		}

		headers[key] = value
	}

	return nil
}

// parseBody reads the request body based on Content-Length.
func parseBody(reader *bufio.Reader, req *Request) error {
	contentLengthStr, ok := req.Headers["content-length"]
	if !ok {
		// No body to parse
		return nil
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil || contentLength < 0 {
		return ErrInvalidContentLength
	}

	if contentLength == 0 {
		return nil
	}

	if contentLength > MaxBodySize {
		return ErrBodyTooLarge
	}

	req.Body = make([]byte, contentLength)
	_, err = io.ReadFull(reader, req.Body)
	if err != nil {
		return err
	}

	return nil
}

// ParseResponse reads from a bufio.Reader and constructs an HTTP Response.
// It parses the Status-Line and Headers, and attaches the reader to Response.Stream
// for zero-allocation body streaming.
func ParseResponse(reader *bufio.Reader) (*Response, error) {
	resp := NewResponse()

	// 1. Parse the status line
	err := parseStatusLine(reader, resp)
	if err != nil {
		return nil, err
	}

	// 2. Parse headers
	err = parseHeaders(reader, resp.Headers)
	if err != nil {
		return nil, err
	}

	// 3. For the response body, we leave it to the caller to stream via resp.Stream
	// But we need to ensure it's available. We can just attach the reader.
	resp.Stream = reader

	return resp, nil
}

// parseStatusLine reads the first line and extracts Version, StatusCode, and StatusText.
func parseStatusLine(reader *bufio.Reader, resp *Response) error {
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// HTTP lines must end with CRLF (\r\n)
	if !strings.HasSuffix(line, "\r\n") {
		return ErrMalformedRequest // We can reuse this or define ErrMalformedResponse
	}

	// Strip the CRLF
	line = line[:len(line)-2]

	// Split by space. Status-Line = HTTP-Version SP Status-Code SP Reason-Phrase CRLF
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 2 {
		return ErrMalformedRequest
	}

	resp.Version = parts[0]
	
	code, err := strconv.Atoi(parts[1])
	if err != nil {
		return ErrMalformedRequest
	}
	resp.StatusCode = StatusCode(code)

	if len(parts) == 3 {
		resp.StatusText = parts[2]
	}

	if !strings.HasPrefix(resp.Version, "HTTP/") {
		return ErrInvalidVersion
	}

	return nil
}
