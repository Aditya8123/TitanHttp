package http

import "errors"

var (
	// ErrMalformedRequest is returned when the HTTP request cannot be parsed.
	ErrMalformedRequest = errors.New("malformed HTTP request")

	// ErrInvalidMethod is returned when the HTTP method is not recognized or missing.
	ErrInvalidMethod = errors.New("invalid HTTP method")

	// ErrInvalidURI is returned when the requested URI is invalid.
	ErrInvalidURI = errors.New("invalid URI")

	// ErrInvalidVersion is returned when the HTTP version is unsupported or malformed.
	ErrInvalidVersion = errors.New("invalid HTTP version")

	// ErrMalformedHeader is returned when an HTTP header is improperly formatted.
	ErrMalformedHeader = errors.New("malformed HTTP header")

	// ErrInvalidContentLength is returned when the Content-Length header is non-numeric or negative.
	ErrInvalidContentLength = errors.New("invalid Content-Length header")

	// ErrBodyTooLarge is returned when the request body exceeds the maximum allowed size.
	ErrBodyTooLarge = errors.New("request body too large")

	// ErrMissingHostHeader is returned when an HTTP/1.1 request lacks the mandatory Host header.
	ErrMissingHostHeader = errors.New("HTTP/1.1 requests must include a Host header")
)
