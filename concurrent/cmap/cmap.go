package cmap

import (
	"fmt"
	"hash/maphash"
	"sync"
	"sync/atomic"
)

// CMap is a generic, thread-safe key-value store with optional capacity hints.
//
// The implementation uses an underlying map protected by a sync.RWMutex.
//
// The zero value of CMap[K,V] is ready for use without initialization:
//
//	var m CMap[string, int]
//	m.Set("Go", 1)  // works without calling New()
//	v, ok := m.Get("Go")
//
// Use New() or NewWithCapacity() if you prefer an explicit constructor
// or want to set an initial capacity. All operations are safe for
// concurrent use by multiple goroutines.

const shardCount = 64

type CMap[K comparable, V any] struct {
	_             noCopy // prevents copying after first use
	once          sync.Once
	shards        [shardCount]shard[K, V]
	seed          maphash.Seed
	shardCapacity int
}

func (c *CMap[K, V]) init(capacity int) {
	c.once.Do(func() {
		c.seed = maphash.MakeSeed()
		for i := 0; i < shardCount; i++ {
			c.shards[i].mu = sync.RWMutex{}
			if capacity <= shardCount {
				c.shards[i].items = make(map[K]V)
			} else {
				c.shardCapacity = capacity / shardCount
				c.shards[i].items = make(map[K]V, c.shardCapacity)
			}
		}
	})
}

type shard[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
	len   atomic.Int64
}

// New returns an empty CMap with no pre-allocated capacity.
func New[K comparable, V any]() *CMap[K, V] {

	var c CMap[K, V]
	c.init(0)
	return &c
}

// NewWithCapacity returns an empty CMap with a capacity hint.
//
// Supplying a capacity reduces allocations if the expected number of
// key-value pairs is known in advance.
func NewWithCapacity[K comparable, V any](capacity int) *CMap[K, V] {

	var c CMap[K, V]
	c.init(capacity)
	return &c
}

// Set associates value with key, creating the map if necessary.
// If key already exists, its value is replaced.
func (c *CMap[K, V]) Set(key K, value V) {
	var zeroSeed maphash.Seed
	if c.seed == zeroSeed {
		c.init(0)
	}
	s := &c.shards[c.shardIndex(key)]
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.items == nil {
		s.items = make(map[K]V, c.shardCapacity)
	}
	s.items[key] = value
	s.len.Add(1)
}

// Get returns the value for key and reports whether it was present.
// Returns the zero value of V if the key does not exist.
func (c *CMap[K, V]) Get(key K) (V, bool) {
	var zeroSeed maphash.Seed
	if c.seed == zeroSeed {
		var zero V
		return zero, false
	}
	s := &c.shards[c.shardIndex(key)]
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.items == nil {
		var zero V
		return zero, false
	}
	val, ok := s.items[key]
	return val, ok
}

// Delete removes key and its value, if present.
// It does nothing if the key is not in the map.
func (c *CMap[K, V]) Delete(key K) {
	var zeroSeed maphash.Seed
	if c.seed == zeroSeed {
		return
	}
	s := &c.shards[c.shardIndex(key)]
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
	s.len.Add(-1)
}

// Contains reports whether key exists in the map.
func (c *CMap[K, V]) Contains(key K) bool {
	var zeroSeed maphash.Seed
	if c.seed == zeroSeed {
		return false
	}
	s := &c.shards[c.shardIndex(key)]
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.items[key]
	return ok
}

// Len returns the number of entries in the map.
func (c *CMap[K, V]) Len() int {
	var zeroSeed maphash.Seed
	if c.seed == zeroSeed {
		return 0
	}
	total := int64(0)
	for i := range c.shards {
		total += c.shards[i].len.Load()
	}
	return int(total)
}

// Keys returns a snapshot of all keys in the map.
// The returned slice does not reflect later modifications.
func (c *CMap[K, V]) Keys() []K {
	var zeroSeed maphash.Seed
	if c.seed == zeroSeed {
		return nil
	}

	keys := make([]K, 0, c.Len())

	for i := range c.shards {
		s := &c.shards[i]
		s.mu.RLock()
		for k := range s.items {
			keys = append(keys, k)
		}
		s.mu.RUnlock()
	}

	return keys
}

// Reset removes all entries while keeping the current allocation.
// Use Reset to reuse the map without triggering new allocations.
func (c *CMap[K, V]) Reset() {
	var zeroSeed maphash.Seed
	if c.seed == zeroSeed {
		c.init(0)
	}

	for i := range c.shards {
		s := &c.shards[i]
		s.mu.Lock()
		clear(s.items)
		s.len.Store(0)
		s.mu.Unlock()
	}
}

// Clear removes all entries and allocates a new underlying map.
// Unlike Reset, Clear releases the old allocation to the runtime.
func (c *CMap[K, V]) Clear() {
	var zeroSeed maphash.Seed
	if c.seed == zeroSeed {
		return
	}
	for i := range c.shards {
		c.shards[i].mu.Lock()
		c.shards[i].items = make(map[K]V, c.shardCapacity)
		c.shards[i].len.Store(0)
		c.shards[i].mu.Unlock()
	}
}

func (c *CMap[K, V]) shardIndex(key K) uint64 {
	h := maphash.Hash{}
	h.SetSeed(c.seed)
	h.Reset()
	_, err := h.WriteString(fmt.Sprintf("%v", key))
	if err != nil {
		panic(err)
	}
	return h.Sum64() % shardCount
}

// noCopy may be added to structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
//
// Note that it must not be embedded, due to the Lock and Unlock methods.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
