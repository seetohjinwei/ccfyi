package rdb

import (
	"testing"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/items"
)

func TestSave(t *testing.T) {
	buf1 := SaveBuffer{}
	ret1 := buf1.Save(nil)
	IsTrue(t, ret1 != nil, "%+v", string(ret1))

	buf2 := SaveBuffer{}
	ret2 := buf2.Save(nil)
	IsTrue(t, ret2 != nil, "%+v", string(ret2))

	buf3 := SaveBuffer{}
	ret3 := buf3.Save(map[string]*items.Value{
		"k1": items.NewValue(items.NewString("v1"), nil),
		"k2": items.NewValue(items.NewString("v2"), nil),
		"k3": items.NewValue(items.NewString("3"), nil),
		"k4": items.NewValue(items.NewList(), nil),
		"k5": items.NewValue(items.NewListBuilder().Add([]string{"1", "2", "3"}).Build(), nil),
	})
	IsTrue(t, ret3 != nil, "%+v", string(ret3))
}

func TestLoad(t *testing.T) {
	buf1 := NewLoadBuffer(nil)
	Equal(t, V(buf1.Load()), V(map[string]*items.Value(nil), AnyError{}))
}

func TestSaveLoad(t *testing.T) {
	contents := []map[string]*items.Value{
		{
			"k1": items.NewValue(items.NewString("v1"), nil),
			"k2": items.NewValue(items.NewString("v2"), nil),
			"k3": items.NewValue(items.NewString("3"), nil),
			"k4": items.NewValue(items.NewList(), nil),
			"k5": items.NewValue(items.NewListBuilder().Add([]string{"1", "2", "3"}).Build(), nil),
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
			IsTrue(t, ok, "key %s was in expected, but not in actual", k)
			IsTrue(t, v1.Equal(v2), "expected %+v, but got %+v", v1, v2)
		}

		NoError(t, err)
	}
}
