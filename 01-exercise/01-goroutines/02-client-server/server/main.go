package main

import (
	"io"
	"log"
	"net"
	"time"
)

func main() {
	// write server program to handle concurrent client connections.
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		// Adding go here means we can handle multiple concurrent client connections
		go handleConn(conn)
	}
}

// handleConn - utility function
func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, "response from server\n")
		if err != nil {
			return
		}
		time.Sleep(time.Second)
	}
}
