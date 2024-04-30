package router

import (
	"reflect"
	"testing"
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

			if test.hasError && ok {
				t.Errorf("expected err, but succeeded with %+v", actual)
			} else if !test.hasError && !ok {
				t.Errorf("expected no err, but got %+v", actual)
			} else if !test.hasError && !reflect.DeepEqual(test.expected, actual) {
				t.Errorf("expected %+v, but got %+v", test.expected, actual)
			}
		})
	}
}
