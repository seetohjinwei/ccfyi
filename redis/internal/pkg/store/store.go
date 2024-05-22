package store

import (
	"context"
	"sync"
	"time"

	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/disk"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/items"
	"github.com/seetohjinwei/ccfyi/redis/internal/pkg/store/rdb"
	"github.com/seetohjinwei/ccfyi/redis/pkg/delay"
)

type Store struct {
	mu        sync.RWMutex
	ctx       context.Context
	ctxCancel context.CancelFunc
	values    map[string]*items.Value
	expirySet map[string]struct{}
}

func New() *Store {
	ret := newNoExpiry()

	go ret.activeExpiry(ret.cleanKeys)

	return ret
}

// newNoExpiry should only be used for tests.
func newNoExpiry() *Store {
	ctx, cancelFunc := context.WithCancel(context.Background())

	ret := &Store{
		mu:        sync.RWMutex{},
		ctx:       ctx,
		ctxCancel: cancelFunc,
		values:    make(map[string]*items.Value),
		expirySet: make(map[string]struct{}),
	}

	return ret
}

func (s *Store) Get(key string) (items.Item, bool) {
	// allows some race condition, but no data races
	s.mu.RLock()
	value := s.values[key]
	s.mu.RUnlock()
	item, ok := value.Item()
	if !ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.values, key)
		delete(s.expirySet, key)
		return nil, false
	}

	return item, ok
}

func (s *Store) Set(key string, item items.Item) error {
	return s.SetWithDelay(key, item, nil)
}

func (s *Store) SetWithDelay(key string, item items.Item, delay *delay.Delay) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = items.NewValue(item, delay)
	if delay != nil {
		s.expirySet[key] = struct{}{}
	}

	return nil
}

// LoadFromDisk **overrides** the values in `store` with the values loaded from disk.
// This method should only be called on application startup / recovery!
func (s *Store) LoadFromDisk() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := disk.Load()
	if err != nil {
		return err
	}
	buf := rdb.NewLoadBuffer(data)
	values, err := buf.Load()
	if err != nil {
		return err
	}

	// overrides existing values!
	s.values = values

	return nil
}

func (s *Store) SaveToDisk() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := (&rdb.SaveBuffer{}).Save(s.values)
	return disk.Save(data)
}

// activeExpiry must be run from a goroutine when the store is constructed.
func (s *Store) activeExpiry(cleanFunc func()) {
	// 10 times per second
	t := time.NewTicker(time.Second / 10)
	for {
		select {
		case <-t.C:
			// clean the keys from expiry set
			cleanFunc()
		case <-s.ctx.Done():
			// allows store to be cleaned up
			return
		}
	}
}

const cleanKeysQuantity = 20

func (s *Store) cleanKeys() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		iterations := min(len(s.expirySet), cleanKeysQuantity)
		expiryCount := 0

		// tests 20 random keys from the set of keys with an associated expiry
		for i := 0; i < iterations; i++ {
			var key string
			for k := range s.expirySet {
				key = k
				break
			}
			// `key` is now a random key from the expiry set

			value, ok := s.values[key]
			if !ok {
				// key has been removed in a previous iteration
				continue
			}
			if value.HasExpired() {
				expiryCount++
				delete(s.values, key)
				delete(s.expirySet, key)
			}
		}

		if expiryCount > (iterations / 4) {
			continue
		}
		break
	}
}

// Stops the store, ensuring it can be garbage collected.
// Stores shouldn't need to be cancelled in production! (only really needed for tests)
func (s *Store) stop() {
	s.ctxCancel()
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

	if store != nil {
		store.stop()
	}

	store = New()
	return store
}
