package main

import (
	"http-server/internal/request"
	"http-server/internal/response"
	"http-server/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandleError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandleError{StatusCode: response.BadRequest, Msg: "Your problem is not my problem\n"}
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandleError{StatusCode: response.InternalServerError, Msg: "Woopsie, my bad\n"}
		} else {
			w.Write([]byte("All good, frfr\n"))
		}
		return nil
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
