package store

import (
	"sync"

	"github.com/gammazero/deque"
	"github.com/rs/zerolog/log"
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

const listEncoding = 1

func (s *List) Serialise() string {

	// TODO:

	return ""
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
