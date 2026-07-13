package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type TestCase struct {
	Name        string
	Request     string
	CheckFunc   func(response string, status string) (bool, string)
}

func main() {
	target := "localhost:8080"
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	fmt.Printf("========================================\n")
	fmt.Printf(" TitanHTTP Compliance Suite\n")
	fmt.Printf(" Target: %s\n", target)
	fmt.Printf("========================================\n\n")

	tests := []TestCase{
		{
			Name: "GET Method",
			Request: "GET /ping HTTP/1.1\r\nHost: localhost\r\n\r\n",
			CheckFunc: expectStatus("200 OK"),
		},
		{
			Name: "POST Method",
			Request: "POST /ping HTTP/1.1\r\nHost: localhost\r\nContent-Length: 0\r\n\r\n",
			CheckFunc: expectStatus("405"), // Depending on /ping routes, but it responds with something valid
		},
		{
			Name: "PUT Method",
			Request: "PUT /ping HTTP/1.1\r\nHost: localhost\r\nContent-Length: 0\r\n\r\n",
			CheckFunc: expectStatus("405"),
		},
		{
			Name: "PATCH Method",
			Request: "PATCH /ping HTTP/1.1\r\nHost: localhost\r\nContent-Length: 0\r\n\r\n",
			CheckFunc: expectStatus("405"),
		},
		{
			Name: "DELETE Method",
			Request: "DELETE /ping HTTP/1.1\r\nHost: localhost\r\n\r\n",
			CheckFunc: expectStatus("405"),
		},
		{
			Name: "OPTIONS Method",
			Request: "OPTIONS /ping HTTP/1.1\r\nHost: localhost\r\n\r\n",
			CheckFunc: expectStatus("405"), // Assuming OPTIONS not explicitly defined
		},
		{
			Name: "HEAD Method",
			Request: "HEAD /ping HTTP/1.1\r\nHost: localhost\r\n\r\n",
			CheckFunc: expectStatus("405"), // Or 200 if handled by router
		},
		{
			Name: "Chunked Transfer-Encoding (Request)",
			Request: "POST /ping HTTP/1.1\r\nHost: localhost\r\nTransfer-Encoding: chunked\r\n\r\n0\r\n\r\n",
			CheckFunc: expectStatus("405"), // TitanHTTP now supports chunked requests. /ping is GET-only, so 405 is correct.
		},
		{
			Name: "Content-Length",
			Request: "POST /ping HTTP/1.1\r\nHost: localhost\r\nContent-Length: 4\r\n\r\nbody",
			CheckFunc: expectStatus("405"), // Standard length check
		},
		{
			Name: "Keep-Alive",
			Request: "GET /ping HTTP/1.1\r\nHost: localhost\r\nConnection: keep-alive\r\n\r\n",
			CheckFunc: expectHeader("Connection: keep-alive"),
		},
		{
			Name: "Connection: close",
			Request: "GET /ping HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n",
			CheckFunc: expectHeader("Connection: close"),
		},
		{
			Name: "100-Continue",
			Request: "POST /ping HTTP/1.1\r\nHost: localhost\r\nExpect: 100-continue\r\nContent-Length: 4\r\n\r\nbody",
			CheckFunc: expectStatus("100"), // Will likely fail
		},
		{
			Name: "Compression (Gzip)",
			Request: "GET /ping HTTP/1.1\r\nHost: localhost\r\nAccept-Encoding: gzip\r\n\r\n",
			CheckFunc: expectHeader("Content-Encoding: gzip"), // Assuming /ping doesn't gzip, but maybe something else does
		},
		{
			Name: "Range Requests",
			Request: "GET / HTTP/1.1\r\nHost: localhost\r\nRange: bytes=0-10\r\n\r\n",
			CheckFunc: expectStatus("206"), // Will likely fail
		},
		{
			Name: "Conditional GET (ETag)",
			Request: "GET / HTTP/1.1\r\nHost: localhost\r\nIf-None-Match: \"fake-etag\"\r\n\r\n",
			CheckFunc: expectStatus("200 OK"), 
		},
		{
			Name: "Invalid HTTP methods",
			Request: "INVALID / HTTP/1.1\r\nHost: localhost\r\n\r\n",
			CheckFunc: expectStatus("400"), // or 405
		},
		{
			Name: "Malformed headers",
			Request: "GET / HTTP/1.1\r\nHost: localhost\r\nBadHeader\r\n\r\n",
			CheckFunc: expectStatus("400"),
		},
		{
			Name: "Invalid request line",
			Request: "GET / HTTP/1.1 EXTRABYTES\r\nHost: localhost\r\n\r\n",
			CheckFunc: expectStatus("400"),
		},
		{
			Name: "Duplicate Content-Length",
			Request: "POST / HTTP/1.1\r\nHost: localhost\r\nContent-Length: 5\r\nContent-Length: 5\r\n\r\nbody!",
			CheckFunc: expectStatus("400"), // Security vulnerability
		},
		{
			Name: "Huge URI",
			Request: fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: localhost\r\n\r\n", strings.Repeat("a", 8000)),
			CheckFunc: expectStatus("414"), // URI Too Long (might fail if not implemented)
		},
	}

	// We'll need a reliable ping endpoint that accepts methods to test properly.
	// But let's build the runner first.
	passed := 0
	failed := 0

	for _, tc := range tests {
		fmt.Printf("Testing %-40s ", tc.Name)
		
		conn, err := net.DialTimeout("tcp", target, 2*time.Second)
		if err != nil {
			fmt.Printf("[ \033[31mFAIL\033[0m ] (Dial error: %v)\n", err)
			failed++
			continue
		}

		conn.SetDeadline(time.Now().Add(2 * time.Second))
		_, err = conn.Write([]byte(tc.Request))
		if err != nil {
			fmt.Printf("[ \033[31mFAIL\033[0m ] (Write error: %v)\n", err)
			conn.Close()
			failed++
			continue
		}

		// Read response
		var respBuffer bytes.Buffer
		reader := bufio.NewReader(conn)
		statusLine, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			// Some tests expect the connection to be dropped immediately
		}
		
		respBuffer.WriteString(statusLine)
		
		// Read headers
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			respBuffer.WriteString(line)
			if line == "\r\n" || line == "\n" {
				break
			}
		}

		status := strings.TrimSpace(statusLine)
		fullResp := respBuffer.String()

		ok, msg := tc.CheckFunc(fullResp, status)
		if ok {
			fmt.Printf("[ \033[32mPASS\033[0m ]\n")
			passed++
		} else {
			fmt.Printf("[ \033[31mFAIL\033[0m ] (%s. Got: %s)\n", msg, status)
			failed++
		}
		conn.Close()
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf(" Results: %d Passed, %d Failed\n", passed, failed)
	fmt.Printf("========================================\n")

	if failed > 0 {
		os.Exit(1)
	}
}

func expectStatus(expected string) func(string, string) (bool, string) {
	return func(resp string, status string) (bool, string) {
		if strings.Contains(status, expected) {
			return true, ""
		}
		return false, fmt.Sprintf("Expected status containing '%s'", expected)
	}
}

func expectHeader(expected string) func(string, string) (bool, string) {
	return func(resp string, status string) (bool, string) {
		// lower case check to be safe
		if strings.Contains(strings.ToLower(resp), strings.ToLower(expected)) {
			return true, ""
		}
		return false, fmt.Sprintf("Expected header containing '%s'", expected)
	}
}
