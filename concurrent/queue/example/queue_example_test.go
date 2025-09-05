package queue_example

import (
	"collections/concurrent/queue"
	"fmt"
	"sync"
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

	var s3 queue.Queue[int]
	s3.Push(1)
	val, ok = s3.Pop()
	fmt.Println(val, ok)

	var wg sync.WaitGroup
	cs := queue.New[int]()
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			cs.Push(v)
		}(i)
	}
	wg.Wait()
	fmt.Println(cs.Len())

	// Output:
	// 1 true
	// 1 true
	// 2
	// 2 true
	// 2
	// 0 false
	// 0 false
	// 1 true
	// 3
}
