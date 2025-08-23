package main

import (
	"fmt"
	"http-server/internal/headers"
)

func main() {
	h := headers.Headers{}
	h.Parse([]byte("Host: localhost:42069    \r\n\r\n"))

	fmt.Println(h)
}
