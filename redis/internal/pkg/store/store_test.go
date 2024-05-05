package store

import (
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/seetohjinwei/ccfyi/redis/pkg/delay"
)

func TestStoreGetSet(t *testing.T) {
	t.Parallel()

	store := New()

	err := store.Set("key", NewString("value"))
	if err != nil {
		t.Errorf("expected no err, but got %+v", err)
	}
}

func TestStoreGetConcurrent(t *testing.T) {
	t.Parallel()

	wait := make(chan struct{})

	store := New()
	store.Set("key", NewString("value"))

	go func() {
		<-wait
		store.Get("key")
	}()
	go func() {
		<-wait
		store.Get("key")
	}()
	go func() {
		<-wait
		store.Get("key2")
	}()

	close(wait)
}

func TestStoreSetConcurrent(t *testing.T) {
	t.Parallel()

	wait := make(chan struct{})

	store := New()

	go func() {
		<-wait
		store.Set("key1", NewString("value"))
	}()
	go func() {
		<-wait
		store.Set("key2", NewString("value"))
	}()
	go func() {
		<-wait
		store.Set("key3", NewString("value"))
	}()

	close(wait)
}

func TestStoreGetSetConcurrent(t *testing.T) {
	t.Parallel()

	wait := make(chan struct{})

	store := New()

	go func() {
		<-wait
		store.Set("key1", NewString("value"))
	}()
	go func() {
		<-wait
		store.Get("key1")
	}()
	go func() {
		<-wait
		store.Set("key2", NewString("value"))
	}()

	close(wait)
}

func TestStoreSetDelay(t *testing.T) {
	t.Parallel()

	var key string
	var value string
	var actual string

	store := New()

	key = "key"
	value = "value"

	store.SetWithDelay(key, NewString(value), delay.NewDelay(time.Now().Add(50*time.Millisecond)))
	item, ok := store.Get(key)
	if !ok {
		t.Errorf("expected to get the value before expiry")
	}
	actual, ok = item.Get()
	if !ok {
		t.Errorf("expected to get the value before expiry")
	}
	if actual != value {
		t.Errorf("expected %q, but got %q", value, actual)
	}

	time.Sleep(50 * time.Millisecond)
	_, ok = store.Get(key)
	if ok {
		t.Errorf("expected to key to have expired")
	}
}

func TestStoreExpiryTrigger(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping long test")
	}

	result := atomic.Int32{}

	store := newNoExpiry()
	go store.activeExpiry(func() {
		result.Add(1)
	})

	time.Sleep(time.Second)
	store.stop()

	actual := result.Load()
	if actual != 10 {
		t.Errorf("expected result to be %v, but got %v", 10, actual)
	}
}

func TestStoreCleanKeys(t *testing.T) {
	t.Parallel()

	store := newNoExpiry()

	// expire immediately
	store.SetWithDelay("k", NewString("v"), delay.NewDelay(time.Now()))
	if len(store.values) != 1 {
		t.Errorf("expected store to contain the value before cleaning")
	}
	store.cleanKeys()
	if len(store.values) != 0 {
		t.Errorf("expected store to not contain the value after cleaning, got %v instead", len(store.values))
	}

	for i := 0; i < 19; i++ {
		store.SetWithDelay(strconv.Itoa(i), NewString("v"), delay.NewDelay(time.Now().Add(time.Hour)))
	}
	for i := 19; i < 22; i++ {
		store.SetWithDelay(strconv.Itoa(i), NewString("v"), delay.NewDelay(time.Now()))
	}
	if len(store.values) != 22 {
		t.Errorf("expected store to contain the values before cleaning")
	}
	store.cleanKeys()

	// if len(store.values) == 22 {
	// 	// it checks 20 keys (19 have not expired; 3 have expired => must check at least one expired key) <- not true, because it can check the same key multiple times
	// 	t.Errorf("expected store to have cleaned some keys, got %v instead", len(store.values))
	// }
}

func TestMapRandomness(t *testing.T) {
	t.Skip("not an actual test")

	m := make(map[int]struct{}, 100)
	for i := 0; i < 100; i++ {
		m[i] = struct{}{}
	}

	for i := 0; i < 20; i++ {
		for k := range m {
			t.Logf("k %d", k)
			break
		}
	}

	// forces Logf logs to be printed
	t.Fail()
}
