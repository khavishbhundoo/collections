package set

import (
	"sync"
	"testing"
)

func TestSet_New(t *testing.T) {
	s := New[int]()
	if s == nil {
		t.Fatal("Expected non-nil Set")
	}
	if s.Len() != 0 {
		t.Errorf("Expected size 0, got %d", s.Len())
	}
}

func TestSet_NewWithCapacity(t *testing.T) {
	capacity := 10
	s := NewWithCapacity[int](capacity)
	if s == nil {
		t.Fatal("Expected non-nil Set")
	}
	if s.Len() != 0 {
		t.Errorf("Expected size 0, got %d", s.Len())
	}

	s.Add(1)
	s.Add(2)
	s.Add(3)
	if s.Len() != 3 {
		t.Errorf("Expected size 3, got %d", s.Len())
	}
}

func TestSet_Add(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	for _, v := range []int{1, 2, 3} {
		if !s.Contains(v) {
			t.Errorf("Expected set to contain %d", v)
		}
	}

	if s.Len() != 3 {
		t.Errorf("Expected size 3, got %d", s.Len())
	}

	// Adding duplicate should not increase size
	s.Add(1)
	if s.Len() != 3 {
		t.Errorf("Expected size 3 after adding duplicate, got %d", s.Len())
	}
}

func TestSet_AddMany(t *testing.T) {
	s := New[int]()
	s.AddMany(1, 2, 2, 3)
	if s.Len() != 3 {
		t.Errorf("Expected size 3, got %d", s.Len())
	}
	for _, v := range []int{1, 2, 3} {
		if !s.Contains(v) {
			t.Errorf("Expected set to contain %d", v)
		}
	}
}

func TestSet_Remove(t *testing.T) {
	s := New[int]()
	s.AddMany(1, 2, 3)

	s.Remove(2)
	if s.Contains(2) {
		t.Errorf("Expected 2 to be removed")
	}
	if s.Len() != 2 {
		t.Errorf("Expected size 2, got %d", s.Len())
	}

	// Removing non-existent element should not panic
	s.Remove(42)
	if s.Len() != 2 {
		t.Errorf("Expected size 2 after removing non-existent element, got %d", s.Len())
	}
}

func TestSet_Contains(t *testing.T) {
	var s Set[string] // zero-value set
	if s.Contains("x") {
		t.Errorf("Zero-value set should not contain any element")
	}

	s.Add("a")
	s.Add("b")

	tests := []struct {
		value    string
		expected bool
	}{
		{"a", true},
		{"b", true},
		{"c", false},
	}

	for _, tt := range tests {
		if got := s.Contains(tt.value); got != tt.expected {
			t.Errorf("Contains(%q) = %v, want %v", tt.value, got, tt.expected)
		}
	}
}

func TestSet_Size(t *testing.T) {
	var s Set[int] // zero-value
	if s.Len() != 0 {
		t.Errorf("Expected size 0, got %d", s.Len())
	}

	s.Add(1)
	s.Add(2)
	if s.Len() != 2 {
		t.Errorf("Expected size 2, got %d", s.Len())
	}

	s.Remove(1)
	if s.Len() != 1 {
		t.Errorf("Expected size 1, got %d", s.Len())
	}
}

func TestSet_Reset(t *testing.T) {
	s := New[int]()
	s.AddMany(1, 2, 3)
	s.Reset()
	if s.Len() != 0 {
		t.Errorf("Expected size 0 after Reset, got %d", s.Len())
	}
	if s.items == nil {
		t.Errorf("Reset should retain underlying map, got nil")
	}
}

func TestSet_Clear(t *testing.T) {
	s := New[int]()
	s.AddMany(1, 2, 3)
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", s.Len())
	}
	if s.items == nil {
		t.Errorf("Clear should allocate a new map, got nil")
	}
}

func TestSet_ConcurrentAdd(t *testing.T) {
	s := New[int]()
	var wg sync.WaitGroup
	n := 1000

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			s.Add(val)
		}(i)
	}
	wg.Wait()

	if s.Len() != n {
		t.Errorf("Expected size %d after concurrent Add, got %d", n, s.Len())
	}
}

func TestSet_ConcurrentAddAndRemove(t *testing.T) {
	s := New[int]()
	var wg sync.WaitGroup
	n := 1000

	// Start by adding everything
	for i := 0; i < n; i++ {
		s.Add(i)
	}

	// Add and remove concurrently
	for i := 0; i < n; i++ {
		wg.Add(2)
		go func(val int) {
			defer wg.Done()
			s.Add(val) // re-adding shouldn't change size
		}(i)
		go func(val int) {
			defer wg.Done()
			s.Remove(val)
		}(i)
	}
	wg.Wait()

	// Len could be anything between 0 and n depending on timing,
	// but it should never be negative or > n.
	size := s.Len()
	if size < 0 || size > n {
		t.Errorf("Unexpected size after concurrent Add/Remove: %d", size)
	}
}

func TestSet_ConcurrentContains(t *testing.T) {
	s := New[int]()
	for i := 0; i < 100; i++ {
		s.Add(i)
	}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			_ = s.Contains(val % 100) // Just ensure no panics or race conditions
		}(i)
	}
	wg.Wait()
}

func TestSet_ConcurrentResetAndClear(t *testing.T) {
	s := New[int]()
	for i := 0; i < 100; i++ {
		s.Add(i)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.Reset() // Should zero size but keep underlying map
	}()
	go func() {
		defer wg.Done()
		s.Clear() // Should zero size and allocate new map
	}()
	wg.Wait()

	// Just ensure it's in a valid state (no panic/race)
	if size := s.Len(); size < 0 {
		t.Errorf("Invalid size after concurrent Reset/Clear: %d", size)
	}
}
