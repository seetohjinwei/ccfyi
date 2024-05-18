package rdb

import (
	"bytes"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/rdb/encoding"
)

// *partially* follows the file format listed here: https://rdb.fnordig.de/file_format.html
// liberties taken because this is a Go program, some of these optimisations don't matter

const redis = "REDIS"
const version = "LITE" // intentionally not an integer to not collide with redis version numbers
const magicString = redis + version

type SaveBuffer struct {
	bytes.Buffer
}

func (buf *SaveBuffer) saveHeader() {
	buf.WriteString(magicString)
}

func (buf *SaveBuffer) saveValue(k string, v *store.Value) {
	item, ok := v.Item()
	if !ok {
		return
	}

	buf.Write(v.SerialiseExpiry())
	buf.WriteByte(byte(item.ValueType()))
	buf.Write(encoding.EncodeString(k))
	buf.Write(item.Serialise())
}

// Make sure to lock the map!
func (buf *SaveBuffer) Save(values map[string]*store.Value) []byte {
	buf.saveHeader()

	for k, v := range values {
		buf.saveValue(k, v)
	}

	return buf.Bytes()
}

func Load(source []byte) map[string]*store.Value {
	ret := make(map[string]*store.Value)
	return ret
}
