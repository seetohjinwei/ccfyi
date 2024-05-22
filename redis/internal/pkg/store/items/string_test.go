package items

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

func TestStringSerialise(t *testing.T) {
	s1 := NewString("0")
	NoPanic(t, func() {
		s1.Serialise()
	})
}
