package router

import (
	"reflect"
	"testing"

	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func isError(message string) bool {
	msg, err := messages.Deserialise(message)
	if err != nil {
		// not a valid message => not an error
		return false
	}

	_, ok := msg.(*messages.Error)
	return ok
}

func TestHandle(t *testing.T) {
	tests := []struct {
		name     string
		request  string
		expected string
		hasError bool
	}{
		{"invalid_message_1", "", "", true},
		{"invalid_message_2", "xd", "", true},
		{"not_array_1", "+OK\r\n", "", true},
	}

	r := Router{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := r.handle(test.request)
			hasError := isError(actual)

			if test.hasError && !hasError {
				t.Errorf("expected err, but succeeded with %+v", actual)
			} else if !test.hasError && hasError {
				t.Errorf("expected no err, but got %+v", actual)
			} else if !test.hasError && !reflect.DeepEqual(test.expected, actual) {
				t.Errorf("expected %+v, but got %+v", test.expected, actual)
			}
		})
	}
}
