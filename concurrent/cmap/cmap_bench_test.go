package cmap

import (
	"strconv"
	"sync"
	"testing"
)

func BenchmarkCMap_Set(b *testing.B) {
	m := New[string, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), i)
	}
}

func BenchmarkCMap_Get(b *testing.B) {
	m := New[string, int]()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Get(strconv.Itoa(i))
	}
}

func BenchmarkCMap_Contains(b *testing.B) {
	m := New[string, int]()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Contains(strconv.Itoa(i))
	}
}

func BenchmarkCMap_ConcurrentSet(b *testing.B) {
	m := New[string, int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := strconv.Itoa(i)
			m.Set(key, i)
			i++
		}
	})
}

func BenchmarkCMap_ConcurrentGet(b *testing.B) {
	m := New[string, int]()
	const N = 100000
	for i := 0; i < N; i++ {
		m.Set(strconv.Itoa(i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := strconv.Itoa(i % N)
			_, _ = m.Get(key)
			i++
		}
	})
}

func BenchmarkCMap_ZeroValueSet(b *testing.B) {
	var m CMap[int, int]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}

func BenchmarkCMap_Reset(b *testing.B) {
	m := New[int, int]()
	for i := 0; i < 10000; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Reset()
	}
}

func BenchmarkCMap_Clear(b *testing.B) {
	m := New[int, int]()
	for i := 0; i < 10000; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Clear()
	}
}

func BenchmarkCMap_Keys(b *testing.B) {
	m := New[int, int]()
	for i := 0; i < 10000; i++ {
		m.Set(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Keys()
	}
}

func BenchmarkCMap_ConcurrentMixed(b *testing.B) {
	m := New[int, int]()
	var wg sync.WaitGroup
	const workers = 8
	b.ResetTimer()

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < b.N/workers; i++ {
				m.Set(i+id*b.N/workers, i)
				_, _ = m.Get(i + id*b.N/workers)
			}
		}(w)
	}
	wg.Wait()
}

func BenchmarkMap_Set(b *testing.B) {
	m := make(map[string]int)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[strconv.Itoa(i)] = i
	}
}

func BenchmarkMap_Get(b *testing.B) {
	m := make(map[string]int)
	for i := 0; i < b.N; i++ {
		m[strconv.Itoa(i)] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[strconv.Itoa(i)]
	}
}

func BenchmarkMap_Contains(b *testing.B) {
	m := make(map[string]int)
	for i := 0; i < b.N; i++ {
		m[strconv.Itoa(i)] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m[strconv.Itoa(i)]
	}
}

func BenchmarkMap_ConcurrentSet(b *testing.B) {
	m := make(map[string]int)
	var mu sync.Mutex
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := strconv.Itoa(i)
			mu.Lock()
			m[key] = i
			mu.Unlock()
			i++
		}
	})
}

func BenchmarkMap_ConcurrentGet(b *testing.B) {
	m := make(map[string]int)
	const N = 100000
	for i := 0; i < N; i++ {
		m[strconv.Itoa(i)] = i
	}
	var mu sync.RWMutex
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := strconv.Itoa(i % N)
			mu.RLock()
			_ = m[key]
			mu.RUnlock()
			i++
		}
	})
}

func BenchmarkMap_ZeroValueSet(b *testing.B) {
	var m map[int]int // zero-value map
	// Need to initialize on first use
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if m == nil {
			m = make(map[int]int)
		}
		m[i] = i
	}
}

func BenchmarkMap_Keys(b *testing.B) {
	m := make(map[int]int)
	for i := 0; i < 10000; i++ {
		m[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		keys := make([]int, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
	}
}
