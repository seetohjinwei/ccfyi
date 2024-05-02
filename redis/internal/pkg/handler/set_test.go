package handler

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/pkg/messages"
)

func TestParseSetArguments(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		expected setArgs
		hasError bool
	}{
		{"simple", strings.Split("SET k v", " "), setArgs{false, false, false, time.Time{}}, false},
		{"simple", strings.Split("SET k v NX", " "), setArgs{true, false, false, time.Time{}}, false},
		{"simple", strings.Split("SET k v XX", " "), setArgs{false, true, false, time.Time{}}, false},
		{"simple", strings.Split("SET k v GET", " "), setArgs{false, false, true, time.Time{}}, false},

		// would have to mock time.Now()...
		// {"simple", strings.Split("SET k v EX 10", " "), setArgs{setNone, false, time.Time{}}, false},
		// {"simple", strings.Split("SET k v PX 10", " "), setArgs{setNone, false, time.Time{}}, false},
		{"simple", strings.Split("SET k v EXAT 1714662500", " "), setArgs{false, false, false, time.Date(2024, time.May, 2, 15, 8, 20, 0, time.UTC)}, false},
		{"simple", strings.Split("SET k v PXAT 1714662500000", " "), setArgs{false, false, false, time.Date(2024, time.May, 2, 15, 8, 20, 0, time.UTC)}, false},
		{"complex", strings.Split("SET k v GET XX PXAT 1714662500000", " "), setArgs{false, true, true, time.Date(2024, time.May, 2, 15, 8, 20, 0, time.UTC)}, false},
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

func assertSet(t *testing.T, command string, expected messages.Message, expectedOk bool) {
	commands := strings.Split(command, " ")

	res, ok := Set(commands)
	if ok != expectedOk {
		t.Errorf("expected %v, but got %v", expectedOk, ok)
	}

	actual, err := messages.Deserialise(res)
	if err != nil {
		t.Errorf("expected a valid result, but got err %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %+v, but got %+v", expected, actual)
	}
}

// Tests the interactions between multiple sets.
func TestSet(t *testing.T) {
	t.Run("no options, NX, XX", func(t *testing.T) {
		store.ResetSingleton()

		assertSet(t, "SET k v", messages.NewSimpleString("OK"), true)
		assertSet(t, "SET k v", messages.NewSimpleString("OK"), true)

		assertSet(t, "SET k v NX", messages.NewNullBulkString(), true)
		// check that repeated sets still won't set
		assertSet(t, "SET k v NX", messages.NewNullBulkString(), true)
		assertSet(t, "SET k v NX", messages.NewNullBulkString(), true)
		assertSet(t, "SET k_nx v NX", messages.NewSimpleString("OK"), true)

		assertSet(t, "SET k v XX", messages.NewSimpleString("OK"), true)
		assertSet(t, "SET k_xx v XX", messages.NewNullBulkString(), true)
		// check that repeated sets still won't set
		assertSet(t, "SET k_xx v XX", messages.NewNullBulkString(), true)
		assertSet(t, "SET k_xx v XX", messages.NewNullBulkString(), true)
	})

	t.Run("tests GET", func(t *testing.T) {
		store.ResetSingleton()

		assertSet(t, "SET k v1 GET", messages.NewNullBulkString(), true)
		assertSet(t, "SET k v2 GET", messages.NewBulkString("v1"), true)
		assertSet(t, "SET k v3 GET", messages.NewBulkString("v2"), true)
		assertSet(t, "SET k v4 GET", messages.NewBulkString("v3"), true)
	})

	t.Run("tests NX/XX + GET", func(t *testing.T) {
		store.ResetSingleton()

		assertSet(t, "SET k1 v1", messages.NewSimpleString("OK"), true)
		// shouldn't set
		assertSet(t, "SET k1 v2 NX GET", messages.NewNullBulkString(), true)
		assertSet(t, "SET k1 v3 GET", messages.NewBulkString("v1"), true)

		assertSet(t, "SET k2 v1 XX GET", messages.NewNullBulkString(), true)
		assertSet(t, "SET k2 v2 GET", messages.NewNullBulkString(), true)
	})
}
