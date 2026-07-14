package main

import (
	"sync"
	"testing"
)

func TestSafeMap_SetAndGet(t *testing.T) {
	sm := NewSafeMap()
	sm.Set("name", "alice")
	val, ok := sm.Get("name")
	if !ok || val != "alice" {
		t.Errorf("Get(\"name\") = %q, %v, want \"alice\", true", val, ok)
	}
}

func TestSafeMap_GetMissingKey(t *testing.T) {
	sm := NewSafeMap()
	_, ok := sm.Get("missing")
	if ok {
		t.Error("Get(\"missing\") returned ok=true, want false")
	}
}

func TestSafeMap_Delete(t *testing.T) {
	sm := NewSafeMap()
	sm.Set("key", "val")
	sm.Delete("key")
	_, ok := sm.Get("key")
	if ok {
		t.Error("key still exists after Delete")
	}
}

func TestSafeMap_List(t *testing.T) {
	sm := NewSafeMap()
	sm.Set("a", "1")
	sm.Set("b", "2")
	keys := sm.List()
	if len(keys) != 2 {
		t.Errorf("List() returned %d keys, want 2", len(keys))
	}
}

func TestSafeMap_ListEmpty(t *testing.T) {
	sm := NewSafeMap()
	keys := sm.List()
	if len(keys) != 0 {
		t.Errorf("List() on empty store returned %d keys, want 0", len(keys))
	}
}

func TestSafeMap_ConcurrentReadWrite(t *testing.T) {
	sm := NewSafeMap()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			sm.Set("key", "value")
		}()
		go func() {
			defer wg.Done()
			sm.Get("key")
		}()
	}

	wg.Wait()
}
