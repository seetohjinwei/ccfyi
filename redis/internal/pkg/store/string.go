package store

import (
	"strconv"
	"sync"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/rdb/encoding"
)

type stringType int

const (
	stringUnknown stringType = iota
	stringString
	stringInteger
)

type String struct {
	mu         sync.RWMutex
	str        string
	integer    int64
	actualType stringType

	// String is either a string or an integer. (according to redis specs)
	// but "Redis stores integers in their integer representation".

	*AbstractItem
}

func NewString(str string) *String {
	ret := &String{
		mu:         sync.RWMutex{},
		str:        "",
		integer:    int64(0),
		actualType: stringUnknown,
	}

	integer, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		ret.integer = integer
		ret.actualType = stringInteger
	} else {
		ret.str = str
		ret.actualType = stringString
	}

	if ret.actualType == stringUnknown {
		panic("ret.ActualType == stringUnknown")
	}

	return ret
}

const stringEncoding = 0

func (s *String) Serialise() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch s.actualType {
	case stringString:
		return encoding.EncodeString(s.str)
	case stringInteger:
		return encoding.EncodeInteger(s.integer)
	}

	panic("ret.ActualType == stringUnknown")
}

func (s *String) Get() (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch s.actualType {
	case stringString:
		return s.str, true
	case stringInteger:
		return strconv.FormatInt(s.integer, 10), true
	}

	panic("ret.ActualType == stringUnknown")
}

func (s *String) Incr() (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.actualType != stringInteger {
		return 0, false
	}

	s.integer++

	return s.integer, true
}

func (s *String) Decr() (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.actualType != stringInteger {
		return 0, false
	}

	s.integer--

	return s.integer, true
}
