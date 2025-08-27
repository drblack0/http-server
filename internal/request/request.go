package request

import (
	"bytes"
	"fmt"
	"http-server/internal/headers"
	"io"
	"unicode"
)

type parserState string

const (
	StateInit           parserState = "init"
	StateDone           parserState = "done"
	StateParsingHeaders parserState = "parsingHeaders"
)

var ERROR_MALFORMED_REQUEST = fmt.Errorf("malformed request")
var ERROR_BAD_REQUEST_LINE = fmt.Errorf("malformed request line recieved")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version used")

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: *headers.NewHeaders(),
	}
}

func (r *Request) done() bool {
	return r.state == StateDone
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	totalBytesRead := 0
outer:
	for totalBytesRead != len(data) {
		partialData := data[totalBytesRead:]

		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(partialData)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break
			}

			r.RequestLine = *rl
			read += n
			totalBytesRead += n

			r.state = StateParsingHeaders
		case StateParsingHeaders:
			n, done, err := r.Headers.Parse(partialData)

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break
			}

			read += n
			totalBytesRead += n

			if done {
				r.state = StateDone
			}
		case StateDone:
			break outer
		}
	}
	return read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	buffLen := 0

	for !request.done() {
		fmt.Println("here in the for loop")
		n, err := reader.Read(buf[buffLen:])
		if err != nil {
			return nil, err
		}

		buffLen += n

		readN, err := request.parse(buf[:buffLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:buffLen])
		buffLen -= readN
	}

	return request, nil
}

func isUpper(s string) bool {
	for _, bit := range s {
		if !unicode.IsUpper(bit) && unicode.IsLetter(bit) {
			return false
		}
	}
	return true
}

func parseRequestLine(r []byte) (*RequestLine, int, error) {
	idx := bytes.Index(r, []byte("\r\n"))

	if idx == -1 {
		return nil, 0, nil
	}
	read := idx + len([]byte("\r\n"))
	parts := bytes.Split(r, []byte("\r\n"))

	requestLineParts := bytes.Split(parts[0], []byte(" "))

	if len(requestLineParts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST
	}

	if !isUpper(string(requestLineParts[0])) {
		return nil, 0, ERROR_MALFORMED_REQUEST
	}

	httpParts := bytes.Split(requestLineParts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_REQUEST
	}

	rl := &RequestLine{
		HttpVersion:   string(httpParts[1]),
		RequestTarget: string(requestLineParts[1]),
		Method:        string(requestLineParts[0]),
	}

	return rl, read, nil
}
