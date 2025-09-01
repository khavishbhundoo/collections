package cmap

import (
	"sync"
	"testing"
)

func TestCMap_BasicOperations(t *testing.T) {
	var m CMap[string, int]

	// Set and Get
	m.Set("one", 1)
	m.Set("two", 2)

	if val, ok := m.Get("one"); !ok || val != 1 {
		t.Errorf("expected 1, got %v, ok=%v", val, ok)
	}

	if val, ok := m.Get("two"); !ok || val != 2 {
		t.Errorf("expected 2, got %v, ok=%v", val, ok)
	}

	// Contains
	if !m.Contains("one") || !m.Contains("two") {
		t.Errorf("expected keys to exist")
	}
	if m.Contains("three") {
		t.Errorf("expected key 'three' to not exist")
	}

	// Len
	if l := m.Len(); l != 2 {
		t.Errorf("expected length 2, got %d", l)
	}

	// Delete
	m.Delete("one")
	if m.Contains("one") {
		t.Errorf("key 'one' should have been deleted")
	}
	if l := m.Len(); l != 1 {
		t.Errorf("expected length 1 after delete, got %d", l)
	}

	// Keys
	keys := m.Keys()
	if len(keys) != 1 || keys[0] != "two" {
		t.Errorf("expected keys ['two'], got %v", keys)
	}
}

func TestCMap_ResetAndClear(t *testing.T) {
	m := NewWithCapacity[string, int](5)
	m.Set("a", 1)
	m.Set("b", 2)

	m.Reset()
	if m.Len() != 0 {
		t.Errorf("expected length 0 after Reset, got %d", m.Len())
	}
	m.Set("c", 3)
	if val, ok := m.Get("c"); !ok || val != 3 {
		t.Errorf("expected key 'c' after Reset, got %v, ok=%v", val, ok)
	}

	m.Clear()
	if m.Len() != 0 {
		t.Errorf("expected length 0 after Clear, got %d", m.Len())
	}
	m.Set("d", 4)
	if val, ok := m.Get("d"); !ok || val != 4 {
		t.Errorf("expected key 'd' after Clear, got %v, ok=%v", val, ok)
	}
}

func TestCMap_ConcurrentAccess(t *testing.T) {
	m := New[string, int]()
	wg := sync.WaitGroup{}
	const n = 1000

	// concurrent writers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Set(string(rune(i)), i)
		}(i)
	}

	wg.Wait()

	// concurrent readers
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, _ = m.Get(string(rune(i)))
		}(i)
	}

	wg.Wait()

	if m.Len() != n {
		t.Errorf("expected length %d after concurrent writes, got %d", n, m.Len())
	}
}

func TestCMap_ZeroValue(t *testing.T) {
	var m CMap[int, string]

	m.Set(1, "one")
	if val, ok := m.Get(1); !ok || val != "one" {
		t.Errorf("expected 'one', got %v, ok=%v", val, ok)
	}

	if l := m.Len(); l != 1 {
		t.Errorf("expected length 1, got %d", l)
	}
}
