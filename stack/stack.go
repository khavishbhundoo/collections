package stack

type Stack[T any] struct {
	items           []T
	initialCapacity int
}

const shrinkCapacityThreshold = 16

func New[T any]() *Stack[T] {
	return &Stack[T]{
		items:           []T{},
		initialCapacity: 0,
	}
}

func NewWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		items:           make([]T, 0, capacity),
		initialCapacity: capacity,
	}
}

func (s *Stack[T]) PushMany(item ...T) {
	s.items = append(s.items, item...)
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	// Shrink if too much extra initialCapacity is unused
	if cap(s.items) > shrinkCapacityThreshold && cap(s.items) > s.initialCapacity && len(s.items) < cap(s.items)/4 {
		newCap := cap(s.items) / 2
		if newCap < s.initialCapacity {
			newCap = s.initialCapacity
		}
		newItems := make([]T, len(s.items), newCap)
		copy(newItems, s.items)
		s.items = newItems
	}

	return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

// Reset clears the stack but keeps the backing array
func (s *Stack[T]) Reset() {
	s.items = s.items[:0]
}

// Clear creates a new backing array with the initialCapacity
func (s *Stack[T]) Clear() {
	s.items = make([]T, 0, s.initialCapacity)
}
