package rate

import (
	"sync"
	"time"
)

// defaultShardCount defines the number of shards used to partition the rate limiters.
// Using 64 shards significantly reduces global mutex contention across parallel
// requests while keeping memory overhead reasonable.
const defaultShardCount = 64

type slidingWindowShard struct {
	mu   sync.Mutex
	logs map[string][]time.Time
}

// SlidingWindow implements the Limiter interface using the sliding window log algorithm
// mapped across multiple shards to eliminate lock contention.
type SlidingWindow struct {
	limit  int
	window time.Duration
	shards []*slidingWindowShard
}

// NewSlidingWindow creates a new sharded SlidingWindow rate limiter.
func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	sw := &SlidingWindow{
		limit:  limit,
		window: window,
		shards: make([]*slidingWindowShard, defaultShardCount),
	}
	for i := 0; i < defaultShardCount; i++ {
		sw.shards[i] = &slidingWindowShard{
			logs: make(map[string][]time.Time),
		}
	}
	return sw
}

// getShard returns the specific shard for a given IP address using FNV-1a hashing.
func (sw *SlidingWindow) getShard(ip string) *slidingWindowShard {
	var hash uint32 = 2166136261
	for i := 0; i < len(ip); i++ {
		hash ^= uint32(ip[i])
		hash *= 16777619
	}
	return sw.shards[hash%defaultShardCount]
}

// StartSweeper begins a background goroutine to periodically clean up stale IPs.
// It iterates through shards independently to minimize lock hold times.
func (sw *SlidingWindow) StartSweeper(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			cutoff := now.Add(-sw.window)
			for _, shard := range sw.shards {
				shard.mu.Lock()
				for ip, timestamps := range shard.logs {
					// If there are no timestamps or the newest timestamp is older than the cutoff
					if len(timestamps) == 0 || timestamps[len(timestamps)-1].Before(cutoff) {
						delete(shard.logs, ip)
					}
				}
				shard.mu.Unlock()
			}
		}
	}()
}

// Allow implements the Limiter interface.
func (sw *SlidingWindow) Allow(ip string) bool {
	shard := sw.getShard(ip)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	timestamps := shard.logs[ip]

	// Purge timestamps older than the cutoff window
	// Since timestamps are appended chronologically, we can just find the first valid one
	validIdx := 0
	for i, t := range timestamps {
		if t.After(cutoff) {
			validIdx = i
			break
		}
		// If we reach the end and all are older, validIdx stays 0 but we want len(timestamps)
		if i == len(timestamps)-1 {
			validIdx = len(timestamps)
		}
	}

	// Slice off the old timestamps to prevent memory growth
	if validIdx > 0 {
		timestamps = timestamps[validIdx:]

		// Shrink the backing array if the capacity is large but only a small fraction is used.
		// This prevents indefinite memory retention (leak) during burst loads.
		if cap(timestamps) > 1024 && len(timestamps) < cap(timestamps)/4 {
			newTs := make([]time.Time, len(timestamps))
			copy(newTs, timestamps)
			timestamps = newTs
		}
	}

	// Check if within limit
	if len(timestamps) < sw.limit {
		// Allow request and append timestamp
		timestamps = append(timestamps, now)
		shard.logs[ip] = timestamps
		return true
	}

	// Reject request, save the purged slice back
	shard.logs[ip] = timestamps
	return false
}
