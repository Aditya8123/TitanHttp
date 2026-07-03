package rate

import (
	"sync"
	"time"
)

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

// TokenBucket implements the Limiter interface using the token bucket algorithm.
type TokenBucket struct {
	mu           sync.Mutex
	capacity     float64
	refillRate   float64 // tokens per second
	buckets      map[string]*bucket
}

// NewTokenBucket creates a new TokenBucket rate limiter.
func NewTokenBucket(capacity int, refillPerSec float64) *TokenBucket {
	return &TokenBucket{
		capacity:   float64(capacity),
		refillRate: refillPerSec,
		buckets:    make(map[string]*bucket),
	}
}

// StartSweeper begins a background goroutine to periodically clean up stale IPs.
// An IP is considered stale if it hasn't made a request in a long time (e.g., 5 minutes).
func (tb *TokenBucket) StartSweeper(interval time.Duration, maxIdle time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			tb.mu.Lock()
			now := time.Now()
			for ip, b := range tb.buckets {
				if now.Sub(b.lastRefill) > maxIdle {
					delete(tb.buckets, ip)
				}
			}
			tb.mu.Unlock()
		}
	}()
}

// Allow implements the Limiter interface.
func (tb *TokenBucket) Allow(ip string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	b, exists := tb.buckets[ip]
	now := time.Now()

	if !exists {
		// New IP: bucket starts full, minus one for this request
		tb.buckets[ip] = &bucket{
			tokens:     tb.capacity - 1,
			lastRefill: now,
		}
		return true
	}

	// Calculate refilled tokens
	elapsed := now.Sub(b.lastRefill).Seconds()
	refill := elapsed * tb.refillRate

	if refill > 0 {
		b.tokens += refill
		if b.tokens > tb.capacity {
			b.tokens = tb.capacity
		}
		// We only update lastRefill if we actually added tokens, 
		// but typically we can just set it to now.
		b.lastRefill = now
	}

	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true
	}

	return false
}
