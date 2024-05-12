package rdb

import (
	"bytes"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/rdb/encoding"
)

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
