// Package server will host the raw TCP listener and connection-accept loop
// that underpin TitanHTTP.
//
// It handles socket connections, manages worker goroutines, and forms the
// core entrypoint for incoming network traffic.
package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"runtime"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

const (
	DefaultIdleTimeout = 30 * time.Second
)

var readerPool = sync.Pool{
	New: func() interface{} {
		return bufio.NewReaderSize(nil, 4096)
	},
}

type WorkerShard struct {
	queue chan net.Conn
}

type Server struct {
	addr         string
	mu           sync.Mutex // Protects listener
	listener     net.Listener
	router       *router.Router
	metrics      *Metrics
	done         chan struct{}
	activeConnWg sync.WaitGroup
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	MaxWorkers   int
	shards       []WorkerShard
	roundRobin   uint64
}

func NewServer(addr string) *Server {
	s := &Server{
		addr:        addr,
		router:      router.NewRouter(),
		metrics:     &Metrics{},
		done:        make(chan struct{}),
		IdleTimeout: DefaultIdleTimeout,
		ReadTimeout: 10 * time.Second,
		MaxWorkers:  10000,
	}
	
	// Create shards to reduce channel contention
	numShards := runtime.GOMAXPROCS(0) * 8
	if numShards < 16 {
		numShards = 16
	}
	s.shards = make([]WorkerShard, numShards)
	
	// Queue size per shard. e.g. 10000 total connections / 16 shards = 625 connections per shard
	queueSize := s.MaxWorkers / numShards
	for i := 0; i < numShards; i++ {
		s.shards[i] = WorkerShard{
			queue: make(chan net.Conn, queueSize),
		}
	}
	
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

	s.startWorkerPool()

	return s.serve()
}

// StartTLS opens a TCP socket and wraps it in a TLS listener.
func (s *Server) StartTLS(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load key pair: %w", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2", "http/1.1"},
	}

	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to bind to address %s: %w", s.addr, err)
	}

	tlsListener := tls.NewListener(l, config)

	s.mu.Lock()
	s.listener = tlsListener
	s.mu.Unlock()

	s.startWorkerPool()

	return s.serve()
}

func (s *Server) serve() error {
	fmt.Printf("TitanHTTP Server successfully started.\n")
	fmt.Printf("Listening on %s...\n", s.Addr())

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

	s.mu.Lock()
	l := s.listener
	s.mu.Unlock()

	if l == nil {
		return nil
	}

	// The main connection accept loop (Multi-Acceptor Model).
	// We spawn multiple goroutines to call Accept() concurrently.
	// This drastically increases the rate at which we drain the OS TCP Backlog,
	// preventing "target machine actively refused it" errors during C10K bursts.
	acceptors := 4 // Optimal for generic multi-core systems without excessive contention
	
	for i := 0; i < acceptors; i++ {
		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					select {
					case <-s.done:
						return // server is shutting down
					default:
						if errors.Is(err, net.ErrClosed) {
							return // Listener closed during shutdown
						}
						// Silently ignore accept errors during burst shutdown
						if !strings.Contains(err.Error(), "use of closed network connection") {
							fmt.Printf("Error accepting connection: %v\n", err)
						}
						continue
					}
				}

				// Increment active connections before dispatching to prevent shutdown races
				s.activeConnWg.Add(1)

				// Dispatch to the worker pool.
				// We use atomic round-robin to pick a shard, ensuring perfect distribution
				// and eliminating the single-channel mutex bottleneck.
				shardIdx := atomic.AddUint64(&s.roundRobin, 1) % uint64(len(s.shards))
				
				select {
				case s.shards[shardIdx].queue <- conn:
					// Dispatched successfully
				default:
					// Capacity exhausted! Backpressure triggered.
					// Reject connection gracefully by closing instantly to protect server memory.
					conn.Close()
					s.metrics.ConnectionRejected()
					s.activeConnWg.Done() // Decrement because we didn't process it
				}
			}
		}()
	}

	// Block the main thread until shutdown is triggered
	<-s.done
	return nil
}

