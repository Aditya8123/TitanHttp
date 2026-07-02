// Package server will host the raw TCP listener and connection-accept loop
// that underpin TitanHTTP.
//
// It handles socket connections, manages worker goroutines, and forms the
// core entrypoint for incoming network traffic.
package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// Server represents the TitanHTTP core server.
type Server struct {
	addr     string
	listener net.Listener
	router   *router.Router
}

// NewServer initializes a new TitanHTTP Server configured to listen on the given address.
func NewServer(addr string) *Server {
	return &Server{
		addr:   addr,
		router: router.NewRouter(),
	}
}

// Router returns the underlying router for registering endpoints.
func (s *Server) Router() *router.Router {
	return s.router
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

		go s.handleConnection(conn)
	}
}

// handleConnection processes an individual client connection.
func (s *Server) handleConnection(conn net.Conn) {
	// Defers are executed when the surrounding function returns (LIFO order).
	// We close the connection last.
	defer conn.Close()

	// Recover from panics to isolate connection failures and prevent server crashes.
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Critical: Connection panic recovered: %v\n", r)
			resp := http.NewResponse500()
			conn.Write(resp.Bytes())
		}
	}()

	reader := bufio.NewReader(conn)

	// Connection loop: continuously read from the socket until EOF or error.
	for {
		// Set a 5-second timeout for reading to prevent hanging connections.
		err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Printf("Failed to set read deadline: %v\n", err)
			return
		}

		req, err := http.ParseRequest(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client disconnected (EOF).\n")
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Connection timed out.\n")
			} else {
				fmt.Printf("Error parsing request: %v\n", err)
				// Send a 400 Bad Request on parser errors
				resp := http.NewResponse400()
				if _, err := conn.Write(resp.Bytes()); err != nil {
					fmt.Printf("Failed to write parser error response: %v\n", err)
				}
			}
			return
		}

		err = req.Validate()
		if err != nil {
			fmt.Printf("Request validation failed: %v\n", err)
			resp := http.NewResponse400()
			if _, err := conn.Write(resp.Bytes()); err != nil {
				fmt.Printf("Failed to write validation error response: %v\n", err)
			}
			return
		}

		// Dispatch request to the router
		resp := s.router.ServeHTTP(req)

		_, err = conn.Write(resp.Bytes())
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			return
		}
	}
}
