package rate

import (
	"testing"
	"time"
)

func TestSlidingWindow(t *testing.T) {
	// Limit to 2 requests per 100ms
	sw := NewSlidingWindow(2, 100*time.Millisecond)

	// 1st request (t=0)
	if !sw.Allow("127.0.0.1") {
		t.Error("Expected 1st request to be allowed")
	}

	// 2nd request (t=0)
	if !sw.Allow("127.0.0.1") {
		t.Error("Expected 2nd request to be allowed")
	}

	// 3rd request (t=0) - should fail because limit is 2
	if sw.Allow("127.0.0.1") {
		t.Error("Expected 3rd request to be rejected")
	}

	// Wait 120ms to pass the 100ms window
	time.Sleep(120 * time.Millisecond)

	// 4th request (t=120) - should pass because previous requests are purged
	if !sw.Allow("127.0.0.1") {
		t.Error("Expected 4th request to be allowed after window passed")
	}

	// Check another IP doesn't share limit
	if !sw.Allow("10.0.0.1") {
		t.Error("Expected request from new IP to be allowed")
	}
}
