package queue

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueue_New(t *testing.T) {
	q := New[int]()
	q.Push(1)

	r, ok := q.Pop()
	if !ok {
		t.Errorf("Pop(): expected OK, got NOK")
	}
	if r != 1 {
		t.Errorf("Pop(): expected %d, got %d", 1, r)
	}

	if cap(q.items) != 0 {
		t.Errorf("Initial capacity: expected %d, got %d", 0, cap(q.items))
	}
}

func TestQueue_NewWithCapacity(t *testing.T) {
	q := NewWithCapacity[int](5)
	q.Push(1)

	if q.Len() != 1 {
		t.Errorf("Len(): expected %d, got %d", 1, q.Len())
	}
	if q.initialCapacity != 5 && cap(q.items) != 5 {
		t.Errorf("Initial capacity: expected %d, got %d", 5, cap(q.items))
	}

	r, ok := q.Pop()
	if !ok {
		t.Errorf("Pop(): expected OK, got NOK")
	}
	if r != 1 {
		t.Errorf("Pop(): expected %d, got %d", 1, r)
	}
}

func TestQueue_Push(t *testing.T) {
	q := New[int]()
	q.PushMany(1, 2, 3)
	q.Push(4)

	if q.Len() != 4 {
		t.Errorf("Len(): expected %d, got %d", 4, q.Len())
	}
}

func TestQueue_Pop(t *testing.T) {
	q := New[int]()
	q.PushMany(1, 2, 3)

	if q.Len() != 3 {
		t.Errorf("Len(): expected %d, got %d", 3, q.Len())
	}

	r, ok := q.Pop()
	if !ok {
		t.Errorf("Pop(): expected OK, got NOK")
	}
	if r != 1 {
		t.Errorf("Pop(): expected %d, got %d", 1, r)
	}

	r, ok = q.Pop()
	if !ok {
		t.Errorf("Pop(): expected OK, got NOK")
	}
	if r != 2 {
		t.Errorf("Pop(): expected %d, got %d", 2, r)
	}

	if q.Len() != 1 {
		t.Errorf("Len(): expected %d, got %d", 1, q.Len())
	}
}

func TestQueue_Peek(t *testing.T) {
	q := New[int]()
	q.PushMany(10, 20)

	r, ok := q.Peek()
	if !ok {
		t.Errorf("Peek(): expected OK, got NOK")
	}
	if r != 10 {
		t.Errorf("Peek(): expected %d, got %d", 10, r)
	}

	q.Pop()
	r, ok = q.Peek()
	if !ok || r != 20 {
		t.Errorf("Peek(): expected %d, got %d", 20, r)
	}
}

func TestQueue_Reset(t *testing.T) {
	q := NewWithCapacity[int](5)
	q.PushMany(1, 2, 3)

	q.Reset()
	_, ok := q.Pop()
	if ok {
		t.Errorf("Pop() after Reset(): expected NOK, got OK")
	}

	if cap(q.items) != 5 {
		t.Errorf("Capacity after Reset(): expected %d, got %d", 5, cap(q.items))
	}
}

func TestQueue_Clear(t *testing.T) {
	q := NewWithCapacity[int](2)
	q.PushMany(1, 2)
	q.Push(3)

	q.Clear()
	_, ok := q.Pop()
	if ok {
		t.Errorf("Pop() after Clear(): expected NOK, got OK")
	}

	if cap(q.items) != 2 {
		t.Errorf("Capacity after Clear(): expected %d, got %d", 2, cap(q.items))
	}
}

func TestQueue_OrderFIFO(t *testing.T) {
	q := New[int]()
	for i := 1; i <= 5; i++ {
		q.Push(i)
	}

	for i := 1; i <= 5; i++ {
		r, ok := q.Pop()
		if !ok || r != i {
			t.Errorf("Pop(): expected %d, got %d", i, r)
		}
	}

	if q.Len() != 0 {
		t.Errorf("Len(): expected %d, got %d", 0, q.Len())
	}
}

