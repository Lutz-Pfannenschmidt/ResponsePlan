package htmx

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ParseError struct {
	Body string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Error parsing body: %s", e.Body)
}

// parse the body of a htmx request into a map of key-value pairs
func ParseBody(body io.ReadCloser) (map[string]string, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	newStr := buf.String()

	pairs := strings.Split(newStr, "&")
	result := make(map[string]string)
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, &ParseError{pair}
		}
		result[kv[0]] = kv[1]
	}

	return result, nil
}
