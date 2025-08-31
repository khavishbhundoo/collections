package set

import "testing"

func BenchmarkSet_Add(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkSet_Add_PreSized(b *testing.B) {
	b.ReportAllocs()
	s := NewWithCapacity[int](b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
	}
}

func BenchmarkSet_AddMany100_FreshEmpty(b *testing.B) {
	b.ReportAllocs()
	values := make([]int, 100)
	for i := 0; i < 100; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := New[int]() // fresh set, capacity = 0
		s.AddMany(values...)
	}
}

func BenchmarkSet_AddMany100_FreshPreSized(b *testing.B) {
	b.ReportAllocs()
	values := make([]int, 100)
	for i := 0; i < 100; i++ {
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := NewWithCapacity[int](len(values)) // fresh set, capacity = 100
		s.AddMany(values...)
	}
}

func BenchmarkSet_AddMany100_AmortizedPreSized(b *testing.B) {
	b.ReportAllocs()
	values := make([]int, 100)
	for i := 0; i < 100; i++ {
		values[i] = i
	}
	s := NewWithCapacity[int](b.N * len(values)) // big reusable set
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
		key := i % 1000
		s.Remove(key)
		s.Add(key) // alternate remove/add to keep size stable
	}
}

func BenchmarkSet_Contains(b *testing.B) {
	b.ReportAllocs()
	const n = 1000
	s := New[int]()
	for i := 0; i < n; i++ {
		s.Add(i)
	}
	keys := make([]int, 2*n)
	for i := 0; i < n; i++ {
		keys[i] = i       // exists
		keys[i+n] = i + n // does not exist
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Contains(keys[i%len(keys)])
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

func BenchmarkSet_Reset_1k(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	vals := make([]int, 1000)
	for i := range vals {
		vals[i] = i
	}
	s.AddMany(vals...) // warm-up
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Reset()          // keeps buckets
		s.AddMany(vals...) // refill
	}
}

func BenchmarkSet_Clear_1k_InitialCapacityZero(b *testing.B) {
	b.ReportAllocs()
	s := New[int]()
	vals := make([]int, 1000)
	for i := range vals {
		vals[i] = i
	}
	s.AddMany(vals...) // warm-up
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Clear()          // new map with capacity 0 → will grow every time
		s.AddMany(vals...) // refill triggers map growth
	}
}

func BenchmarkSet_Clear_1k_PreSized(b *testing.B) {
	b.ReportAllocs()
	s := NewWithCapacity[int](1000)
	vals := make([]int, 1000)
	for i := range vals {
		vals[i] = i
	}
	s.AddMany(vals...) // warm-up
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Clear()          // new map with cap 1000 → minimal growth
		s.AddMany(vals...) // refill avoids extra allocations
	}
}
