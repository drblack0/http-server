package request

import (
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal("error while reading request: ", err)
	}

	requestLine, err := parseRequestLine(request)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(r []byte) (*RequestLine, error) {

	// we are just going to be using the first part which is the request line
	parts := strings.Split(string(r), "\r\n")

	requestLineParts := strings.Split(parts[0], " ")

	if len(requestLineParts) != 3 {
		log.Println("invalid request line, number of args more than 3")
		return nil, fmt.Errorf("invalid request line, number of args more than 3")
	}

	if !isUpper(requestLineParts[0]) {
		log.Println("http method not proper, should be in caps, recieved: ", requestLineParts[0])
		return nil, fmt.Errorf("http method not proper, should be in caps, recieved: %s", requestLineParts[0])
	}

	if !strings.HasPrefix(requestLineParts[1], "/") {
		log.Println("the target should start with /")
		return nil, fmt.Errorf("the target should start with /, recieved: %s", requestLineParts[1])
	}

	if strings.Split(requestLineParts[2], "/")[1] != "1.1" {
		log.Println("only http 1.1 is supported right now and recieved: ", requestLineParts[1])
		return nil, fmt.Errorf("only http 1.1 is supported right now and recieved: %s", requestLineParts[1])
	}

	fmt.Println("here si the request line parts: ", requestLineParts)

	return &RequestLine{
		HttpVersion:   strings.Split(requestLineParts[2], "/")[1],
		RequestTarget: requestLineParts[1],
		Method:        requestLineParts[0],
	}, nil
}

func isUpper(s string) bool {
	for _, bit := range s {
		if !unicode.IsUpper(bit) && unicode.IsLetter(bit) {
			return false
		}
	}
	return true
}
