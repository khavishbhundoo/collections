package set_test

import (
	"fmt"
	"sync"

	"github.com/khavishbhundoo/collections/concurrent/set"
)

func ExampleSet() {
	s := set.New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)
	fmt.Println("After Add:", s.Len())

	//The zero value of Set[T] is ready to use without initialization
	var s3 set.Set[int]
	s3.Add(1)
	fmt.Println("After Add:", s3.Len())

	// Add multiple elements at once
	s.AddMany(3, 4, 5)
	fmt.Println("After AddMany:", s.Len())

	// Check if an element exists
	fmt.Println("Contains 3?", s.Contains(3))
	fmt.Println("Contains 10?", s.Contains(10))

	// Remove an element
	s.Remove(2)
	fmt.Println("After Remove 2:", s.Len())

	// Reset the set (keeps capacity)
	s.Reset()
	fmt.Println("After Reset:", s.Len())

	// Clear the set (resets map to initial capacity)
	s.Clear()
	fmt.Println("After Clear:", s.Len())

	var wg sync.WaitGroup
	cs := set.New[int]()
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			cs.Add(v)
		}(i)
	}
	wg.Wait()
	fmt.Println(cs.Len())

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			cs.Add(v)
			cs.Remove(v)
		}(i)
	}
	wg.Wait()
	fmt.Println(cs.Len())

	// Output:
	// After Add: 3
	// After Add: 1
	// After AddMany: 5
	// Contains 3? true
	// Contains 10? false
	// After Remove 2: 4
	// After Reset: 0
	// After Clear: 0
	// 3
	// 0
}
