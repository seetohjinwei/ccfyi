package assert

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/logging"
)

func init() {
	logging.Init()
}

func V(v ...any) []any {
	return v
}

func Equal(t testing.TB, actual, expected []any) {
	t.Helper()

	if len(expected) != len(actual) {
		// something wrong with the test case
		t.Fatalf("expected len(expected) == %d, but got len(actual) == %d", len(expected), len(actual))
	}

	for i := 0; i < len(expected); i++ {
		e := expected[i]
		a := actual[i]

		if !cmp.Equal(e, a) {
			if reflect.TypeOf(e) != reflect.TypeOf(a) {
				t.Errorf("expected %+v (%v), but got %+v (%v) - type mismatch!", e, reflect.TypeOf(e), a, reflect.TypeOf(a))
			} else {
				t.Errorf("expected %+v, but got %+v", e, a)
			}
		}
	}
}

func EqualO[T any](t testing.TB, actual, expected T) {
	t.Helper()

	Equal(t, []any{actual}, []any{expected})
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

// TODO: do similar changes as done in router_test.go
func IsTrue(t testing.TB, value bool, format string, args ...any) {
	t.Helper()

	if value {
		return
	}

	if len(format) > 0 {
		fullArgs := make([]any, 0, len(args)+1)
		fullArgs = append(fullArgs, value)
		fullArgs = append(fullArgs, args...)
		t.Errorf("expected true, but got %v ("+format+")", fullArgs...)
	} else {
		t.Errorf("expected true, but got %v", value)
	}
}

func IsFalse(t testing.TB, value bool, format string, args ...any) {
	t.Helper()

	IsTrue(t, !value, format, args...)
}

func NoPanic(t testing.TB, f func()) {
	t.Helper()

	defer func() {
		if ret := recover(); ret != nil {
			t.Errorf("expected no panic, but got %+v", ret)
		}
	}()

	f()
}
