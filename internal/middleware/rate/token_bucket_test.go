package rate

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	// 2 tokens capacity, 10 tokens per second refill
	tb := NewTokenBucket(2, 10.0)

	// First request: gets token (1 remaining)
	if !tb.Allow("127.0.0.1") {
		t.Error("Expected first request to be allowed")
	}

	// Second request: gets token (0 remaining)
	if !tb.Allow("127.0.0.1") {
		t.Error("Expected second request to be allowed")
	}

	// Third request: rejected (0 remaining)
	if tb.Allow("127.0.0.1") {
		t.Error("Expected third request to be rejected")
	}

	// Wait 150ms. Since rate is 10/sec, 150ms should yield 1.5 tokens.
	// So 1 token is available.
	time.Sleep(150 * time.Millisecond)

	// Fourth request: gets token (0.5 remaining)
	if !tb.Allow("127.0.0.1") {
		t.Error("Expected fourth request to be allowed after refill")
	}

	// Fifth request: rejected (0.5 remaining, need 1.0)
	if tb.Allow("127.0.0.1") {
		t.Error("Expected fifth request to be rejected")
	}
}
