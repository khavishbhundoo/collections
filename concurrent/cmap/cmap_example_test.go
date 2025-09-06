package cmap_test

import (
	"collections/concurrent/cmap"
	"fmt"
	"sort"
	"sync"
)

func ExampleCMap() {

	var m = cmap.New[string, int]()

	// Insert values
	m.Set("Go", 1)
	m.Set("C#", 2)

	// Retrieve values
	if v, ok := m.Get("Go"); ok {
		fmt.Println("Go =", v)
	}

	// The zero value of CMap[K,V] is ready to use without initialization
	var n cmap.CMap[string, int]
	n.Set("Go", 1)
	_, _ = n.Get("Go")

	// Check existence
	fmt.Println("Contains C#?", m.Contains("C#"))
	fmt.Println("Len =", m.Len())

	// Get all keys
	keys := m.Keys()
	sort.Strings(keys)
	fmt.Println("Keys:", keys)

	// Delete a key
	m.Delete("Go")
	fmt.Println("Contains Go after delete?", m.Contains("Go"))

	// Reset map
	m.Reset()
	fmt.Println("Len after Reset:", m.Len())

	// Clear map
	m.Set("Rust", 3)
	m.Clear()
	fmt.Println("Len after Clear:", m.Len())

	// Output:
	// Go = 1
	// Contains C#? true
	// Len = 2
	// Keys: [C# Go]
	// Contains Go after delete? false
	// Len after Reset: 0
	// Len after Clear: 0
}

func ExampleCMap_concurrent() {
	m := cmap.NewWithCapacity[string, int](50)
	var wg sync.WaitGroup

	// Start 5 goroutines that write to the map
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("w%d_%d", id, j)
				m.Set(key, j)
			}
		}(i)
	}

	// Wait for all writers to finish before inspecting the map
	wg.Wait()

	// Now read concurrently after all writes are done
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			keys := m.Keys()
			fmt.Printf("Reader %d sees %d keys\n", id, len(keys))
		}(i)
	}

	wg.Wait()

	// Final state of the map
	fmt.Println("Final map length:", m.Len())
}
