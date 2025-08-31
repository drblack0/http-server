package main

import (
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"http-server/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		if req.RequestLine.RequestTarget == "/yourproblem" {
			responseString := `<!DOCTYPE html>
					<html>
					<head>
						<title>400 Bad Request</title>
					</head>
					<body>
						<h1>Bad Request</h1>
						<p>Your request honestly kinda sucked.</p>
					</body>
					</html>`

			w.WriteStatusLine(response.BadRequest)
			h.Replace("content-length", fmt.Sprintf("%d", len(responseString)))
			h.Replace("content-type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody([]byte(responseString))

		} else if req.RequestLine.RequestTarget == "/myproblem" {
			responseString := `<html>
					<head>
						<title>500 Internal Server Error</title>
					</head>
					<body>
						<h1>Internal Server Error</h1>
						<p>Okay, you know what? This one is on me.</p>
					</body>
					</html>`

			w.WriteStatusLine(response.InternalServerError)
			h.Replace("content-length", fmt.Sprintf("%d", len(responseString)))
			h.Replace("content-type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody([]byte(responseString))

		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream/") {
			target := "https://httpbin.org/" + req.RequestLine.RequestTarget[len("httpbin/"):]
			resp, err := http.Get(target)

			if err != nil {
				w.WriteStatusLine(response.InternalServerError)
				w.WriteHeaders(h)
				return
			}

			h.Delete("content-length")
			h.Set("transfer-encoding", "chunnked")

			for {
				data := make([]byte, 32)

				n, err := resp.Body.Read(data)

				if err != nil {
					break
				}

				w.WriteBody([]byte(fmt.Sprintf("%x", n)))
				w.WriteBody(data)
				w.WriteBody([]byte("\r\n"))
			}

			w.WriteBody([]byte("0\r\n\r\n"))
		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
			videoData, err := os.Open("/mnt/c/Users/karti/projects/http-server/assets/vim.mp4")

			if err != nil {
				w.WriteStatusLine(response.InternalServerError)
				return
			}

			h.Replace("content-type", "video/mp4")
			totalLength := 0
			totalData := make([]byte, 0)
			for {
				data := make([]byte, 32)

				n, err := videoData.Read(data)

				if err != nil {
					break
				}

				totalData = append(totalData, data...)
				totalLength += n
			}

			h.Replace("content-length", fmt.Sprintf("%d", totalLength))
			w.WriteHeaders(h)
			w.WriteBody(totalData)
		} else {
			responseString := `<html>
					<head>
						<title>200 OK</title>
					</head>
					<body>
						<h1>Success!</h1>
						<p>Your request was an absolute banger.</p>
					</body>
					</html>`
			w.WriteStatusLine(response.SuccessfulRequest)
			h.Replace("content-length", fmt.Sprintf("%d", len(responseString)))
			h.Replace("content-type", "text/html")
			w.WriteHeaders(h)
			w.WriteBody([]byte(responseString))
		}
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