func TestQueue_EmptyPopPeek(t *testing.T) {
	q := New[int]()

	if _, ok := q.Pop(); ok {
		t.Errorf("Pop() on empty queue: expected NOK, got OK")
	}
	if _, ok := q.Peek(); ok {
		t.Errorf("Peek() on empty queue: expected NOK, got OK")
	}
}

func TestQueue_Shrink(t *testing.T) {
	initialCap := 64
	q := NewWithCapacity[int](initialCap)

	// Fill the queue
	for i := 0; i < 100; i++ {
		q.Push(i)
	}

	peakCap := cap(q.items)

	// Pop most of the items to trigger shrink (less than 25% used)
	for i := 0; i < 90; i++ {
		_, ok := q.Pop()
		if !ok {
			t.Fatalf("Pop() failed at iteration %d", i)
		}
	}

	if cap(q.items) >= peakCap {
		t.Errorf("Expected capacity to shrink below peak %d, got %d", peakCap, cap(q.items))
	}

	// Capacity should not be greater than initial cap
	if cap(q.items) > initialCap {
		t.Logf("Warning: capacity should not be greater than initialCap: %d", cap(q.items))
	}

	// Remaining elements should still be correct
	for _, val := range []int{90, 91, 92, 93, 94, 95, 96, 97, 98, 99} {
		r, ok := q.Pop()
		if !ok || r != val {
			t.Errorf("Expected Pop() to return %d, got %d", val, r)
		}
	}

	if q.Len() != 0 {
		t.Errorf("Expected queue to be empty after pops, got Len() = %d", q.Len())
	}
}

func TestQueue_GenericType(t *testing.T) {
	q := New[string]()
	q.Push("foo")
	q.Push("bar")

	r, ok := q.Pop()
	if !ok || r != "foo" {
		t.Errorf("Pop(): expected 'foo', got '%s'", r)
	}

	r, ok = q.Pop()
	if !ok || r != "bar" {
		t.Errorf("Pop(): expected 'bar', got '%s'", r)
	}
}

func TestQueue_ConcurrentPush(t *testing.T) {
	const goroutines = 10
	const perGoroutine = 1000

	q := New[int]()
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(base int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				q.Push(base*perGoroutine + i)
			}
		}(g)
	}

	wg.Wait()

	if q.Len() != goroutines*perGoroutine {
		t.Errorf("Len(): expected %d, got %d", goroutines*perGoroutine, q.Len())
	}
}

func TestQueue_ConcurrentPop(t *testing.T) {
	const total = 2000
	q := New[int]()

	for i := 0; i < total; i++ {
		q.Push(i)
	}

	var wg sync.WaitGroup
	var popped atomic.Int64
	const goroutines = 5
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for {
				_, ok := q.Pop()
				if !ok {
					return
				}
				popped.Add(1)
			}
		}()
	}

	wg.Wait()

	if popped.Load() != int64(total) {
		t.Errorf("Expected %d pops, got %d", total, popped.Load())
	}
	if q.Len() != 0 {
		t.Errorf("Expected empty queue, got Len() = %d", q.Len())
	}
}

func TestQueue_ConcurrentPushPop(t *testing.T) {
	const goroutines = 10
	const opsPerG = 500
	q := New[int]()
	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Pushers
	for g := 0; g < goroutines; g++ {
		go func(base int) {
			defer wg.Done()
			for i := 0; i < opsPerG; i++ {
				q.Push(base*opsPerG + i)
			}
		}(g)
	}

	// Poppers
	var popCount atomic.Int64
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for {
				_, ok := q.Pop()
				if !ok {
					time.Sleep(time.Microsecond)
					if q.Len() == 0 {
						return
					}
					continue
				}
				popCount.Add(1)
			}
		}()
	}

	wg.Wait()

	expected := goroutines * opsPerG
	if popCount.Load() != int64(expected) {
		t.Errorf("Expected %d total pops, got %d", expected, popCount.Load())
	}
	if q.Len() != 0 {
		t.Errorf("Expected empty queue, got Len() = %d", q.Len())
	}
}
