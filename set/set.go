package set

// Set is a generic, non-thread-safe set implementation backed by a map[T]struct{}.
// It stores unique elements of type T. The zero value of Set[T] is ready to use
// without initialization.
//
// Use New() or NewWithCapacity() to explicitly create a set or provide an initial capacity.
// For a thread-safe set, see collections/concurrent/set.
type Set[T comparable] struct {
	items           map[T]struct{}
	initialCapacity int
}

// New creates an empty set of type T with no pre-allocated capacity.
// Equivalent to declaring `var s set.Set[int]`.
func New[T comparable]() *Set[T] {
	return &Set[T]{
		items:           make(map[T]struct{}),
		initialCapacity: 0,
	}
}

// NewWithCapacity creates an empty set with a capacity hint for the underlying map.
// Useful when you know approximately how many elements the set will contain.
func NewWithCapacity[T comparable](capacity int) *Set[T] {
	return &Set[T]{
		items:           make(map[T]struct{}, capacity),
		initialCapacity: capacity,
	}
}

// Add inserts a value into the set. If the value already exists, it does nothing.
// Initializes the underlying map if it is nil.
func (s *Set[T]) Add(value T) {
	if s.items == nil {
		s.items = make(map[T]struct{}, s.initialCapacity)
	}
	s.items[value] = struct{}{}
}

// AddMany inserts multiple values into the set. Duplicates are ignored.
// Initializes the underlying map if it is nil, sizing it to hold all values.
func (s *Set[T]) AddMany(values ...T) {
	if s.items == nil {
		s.items = make(map[T]struct{}, max(s.initialCapacity, len(values)))
	}
	for _, v := range values {
		s.items[v] = struct{}{}
	}
}

// Remove deletes a value from the set if it exists. Safe on a zero-value Set.
func (s *Set[T]) Remove(value T) {
	if s.items == nil {
		return
	}
	delete(s.items, value)
}

// Contains reports whether a value exists in the set.
// Safe to call on a zero-value Set; returns false without allocating.
func (s *Set[T]) Contains(value T) bool {
	if s.items == nil {
		return false
	}
	_, exists := s.items[value]
	return exists
}

// Len returns the number of elements in the set.
// Safe to call on a zero-value Set; returns 0 without allocating.
func (s *Set[T]) Len() int {
	if s.items == nil {
		return 0
	}
	return len(s.items)
}

// Reset removes all elements from the set but retains the underlying map capacity.
// Initializes the map if it is nil.
func (s *Set[T]) Reset() {
	if s.items == nil {
		s.items = make(map[T]struct{}, s.initialCapacity)
		return
	}
	clear(s.items)
}

// Clear removes all elements and resets the underlying map to the initial capacity.
// Always allocates a new map.
func (s *Set[T]) Clear() {
	s.items = make(map[T]struct{}, s.initialCapacity)
}
