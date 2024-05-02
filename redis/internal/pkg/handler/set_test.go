package handler

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseSetArguments(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		expected setArgs
		hasError bool
	}{
		{"simple", strings.Split("SET k v", " "), setArgs{setNone, false, time.Time{}}, false},
		{"simple", strings.Split("SET k v NX", " "), setArgs{setNX, false, time.Time{}}, false},
		{"simple", strings.Split("SET k v XX", " "), setArgs{setXX, false, time.Time{}}, false},
		{"simple", strings.Split("SET k v GET", " "), setArgs{setNone, true, time.Time{}}, false},

		// would have to mock time.Now()...
		// {"simple", strings.Split("SET k v EX 10", " "), setArgs{setNone, false, time.Time{}}, false},
		// {"simple", strings.Split("SET k v PX 10", " "), setArgs{setNone, false, time.Time{}}, false},
		{"simple", strings.Split("SET k v EXAT 1714662500", " "), setArgs{setNone, false, time.Date(2024, time.May, 2, 15, 8, 20, 0, time.UTC)}, false},
		{"simple", strings.Split("SET k v PXAT 1714662500000", " "), setArgs{setNone, false, time.Date(2024, time.May, 2, 15, 8, 20, 0, time.UTC)}, false},
		{"complex", strings.Split("SET k v GET XX PXAT 1714662500000", " "), setArgs{setXX, true, time.Date(2024, time.May, 2, 15, 8, 20, 0, time.UTC)}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseSetArguments(test.commands)

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
