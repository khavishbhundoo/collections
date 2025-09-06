# Collections
A collection of common data structures for Go in both thread safe and non-thread safe variants.The documentation, 
example usage are available in their own dedicated page. [Benchmark](benchmark/benchmark.txt) results are available for 
all the data structures. 

## Installation 

```bash
go get -u github.com/khavishbhundoo/collections
```

## Design Principles

- General-Purpose Design
  These data structures are designed to perform well enough across a wide range of use cases rather than being optimized 
  for a type of operation.

- Zero-Value Usability

All data structures are immediately usable without explicit initialization.
For example, `var q Queue[int]` is valid and ready to accept elements.

- Explicit Construction
Constructors are provided for flexibility:

`New` creates a structure with default capacity.
`NewWithCapacity` allows pre-allocation when the expected size is known.

- Efficient Memory Management

Slice-backed structures (e.g. Stack, Queue) grow automatically and shrink when appropriate to reclaim memory, minimizing 
long-term footprint.

- Concurrency by Design

Concurrent variants are built with minimal locking to ensure safety across goroutines while prioritizing throughput and 
reducing contention.

## Non Thread safe

[Stack](stack/)

[Queue](queue/)

[Set](set/)

## Thread safe

[Stack](concurrent/stack/)

[Queue](concurrent/queue/)

[Set](concurrent/set/)

[CMap](concurrent/cmap/)