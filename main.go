package main

import (
	"http-server/internal/request"
	"strings"
)

func main() {
	request.RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
}
