package stack

import (
	"sync"
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

func TestStack_ConcurrentPush(t *testing.T) {
	s := New[int]()
	const goroutines = 10
	const pushesPerGoroutine = 100
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < pushesPerGoroutine; j++ {
				s.Push(id)
			}
		}(i)
	}

	wg.Wait()

	expectedSize := goroutines * pushesPerGoroutine
	if s.Size() != expectedSize {
		t.Errorf("Expected size %d after concurrent Push, got %d", expectedSize, s.Size())
	}
}

func TestStack_ConcurrentPushMany(t *testing.T) {
	s := New[int]()
	const goroutines = 5
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			s.PushMany(id*10+1, id*10+2, id*10+3)
		}(i)
	}

	wg.Wait()

	if s.Size() != goroutines*3 {
		t.Errorf("Expected size %d after concurrent PushMany, got %d", goroutines*3, s.Size())
	}
}

func TestStack_ConcurrentPop(t *testing.T) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Push(i)
	}

	const goroutines = 5
	const popsPerGoroutine = 200
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < popsPerGoroutine; j++ {
				s.Pop()
			}
		}()
	}

	wg.Wait()

	expectedSize := 1000 - (goroutines * popsPerGoroutine)
	if s.Size() != expectedSize {
		t.Errorf("Expected size %d after concurrent Pop, got %d", expectedSize, s.Size())
	}
}

func TestStack_ConcurrentPushPop(t *testing.T) {
	s := New[int]()
	const goroutines = 5
	const opsPerGoroutine = 200
	wg := sync.WaitGroup{}
	wg.Add(2 * goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				s.Push(j + id*opsPerGoroutine)
			}
		}(i)

		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				s.Pop() // ignore empty pops
			}
		}()
	}

	wg.Wait()

	size := s.Size()
	if size < 0 || size > goroutines*opsPerGoroutine {
		t.Errorf("Unexpected size %d after concurrent Push/Pop", size)
	}
}

func TestStack_ConcurrentResetClear(t *testing.T) {
	s := New[int]()
	const goroutines = 5
	const opsPerGoroutine = 200
	wg := sync.WaitGroup{}
	wg.Add(2 * goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				s.Push(j + id*opsPerGoroutine)
			}
		}(i)

		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				s.Pop()
			}
		}()
	}

	wg.Wait()

	s.Reset()
	if s.Size() != 0 {
		t.Error("Expected size 0 after Reset")
	}

	for i := 0; i < goroutines; i++ {
		s.PushMany(i, i+1)
	}
	s.Clear()
	if s.Size() != 0 {
		t.Error("Expected size 0 after Clear")
	}
}
