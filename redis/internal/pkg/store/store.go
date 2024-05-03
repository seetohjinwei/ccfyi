package store

import (
	"sync"

	"github.com/seetohjinwei/ccfyi/redis/pkg/delay"
)

type Store struct {
	mu     sync.RWMutex
	values map[string]*Value
}

func New() *Store {
	ret := &Store{
		mu:     sync.RWMutex{},
		values: make(map[string]*Value),
	}
	return ret
}

func (s *Store) Get(key string) (Item, bool) {
	// allows some race condition, but no data races
	s.mu.RLock()
	value := s.values[key]
	s.mu.RUnlock()
	item, ok := value.Item()
	if !ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.values, key)
		return nil, false
	}

	return item, ok
}

func (s *Store) Set(key string, item Item) error {
	return s.SetDelay(key, item, nil)
}

func (s *Store) SetDelay(key string, item Item, delay *delay.Delay) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = NewValue(item, delay)

	return nil
}

var (
	store *Store
	once  sync.Once
	mu    sync.Mutex
)

func GetSingleton() *Store {
	once.Do(func() {
		ResetSingleton()
	})

	return store
}

// ResetSingleton should only be used in tests.
func ResetSingleton() *Store {
	mu.Lock()
	defer mu.Unlock()

	store = New()
	return store
}
