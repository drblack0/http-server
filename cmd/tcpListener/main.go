package main

import (
	"fmt"
	"http-server/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal("error starting server:", err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("channel closed")
			break
		}

		rl, err := request.RequestFromReader(conn)

		fmt.Println("Request line: ")
		fmt.Printf("- Method: %s\n", rl.RequestLine.Method)
		fmt.Printf("- Target: %s\n", rl.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", rl.RequestLine.HttpVersion)
	}
}
