package rdb

// import (
// 	"testing"

// 	. "github.com/seetohjinwei/ccfyi/redis/internal/pkg/assert"
// 	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
// )

// TODO: continue here

// func TestSave(t *testing.T) {
// 	buf1 := SaveBuffer{}
// 	EqualO(t, buf1.Save(nil), []byte("REDISLITE"))

// 	buf2 := SaveBuffer{}
// 	ret2 := buf2.Save(nil)                      // TODO: try this with some values
// 	IsTrue(t, ret2 == nil, "%+v", string(ret2)) // change this back to ret2 != nil
// }

// func TestLoad(t *testing.T) {
// 	Equal(t, V(Load([]byte{})), V(map[string]*store.Value{}))
// }
