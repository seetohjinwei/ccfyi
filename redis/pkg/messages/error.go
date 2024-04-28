package messages

import (
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	str string
}

func NewError(message string) *Error {
	return &Error{message}
}

func (r *Error) Serialise() string {
	// -Error message\r\n

	return fmt.Sprintf("-%s\r\n", r.str)
}

func deserialiseError(message string) (*Error, string, error) {
	// -Error message\r\n

	curr := message
	for i := 0; i < len(message); i++ {
		if strings.HasPrefix(curr, CRLF) {
			// found the end
			ret := message[:i]
			remaining := curr[2:]
			return &Error{ret}, remaining, nil
		}

		c := curr[0]
		if c == CR || c == LF {
			return nil, "", errors.New("error must not contain CR (\\r) or LF (\\n)")
		}
		curr = curr[1:]
	}

	return nil, "", errors.New("error must end with CRLF")
}
