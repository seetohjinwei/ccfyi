package messages

import (
	"errors"
	"fmt"
	"strings"
)

type Array struct {
	len   uint
	items []Message
}

func (r *Array) Serialise() string {
	// *<number-of-elements>\r\n<element-1>...<element-n>

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("*%d\r\n", r.len))
	for _, item := range r.items {
		builder.WriteString(item.Serialise())
	}

	return builder.String()
}

func deserialiseArray(message string) (*Array, string, error) {
	// *<number-of-elements>\r\n<element-1>...<element-n>

	integerLength, message, err := deserialiseInteger(message)
	if err != nil {
		return nil, "", errors.New("array must contain a valid integer length")
	}
	if integerLength.value < 0 {
		return nil, "", errors.New("array length must be non-negative")
	}

	length := uint(integerLength.value)

	items := make([]Message, length)
	for i := 0; i < len(items); i++ {
		if message == "" {
			return nil, "", errors.New("array length is incorrect")
		}

		var item Message
		item, message, err = deserialise(message)
		if err != nil {
			return nil, "", err
		}
		items[i] = item
	}

	return &Array{length, items}, "", nil
}