func (s *Server) startWorkerPool() {
	workersPerShard := s.MaxWorkers / len(s.shards)
	
	for i := 0; i < len(s.shards); i++ {
		shard := s.shards[i]
		for j := 0; j < workersPerShard; j++ {
			go func(q chan net.Conn) {
				for conn := range q {
					s.handleConnection(conn)
				}
			}(shard.queue)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// Defers are executed when the surrounding function returns (LIFO order).
	// We close the connection last, and signal the WaitGroup.
	defer s.activeConnWg.Done()
	defer conn.Close()

	s.metrics.ConnectionOpened()
	defer func() {
		s.metrics.ConnectionClosed()
	}()

	// Recover from panics to isolate connection failures and prevent server crashes.
	defer func() {
		if r := recover(); r != nil {
			s.metrics.PanicRecovered()
			fmt.Printf("Critical: Connection panic recovered: %v\n", r)
			importDebug := "runtime/debug" // to ensure we can print
			_ = importDebug
			// Print stack trace directly instead of importing
			fmt.Printf("Stack trace:\n%s\n", string(func() []byte {
				buf := make([]byte, 10240)
				n := runtime.Stack(buf, false)
				return buf[:n]
			}()))
			resp := http.NewResponse500()
			n, _ := resp.WriteTo(conn)
			s.metrics.RequestServed(n)
		}
	}()

	// Check for HTTP/2 ALPN negotiation
	if tlsConn, ok := conn.(*tls.Conn); ok {
		// Force TLS Handshake to read the negotiated protocol before parsing
		if err := tlsConn.Handshake(); err == nil {
			if tlsConn.ConnectionState().NegotiatedProtocol == "h2" {
				s.handleHTTP2(conn)
				return
			}
		} else {
			fmt.Printf("TLS Handshake error: %v\n", err)
			return
		}
	}

	r := readerPool.Get().(*bufio.Reader)
	r.Reset(conn)
	defer readerPool.Put(r)

	requestsServed := 0

	// Connection loop: continuously read from the socket until EOF or error.
	for {
		// Apply ReadTimeout for the first request, and IdleTimeout for subsequent keep-alive requests.
		timeout := s.ReadTimeout
		if requestsServed > 0 {
			timeout = s.IdleTimeout
		}

		if timeout > 0 {
			err := conn.SetReadDeadline(time.Now().Add(timeout))
			if err != nil {
				fmt.Printf("Failed to set read deadline: %v\n", err)
				return
			}
		}

		req, err := http.ParseRequest(r)
		
		// Clear the read deadline while processing the request so long handlers aren't killed
		if timeout > 0 {
			conn.SetReadDeadline(time.Time{})
		}
		if err != nil {
			if err == io.EOF {
				// Silently handle normal disconnects to avoid log noise in concurrent environments.
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Silently handle read timeouts to avoid log noise.
			} else if strings.Contains(err.Error(), "forcibly closed") || strings.Contains(err.Error(), "connection reset") {
				// Silently handle abrupt client disconnects (e.g., bombardier shutting down)
			} else {
				fmt.Printf("Error parsing request: %v\n", err)
				// Send specific HTTP errors on parser errors
				var resp *http.Response
				switch {
				case errors.Is(err, http.ErrURITooLong):
					resp = http.NewResponse414()
				case errors.Is(err, http.ErrHeaderFieldsTooLarge):
					resp = http.NewResponse431()
				case errors.Is(err, http.ErrNotImplemented):
					resp = http.NewResponse501()
				case errors.Is(err, http.ErrMethodNotAllowed):
					resp = http.NewResponse405()
				case errors.Is(err, http.ErrInvalidMethod):
					resp = http.NewResponse400()
				default:
					resp = http.NewResponse400()
				}
				if _, wErr := resp.WriteTo(conn); wErr != nil {
					fmt.Printf("Failed to write parser error response: %v\n", wErr)
				}
			}
			return
		}

		// Provide a default background context. 
		// We avoid per-request context allocations under high churn workloads unless specifically needed.
		req.WithContext(context.Background())

		// Populate network-level details on the Request
		if addr := conn.RemoteAddr(); addr != nil {
			req.RemoteAddr = addr.String()
		} else {
			req.RemoteAddr = "127.0.0.1:0"
		}
		
		if _, isTLS := conn.(*tls.Conn); isTLS {
			req.Scheme = "https"
		} else {
			req.Scheme = "http"
		}

		if req.Headers.Get("expect") == "100-continue" {
			if _, wErr := conn.Write([]byte("HTTP/1.1 100 Continue\r\n\r\n")); wErr != nil {
				return
			}
			if err := http.ParseBody(r, req); err != nil {
				resp := http.NewResponse400()
				resp.WriteTo(conn)
				return
			}
		}

		err = req.Validate()
		if err != nil {
			fmt.Printf("Request validation failed: %v\n", err)
			resp := http.NewResponse400()
			if _, err := resp.WriteTo(conn); err != nil {
				fmt.Printf("Failed to write validation error response: %v\n", err)
			}
			return
		}

		// Dispatch request to the router
		resp := s.router.ServeHTTP(req)

		requestsServed++
		wantsKeepAlive := req.WantsKeepAlive()

		if wantsKeepAlive {
			resp.Headers.Set("Connection", "keep-alive")
			if s.IdleTimeout > 0 {
				resp.Headers.Set("Keep-Alive", fmt.Sprintf("timeout=%d", int(s.IdleTimeout.Seconds())))
			}
		} else {
			resp.Headers.Set("Connection", "close")
		}

		n, err := resp.WriteTo(conn)
		s.metrics.RequestServed(n)

		// Return pooled objects to reduce GC pressure
		http.ReleaseRequest(req)
		if resp != nil {
			http.ReleaseResponse(resp)
		}

		if err != nil {
			errStr := err.Error()
			// Silently ignore normal benchmark disconnections
			if !strings.Contains(errStr, "aborted by the software") &&
				!strings.Contains(errStr, "forcibly closed") &&
				!strings.Contains(errStr, "connection reset") &&
				!strings.Contains(errStr, "broken pipe") &&
				!strings.Contains(errStr, "closed pipe") &&
				!strings.Contains(errStr, "use of closed") {
				fmt.Printf("Error writing to connection: %v\n", err)
			}
			return
		}

		// Close connection if keep-alive is not desired or max requests reached
		if !wantsKeepAlive {
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

	// Wait for active connections to drain, racing against ctx.Done()
	stopCh := make(chan struct{})
	go func() {
		s.activeConnWg.Wait()
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCh:
		return nil
	}
}
