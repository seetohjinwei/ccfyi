package rdb

import (
	"testing"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
)

func TestSave(t *testing.T) {
	// Equal(t, V(Save(nil)), V([]byte(nil)))
	// buf := Save(nil)
	// t.Logf(string(buf))
	// t.Fail()
}

func TestLoad(t *testing.T) {
	Equal(t, V(Load([]byte{})), V(map[string]*store.Value{}))
}
