package store

import (
	"testing"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
)

func TestList(t *testing.T) {
	list := NewList()

	Equal(t, V(list.LPush([]string{"a", "b", "c"})), V(int64(3), true))
	Equal(t, V(list.LLen()), V(int64(3), true))
	Equal(t, V(list.LRange(0, -1)), V([]string{"c", "b", "a"}, true))
}

func TestListLRange(t *testing.T) {
	list := NewList()
	list.RPush([]string{"a", "b", "c"})

	Equal(t, V(list.LRange(0, 0)), V([]string{"a"}, true))
	Equal(t, V(list.LRange(-3, 2)), V([]string{"a", "b", "c"}, true))
	Equal(t, V(list.LRange(-100, 100)), V([]string{"a", "b", "c"}, true))
	Equal(t, V(list.LRange(5, 10)), V([]string{}, true))
	Equal(t, V(list.LRange(0, 1)), V([]string{"a", "b"}, true))
	Equal(t, V(list.LRange(-4, 0)), V([]string{"a"}, true))
	Equal(t, V(list.LRange(-100, 1)), V([]string{"a", "b"}, true))
	Equal(t, V(list.LRange(-1, -2)), V([]string{}, true))
	Equal(t, V(list.LRange(-1, -1)), V([]string{"c"}, true))
	Equal(t, V(list.LRange(1, 0)), V([]string{}, true))
}

func TestListSerialise(t *testing.T) {
	l1 := NewList()
	Equal(t, V(l1.RPush([]string{"", "1", "2"})), V(int64(3), true))

	NoPanic(t, func() {
		l1.Serialise()
	})
}
