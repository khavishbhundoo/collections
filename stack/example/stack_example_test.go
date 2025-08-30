package stack_example

import (
	"collections/stack"
	"fmt"
)

func ExampleStack() {
	s := stack.New[int]()
	s.PushMany(1, 2)
	s.Push(3)

	val, ok := s.Pop()
	fmt.Println(val, ok)
	fmt.Println(s.Size())
	peek, ok := s.Peek()
	fmt.Println(peek, ok)
	fmt.Println(s.Size())
	s.Pop()
	s.Pop()
	val, ok = s.Pop()
	fmt.Println(val, ok)
	peek, ok = s.Peek()
	fmt.Println(peek, ok)

	// Output:
	// 3 true
	// 2
	// 2 true
	// 2
	// 0 false
	// 0 false
}
