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
	err = parseHeaders(reader, req)
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

// parseHeaders reads the headers and adds them to the Request struct.
func parseHeaders(reader *bufio.Reader, req *Request) error {
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

		req.Headers[key] = value
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
