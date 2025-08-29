package request

import (
	"bytes"
	"fmt"
	"http-server/internal/headers"
	"io"
	"strconv"
	"unicode"
)

type parserState string

const (
	StateInit           parserState = "init"
	StateDone           parserState = "done"
	StateParsingHeaders parserState = "headers"
	StateParsingBody    parserState = "body"
)

var ERROR_MALFORMED_REQUEST = fmt.Errorf("malformed request")
var ERROR_BAD_REQUEST_LINE = fmt.Errorf("malformed request line recieved")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version used")

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: *headers.NewHeaders(),
		Body:    make([]byte, 0),
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

func (r *Request) hasBody() bool {
	return r.Headers.Get("content-length") != ""
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for read != len(data) {
		partialData := data[read:]

		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(partialData)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateParsingHeaders
		case StateParsingHeaders:
			n, done, err := r.Headers.Parse(partialData)

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				if r.hasBody() {
					r.state = StateParsingBody
				} else {
					r.state = StateDone
				}
			}

		case StateParsingBody:
			lengthStr := r.Headers.Get("content-length")

			contentLenght, err := strconv.Atoi(lengthStr)

			if err != nil {
				return 0, err
			}

			if contentLenght == 0 {
				r.state = StateDone
				break outer
			}

			remaining := min(contentLenght-len(r.Body), len(partialData))

			r.Body = append(r.Body, partialData[:remaining]...)

			read += remaining

			if len(r.Body) == contentLenght {
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
		n, err := reader.Read(buf[buffLen:])
		if err != nil {
			if err == io.EOF {
				if request.done() {
					return request, nil
				} else {
					return nil, ERROR_MALFORMED_REQUEST
				}
			}
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
