package proxy

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
	"github.com/Aditya8123/TitanHttp/internal/router"
)

// Backend represents a single upstream server in the load balancer pool.
type Backend struct {
	URL          string
	Alive        bool
	mux          sync.RWMutex
	ActiveConns  int64 // For future features like least-connections
}

// SetAlive safely updates the backend's health status.
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive safely checks the backend's health status.
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

// StartHealthCheck starts a background goroutine to ping the backend periodically.
// It performs a simple TCP connection test.
func (b *Backend) StartHealthCheck(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			conn, err := net.DialTimeout("tcp", b.URL, 2*time.Second)
			if err != nil {
				b.SetAlive(false)
			} else {
				conn.Close()
				b.SetAlive(true)
			}
		}
	}()
}

// LoadBalancer manages a pool of backends and distributes incoming requests.
type LoadBalancer struct {
	backends []*Backend
	current  uint32 // Atomic counter for Round-Robin
}

// NewLoadBalancer creates a new LoadBalancer handler routing to multiple targets.
// targets should be in the format "host:port", e.g., "localhost:8081".
func NewLoadBalancer(targets []string) router.Handler {
	lb := &LoadBalancer{
		backends: make([]*Backend, 0, len(targets)),
	}

	for _, target := range targets {
		backend := &Backend{
			URL:   target,
			Alive: true,
		}
		backend.StartHealthCheck(10 * time.Second)
		lb.backends = append(lb.backends, backend)
	}

	return lb.ServeHTTP
}

// ServeHTTP implements router.Handler, forwarding the request to the next available backend.
func (lb *LoadBalancer) ServeHTTP(req *http.Request) *http.Response {
	target := lb.NextBackend()
	if target == nil {
		// No healthy backends available
		return http.NewResponse503() // Service Unavailable
	}

	return ForwardRequest(req, target.URL)
}

// NextBackend selects the next healthy backend using Round-Robin.
func (lb *LoadBalancer) NextBackend() *Backend {
	l := uint32(len(lb.backends))
	if l == 0 {
		return nil
	}

	// We loop to find a healthy backend, but we don't want to loop forever.
	// We'll loop at most 'l' times to check all backends.
	for i := uint32(0); i < l; i++ {
		idx := atomic.AddUint32(&lb.current, 1) % l
		backend := lb.backends[idx]
		if backend.IsAlive() {
			return backend
		}
	}

	return nil
}
