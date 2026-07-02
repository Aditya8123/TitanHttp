package server

import (
	"bytes"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// mockConn implements net.Conn for testing purposes
type mockConn struct {
	net.Conn
	writeBuf bytes.Buffer
	closed   bool
	mu       sync.Mutex
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.writeBuf.Write(b)
}

func (m *mockConn) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1234}
}

func TestWorkerPool_Execution(t *testing.T) {
	var handledCount int32
	var wg sync.WaitGroup

	handler := func(conn net.Conn) {
		atomic.AddInt32(&handledCount, 1)
		wg.Done()
	}

	pool := NewWorkerPool(2, 5, handler)
	pool.Start()

	// Submit 5 jobs
	for i := 0; i < 5; i++ {
		wg.Add(1)
		success := pool.Submit(&mockConn{})
		if !success {
			t.Errorf("Expected job to be submitted successfully")
		}
	}

	// Wait for jobs to be handled
	wg.Wait()
	pool.Stop()

	if handledCount != 5 {
		t.Errorf("Expected 5 jobs to be handled, got %d", handledCount)
	}
}

func TestWorkerPool_LoadShedding(t *testing.T) {
	// A handler that blocks so we can fill up the queue
	blockCh := make(chan struct{})
	handler := func(conn net.Conn) {
		<-blockCh
	}

	// 1 worker, queue size 1
	pool := NewWorkerPool(1, 1, handler)
	pool.Start()

	// Job 1: picked up by the 1 worker (blocking)
	pool.Submit(&mockConn{})

	// Wait a tiny bit to ensure worker picks it up
	time.Sleep(10 * time.Millisecond)

	// Job 2: goes into the queue of size 1 (queue is now full)
	pool.Submit(&mockConn{})

	// Job 3: should be rejected due to load shedding
	conn3 := &mockConn{}
	success := pool.Submit(conn3)
	if success {
		t.Errorf("Expected job 3 to be rejected (load shed)")
	}

	// Verify the 503 response was written to the rejected connection
	conn3.mu.Lock()
	output := conn3.writeBuf.String()
	isClosed := conn3.closed
	conn3.mu.Unlock()

	expectedResponse := "HTTP/1.1 503 Service Unavailable"
	if !bytes.Contains([]byte(output), []byte(expectedResponse)) {
		t.Errorf("Expected 503 response, got: %s", output)
	}

	if !isClosed {
		t.Errorf("Expected rejected connection to be closed")
	}

	// Unblock workers and stop pool
	close(blockCh)
	pool.Stop()
}
