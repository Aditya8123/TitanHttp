package rate

import (
	"sync"
	"time"
)

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

type tokenBucketShard struct {
	mu      sync.Mutex
	buckets map[string]*bucket
}

// TokenBucket implements the Limiter interface using the token bucket algorithm
// mapped across multiple shards to eliminate lock contention.
type TokenBucket struct {
	capacity   float64
	refillRate float64 // tokens per second
	shards     []*tokenBucketShard
}

// NewTokenBucket creates a new sharded TokenBucket rate limiter.
func NewTokenBucket(capacity int, refillPerSec float64) *TokenBucket {
	tb := &TokenBucket{
		capacity:   float64(capacity),
		refillRate: refillPerSec,
		shards:     make([]*tokenBucketShard, defaultShardCount),
	}
	for i := 0; i < defaultShardCount; i++ {
		tb.shards[i] = &tokenBucketShard{
			buckets: make(map[string]*bucket),
		}
	}
	return tb
}

// getShard returns the specific shard for a given IP address using FNV-1a hashing.
func (tb *TokenBucket) getShard(ip string) *tokenBucketShard {
	var hash uint32 = 2166136261
	for i := 0; i < len(ip); i++ {
		hash ^= uint32(ip[i])
		hash *= 16777619
	}
	return tb.shards[hash%defaultShardCount]
}

// StartSweeper begins a background goroutine to periodically clean up stale IPs.
// An IP is considered stale if it hasn't made a request in a long time (e.g., 5 minutes).
func (tb *TokenBucket) StartSweeper(interval time.Duration, maxIdle time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			for _, shard := range tb.shards {
				shard.mu.Lock()
				for ip, b := range shard.buckets {
					if now.Sub(b.lastRefill) > maxIdle {
						delete(shard.buckets, ip)
					}
				}
				shard.mu.Unlock()
			}
		}
	}()
}

// Allow implements the Limiter interface.
func (tb *TokenBucket) Allow(ip string) bool {
	shard := tb.getShard(ip)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	b, exists := shard.buckets[ip]
	now := time.Now()

	if !exists {
		// New IP: bucket starts full, minus one for this request
		shard.buckets[ip] = &bucket{
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
