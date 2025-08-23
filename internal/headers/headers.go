package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const clrf = "\r\n"

var MalformedHeadersError = fmt.Errorf("malformed headers")

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(clrf))

	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 0, true, nil
	}

	part := data[:idx]

	colonIdx := bytes.Index(part, []byte(":"))

	if colonIdx == -1 {
		return 0, false, nil
	}

	if data[colonIdx-1] == ' ' {
		return 0, false, MalformedHeadersError
	}

	clrfIdx := bytes.Index(data[idx+1:], []byte(clrf))
	key := data[:colonIdx]
	value := data[colonIdx+1 : idx+clrfIdx - 1]

	h[string(key)] = strings.Trim(string(value), " ")
	n = idx + len(clrf)
	return n, false, nil
}
