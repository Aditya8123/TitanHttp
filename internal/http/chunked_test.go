package http

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestChunkedReader_Basic(t *testing.T) {
	input := []byte("4\r\nWiki\r\n5\r\npedia\r\nd\r\n in\r\n\r\nchunks\r\n0\r\n\r\n")
	r := bufio.NewReader(bytes.NewReader(input))
	cr := NewChunkedReader(r)
	defer cr.Close()

	output, err := io.ReadAll(cr)
	if err != nil {
		t.Fatalf("Failed to read chunked data: %v", err)
	}

	expected := "Wikipedia in\r\n\r\nchunks"
	if string(output) != expected {
		t.Errorf("Expected %q, got %q", expected, string(output))
	}
}

func TestChunkedReader_MalformedHeader(t *testing.T) {
	// Missing CRLF in chunk header
	input := []byte("4Wiki\r\n0\r\n\r\n")
	r := bufio.NewReader(bytes.NewReader(input))
	cr := NewChunkedReader(r)
	defer cr.Close()

	_, err := io.ReadAll(cr)
	if err == nil {
		t.Fatal("Expected error for malformed chunk header")
	}
}

func TestChunkedReader_InvalidSize(t *testing.T) {
	// Non-hex characters in size
	input := []byte("z\r\nWiki\r\n0\r\n\r\n")
	r := bufio.NewReader(bytes.NewReader(input))
	cr := NewChunkedReader(r)
	defer cr.Close()

	_, err := io.ReadAll(cr)
	if err != ErrInvalidChunkSize {
		t.Errorf("Expected ErrInvalidChunkSize, got %v", err)
	}
}

func TestChunkedReader_ChunkTooLarge(t *testing.T) {
	// Size exceeds limit (MaxChunkSize is 10MB)
	// Let's request 11MB chunk size: 11 * 1024 * 1024 bytes = 11534336 bytes = 0xB00000
	input := []byte("B00000\r\nWiki\r\n0\r\n\r\n")
	r := bufio.NewReader(bytes.NewReader(input))
	cr := NewChunkedReader(r)
	defer cr.Close()

	_, err := io.ReadAll(cr)
	if err != ErrChunkTooLarge {
		t.Errorf("Expected ErrChunkTooLarge, got %v", err)
	}
}

func TestChunkedReader_UnexpectedEOF(t *testing.T) {
	// Chunk claims to be 5 bytes but stream ends early
	input := []byte("5\r\nWiki")
	r := bufio.NewReader(bytes.NewReader(input))
	cr := NewChunkedReader(r)
	defer cr.Close()

	_, err := io.ReadAll(cr)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("Expected io.ErrUnexpectedEOF, got %v", err)
	}
}
