package messages

import (
	"errors"
	"fmt"
	"strings"
)

type SimpleString struct {
	str string
}

func (r *SimpleString) Serialise() string {
	// +OK\r\n

	return fmt.Sprintf("+%s\r\n", r.str)
}

func deserialiseSimpleString(message string) (*SimpleString, string, error) {
	// +OK\r\n

	curr := message
	for i := 0; i < len(message); i++ {
		if strings.HasPrefix(curr, CRLF) {
			// found the end
			ret := message[:i]
			remaining := curr[2:]
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
