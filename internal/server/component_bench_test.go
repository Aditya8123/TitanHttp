package server

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/Aditya8123/TitanHttp/internal/http"
)

// mockConn implements net.Conn to avoid OS pipe overhead during benchmarking
type mockConn struct {
	reqData    []byte
	readOffset int
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	if m.readOffset >= len(m.reqData) {
		return 0, io.EOF
	}
	n = copy(b, m.reqData[m.readOffset:])
	m.readOffset += n
	return n, nil
}
func (m *mockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

// BenchmarkServer_Component simulates a full request-response lifecycle
// using an in-memory mock net.Conn to bypass OS networking, focusing strictly
// on the performance of our Router, Parser, and connection handlers.
func BenchmarkServer_Component(b *testing.B) {
	srv := NewServer(":0")

	srv.Router().Get("/ping", func(req *http.Request) *http.Response {
		resp := http.NewResponse200()
		resp.Body = []byte("pong")
		return resp
	})

	reqBytes := []byte("GET /ping HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n")

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn := &mockConn{reqData: reqBytes}
			srv.activeConnWg.Add(1)
			srv.handleConnection(conn)
		}
	})
}
