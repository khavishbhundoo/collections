package queue

// Queue is a generic, non-thread-safe FIFO (first-in-first-out) queue
// implementation backed by a dynamically resizing slice.The zero value
// of Queue[T] is ready to use without initialization
//
// Use New() or NewWithCapacity() if you prefer an explicit constructor
// or want to set an initial capacity.
// If you do need thread-safety, use the collections/concurrent/queue package instead.
type Queue[T any] struct {
	items           []T
	initialCapacity int
}

// shrinkCapacityThreshold defines the minimum slice capacity before
// shrink operations are considered. Avoids aggressive shrinking for
// small queues that would just grow again.
// shrinkCapacityThreshold is just one parameter when deciding to shrink
// the underlying array
const shrinkCapacityThreshold = 16

// New creates an empty queue of type T with no pre-allocated capacity.
// Use this when you don't know in advance how many elements you will push.
// This is equivalent to creating a queue as `var q queue.Queue[int]`
//
// Example:
//
//	q := queue.New[int]()
func New[T any]() *Queue[T] {
	return &Queue[T]{
		items:           []T{},
		initialCapacity: 0,
	}
}

// NewWithCapacity creates an empty queue of type T with a pre-allocated
// capacity. This avoids repeated allocations if you know roughly how
// many elements you’ll push.
//
// Example:
//
//	s := queue.NewWithCapacity[int](10)
func NewWithCapacity[T any](capacity int) *Queue[T] {
	return &Queue[T]{
		items:           make([]T, 0, capacity),
		initialCapacity: capacity,
	}
}

// PushMany pushes one or more items onto the queue in order.
// Equivalent to calling Push repeatedly but more efficient
// when adding multiple elements.
//
// Example:
//
//	q.PushMany(1, 2, 3)
func (q *Queue[T]) PushMany(item ...T) {
	q.items = append(q.items, item...)
}

// Push adds a single item to the end of the queue.
//
// Example:
//
//	q.Push(42)
func (q *Queue[T]) Push(item T) {
	q.items = append(q.items, item)
}

// Pop removes and returns the element in front of the queue.
// The boolean return is false if the queue is empty.
// The queue may shrink its capacity automatically if
// it has grown significantly and is mostly empty.
//
// Example:
//
//	value, ok := q.Pop()
//	if ok { fmt.Println(value) }
func (q *Queue[T]) Pop() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]

	// Reduce capacity if:
	//   - slice is larger than the shrink threshold (avoid tiny slice reallocations),
	//   - current capacity exceeds 2× the initial capacity (if any),
	//   - and fewer than 12.5% of elements are in use (cap/8).
	//
	// Why 1/8 instead of 1/4?
	//   Using 1/4 is fine for general use, but in tight push/pop workloads
	//   it may trigger frequent grow/shrink oscillations. Using 1/8 shrinks
	//   only when the queue is significantly underutilized.
	//
	// Why halve capacity?
	//   Halving avoids repeated reallocations while still reclaiming
	//   unused memory proportionally. It balances memory efficiency and speed.
	capNow := cap(q.items)
	if capNow > shrinkCapacityThreshold &&
		(q.initialCapacity == 0 || capNow > q.initialCapacity*2) &&
		len(q.items) < capNow/8 {

		newCap := capNow / 2
		if q.initialCapacity > 0 && newCap < q.initialCapacity {
			newCap = q.initialCapacity
		}
		if newCap != capNow { // only shrink if capacity actually changes
			newItems := make([]T, len(q.items), newCap)
			copy(newItems, q.items)
			q.items = newItems
		}
	}

	return item, true
}

// Peek returns the front of the queue without removing it.
// The boolean return is false if the queue is empty.
//
// Example:
//
//	value, ok := q.Peek()
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	return q.items[0], true
}

// Len returns the current number of items in the queue.
//
// Example:
//
//	n := s.Len()
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// Reset clears all items but keeps the current capacity
// of the underlying slice. This is faster than Clear()
// when you expect to reuse the same queue size.
//
// Example:
//
//	q.Reset()
func (q *Queue[T]) Reset() {
	q.items = q.items[:0]
}

// Clear removes all items and reallocates a slice with
// the initial capacity (if any). Use this to shrink the
// backing array explicitly.
//
// Example:
//
//	q.Clear()
func (q *Queue[T]) Clear() {
	q.items = make([]T, 0, q.initialCapacity)
}
