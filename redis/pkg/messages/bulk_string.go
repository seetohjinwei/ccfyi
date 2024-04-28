package messages

import (
	"errors"
	"strings"
)

// A null bulk string is represented by `nil`.
type BulkString struct {
	len uint
	str string
}

func (r *BulkString) Serialise() string {
	return r.str
}

func deserialiseBulkString(message string) (*BulkString, string, error) {
	// $<length>\r\n<data>\r\n

	integerLength, message, err := deserialiseInteger(message)
	if err != nil {
		return nil, "", errors.New("bulk string must contain a valid integer length")
	}

	if integerLength.value == -1 {
		// null bulk string is represented by a nil object
		return nil, "", nil
	} else if integerLength.value < 0 {
		return nil, "", errors.New("bulk string length must be either -1 (null string) or non-negative")
	}

	length := uint(integerLength.value)

	ret := message[:length]

	if !strings.HasPrefix(message[length:], CRLF) {
		return nil, "", errors.New("bulk string does not have CRLF after the specified length")
	}

	remaining := message[length+2:]

	return &BulkString{length, ret}, remaining, nil
}
