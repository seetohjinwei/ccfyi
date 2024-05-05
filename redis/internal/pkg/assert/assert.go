package assert

import (
	"reflect"
	"testing"
)

func V(v ...interface{}) []interface{} {
	return v
}

func Equal(t testing.TB, expected, actual []interface{}) {
	if len(expected) != len(actual) {
		// something wrong with the test case
		t.Errorf("expected len(expected) == %d, but got len(actual) == %d", len(expected), len(actual))
	}

	for i := 0; i < len(expected); i++ {
		e := expected[i]
		a := actual[i]

		if !reflect.DeepEqual(e, a) {
			t.Errorf("expected %+v (%v), but got %+v (%v)", e, reflect.TypeOf(e), a, reflect.TypeOf(a))
		}
	}
}

func NoError(t testing.TB, err error) {
	if err == nil {
		return
	}

	t.Errorf("expected no err, but got %v", err)
}
