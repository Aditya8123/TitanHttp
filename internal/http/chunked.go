package http

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

var (
	ErrInvalidChunkSize = errors.New("invalid chunk size")
	ErrChunkTooLarge    = errors.New("chunk size exceeds limits")
)

const MaxChunkSize = 10 * 1024 * 1024 // 10MB limit per chunk

// chunkedReader implements io.ReadCloser for HTTP/1.1 chunked transfer encoding.
type chunkedReader struct {
	r         *bufio.Reader
	chunkLeft int
	done      bool
}

// NewChunkedReader creates a new io.ReadCloser for chunked decoding.
func NewChunkedReader(r *bufio.Reader) io.ReadCloser {
	return &chunkedReader{
		r: r,
	}
}

func (cr *chunkedReader) Read(p []byte) (n int, err error) {
	if cr.done {
		return 0, io.EOF
	}

	if cr.chunkLeft == 0 {
		err = cr.readChunkHeader()
		if err != nil {
			return 0, err
		}
		if cr.chunkLeft == 0 {
			cr.done = true
			return 0, io.EOF
		}
	}

	toRead := cr.chunkLeft
	if len(p) < toRead {
		toRead = len(p)
	}

	n, err = cr.r.Read(p[:toRead])
	cr.chunkLeft -= n

	if cr.chunkLeft == 0 {
		// Read the trailing CRLF after the chunk data
		errCRLF := cr.readCRLF()
		if errCRLF != nil {
			return n, errCRLF
		}
	}

	if err == io.EOF && cr.chunkLeft > 0 {
		return n, io.ErrUnexpectedEOF
	}

	// Mask inner EOFs until we hit the zero-chunk
	if err == io.EOF {
		err = nil
	}

	return n, err
}

func (cr *chunkedReader) readChunkHeader() error {
	line, err := cr.r.ReadSlice('\n')
	if err != nil {
		return err
	}

	if len(line) < 2 || line[len(line)-2] != '\r' || line[len(line)-1] != '\n' {
		return errors.New("malformed chunk header")
	}

	// Remove CRLF
	line = line[:len(line)-2]

	// Ignore chunk extensions
	idx := bytes.IndexByte(line, ';')
	if idx != -1 {
		line = line[:idx]
	}
	line = bytes.TrimSpace(line)

	size, err := strconv.ParseInt(string(line), 16, 64)
	if err != nil || size < 0 {
		return ErrInvalidChunkSize
	}

	if size > MaxChunkSize {
		return ErrChunkTooLarge
	}

	cr.chunkLeft = int(size)
	return nil
}

func (cr *chunkedReader) readCRLF() error {
	b := make([]byte, 2)
	_, err := io.ReadFull(cr.r, b)
	if err != nil {
		return err
	}
	if b[0] != '\r' || b[1] != '\n' {
		return errors.New("missing CRLF after chunk data")
	}
	return nil
}

func (cr *chunkedReader) Close() error {
	cr.done = true
	return nil
}
