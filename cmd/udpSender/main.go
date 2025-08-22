package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")

	if err != nil {
		log.Fatal("error while resolving address: ", err)
	}

	// the second arg is the server's own address and the third arg is the remote address (where we want to send the udp msgs)
	conn, err := net.DialUDP("udp", nil, addr)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println(">")
	for {
		input, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("error while reading from stdin: ", err)
		}

		conn.Write([]byte(input))
	}
}
