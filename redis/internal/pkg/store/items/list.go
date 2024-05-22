package items

import (
	"sync"

	"github.com/gammazero/deque"
	"github.com/rs/zerolog/log"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/rdb/encoding"
)

type List struct {
	mu   sync.RWMutex
	list *deque.Deque[string]

	*AbstractItem
}

func NewList() *List {
	ret := &List{
		mu:   sync.RWMutex{},
		list: deque.New[string](),
	}
	return ret
}

func (l *List) ValueType() encoding.ValueType {
	return encoding.ValueList
}

func (l *List) Serialise() []byte {
	l.mu.RLock()
	defer l.mu.RUnlock()

	list := make([]string, l.list.Len())
	for i := 0; i < l.list.Len(); i++ {
		list[i] = l.list.At(i)
	}

	return encoding.EncodeList(list)
}

func DeserialiseList(b []byte) (*List, []byte, error) {
	strs, remaining, err := encoding.DecodeList(b)
	if err != nil {
		return nil, b, err
	}
	list := NewList()
	list.RPush(strs)

	return list, remaining, nil
}

func (l *List) LPush(strs []string) (int64, bool) {
	for _, s := range strs {
		l.list.PushFront(s)
	}
	return int64(l.list.Len()), true
}

func (l *List) RPush(strs []string) (int64, bool) {
	for _, s := range strs {
		l.list.PushBack(s)
	}
	return int64(l.list.Len()), true
}

func (l *List) LRange(start, stop int) ([]string, bool) {
	// handle negative indexes
	if start < 0 {
		start = l.list.Len() + start
	}
	if stop < 0 {
		stop = l.list.Len() + stop
	}
	start = max(0, start)
	stop = min(l.list.Len()-1, stop)

	if start > stop {
		return []string{}, true
	}

	// capacity == final length
	log.Debug().Int("start", start).Int("stop", stop).Int("capacity", stop-start+1).Msg("LRange return slice")
	ret := make([]string, 0, stop-start+1)
	// both start & stop are INCLUSIVE
	for i := start; i <= stop; i++ {
		ret = append(ret, l.list.At(i))
	}

	return ret, true
}

func (l *List) LLen() (int64, bool) {
	return int64(l.list.Len()), true
}

func (l *List) Equal(other any) bool {
	o, ok := other.(*List)
	if !ok {
		return false
	}

	if l == nil || o == nil {
		return (l == nil) && (o == nil)
	}

	if l.list.Len() != o.list.Len() {
		return false
	}
	for i := 0; i < l.list.Len(); i++ {
		if l.list.At(i) != o.list.At(i) {
			return false
		}
	}

	return true
}

type ListBuilder struct {
	*List
}

func NewListBuilder() *ListBuilder {
	return &ListBuilder{
		NewList(),
	}
}

func (b *ListBuilder) Add(strs []string) *ListBuilder {
	b.List.RPush(strs)
	return b
}

func (b *ListBuilder) Build() *List {
	return b.List
}
