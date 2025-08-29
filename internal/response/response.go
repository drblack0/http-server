package response

import (
	"fmt"
	"http-server/internal/headers"
	"io"
)

type StatusCode int

const (
	SuccessfulRequest   StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func getStatusCodeLine(statusCode StatusCode) string {
	response := ""
	switch statusCode {
	case SuccessfulRequest:
		response = "HTTP/1.1 200 OK\r\n"
	case BadRequest:
		response = "HTTP/1.1 400 Bad Request\r\n"
	case InternalServerError:
		response = "HTTP/1.1 500 Internal Server Error\r\n"
	}
	return response
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("content-length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("content-type", "text/plain")

	return *headers
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	out := getStatusCodeLine(statusCode)
	_, err := w.Write([]byte(out))

	if err != nil {
		return err
	}

	return nil
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for k, v := range headers.Headers {
		out := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(out))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}
