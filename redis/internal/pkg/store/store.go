package store

import (
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	items map[string]Item
}

func New() *Store {
	return &Store{
		mu:    sync.RWMutex{},
		items: make(map[string]Item),
	}
}

func (s *Store) Get(key string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[key]
	return item, ok
}

func (s *Store) Set(key string, value Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = value

	return nil
}
