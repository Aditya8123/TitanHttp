package server

import (
	"fmt"
	"sync/atomic"
)

// Metrics tracks the health and performance of the server using lock-free atomic operations.
type Metrics struct {
	totalRequests atomic.Int64
	activeConns   atomic.Int64
	totalBytes    atomic.Int64
	panicCount    atomic.Int32
}

// RequestStarted should be called when a new connection begins processing.
func (m *Metrics) RequestStarted() {
	m.totalRequests.Add(1)
	m.activeConns.Add(1)
}

// RequestFinished should be called when a connection completes, providing the bytes processed.
func (m *Metrics) RequestFinished(bytes int64) {
	m.activeConns.Add(-1)
	m.totalBytes.Add(bytes)
}

// PanicRecovered increments the counter of panics gracefully handled by the server.
func (m *Metrics) PanicRecovered() {
	m.panicCount.Add(1)
}

// Report prints the current snapshot of server metrics.
func (m *Metrics) Report() {
	fmt.Printf("[Metrics] Total: %d | Active: %d | Bytes: %d | Panics: %d\n",
		m.totalRequests.Load(),
		m.activeConns.Load(),
		m.totalBytes.Load(),
		m.panicCount.Load(),
	)
}
