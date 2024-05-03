package store

import (
	"sync"
)

type Store struct {
	mu     sync.RWMutex
	values map[string]*Value
}

func New() *Store {
	return &Store{
		mu:     sync.RWMutex{},
		values: make(map[string]*Value),
	}
}

func (s *Store) Get(key string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.values[key]
	if !ok {
		return nil, false
	}
	item := value.Item()

	return item, ok
}

func (s *Store) Set(key string, item Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = NewValue(item)

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
