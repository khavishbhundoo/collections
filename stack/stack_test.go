package stack

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New[int]()
	s.Push(1)
	expected := 1
	r, ok := s.Pop()
	if !ok {
		t.Error("Expected", "OK", "Actual", "NOK")
	}
	if r != expected {
		t.Error("Expected", expected, "got", r)
	}
	expected = 1
	if cap(s.items) != expected {
		t.Error("Expected", expected, "Actual", cap(s.items))
	}
}

func TestNewWithCapacity(t *testing.T) {
	s := NewWithCapacity[int](5)
	s.Push(1)
	expected := 1
	if s.Size() != 1 {
		t.Error("Expected", expected, "Actual", s.Size())
	}
	if s.initialCapacity != 5 && cap(s.items) != 5 {
		t.Error("Expected", 5, "Actual", cap(s.items))
	}
	r, ok := s.Pop()
	if !ok {
		t.Error("Expected", "OK", "Actual", "NOK")
	}
	if r != expected {
		t.Error("Expected", expected, "got", r)
	}

}

func TestStack_Push(t *testing.T) {
	s := New[int]()
	s.PushMany(1, 2, 3)
	s.Push(1)
	if s.Size() != 4 {
		t.Error("Expected", 3, "Actual", s.Size())
	}
}

func TestStack_Pop(t *testing.T) {
	s := New[int]()
	s.PushMany(1, 2, 3)
	if s.Size() != 3 {
		t.Error("Expected", 3, "Actual", s.Size())
	}
	expected := 3
	r, ok := s.Pop()
	if !ok {
		t.Error("Expected", "OK", "Actual", "NOK")
	}
	if r != expected {
		t.Error("Expected", expected, "got", r)
	}

	expected = 2
	r, ok = s.Pop()
	if !ok {
		t.Error("Expected", "OK", "Actual", "NOK")
	}
	if r != expected {
		t.Error("Expected", expected, "got", r)
	}
}

func TestStack_Size(t *testing.T) {
	s := New[int]()
	s.PushMany(1, 2, 3)
	if s.Size() != 3 {
		t.Error("Expected", 3, "Actual", s.Size())
	}
}

func TestStack_Reset(t *testing.T) {
	capacity := 5
	s := NewWithCapacity[int](capacity)
	s.PushMany(1, 2)
	s.Push(3)
	s.Reset()
	_, ok := s.Pop()
	if ok {
		t.Error("Expected", "NOK", "Actual", "OK")
	}
	if cap(s.items) != capacity {
		t.Error("Expected", capacity, "Actual", cap(s.items))
	}
}

func TestStack_Clear(t *testing.T) {
	capacity := 2
	s := NewWithCapacity[int](capacity)
	s.PushMany(1, 2)
	s.Push(3)
	s.Clear()
	_, ok := s.Pop()
	if ok {
		t.Error("Expected", "NOK", "Actual", "OK")
	}
	if cap(s.items) != capacity {
		t.Error("Expected", capacity, "Actual", cap(s.items))
	}
}
