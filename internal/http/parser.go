package http

import (
	"bufio"
	"bytes"
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
		ReleaseRequest(req)
		return nil, err
	}

	// 2. Parse headers
	err = parseHeaders(reader, &req.Headers)
	if err != nil {
		ReleaseRequest(req)
		return nil, err
	}

	// If Expect: 100-continue is set, we defer body parsing to the server layer
	// which will send a 100 Continue response before proceeding.
	if req.Headers.Get("expect") == "100-continue" {
		return req, nil
	}

	// 3. Parse body
	err = ParseBody(reader, req)
	if err != nil {
		ReleaseRequest(req)
		return nil, err
	}

	return req, nil
}

func parseRequestLine(reader *bufio.Reader, req *Request) error {
	line, err := reader.ReadSlice('\n')
	if err != nil {
		if err == bufio.ErrBufferFull {
			return ErrURITooLong
		}
		return err
	}

	// HTTP lines must end with CRLF (\r\n)
	if len(line) < 2 || line[len(line)-2] != '\r' || line[len(line)-1] != '\n' {
		return ErrMalformedRequest
	}

	// Strip the CRLF
	line = line[:len(line)-2]

	// Request-Line = Method SP Request-URI SP HTTP-Version CRLF
	idx1 := bytes.IndexByte(line, ' ')
	if idx1 == -1 {
		return ErrMalformedRequest
	}
	methodBytes := line[:idx1]
	if bytes.Equal(methodBytes, []byte("GET")) {
		req.Method = MethodGet
	} else if bytes.Equal(methodBytes, []byte("POST")) {
		req.Method = MethodPost
	} else if bytes.Equal(methodBytes, []byte("PUT")) {
		req.Method = MethodPut
	} else if bytes.Equal(methodBytes, []byte("DELETE")) {
		req.Method = MethodDelete
	} else if bytes.Equal(methodBytes, []byte("OPTIONS")) {
		req.Method = MethodOptions
	} else if bytes.Equal(methodBytes, []byte("HEAD")) {
		req.Method = MethodHead
	} else if bytes.Equal(methodBytes, []byte("PATCH")) {
		req.Method = MethodPatch
	} else {
		return ErrInvalidMethod
	}

	idx2 := bytes.IndexByte(line[idx1+1:], ' ')
	if idx2 == -1 {
		return ErrMalformedRequest
	}
	idx2 += idx1 + 1

	req.Path = string(line[idx1+1 : idx2])

	versionBytes := line[idx2+1:]
	if bytes.Equal(versionBytes, []byte("HTTP/1.1")) {
		req.Version = "HTTP/1.1"
	} else if bytes.Equal(versionBytes, []byte("HTTP/1.0")) {
		req.Version = "HTTP/1.0"
	} else {
		return ErrInvalidVersion
	}

	// Basic validation
	if req.Path == "" || !strings.HasPrefix(req.Path, "/") {
		return ErrInvalidURI
	}

	return nil
}

const MaxHeadersCount = 100

