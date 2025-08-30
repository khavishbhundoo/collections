package stack

import "testing"

func BenchmarkStack_Push(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	for b.Loop() {
		s.Push(1)
	}
}

func BenchmarkStack_PushWithCapacity(b *testing.B) {
	b.ReportAllocs()
	s := NewWithCapacity[int](b.N)
	for b.Loop() {
		s.Push(1)
	}
}

func BenchmarkStack_Pop(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	for i := 0; i < b.N; i++ {
		s.Push(1)
	}
	for b.Loop() {
		s.Pop()
	}
}

func BenchmarkStack_PopWithCapacity(b *testing.B) {
	b.ReportAllocs()
	s := NewWithCapacity[int](b.N)
	for i := 0; i < b.N; i++ {
		s.Push(1)
	}
	for b.Loop() {
		s.Pop()
	}
}

func BenchmarkStack_ClearAndPush10M(b *testing.B) {
	b.ReportAllocs()
	const stackSize = 10_000_000 // 10M elements
	s := New[int]()

	for i := 0; i < stackSize; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for b.Loop() {
		s.Clear()

		// Rebuild the stack with 10M elements
		for j := 0; j < stackSize; j++ {
			s.Push(j)
		}
	}
}

func BenchmarkStack_ResetAndPush10M(b *testing.B) {
	b.ReportAllocs()
	const stackSize = 10_000_000 // 10M elements
	s := New[int]()

	for i := 0; i < stackSize; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for b.Loop() {
		s.Reset()

		// Rebuild the stack with 10M elements
		for j := 0; j < stackSize; j++ {
			s.Push(j)
		}
	}
}
