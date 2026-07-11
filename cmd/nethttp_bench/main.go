package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	// Start pprof on 6061 for net/http reference server
	go func() {
		fmt.Println("Starting net/http pprof debug server on localhost:6061")
		if err := http.ListenAndServe("localhost:6061", nil); err != nil {
			log.Printf("net/http pprof server failed: %v", err)
		}
	}()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	heavyPayload := []byte(strings.Repeat("junk_data_", 10000))
	http.HandleFunc("/heavy", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(heavyPayload)
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/users/")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello, %s!\n", name)))
	})

	http.HandleFunc("/payload/", func(w http.ResponseWriter, r *http.Request) {
		sizeStr := strings.TrimPrefix(r.URL.Path, "/payload/")
		size, _ := strconv.Atoi(sizeStr)
		if size < 0 || size > 10_000_000 {
			http.Error(w, "Size must be between 0 and 10MB", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(make([]byte, size))
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	http.HandleFunc("/debug/gc", func(w http.ResponseWriter, r *http.Request) {
		runtime.GC()
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GC complete"))
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	fmt.Println("Starting net/http reference server...")
	if os.Getenv("TLS_MODE") == "1" {
		fmt.Println("Listening on :8444 with TLS")
		if err := http.ListenAndServeTLS(":8444", "certs/cert.pem", "certs/key.pem", nil); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Listening on :8081...")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatal(err)
		}
	}
}
