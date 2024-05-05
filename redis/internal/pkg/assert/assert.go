package assert

import (
	"reflect"
	"testing"
)

func V(v ...interface{}) []interface{} {
	return v
}

func Equal(t testing.TB, actual, expected []interface{}) {
	t.Helper()

	if len(expected) != len(actual) {
		// something wrong with the test case
		t.Errorf("expected len(expected) == %d, but got len(actual) == %d", len(expected), len(actual))
	}

	for i := 0; i < len(expected); i++ {
		e := expected[i]
		a := actual[i]

		if !reflect.DeepEqual(e, a) {
			if reflect.TypeOf(e) != reflect.TypeOf(a) {
				t.Errorf("expected %+v (%v), but got %+v (%v) - type mismatch!", e, reflect.TypeOf(e), a, reflect.TypeOf(a))
			}
			t.Errorf("expected %+v, but got %+v", e, a)
		}
	}
}

func HasError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		return
	}

	t.Errorf("expected err, but got no error")
}

func NoError(t testing.TB, err error) {
	t.Helper()

	if err == nil {
		return
	}

	t.Errorf("expected no err, but got %v", err)
}
