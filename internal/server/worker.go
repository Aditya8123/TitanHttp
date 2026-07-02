package server

import (
	"fmt"
	"net"
	"sync"
)

// WorkerPool manages a pool of goroutines to handle incoming TCP connections
// concurrently while bounding the total number of active workers.
type WorkerPool struct {
	workers int
	jobs    chan net.Conn
	wg      sync.WaitGroup
	handler func(net.Conn)
}

// NewWorkerPool initializes a WorkerPool with the specified number of workers
// and the maximum size of the pending job queue.
func NewWorkerPool(workers, queueSize int, handler func(net.Conn)) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan net.Conn, queueSize),
		handler: handler,
	}
}

// Start spins up the fixed number of worker goroutines. Each worker blocks
// on the jobs channel, waiting to process incoming connections.
func (p *WorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(id int) {
			defer p.wg.Done()
			for conn := range p.jobs {
				p.handler(conn)
			}
		}(i)
	}
}

// Submit attempts to enqueue a new connection into the worker pool.
// If the job queue is full, it performs load shedding by immediately
// sending a 503 response, closing the connection, and returning false.
func (p *WorkerPool) Submit(conn net.Conn) bool {
	select {
	case p.jobs <- conn:
		return true // Job successfully enqueued
	default:
		// Load shedding: the queue is full.
		// Send a raw HTTP 503 response and drop the connection.
		fmt.Printf("Warning: Worker pool queue full. Shedding load for %s\n", conn.RemoteAddr().String())
		conn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\nContent-Length: 0\r\nConnection: close\r\n\r\n"))
		conn.Close()
		return false
	}
}

// Stop gracefully shuts down the worker pool. It closes the jobs channel,
// signaling all workers to exit after they finish their current connection,
// and blocks until all workers have completely stopped.
func (p *WorkerPool) Stop() {
	close(p.jobs)
	p.wg.Wait()
}
