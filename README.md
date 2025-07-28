# go-extras

A Go library that provides three powerful utilities:
1. **Rich Parameterized Enums**: Type-safe enum patterns inspired by Java
2. **Optional Pattern**: A robust way to handle nullable/absent values
3. **Stream Processing**: Functional-style operations on collections inspired by Java Streams

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Packages](#packages)
  - [Enum Package](#enum-package)
  - [Optional Package](#optional-package)
  - [Stream Package](#stream-package)
- [Examples](#examples)
- [Documentation](#documentation)
- [License](#license)

## Features

### Enum Features
- Type-safe parameterized enums with any associated values
- Association of functions to enum values
- Custom search and ordering capabilities
- JSON serialization support

### Optional Features
- Type-safe generic implementation for optional values
- Functional programming patterns (map, orElse, ifPresent)
- Zero-allocation empty optionals
- Chainable operations
- Clean error handling

### Stream Features
- Generic functional-style collection processing
- Fluent API for chaining operations (filter, map, reduce, etc.)
- Parallel stream processing with controllable concurrency
- Grouping, sorting, and aggregation operations
- Java-like stream pattern for Go

## Installation

```bash
go get github.com/tiagods/go-extras
```

## Packages

### Enum Package

`go-extras` provides a flexible way to create rich enums in Go with multiple attributes, associated functions, and custom behavior. It solves the limitation of Go's standard `const` with `iota` approach.

#### Quick Example - Enum

```go
// Define enum structure
type OperationValue struct {
    Symbol string
    Apply  func(a, b float64) float64
}

// Create enum instances
var (
    SUM = Enum[OperationValue]{
        Name: "SUM",
        Value: OperationValue{
            Symbol: "+",
            Apply: func(a, b float64) float64 { return a + b },
        },
    }
    
    MULTIPLY = Enum[OperationValue]{
        Name: "MULTIPLY",
        Value: OperationValue{
            Symbol: "*",
            Apply: func(a, b float64) float64 { return a * b },
        },
    }
)

// Create and use an EnumSet
operations := FromValues([]Enum[OperationValue]{SUM, MULTIPLY})
result := operations.FindByName("SUM").OrElse(SUM).Value.Apply(5, 3)
// result = 8
```

### Optional Package

The `optional` package implements the Optional pattern for handling values that may or may not be present, similar to Java's Optional or Rust's Option. It provides a cleaner way to handle nullable values in Go.

#### Quick Example - Optional

```go
// Create Optional values
present := optional.Of("hello")
empty := optional.Empty[string]()

// Check presence
if present.IsPresent() {
    // Do something
}

// Get value with default
value := empty.OrElse("default")

// Execute code conditionally 
present.IfPresent(func(v string) {
    fmt.Println("Value is:", v)
})

// With error handling
value, err := empty.OrElseThrow(errors.New("value required"))
```

### Stream Package

The `stream` package provides a functional-style collection processing API inspired by Java's Stream API. It allows for a fluent, chainable approach to data transformations, filtering, and aggregation.

#### Quick Example - Stream

```go
import "github.com/tiagods/go-extras/stream"

// Create a stream from a slice
s := stream.NewStream([]int{1, 2, 3, 4, 5})

// Chain operations: filter even numbers, double them, and collect
result := s.
    Filter(func(n int) bool {
        return n%2 == 0 // Only even numbers
    }).
    Sort(func(a, b int) bool {
        return a > b // Descending order
    }).
    ToSlice() // [4, 2]

// Process elements in parallel
parallelResult := stream.NewStream([]int{1, 2, 3, 4, 5}).
    ParallelStream(func(n int) interface{} {
        return n * 2 // Double each number in parallel
    }, 4) // Using 4 goroutines
```

## Examples

### Enum Use Cases

#### Rich Enums with Associated Functions

```go
// Color enum with RGB values and utility methods
var (
    RED = Enum[ColorValue]{
        Name: "RED",
        Value: ColorValue{
            Hex: "#FF0000",
            RGB: [3]int{255, 0, 0},
            Darken: func(amount float64) string {
                // Implementation
            },
        },
    }
)
```

#### Sorting and Filtering

```go
// Sort enums by a custom order
sortedOperations := operationSet.SortByOrder(func(op OperationValue) int {
    return op.Order
})

// Find an enum by name
multiplyOp := operationSet.FindByName("MULTIPLY")
```

### Optional Use Cases

#### Repository Pattern

```go
// Returns an Optional instead of (User, bool) or (User, error)
func (r *UserRepository) FindUserByID(id int) optional.Optional[User] {
    if user, ok := r.users[id]; ok {
        return optional.Of(user)
    }
    return optional.Empty[User]()
}

// Usage
userOpt := repo.FindUserByID(123)
userOpt.IfPresent(func(user User) {
    processUser(user)
})
```

#### Safe Parsing

```go
func ParseInt(s string) optional.Optional[int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return optional.Empty[int]()
    }
    return optional.Of(n)
}

// Usage 
result := ParseInt("42").OrElse(0)
```

### Stream Use Cases

#### Data Transformation Pipeline

```go
// Process a collection of users
users := []User{...}
result := stream.NewStream(users).
    Filter(func(u User) bool {
        return u.Age >= 18 // Only adults
    }).
    Sort(func(a, b User) bool {
        return a.Name < b.Name // Sort by name
    }).
    ToSlice()
```

#### Parallel Data Processing

```go
// Process items in parallel with controllable concurrency
items := []Item{...}
results := stream.NewStream(items).
    ParallelStream(func(item Item) interface{} {
        return processItemIntensively(item) // CPU intensive work
    }, runtime.NumCPU()) // Use all available CPUs

// Group results
groupedByCategory := stream.NewStream(results.ToSlice()).
    GroupByString(func(r interface{}) string {
        return r.(ProcessedItem).Category
    })
```

## Documentation

For detailed documentation about each package:

- [Enum Package Documentation](enum/README.md)
- [Optional Package Documentation](optional/README.md)
- [Stream Package Documentation](stream/README.md)

## License

MIT License 