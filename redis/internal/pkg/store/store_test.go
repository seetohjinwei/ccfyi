package store

import (
	"testing"
)

func TestStoreGetSet(t *testing.T) {
	store := New()

	err := store.Set("key", NewString("value"))
	if err != nil {
		t.Errorf("expected no err, but got %+v", err)
	}
}

func TestStoreGetConcurrent(t *testing.T) {
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
