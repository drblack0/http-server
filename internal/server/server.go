package server

import (
	"bytes"
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
	"net"
)

type Server struct {
	closed   bool
	listener net.Listener
}

type HandleError struct {
	StatusCode response.StatusCode
	Msg        string
}

type Handler func(w io.Writer, req *request.Request) *HandleError

func newServer() *Server {
	return &Server{
		closed: false,
	}
}

func Serve(port uint16, handleFunc Handler) (*Server, error) {
	server := newServer()
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}
	server.listener = listner
	go server.listen(handleFunc)
	return server, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) handle(conn io.ReadWriteCloser, handleFunc Handler) {
	defer conn.Close()
	responseBody := bytes.Buffer{}

	req, err := request.RequestFromReader(conn)

	if err != nil {
		writeErrorMsg(conn, &HandleError{StatusCode: response.BadRequest, Msg: "Malformed request\n"})
		return
	}

	handleFuncError := handleFunc(&responseBody, req)

	if handleFuncError != nil {
		writeErrorMsg(conn, handleFuncError)
		return
	}

	statusLineError := response.WriteStatusLine(conn, response.SuccessfulRequest)

	if statusLineError != nil {
		return
	}

	headersError := response.WriteHeaders(conn, response.GetDefaultHeaders(len(responseBody.String())))

	if headersError != nil {
		return
	}

	conn.Write(responseBody.Bytes())
}

func (s *Server) listen(handleFunc Handler) {
	for {
		conn, err := s.listener.Accept()

		if s.closed {
			return
		}

		if err != nil {
			return
		}

		go s.handle(conn, handleFunc)
	}
}

func writeErrorMsg(w io.Writer, err *HandleError) {
	statusLineError := response.WriteStatusLine(w, err.StatusCode)

	if statusLineError != nil {
		return
	}

	headersErr := response.WriteHeaders(w, response.GetDefaultHeaders(len(err.Msg)))

	if headersErr != nil {
		return
	}

	w.Write([]byte(err.Msg))
}
