package rdb

import (
	"bytes"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store"
)

func saveHeader(buf bytes.Buffer) {
	buf.WriteString(magicString)
}

func saveValue(buf bytes.Buffer, k string, v *store.Value) {
	// TODO: utilise string encoding
}

func Save(values map[string]*store.Value) []byte {
	buf := bytes.Buffer{}

	saveHeader(buf)

	for k, v := range values {
		saveValue(buf, k, v)
	}

	return buf.Bytes()
}
