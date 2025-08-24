package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

var clrf = []byte("\r\n")

var MalformedHeadersError = fmt.Errorf("malformed headers")

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{headers: map[string]string{}}
}

func (h *Headers) Get(key string) string {
	return h.headers[strings.ToLower(key)]
}

func (h *Headers) Set(key string, value string) {
	h.headers[strings.ToLower(key)] = value
}

func validateFieldName(key string) bool {
	for _, s := range key {
		// todo: put the special character check here
		if !unicode.IsLetter(s) && !unicode.IsNumber(s) && !strings.ContainsRune("!#$%&'*+-.^_`|~", s) {
			return false
		}
	}
	return true
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	key := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(key, []byte(" ")) {
		return "", "", fmt.Errorf("malformed field name")
	}

	if !validateFieldName(string(key)) {
		return "", "", fmt.Errorf("invalid character in field name")
	}

	return string(key), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {

	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], clrf)

		if idx == -1 {
			break
		}

		// Empty line indicating end of headers
		if idx == 0 {
			done = true
			read += len(clrf)
			break
		}

		key, value, err := parseHeader(data[read : read+idx])

		if err != nil {
			return 0, false, err
		}

		h.Set(key, value)
		fmt.Println(h)
		read += idx + len(clrf)
	}

	return read, done, nil
}
