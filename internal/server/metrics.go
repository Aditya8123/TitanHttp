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
	rejectedConns atomic.Int64
}

// ConnectionOpened should be called when a new connection begins processing.
func (m *Metrics) ConnectionOpened() {
	m.activeConns.Add(1)
}

// ConnectionClosed should be called when a connection completes.
func (m *Metrics) ConnectionClosed() {
	m.activeConns.Add(-1)
}

// RequestServed should be called when an HTTP request completes, providing the bytes processed.
func (m *Metrics) RequestServed(bytes int64) {
	m.totalRequests.Add(1)
	m.totalBytes.Add(bytes)
}

// PanicRecovered increments the counter of panics gracefully handled by the server.
func (m *Metrics) PanicRecovered() {
	m.panicCount.Add(1)
}

// ConnectionRejected increments the counter for connections shed under backpressure.
func (m *Metrics) ConnectionRejected() {
	m.rejectedConns.Add(1)
}

// Report prints the current snapshot of server metrics.
func (m *Metrics) Report() {
	fmt.Printf(
		"[Metrics] Total: %d | Active: %d | Bytes: %d | Panics: %d | Rejected: %d\n",
		m.totalRequests.Load(),
		m.activeConns.Load(),
		m.totalBytes.Load(),
		m.panicCount.Load(),
		m.rejectedConns.Load(),
	)
}
