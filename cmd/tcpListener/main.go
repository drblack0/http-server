package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		sentenceChan := getLinesChannel(conn)

		for value := range sentenceChan {
			fmt.Printf("%s\n", value)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	sentenceChan := make(chan string)
	go func() {
		defer f.Close()
		defer close(sentenceChan)

		msg := make([]byte, 8)
		curr := ""

		for {
			_, err := f.Read(msg)

			if err == io.EOF {
				fmt.Println("read: end")
				return
			}

			parts := strings.Split(string(msg), "\n")

			if len(parts) > 1 {
				curr += parts[0]
				sentenceChan <- curr

				// there is no remove function to use so we are using the indexing to remove the first element
				subParts := parts[1:]
				parts = subParts

				curr = ""
			}

			curr += parts[0]
		}
	}()
	return sentenceChan
}
