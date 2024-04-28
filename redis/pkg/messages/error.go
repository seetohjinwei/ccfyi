package messages

import (
	"errors"
	"strings"
)

type Error struct {
	str string
}

func deserialiseError(message string) (*Error, string, error) {
	curr := message
	for i := 0; i < len(message); i++ {
		if strings.HasPrefix(curr, CRLF) {
			// found the end
			ret := message[:i]
			length := len(curr)
			remaining := curr[:length-2]
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

func (r *Error) Serialise() string {
	return r.str
}
