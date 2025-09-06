package queue

import (
	"runtime"
	"testing"
)

func BenchmarkQueue_Push(b *testing.B) {
	b.ReportAllocs()
	b.Cleanup(func() { runtime.GC() })
	s := New[int]()
	for b.Loop() {
		s.Push(1)
	}
}

func BenchmarkQueue_PushWithCapacity(b *testing.B) {
	b.ReportAllocs()
	b.Cleanup(func() { runtime.GC() })
	s := NewWithCapacity[int](b.N)
	for b.Loop() {
		s.Push(1)
	}
}

func BenchmarkQueue_Pop(b *testing.B) {
	b.ReportAllocs()
	b.Cleanup(func() { runtime.GC() })
	s := New[int]()
	for i := 0; i < b.N; i++ {
		s.Push(1)
	}
	for b.Loop() {
		s.Pop()
	}
}

func BenchmarkQueue_PopWithCapacity(b *testing.B) {
	b.ReportAllocs()
	b.Cleanup(func() { runtime.GC() })
	s := NewWithCapacity[int](b.N)
	for i := 0; i < b.N; i++ {
		s.Push(1)
	}
	for b.Loop() {
		s.Pop()
	}
}

func BenchmarkQueue_ClearAndPush10M(b *testing.B) {
	b.ReportAllocs()
	b.Cleanup(func() { runtime.GC() })
	const queueSize = 10_000_000 // 10M elements
	s := New[int]()

	for i := 0; i < queueSize; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for b.Loop() {
		s.Clear()

		// Rebuild the stack with 10M elements
		for j := 0; j < queueSize; j++ {
			s.Push(j)
		}
	}
}

func BenchmarkQueue_ResetAndPush10M(b *testing.B) {
	b.ReportAllocs()
	b.Cleanup(func() { runtime.GC() })
	const queueSize = 10_000_000 // 10M elements
	s := New[int]()

	for i := 0; i < queueSize; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for b.Loop() {
		s.Reset()

		// Rebuild the stack with 10M elements
		for j := 0; j < queueSize; j++ {
			s.Push(j)
		}
	}
}

func BenchmarkQueue_ClearAndPush10MWithCapacity(b *testing.B) {
	b.ReportAllocs()
	b.Cleanup(func() { runtime.GC() })
	const queueSize = 10_000_000 // 10M elements
	s := NewWithCapacity[int](queueSize)

	for i := 0; i < queueSize; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for b.Loop() {
		s.Clear()

		// Rebuild the stack with 10M elements
		for j := 0; j < queueSize; j++ {
			s.Push(j)
		}
	}
}

func BenchmarkQueue_ResetAndPush10MWithCapacity(b *testing.B) {
	b.ReportAllocs()
	b.Cleanup(func() { runtime.GC() })
	const queueSize = 10_000_000 // 10M elements
	s := NewWithCapacity[int](queueSize)

	for i := 0; i < queueSize; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for b.Loop() {
		s.Reset()

		// Rebuild the stack with 10M elements
		for j := 0; j < queueSize; j++ {
			s.Push(j)
		}
	}
}
