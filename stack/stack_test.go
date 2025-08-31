package stack

import (
	"testing"
)

func TestStack_New(t *testing.T) {
	s := New[int]()
	s.Push(1)

	r, ok := s.Pop()
	if !ok {
		t.Errorf("Pop(): expected OK, got NOK")
	}
	if r != 1 {
		t.Errorf("Pop(): expected %d, got %d", 1, r)
	}

	if cap(s.items) != 1 {
		t.Errorf("Initial capacity: expected %d, got %d", 1, cap(s.items))
	}
}

func TestStack_NewWithCapacity(t *testing.T) {
	s := NewWithCapacity[int](5)
	s.Push(1)

	if s.Size() != 1 {
		t.Errorf("Size(): expected %d, got %d", 1, s.Size())
	}
	if s.initialCapacity != 5 && cap(s.items) != 5 {
		t.Errorf("Initial capacity: expected %d, got %d", 5, cap(s.items))
	}

	r, ok := s.Pop()
	if !ok {
		t.Errorf("Pop(): expected OK, got NOK")
	}
	if r != 1 {
		t.Errorf("Pop(): expected %d, got %d", 1, r)
	}
}

func TestStack_Push(t *testing.T) {
	s := New[int]()
	s.PushMany(1, 2, 3)
	s.Push(1)

	if s.Size() != 4 {
		t.Errorf("Size(): expected %d, got %d", 4, s.Size())
	}
}

func TestStack_Pop(t *testing.T) {
	s := New[int]()
	s.PushMany(1, 2, 3)

	if s.Size() != 3 {
		t.Errorf("Size(): expected %d, got %d", 3, s.Size())
	}

	r, ok := s.Pop()
	if !ok {
		t.Errorf("Pop(): expected OK, got NOK")
	}
	if r != 3 {
		t.Errorf("Pop(): expected %d, got %d", 3, r)
	}

	r, ok = s.Pop()
	if !ok {
		t.Errorf("Pop(): expected OK, got NOK")
	}
	if r != 2 {
		t.Errorf("Pop(): expected %d, got %d", 2, r)
	}
}

func TestStack_Size(t *testing.T) {
	s := New[int]()
	s.PushMany(1, 2, 3)

	if s.Size() != 3 {
		t.Errorf("Size(): expected %d, got %d", 3, s.Size())
	}
}

func TestStack_Reset(t *testing.T) {
	s := NewWithCapacity[int](5)
	s.PushMany(1, 2)
	s.Push(3)

	s.Reset()
	_, ok := s.Pop()
	if ok {
		t.Errorf("Pop() after Reset(): expected NOK, got OK")
	}

	if cap(s.items) != 5 {
		t.Errorf("Capacity after Reset(): expected %d, got %d", 5, cap(s.items))
	}
}

func TestStack_Clear(t *testing.T) {
	s := NewWithCapacity[int](2)
	s.PushMany(1, 2)
	s.Push(3)

	s.Clear()
	_, ok := s.Pop()
	if ok {
		t.Errorf("Pop() after Clear(): expected NOK, got OK")
	}

	if cap(s.items) != 2 {
		t.Errorf("Capacity after Clear(): expected %d, got %d", 2, cap(s.items))
	}
}
