package store

import (
	"testing"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
)

func TestString(t *testing.T) {
	s := NewString("0")

	Equal(t, V(s.Incr()), V(int64(1), true))
	Equal(t, V(s.Incr()), V(int64(2), true))
	Equal(t, V(s.Incr()), V(int64(3), true))

	Equal(t, V(s.Get()), V("3", true))

	Equal(t, V(s.Decr()), V(int64(2), true))
	Equal(t, V(s.Decr()), V(int64(1), true))

	Equal(t, V(s.Get()), V("1", true))
}
