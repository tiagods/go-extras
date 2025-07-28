# Stream

A Go implementation of Java-like Stream API for functional collection processing.

## Overview

The Stream package provides a powerful and flexible way to process collections of data in a functional style, inspired by Java's Stream API. It allows for fluent operations like filtering, mapping, reducing, and more.

## Features

- Generic implementation with Go 1.18+ generics
- Fluent API for chaining operations
- Parallel stream processing with controllable concurrency
- Rich set of operations:
  - Filter
  - Map/FlatMap
  - Reduce
  - Sort
  - Distinct
  - Limit
  - Find operations
  - GroupBy with various options
  - Join

## Installation

```bash
go get github.com/tiagods/go-extras/stream
```

## Usage

### Basic Operations

```go
import "github.com/tiagods/go-extras/stream"

// Create a stream from a slice
s := stream.NewStream([]int{1, 2, 3, 4, 5})

// Filter even numbers
evenNumbers := s.Filter(func(n int) bool {
    return n%2 == 0
})

// Get result as slice
result := evenNumbers.ToSlice() // [2, 4]
```

### Chaining Operations

```go
result := stream.NewStream([]int{1, 2, 3, 4, 5}).
    Filter(func(n int) bool {
        return n%2 != 0 // Odd numbers
    }).
    Filter(func(n int) bool {
        return n > 2 // Greater than 2
    }).
    ToSlice() // [3, 5]
```

### Map and Transform

```go
// Map integers to strings
result := stream.Map(stream.NewStream([]int{1, 2, 3}), func(n int) string {
    return fmt.Sprintf("Number: %d", n)
}).ToSlice() // ["Number: 1", "Number: 2", "Number: 3"]
```

### Reduce

```go
// Sum all numbers
sum := stream.NewStream([]int{1, 2, 3, 4, 5}).
    Reduce(func(a, b int) int {
        return a + b
    }, 0) // 15
```

### Sort

```go
// Sort strings
sorted := stream.NewStream([]string{"banana", "apple", "cherry"}).
    Sort(func(a, b string) bool {
        return a < b
    }).ToSlice() // ["apple", "banana", "cherry"]
```

### Find Operations

```go
// Find first element
firstElement := stream.NewStream([]int{1, 2, 3}).FindFirst()
if value, found := firstElement.GetIfPresent(); found {
    fmt.Println("First element:", value)
}
```

### GroupBy Operations

```go
// Group by a key
people := []Person{{"Alice", 25}, {"Bob", 30}, {"Charlie", 25}}
byAge := stream.NewStream(people).GroupByString(func(p Person) string {
    return fmt.Sprintf("%d", p.Age)
})
// Result: {"25": [{"Alice", 25}, {"Charlie", 25}], "30": [{"Bob", 30}]}
```

### Parallel Processing

```go
// Process in parallel with 4 goroutines
results := stream.NewStream(data).ParallelStream(func(item int) interface{} {
    return expensiveOperation(item)
}, 4)
```

## Complete Example

```go
package main

import (
    "fmt"
    "strings"

    "github.com/tiagods/go-extras/stream"
)

type Person struct {
    Name   string
    Age    int
    City   string
}

func main() {
    people := []Person{
        {"John", 30, "New York"},
        {"Alice", 25, "Chicago"},
        {"Bob", 35, "New York"},
        {"Carol", 28, "Chicago"},
        {"Dave", 40, "Boston"},
    }
    
    // Create a stream from the slice of people
    s := stream.NewStream(people)
    
    // Find people over 30, sort by name, limit to 2
    result := s.
        Filter(func(p Person) bool {
            return p.Age > 30
        }).
        Sort(func(a, b Person) bool {
            return a.Name < b.Name
        }).
        ToSlice()
    
    fmt.Println("People over 30:")
    for _, p := range result {
        fmt.Printf("- %s (%d) from %s\n", p.Name, p.Age, p.City)
    }
    
    // Group people by city
    byCity := s.GroupByString(func(p Person) string {
        return p.City
    })
    
    fmt.Println("\nPeople by city:")
    for city, cityPeople := range byCity {
        names := stream.Map(stream.NewStream(cityPeople), 
            func(p Person) string { return p.Name }).ToSlice()
        fmt.Printf("%s: %s\n", city, strings.Join(names, ", "))
    }
}
```

## API Reference

### Creation

- `NewStream[T any](elements []T) *Stream[T]`
  Creates a new stream from a slice.

### Intermediate Operations

- `Filter(predicate func(T) bool) *Stream[T]`
  Filters elements based on a predicate.

- `Map[R any](mapper func(T) R) *Stream[R]`
  Maps elements to a new type.

- `FlatMap[R any](mapper func(T) []R) *Stream[R]`
  Maps each element to multiple elements and flattens the result.

- `Sort(less func(T, T) bool) *Stream[T]`
  Sorts elements based on a comparator.

- `Distinct() *Stream[T]`
  Removes duplicate elements.

- `Limit(n int) *Stream[T]`
  Limits the stream to at most n elements.

### Terminal Operations

- `ToSlice() []T`
  Collects stream elements into a slice.

- `ForEach(action func(T))`
  Performs an action on each element.

- `Reduce(reducer func(T, T) T, initialValue T) T`
  Reduces elements to a single value.

- `FindFirst() Optional[T]`
  Returns the first element if present.

- `FindAny() Optional[T]`
  Returns any element if present.

- `Count() int`
  Returns the number of elements.

- `Join(separator string) string`
  Joins elements into a string.

- `GroupBy(keyMapper func(T) interface{}) map[interface{}][]T`
  Groups elements by a key.

### Parallel Processing

- `ParallelStream(mapper func(T) interface{}, maxGoroutines int) *Stream[interface{}]`
  Processes elements in parallel.

## License

MIT License 