func parseHeaders(reader *bufio.Reader, headers *Header) error {
	headerCount := 0
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			if err == bufio.ErrBufferFull {
				return ErrHeaderFieldsTooLarge
			}
			return err
		}

		headerCount++
		if headerCount > MaxHeadersCount {
			return ErrHeaderFieldsTooLarge
		}

		if len(line) < 2 || line[len(line)-2] != '\r' || line[len(line)-1] != '\n' {
			return ErrMalformedHeader
		}

		// Strip CRLF
		line = line[:len(line)-2]

		// An empty line signifies the end of headers
		if len(line) == 0 {
			break
		}

		idx := bytes.IndexByte(line, ':')
		if idx == -1 {
			return ErrMalformedHeader
		}

		keyBytes := bytes.TrimSpace(line[:idx])
		// Lowercase the key in-place to avoid strings.ToLower allocations
		for i := 0; i < len(keyBytes); i++ {
			if keyBytes[i] >= 'A' && keyBytes[i] <= 'Z' {
				keyBytes[i] += 'a' - 'A'
			}
		}

		if len(keyBytes) == 0 {
			return ErrMalformedHeader
		}

		valueBytes := bytes.TrimSpace(line[idx+1:])

		var keyStr string
		switch {
		case bytes.Equal(keyBytes, []byte("host")):
			keyStr = "host"
		case bytes.Equal(keyBytes, []byte("user-agent")):
			keyStr = "user-agent"
		case bytes.Equal(keyBytes, []byte("accept")):
			keyStr = "accept"
		case bytes.Equal(keyBytes, []byte("connection")):
			keyStr = "connection"
		case bytes.Equal(keyBytes, []byte("content-length")):
			keyStr = "content-length"
			if headers.Get(keyStr) != "" {
				return ErrDuplicateHeader
			}
		case bytes.Equal(keyBytes, []byte("content-type")):
			keyStr = "content-type"
		case bytes.Equal(keyBytes, []byte("accept-encoding")):
			keyStr = "accept-encoding"
		case bytes.Equal(keyBytes, []byte("authorization")):
			keyStr = "authorization"
		case bytes.Equal(keyBytes, []byte("transfer-encoding")):
			keyStr = "transfer-encoding"
		default:
			keyStr = string(keyBytes)
		}

		headers.Add(keyStr, string(valueBytes))
	}

	if headers.Get("content-length") != "" {
		if headers.Get("transfer-encoding") != "" {
			return ErrConflictingHeaders
		}
	}

	return nil
}

// ParseBody reads the request body based on Content-Length or Transfer-Encoding.
func ParseBody(reader *bufio.Reader, req *Request) error {
	isChunked := false
	if strings.Contains(strings.ToLower(req.Headers.Get("transfer-encoding")), "chunked") {
		isChunked = true
	}

	if isChunked {
		req.Body = NewChunkedReader(reader)
		return nil
	}

	contentLengthStr := req.Headers.Get("content-length")
	if contentLengthStr == "" {
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

	// Fast Path: Load into RAM if payload is <= 64KB
	if contentLength <= 64*1024 {
		req.RawBody = make([]byte, contentLength)
		_, err = io.ReadFull(reader, req.RawBody)
		if err != nil {
			return err
		}
		// Wrap RawBody in a reader so standard middleware can still read from Body
		req.Body = io.NopCloser(bytes.NewReader(req.RawBody))
		return nil
	}

	// Slow Path (Streaming): Expose the underlying network stream for large payloads
	req.Body = io.NopCloser(io.LimitReader(reader, int64(contentLength)))
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
		ReleaseResponse(resp)
		return nil, err
	}

	// 2. Parse headers
	err = parseHeaders(reader, &resp.Headers)
	if err != nil {
		ReleaseResponse(resp)
		return nil, err
	}

	// 3. For the response body, we leave it to the caller to stream via resp.Stream
	// But we need to ensure it's available. We can just attach the reader.
	resp.Stream = reader

	return resp, nil
}

func parseStatusLine(reader *bufio.Reader, resp *Response) error {
	line, err := reader.ReadSlice('\n')
	if err != nil {
		return err
	}

	// HTTP lines must end with CRLF (\r\n)
	if len(line) < 2 || line[len(line)-2] != '\r' || line[len(line)-1] != '\n' {
		return ErrMalformedRequest // We can reuse this or define ErrMalformedResponse
	}

	// Strip the CRLF
	line = line[:len(line)-2]

	// Status-Line = HTTP-Version SP Status-Code SP Reason-Phrase CRLF
	idx1 := bytes.IndexByte(line, ' ')
	if idx1 == -1 {
		return ErrMalformedRequest
	}
	resp.Version = string(line[:idx1])

	idx2 := bytes.IndexByte(line[idx1+1:], ' ')
	var codeStr string
	if idx2 == -1 {
		codeStr = string(line[idx1+1:])
	} else {
		idx2 += idx1 + 1
		codeStr = string(line[idx1+1 : idx2])
		resp.StatusText = string(line[idx2+1:])
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		return ErrMalformedRequest
	}
	resp.StatusCode = StatusCode(code)

	if !strings.HasPrefix(resp.Version, "HTTP/") {
		return ErrInvalidVersion
	}

	return nil
}
