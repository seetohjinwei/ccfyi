package router

import (
	"testing"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
)

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
			actual, ok := r.Handle(test.request)

			if test.hasError {
				IsFalse(t, ok, "actual=%q", actual)
			} else {
				IsTrue(t, ok, "actual=%q", actual)
				EqualO(t, test.expected, actual)
			}
		})
	}
}
