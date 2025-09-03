package queue

import (
	"testing"
)

// BenchmarkQueue_Push benchmarks pushing N items into the queue
func BenchmarkQueue_Push(b *testing.B) {
	b.ReportAllocs()
	q := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

// BenchmarkQueue_PushWithCapacity benchmarks pushing N items into a queue
// pre-allocated with enough capacity to avoid reallocations.
func BenchmarkQueue_PushWithCapacity(b *testing.B) {
	b.ReportAllocs()
	q := NewWithCapacity[int](b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

// BenchmarkQueue_PushMany benchmarks pushing N items at once using PushMany
func BenchmarkQueue_PushMany(b *testing.B) {
	b.ReportAllocs()
	q := New[int]()
	items := make([]int, 1000)
	for i := 0; i < len(items); i++ {
		items[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.PushMany(items...)
	}
}

// BenchmarkQueue_Pop benchmarks popping N items from the queue
func BenchmarkQueue_Pop(b *testing.B) {
	b.ReportAllocs()
	q := New[int]()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Pop()
	}
}

// BenchmarkQueue_PopWithCapacity benchmarks popping N items with pre-allocated capacity
func BenchmarkQueue_PopWithCapacity(b *testing.B) {
	b.ReportAllocs()
	q := NewWithCapacity[int](b.N)
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Pop()
	}
}

// BenchmarkQueue_PushPop benchmarks mixed push/pop operations
func BenchmarkQueue_PushPop(b *testing.B) {
	b.ReportAllocs()
	q := New[int]()
	for i := 0; i < b.N; i++ {
		q.Push(i)
		if i%2 == 0 {
			_, _ = q.Pop()
		}
	}
}

// BenchmarkQueue_PushPopWithCapacity benchmarks mixed push/pop with pre-allocated capacity
func BenchmarkQueue_PushPopWithCapacity(b *testing.B) {
	b.ReportAllocs()
	q := NewWithCapacity[int](b.N)
	for i := 0; i < b.N; i++ {
		q.Push(i)
		if i%2 == 0 {
			_, _ = q.Pop()
		}
	}
}

// BenchmarkQueue_Peek benchmarks repeated peeks
func BenchmarkQueue_Peek(b *testing.B) {
	b.ReportAllocs()
	q := New[int]()
	for i := 0; i < 1000; i++ {
		q.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Peek()
	}
}

// BenchmarkQueue_Reset benchmarks repeated resets
func BenchmarkQueue_Reset(b *testing.B) {
	b.ReportAllocs()
	q := NewWithCapacity[int](b.N)
	for i := 0; i < 1000; i++ {
		q.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Reset()
	}
}

// BenchmarkQueue_Clear benchmarks repeated clears
func BenchmarkQueue_Clear(b *testing.B) {
	b.ReportAllocs()
	q := NewWithCapacity[int](b.N)
	for i := 0; i < 1000; i++ {
		q.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Clear()
	}
}
