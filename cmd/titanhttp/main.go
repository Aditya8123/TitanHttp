// Command titanhttp launches the TitanHTTP server.
//
// TitanHTTP is an HTTP/1.1 server implemented from scratch over raw TCP
// sockets, built to master networking, concurrency, and systems
// architecture. The HTTP core deliberately avoids Go's net/http package;
// see phases.md for the project roadmap.
package main

import "fmt"

func main() {
	fmt.Println("TitanHTTP")
}
