package stack

// Stack is a generic, non-thread-safe LIFO (last-in-first-out) stack
// implementation backed by a dynamically resizing slice.
//
// The zero value of Stack[T] is ready to use without initialization:
//
//	var s stack.Stack[int]
//	s.Push(1)
//
// Use New() or NewWithCapacity() if you prefer an explicit constructor
// or want to set an initial capacity.
// If you do need thread-safety, use the collections/concurrent/stack package instead.
type Stack[T any] struct {
	items           []T
	initialCapacity int
}

// shrinkCapacityThreshold defines the minimum slice capacity before
// shrink operations are considered. Avoids aggressive shrinking for
// small stacks that would just grow again.
// shrinkCapacityThreshold is just one parameter when deciding to shrink
// the underlying array
const shrinkCapacityThreshold = 16

// New creates an empty stack of type T with no pre-allocated capacity.
// Use this when you don't know in advance how many elements you will push.
// This is equivalent to creating a stack as `var s stack.Stack[int]`
//
// Example:
//
//	s := stack.New[int]()
func New[T any]() *Stack[T] {
	return &Stack[T]{
		items:           []T{},
		initialCapacity: 0,
	}
}

// NewWithCapacity creates an empty stack of type T with a pre-allocated
// capacity. This avoids repeated allocations if you know roughly how
// many elements you’ll push.
//
// Example:
//
//	s := stack.NewWithCapacity[int](10)
func NewWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		items:           make([]T, 0, capacity),
		initialCapacity: capacity,
	}
}

// PushMany pushes one or more items onto the stack in order.
// Equivalent to calling Push repeatedly but more efficient
// when adding multiple elements.
//
// Example:
//
//	s.PushMany(1, 2, 3)
func (s *Stack[T]) PushMany(item ...T) {
	s.items = append(s.items, item...)
}

// Push adds a single item to the top of the stack.
//
// Example:
//
//	s.Push(42)
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top element of the stack.
// The boolean return is false if the stack is empty.
// The stack may shrink its capacity automatically if
// it has grown significantly and is mostly empty.
//
// Example:
//
//	value, ok := s.Pop()
//	if ok { fmt.Println(value) }
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	// Reduce capacity if:
	//   - slice is larger than the shrink threshold (avoid tiny slices reallocation),
	//   - current capacity exceeds the initial capacity (never shrink below the original size),
	//   - and fewer than 25% of elements are in use.
	//
	// Why 1/4 instead of 1/2?
	//   Using 1/2 would trigger frequent grow/shrink oscillations if the
	//   stack size hovers around the boundary. A 1/4 threshold provides
	//   a buffer: it shrinks only when the slice is significantly underused.
	//
	// Why halve capacity?
	//   Halving avoids repeated reallocations while still reclaiming
	//   unused memory proportionally. It’s a balance between performance
	//   (fewer allocations) and memory efficiency.
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

// Peek returns the top element of the stack without removing it.
// The boolean return is false if the stack is empty.
//
// Example:
//
//	value, ok := s.Peek()
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the current number of items in the stack.
//
// Example:
//
//	n := s.Len()
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Reset clears all items but keeps the current capacity
// of the underlying slice. This is faster than Clear()
// when you expect to reuse the same stack size.
//
// Example:
//
//	s.Reset()
func (s *Stack[T]) Reset() {
	s.items = s.items[:0]
}

// Clear removes all items and reallocates a slice with
// the initial capacity (if any). Use this to shrink the
// backing array explicitly.
//
// Example:
//
//	s.Clear()
func (s *Stack[T]) Clear() {
	s.items = make([]T, 0, s.initialCapacity)
}
