// Package server will host the raw TCP listener and connection-accept loop
// that underpin TitanHTTP.
//
// It handles socket connections, manages worker goroutines, and forms the
// core entrypoint for incoming network traffic.
package server

import (
	"fmt"
	"io"
	"net"
	"time"
)

// Server represents the TitanHTTP core server.
type Server struct {
	addr     string
	listener net.Listener
}

// NewServer initializes a new TitanHTTP Server configured to listen on the given address.
func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

// Start opens a TCP socket on the configured address and begins listening for connections.
// It blocks indefinitely to keep the process alive (for this subtask).
func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to bind to address %s: %w", s.addr, err)
	}

	s.listener = l
	fmt.Printf("TitanHTTP Server successfully started.\n")
	fmt.Printf("Listening on %s...\n", s.addr)

	// The main connection accept loop.
	// This blocks waiting for new connections.
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		fmt.Printf("Accepted new connection from %s\n", conn.RemoteAddr().String())

		s.handleConnection(conn)
	}
}

// handleConnection processes an individual client connection.
func (s *Server) handleConnection(conn net.Conn) {
	// Defers are executed when the surrounding function returns.
	// This guarantees the socket is closed even if a panic occurs or we return early.
	defer conn.Close()

	buf := make([]byte, 1024)

	// Connection loop: continuously read from the socket until EOF or error.
	for {
		// Set a 5-second timeout for reading to prevent hanging connections.
		err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Printf("Failed to set read deadline: %v\n", err)
			return
		}

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client disconnected (EOF).\n")
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Connection timed out.\n")
			} else {
				fmt.Printf("Error reading from connection: %v\n", err)
			}
			return
		}

		fmt.Printf("--- Received %d bytes ---\n", n)
		fmt.Print(string(buf[:n]))
		fmt.Printf("-------------------------\n")

		// Write a simple HTTP response back to the client.
		// In a true HTTP loop, this would happen after parsing a complete request.
		response := "HTTP/1.1 200 OK\r\nContent-Length: 20\r\n\r\nHello from TitanHTTP"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			return
		}
	}
}
