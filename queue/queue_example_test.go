package queue_test

import (
	"collections/queue"
	"fmt"
)

func ExampleQueue() {
	s := queue.New[int]()
	s.PushMany(1, 2)
	s.Push(3)

	s2 := queue.NewWithCapacity[int](4)
	s2.PushMany(1, 2, 3, 4)
	s2.Push(3)
	val, ok := s2.Pop()
	fmt.Println(val, ok)

	val, ok = s.Pop()
	fmt.Println(val, ok)
	fmt.Println(s.Len())
	peek, ok := s.Peek()
	fmt.Println(peek, ok)
	fmt.Println(s.Len())
	s.Pop()
	s.Pop()
	val, ok = s.Pop()
	fmt.Println(val, ok)
	peek, ok = s.Peek()
	fmt.Println(peek, ok)

	// The zero value of Queue[T] is ready to use without initialization
	var s3 queue.Queue[int]
	s3.Push(1)
	val, ok = s3.Pop()
	fmt.Println(val, ok)

	// Output:
	// 1 true
	// 1 true
	// 2
	// 2 true
	// 2
	// 0 false
	// 0 false
	// 1 true
}
