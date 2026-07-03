package rate

// Limiter defines the interface for rate limiting algorithms.
type Limiter interface {
	// Allow returns true if a request from the given IP is permitted, false otherwise.
	Allow(ip string) bool
}
