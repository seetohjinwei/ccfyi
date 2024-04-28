package messages

import (
	"errors"
	"fmt"
)

func Deserialise(message string) (Message, error) {
	ret, remaining, err := deserialise(message)
	if err != nil {
		return nil, err
	}
	if remaining != "" {
		fmt.Printf("ret=%v, remaining=%q\n", ret, remaining)
		return nil, errors.New("deserialise invalid pattern")
	}
	return ret, nil
}

func deserialise(message string) (Message, string, error) {
	if len(message) == 0 {
		return nil, message, errors.New("deserialise empty message")
	}

	remaining := message[1:]

	switch message[0] {
	case '+':
		return deserialiseSimpleString(remaining)
	case '-':
		return deserialiseError(remaining)
	case ':':
		return deserialiseInteger(remaining)
	case '$':
		return deserialiseBulkString(remaining)
	case '*':
		return deserialiseArray(remaining)
	default:
		return nil, message, errors.New("deserialise invalid data type")
	}
}
