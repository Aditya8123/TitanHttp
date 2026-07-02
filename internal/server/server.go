// Package server will host the raw TCP listener and connection-accept loop
// that underpin TitanHTTP.
//
// It handles socket connections, manages worker goroutines, and forms the
// core entrypoint for incoming network traffic.
package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

type Server struct {
	addr       string
	mu         sync.Mutex // Protects listener
	listener   net.Listener
	router     *router.Router
	workerPool *WorkerPool
	metrics    *Metrics
	done       chan struct{}
}

func NewServer(addr string) *Server {
	s := &Server{
		addr:    addr,
		router:  router.NewRouter(),
		metrics: &Metrics{},
		done:    make(chan struct{}),
	}
	// Initialize the worker pool with 100 workers and a queue size of 1024
	s.workerPool = NewWorkerPool(100, 1024, s.handleConnection)
	return s
}

// Router returns the underlying router for registering endpoints.
func (s *Server) Router() *router.Router {
	return s.router
}

// Addr returns the address the server is listening on. It is safe for concurrent use.
func (s *Server) Addr() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return s.addr
}

// Start opens a TCP socket on the configured address and begins listening for connections.
// It blocks indefinitely to keep the process alive (for this subtask).
func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to bind to address %s: %w", s.addr, err)
	}

	s.mu.Lock()
	s.listener = l
	s.mu.Unlock()

	fmt.Printf("TitanHTTP Server successfully started.\n")
	fmt.Printf("Listening on %s...\n", s.Addr())

	// Start the worker pool before accepting connections
	s.workerPool.Start()

	// Start the background metrics reporter
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.metrics.Report()
			case <-s.done:
				return
			}
		}
	}()

	// The main connection accept loop.
	// This blocks waiting for new connections.
	for {
		s.mu.Lock()
		l := s.listener
		s.mu.Unlock()
		
		if l == nil {
			return nil
		}

		conn, err := l.Accept()
		if err != nil {
			select {
			case <-s.done:
				return nil // server is shutting down
			default:
				if errors.Is(err, net.ErrClosed) {
					return nil // Listener closed during shutdown
				}
				fmt.Printf("Error accepting connection: %v\n", err)
				continue
			}
		}

		// Submit the connection to the worker pool instead of spawning a new goroutine
		s.workerPool.Submit(conn)
	}
}

// handleConnection processes an individual client connection.
func (s *Server) handleConnection(conn net.Conn) {
	// Defers are executed when the surrounding function returns (LIFO order).
	// We close the connection last.
	defer conn.Close()

	var totalConnBytes int64
	s.metrics.RequestStarted()
	defer func() {
		s.metrics.RequestFinished(totalConnBytes)
	}()

	// Recover from panics to isolate connection failures and prevent server crashes.
	defer func() {
		if r := recover(); r != nil {
			s.metrics.PanicRecovered()
			fmt.Printf("Critical: Connection panic recovered: %v\n", r)
			resp := http.NewResponse500()
			n, _ := conn.Write(resp.Bytes())
			totalConnBytes += int64(n)
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
				// Silently handle normal disconnects to avoid log noise in concurrent environments.
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Silently handle read timeouts to avoid log noise.
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

		n, err := conn.Write(resp.Bytes())
		totalConnBytes += int64(n)
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			return
		}

		// Close connection if the client requested it
		if req.Headers["connection"] == "close" {
			return
		}
	}
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	// Signal all loops to stop
	close(s.done)
	
	// Close listener to unblock Accept()
	s.mu.Lock()
	if s.listener != nil {
		s.listener.Close()
	}
	s.mu.Unlock()

	// Wait for worker pool to drain, racing against ctx.Done()
	stopCh := make(chan struct{})
	go func() {
		s.workerPool.Stop()
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCh:
		return nil
	}
}
