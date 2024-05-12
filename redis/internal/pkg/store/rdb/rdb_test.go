package rdb

import (
	"testing"

	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
)

func TestSave(t *testing.T) {
	buf := SaveBuffer{}
	Equal(t, V(buf.Save(nil)), V([]byte("REDISLITE")))
}

func TestLoad(t *testing.T) {
	Equal(t, V(Load([]byte{})), V(map[string]*store.Value{}))
}
