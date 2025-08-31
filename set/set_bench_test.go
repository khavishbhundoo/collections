package set

import (
	"testing"
)

func BenchmarkSet_Add(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	b.ResetTimer()
	for b.Loop() {
		s.Add(b.N)
	}
}

func BenchmarkSet_AddMany(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	values := make([]int, 100)
	for i := 0; i < 100; i++ {
		values[i] = i
	}

	for b.Loop() {
		s.AddMany(values...)
	}
}

func BenchmarkSet_Remove(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Remove(i % 1000)
		s.Add(i % 1000)
	}
}

func BenchmarkSet_Contains(b *testing.B) {
	b.ReportAllocs()

	const n = 1000
	s := New[int]()
	for i := 0; i < n; i++ {
		s.Add(i)
	}

	// Precompute lookup keys: half exist, half do not exist
	keys := make([]int, 2*n)
	for i := 0; i < n; i++ {
		keys[i] = i       // exists
		keys[i+n] = i + n // does not exist
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i%len(keys)]
		_ = s.Contains(key)
	}
}

func BenchmarkSet_Size(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Size()
	}
}
