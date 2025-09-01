package set_example

import (
	"collections/set"
	"fmt"
)

func ExampleSet() {
	s := set.New[int]()

	s.Add(1)
	s.Add(2)
	s.Add(3)
	fmt.Println("After Add:", s.Len())

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

	// Output:
	// After Add: 3
	// After AddMany: 5
	// Contains 3? true
	// Contains 10? false
	// After Remove 2: 4
	// After Reset: 0
	// After Clear: 0
}
