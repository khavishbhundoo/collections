package cmap

import "sync"

// CMap is a generic, thread-safe key-value store with optional capacity hints.
// The implementation uses an underlying map protected by a sync.RWMutex.
// The zero value of CMap[K,V] is ready for use without initialization.
//
// Use New() or NewWithCapacity() if you prefer an explicit constructor
// or want to set an initial capacity. All operations are safe for
// concurrent use by multiple goroutines.
type CMap[K comparable, V any] struct {
	_               noCopy // prevents copying after first use
	items           map[K]V
	initialCapacity int
	mu              sync.RWMutex
}

// New returns an empty CMap with no pre-allocated capacity.
func New[K comparable, V any]() *CMap[K, V] {
	return &CMap[K, V]{
		items:           make(map[K]V),
		initialCapacity: 0,
	}
}

// NewWithCapacity returns an empty CMap with a capacity hint.
//
// Supplying a capacity reduces allocations if the expected number of
// key-value pairs is known in advance.
func NewWithCapacity[K comparable, V any](capacity int) *CMap[K, V] {
	return &CMap[K, V]{
		items:           make(map[K]V, capacity),
		initialCapacity: capacity,
	}
}

// Set associates value with key, creating the map if necessary.
// If key already exists, its value is replaced.
func (c *CMap[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.items == nil {
		c.items = make(map[K]V, c.initialCapacity)
	}
	c.items[key] = value
}

// Get returns the value for key and reports whether it was present.
// Returns the zero value of V if the key does not exist.
func (c *CMap[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.items == nil {
		var zero V
		return zero, false
	}
	val, ok := c.items[key]
	return val, ok
}

// Delete removes key and its value, if present.
// It does nothing if the key is not in the map.
func (c *CMap[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Contains reports whether key exists in the map.
func (c *CMap[K, V]) Contains(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.items[key]
	return ok
}

// Len returns the number of entries in the map.
func (c *CMap[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Keys returns a snapshot of all keys in the map.
// The returned slice does not reflect later modifications.
func (c *CMap[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]K, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

// Reset removes all entries while keeping the current allocation.
// Use Reset to reuse the map without triggering new allocations.
func (c *CMap[K, V]) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.items == nil {
		c.items = make(map[K]V, c.initialCapacity)
		return
	}
	clear(c.items)
}

// Clear removes all entries and allocates a new underlying map.
// Unlike Reset, Clear releases the old allocation to the runtime.
func (c *CMap[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[K]V, c.initialCapacity)
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
