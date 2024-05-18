package store

import (
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/rdb/encoding"
)

type Item interface {
	ValueType() encoding.ValueType
	Serialise() []byte
	Get() (string, bool)
	Incr() (int64, bool)
	Decr() (int64, bool)
	LPush(strs []string) (int64, bool)
	RPush(strs []string) (int64, bool)
	LRange(start, stop int) ([]string, bool)
	LLen() (int64, bool)
}
