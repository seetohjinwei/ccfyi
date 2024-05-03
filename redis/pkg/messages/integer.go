package messages

import (
	"errors"
	"fmt"
	"strings"
)

type Integer struct {
	value int64
}

func (r *Integer) Serialise() string {
	// :[<+|->]<value>\r\n

	sign := "" // this is optional
	value := r.value
	if r.value < 0 {
		sign = "-"
		value *= -1
	}

	return fmt.Sprintf(":%s%d\r\n", sign, value)
}

func NewInteger(value int64) *Integer {
	return &Integer{value: value}
}

func deserialiseInteger(message string) (*Integer, string, error) {
	// :[<+|->]<value>\r\n

	if len(message) == 0 {
		return nil, "", errors.New("integer must not be empty")
	}

	isNegative := false
	if message[0] == '-' {
		isNegative = true
		message = message[1:]
	} else if message[0] == '+' {
		message = message[1:]
	}

	value := int64(0)
	curr := message
	for i := 0; i < len(message); i++ {
		if strings.HasPrefix(curr, CRLF) {
			// found the end
			if isNegative {
				value *= -1
			}
			remaining := curr[2:]
			return &Integer{value}, remaining, nil
		}

		c := curr[0]
		if c < '0' || c > '9' {
			return nil, "", errors.New("integer is invalid")
		}
		value *= 10
		value += int64(c - '0')
		curr = curr[1:]
	}

	return nil, "", errors.New("integer must end with CRLF")
}
