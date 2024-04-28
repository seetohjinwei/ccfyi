package messages

import (
	"errors"
	"strings"
)

type SimpleString struct {
	str string
}

func deserialiseSimpleString(message string) (*SimpleString, string, error) {
	curr := message
	for i := 0; i < len(message); i++ {
		if strings.HasPrefix(curr, CRLF) {
			// found the end
			ret := message[:i]
			length := len(curr)
			remaining := curr[:length-2]
			return &SimpleString{ret}, remaining, nil
		}

		c := curr[0]
		if c == CR || c == LF {
			return nil, "", errors.New("simple string must not contain CR (\\r) or LF (\\n)")
		}
		curr = curr[1:]
	}

	return nil, "", errors.New("simple string must end with CRLF")
}

func (r *SimpleString) Serialise() string {
	return r.str
}
