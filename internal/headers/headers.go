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
	Headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{Headers: map[string]string{}}
}

func (h *Headers) Get(key string) string {
	return h.Headers[strings.ToLower(key)]
}

func (h *Headers) Set(key string, value string) {
	if v, ok := h.Headers[strings.ToLower(key)]; ok {
		h.Headers[strings.ToLower(key)] = fmt.Sprintf("%s,%s", v, value)
		return
	} else {
		h.Headers[strings.ToLower(key)] = value
	}
}

func (h *Headers) ForEach(f func(key string, value string)) {
	for k, v := range h.Headers {
		f(k, v)
	}
}

func (h *Headers) Replace(k, v string) {
	if _, ok := h.Headers[strings.ToLower(k)]; ok {
		h.Headers[strings.ToLower(k)] = v
	} else {
		h.Headers[strings.ToLower(k)] = v
	}

	fmt.Println(h.Headers)
}

func validateFieldName(key string) bool {
	for _, s := range key {
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
		read += idx + len(clrf)
	}

	return read, done, nil
}
