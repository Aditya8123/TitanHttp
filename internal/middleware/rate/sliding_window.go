package rate

import (
	"sync"
	"time"
)

// SlidingWindow implements the Limiter interface using the sliding window log algorithm.
type SlidingWindow struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	logs   map[string][]time.Time
}

// NewSlidingWindow creates a new SlidingWindow rate limiter.
func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		limit:  limit,
		window: window,
		logs:   make(map[string][]time.Time),
	}
}

// StartSweeper begins a background goroutine to periodically clean up stale IPs.
// An IP is considered stale if its most recent request is older than the window.
func (sw *SlidingWindow) StartSweeper(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			sw.mu.Lock()
			now := time.Now()
			cutoff := now.Add(-sw.window)
			for ip, timestamps := range sw.logs {
				// If there are no timestamps or the newest timestamp is older than the cutoff
				if len(timestamps) == 0 || timestamps[len(timestamps)-1].Before(cutoff) {
					delete(sw.logs, ip)
				}
			}
			sw.mu.Unlock()
		}
	}()
}

// Allow implements the Limiter interface.
func (sw *SlidingWindow) Allow(ip string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	timestamps := sw.logs[ip]
	
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
	}

	// Check if within limit
	if len(timestamps) < sw.limit {
		// Allow request and append timestamp
		timestamps = append(timestamps, now)
		sw.logs[ip] = timestamps
		return true
	}

	// Reject request, save the purged slice back
	sw.logs[ip] = timestamps
	return false
}
