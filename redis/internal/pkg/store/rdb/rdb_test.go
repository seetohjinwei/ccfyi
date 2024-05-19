package rdb

import (
	"testing"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
)

func TestSave(t *testing.T) {
	buf1 := SaveBuffer{}
	EqualO(t, buf1.Save(nil), []byte("REDISLITEFF"))

	buf2 := SaveBuffer{}
	ret2 := buf2.Save(nil)
	IsTrue(t, ret2 != nil, "%+v", string(ret2))

	buf3 := SaveBuffer{}
	ret3 := buf3.Save(map[string]*store.Value{
		"k1": store.NewValue(store.NewString("v1"), nil),
		"k2": store.NewValue(store.NewString("v2"), nil),
		"k3": store.NewValue(store.NewString("3"), nil),
		"k4": store.NewValue(store.NewList(), nil),
		"k5": store.NewValue(store.NewListBuilder().Add([]string{"1", "2", "3"}).Build(), nil),
	})
	IsTrue(t, ret3 != nil, "%+v", string(ret3))
}

func TestLoad(t *testing.T) {
	buf1 := NewLoadBuffer(nil)
	Equal(t, V(buf1.Load()), V(map[string]*store.Value(nil), AnyError{}))

	buf2 := NewLoadBuffer([]byte("REDISLITEFF"))
	Equal(t, V(buf2.Load()), V(map[string]*store.Value{}, nil))
}

func TestSaveLoad(t *testing.T) {
	contents := []map[string]*store.Value{
		{
			"k1": store.NewValue(store.NewString("v1"), nil),
			"k2": store.NewValue(store.NewString("v2"), nil),
			"k3": store.NewValue(store.NewString("3"), nil),
			"k4": store.NewValue(store.NewList(), nil),
			"k5": store.NewValue(store.NewListBuilder().Add([]string{"1", "2", "3"}).Build(), nil),
		},
		{},
	}

	for _, content := range contents {
		save := SaveBuffer{}
		encoded := save.Save(content)

		load := NewLoadBuffer(encoded)
		actual, err := load.Load()

		EqualO(t, len(actual), len(content))
		for k, v1 := range content {
			v2, ok := actual[k]
			IsTrue(t, ok, "key %s was in expected, but not in actual")
			IsTrue(t, v1.Equal(v2), "expected %+v, but got %+v", v1, v2)
		}

		NoError(t, err)
	}
}
