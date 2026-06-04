package main

import (
	"sync"
)

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]string),
	}
}

func (sm *SafeMap) Set(key string, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

func (sm *SafeMap) Get(key string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, ok := sm.m[key]
	return val, ok
}

func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.m, key)
}

func (sm *SafeMap) List() []string {
	var keyStore []string

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	for key := range sm.m {
		keyStore = append(keyStore, key)
	}

	return keyStore

}
