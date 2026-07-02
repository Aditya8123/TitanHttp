package main

import (
	"fmt"
	"os"

	"github.com/Aditya8123/TitanHttp/internal/server"
)

func main() {
	fmt.Println("Initializing TitanHTTP Server...")

	srv := server.NewServer(":8080")
	if err := srv.Start(); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}
