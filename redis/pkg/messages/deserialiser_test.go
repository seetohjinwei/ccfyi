package messages

import (
	"reflect"
	"testing"
)

func TestSerialise(t *testing.T) {
	tests := []struct {
		name     string
		message  Message
		expected string
	}{
		{"simple_string_1", &SimpleString{"OK"}, "+OK\r\n"},
		{"simple_string_2", &SimpleString{"hello world"}, "+hello world\r\n"},
		{"simple_string_empty", &SimpleString{""}, "+\r\n"},

		{"bulk_string_1", (*BulkString)(nil), "$-1\r\n"},
		{"bulk_string_2", &BulkString{0, ""}, "$0\r\n\r\n"},

		{
			"array_1",
			&Array{
				1,
				[]Message{
					&BulkString{4, "ping"},
				},
			},
			"*1\r\n$4\r\nping\r\n",
		},
		{
			"array_2",
			&Array{
				2,
				[]Message{
					&BulkString{4, "echo"},
					&BulkString{11, "hello world"},
				},
			},
			"*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n",
		},
		{
			"array_3",
			&Array{
				2,
				[]Message{
					&BulkString{3, "get"},
					&BulkString{3, "key"},
				},
			},
			"*2\r\n$3\r\nget\r\n$3\r\nkey\r\n",
		},

		{"error_1", &Error{"Error message"}, "-Error message\r\n"},

		{"integer_1", &Integer{420}, ":420\r\n"},
		{"integer_2", &Integer{-420}, ":-420\r\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.message.Serialise()
			if !reflect.DeepEqual(test.expected, actual) {
				t.Errorf("expected %+v, but got %+v", test.expected, actual)
			}
		})
	}
}

func TestDeserialise(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected Message
		hasError bool
	}{
		{"simple_string_1", "+OK\r\n", &SimpleString{"OK"}, false},
		{"simple_string_2", "+hello world\r\n", &SimpleString{"hello world"}, false},
		{"simple_string_empty", "+\r\n", &SimpleString{""}, false},

		{"bulk_string_1", "$-1\r\n", (*BulkString)(nil), false},
		{"bulk_string_2", "$0\r\n\r\n", &BulkString{0, ""}, false},

		{
			"array_1",
			"*1\r\n$4\r\nping\r\n",
			&Array{
				1,
				[]Message{
					&BulkString{4, "ping"},
				},
			},
			false,
		},
		{
			"array_2",
			"*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n",
			&Array{
				2,
				[]Message{
					&BulkString{4, "echo"},
					&BulkString{11, "hello world"},
				},
			},
			false,
		},
		{
			"array_3",
			"*2\r\n$3\r\nget\r\n$3\r\nkey\r\n",
			&Array{
				2,
				[]Message{
					&BulkString{3, "get"},
					&BulkString{3, "key"},
				},
			},
			false,
		},

		{"error_1", "-Error message\r\n", &Error{"Error message"}, false},

		{"integer_1", ":420\r\n", &Integer{420}, false},
		{"integer_2", ":+420\r\n", &Integer{420}, false},
		{"integer_3", ":-420\r\n", &Integer{-420}, false},

		{"invalid_1", "x_invalid_first_byte\r\n", nil, true},
		{"invalid_2", ":+-420\r\n", nil, true},
		{"invalid_3_wrong_len", "$0\r\nwronglen\r\n", nil, true},
		{"invalid_4_wrong_len", "*2\r\n+str\r\n", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Deserialise(test.message)
			if test.hasError && err == nil {
				t.Errorf("expected err, but succeeded with %+v", actual)
			} else if !test.hasError && err != nil {
				t.Errorf("expected no err, but got %+v", err)
			}

			if !reflect.DeepEqual(test.expected, actual) {
				t.Errorf("expected %+v, but got %+v", test.expected, actual)
			}
		})
	}
}
