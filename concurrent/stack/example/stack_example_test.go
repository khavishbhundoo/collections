package stack_example

import (
	"collections/concurrent/stack"
	"fmt"
	"sync"
)

func ExampleStack() {
	s := stack.New[int]()
	s.PushMany(1, 2)
	s.Push(3)

	s2 := stack.NewWithCapacity[int](4)
	s2.PushMany(1, 2, 3, 4)
	s2.Push(3)
	val, ok := s2.Pop()
	fmt.Println(val, ok)

	val, ok = s.Pop()
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

	var s3 stack.Stack[int]
	s3.Push(1)
	val, ok = s3.Pop()
	fmt.Println(val, ok)

	var wg sync.WaitGroup
	cs := stack.New[int]()
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			cs.Push(v)
		}(i)
	}
	wg.Wait()
	fmt.Println(cs.Size())

	// Output:
	// 3 true
	// 3 true
	// 2
	// 2 true
	// 2
	// 0 false
	// 0 false
	// 1 true
	// 3
}
