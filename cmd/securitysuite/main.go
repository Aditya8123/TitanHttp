package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type SecurityTest struct {
	Name string
	Run  func(target string) (bool, string)
}

func main() {
	target := "localhost:8080"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	fmt.Printf("========================================\n")
	fmt.Printf(" TitanHTTP Security Suite\n")
	fmt.Printf(" Target: %s\n", target)
	fmt.Printf("========================================\n\n")

	tests := []SecurityTest{
		{
			Name: "Slowloris",
			Run: func(target string) (bool, string) {
				conn, err := net.Dial("tcp", target)
				if err != nil {
					return false, err.Error()
				}
				defer conn.Close()

				conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\n"))

				// Trickle headers slowly
				for i := 0; i < 15; i++ {
					_, err := conn.Write([]byte(fmt.Sprintf("X-Slow-%d: 1\r\n", i)))
					if err != nil {
						// Connection closed by server (mitigated)
						return true, ""
					}
					time.Sleep(500 * time.Millisecond)
				}
				return false, "Server did not drop connection after 7.5 seconds of Slowloris"
			},
		},
		{
			Name: "Slow POST",
			Run: func(target string) (bool, string) {
				conn, err := net.Dial("tcp", target)
				if err != nil {
					return false, err.Error()
				}
				defer conn.Close()

				conn.Write([]byte("POST / HTTP/1.1\r\nHost: localhost\r\nContent-Length: 100\r\n\r\n"))

				// Trickle body
				for i := 0; i < 15; i++ {
					_, err := conn.Write([]byte("A"))
					if err != nil {
						return true, ""
					}
					time.Sleep(500 * time.Millisecond)
				}
				return false, "Server did not drop connection during Slow POST"
			},
		},
		{
			Name: "Header Flood (2000 headers)",
			Run: func(target string) (bool, string) {
				conn, err := net.Dial("tcp", target)
				if err != nil {
					return false, err.Error()
				}
				defer conn.Close()

				conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\n"))

				for i := 0; i < 2000; i++ {
					_, err := conn.Write([]byte(fmt.Sprintf("X-Junk-%d: AAAAA\r\n", i)))
					if err != nil {
						return true, ""
					}
				}

				conn.Write([]byte("\r\n"))

				reader := bufio.NewReader(conn)
				resp, err := reader.ReadString('\n')
				if err != nil || strings.Contains(resp, "400") || strings.Contains(resp, "431") {
					return true, ""
				}
				return false, fmt.Sprintf("Server accepted massive headers: %s", strings.TrimSpace(resp))
			},
		},
		{
			Name: "Invalid HTTP",
			Run: func(target string) (bool, string) {
				conn, err := net.Dial("tcp", target)
				if err != nil {
					return false, err.Error()
				}
				defer conn.Close()

				conn.Write([]byte("JUST_GARBAGE_NO_PROTOCOL\r\n\r\n"))
				reader := bufio.NewReader(conn)
				resp, err := reader.ReadString('\n')
				if err != nil || strings.Contains(resp, "400") {
					return true, ""
				}
				return false, fmt.Sprintf("Unexpected response to garbage: %s", strings.TrimSpace(resp))
			},
		},
		{
			Name: "Connection Flood (200 Idle)",
			Run: func(target string) (bool, string) {
				var conns []net.Conn
				for i := 0; i < 200; i++ {
					conn, err := net.Dial("tcp", target)
					if err == nil {
						conns = append(conns, conn)
					}
				}

				time.Sleep(6 * time.Second) // IdleTimeout is usually 5s

				active := 0
				for _, conn := range conns {
					conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
					buf := make([]byte, 1)
					_, err := conn.Read(buf)
					if err == nil {
						active++
					} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						active++
					}
					conn.Close()
				}

				if active == 0 {
					return true, ""
				}
				return false, fmt.Sprintf("%d connections still alive after timeout", active)
			},
		},
		{
			Name: "Oversized Headers",
			Run: func(target string) (bool, string) {
				conn, err := net.Dial("tcp", target)
				if err != nil {
					return false, err.Error()
				}
				defer conn.Close()

				// A single massive header (e.g., 2MB)
				hugeHeader := "X-Huge: " + strings.Repeat("A", 2*1024*1024) + "\r\n"
				conn.Write([]byte("GET / HTTP/1.1\r\nHost: localhost\r\n"))
				conn.Write([]byte(hugeHeader))
				conn.Write([]byte("\r\n"))

				reader := bufio.NewReader(conn)
				resp, err := reader.ReadString('\n')
				// Depending on server limits, it could drop connection or return 431/400
				if err != nil || strings.Contains(resp, "4") { // Any 4xx
					return true, ""
				}
				return false, "Server accepted 2MB single header"
			},
		},
		{
			Name: "Request Smuggling (CL-TE)",
			Run: func(target string) (bool, string) {
				conn, err := net.Dial("tcp", target)
				if err != nil {
					return false, err.Error()
				}
				defer conn.Close()

				// Send conflicting headers
				req := "POST / HTTP/1.1\r\nHost: localhost\r\nContent-Length: 4\r\nTransfer-Encoding: chunked\r\n\r\n0\r\n\r\n"
				conn.Write([]byte(req))

				reader := bufio.NewReader(conn)
				resp, err := reader.ReadString('\n')
				if err != nil || strings.Contains(resp, "400") {
					return true, ""
				}
				return false, "Server did not reject conflicting CL and TE headers"
			},
		},
		{
			Name: "Path Traversal",
			Run: func(target string) (bool, string) {
				conn, err := net.Dial("tcp", target)
				if err != nil {
					return false, err.Error()
				}
				defer conn.Close()

				req := "GET /../../../etc/passwd HTTP/1.1\r\nHost: localhost\r\n\r\n"
				conn.Write([]byte(req))

				var buf bytes.Buffer
				reader := bufio.NewReader(conn)
				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					buf.WriteString(line)
					if line == "\r\n" {
						break
					}
				}

				resp := buf.String()
				if strings.Contains(resp, "400") || strings.Contains(resp, "403") || strings.Contains(resp, "404") {
					return true, ""
				}
				return false, fmt.Sprintf("Unexpected response to traversal: %s", strings.TrimSpace(resp))
			},
		},
	}

	passed := 0
	failed := 0

	for _, tc := range tests {
		fmt.Printf("Testing %-40s ", tc.Name)
		ok, msg := tc.Run(target)
		if ok {
			fmt.Printf("[ \033[32mPASS\033[0m ]\n")
			passed++
		} else {
			fmt.Printf("[ \033[31mFAIL\033[0m ] (%s)\n", msg)
			failed++
		}
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf(" Results: %d Passed, %d Failed\n", passed, failed)
	fmt.Printf("========================================\n")

	if failed > 0 {
		os.Exit(1)
	}
}
