# Optional

A Go package that implements the Optional pattern for handling values that may or may not be present. Similar to Java's Optional or Rust's Option.

## Features

- Create Optional with present value
- Create empty Optional
- Check if value is present
- Safe value retrieval
- Default values for absent values
- Functional programming patterns
- Generic type support

## Usage

```go
import "github.com/tiagods/go-extras/optional"

// Create an Optional with a value
opt1 := optional.Of("Hello, World!")
if opt1.IsPresent() {
    // Value is present
}

// Get value if present
if val, ok := opt1.GetIfPresent(); ok {
    fmt.Println("Value:", val)
}

// Create an empty Optional
opt2 := optional.Empty[string]()

// Use OrElse to provide a default value
value := opt2.OrElse("default value")

// Example of OfNullable
isZero := func(s string) bool { return s == "" }
nullableOpt := optional.OfNullable("text", isZero)
```

## API

### Creation
- `Of[T any](value T) Optional[T]` - Create an Optional with a present value
- `Empty[T any]() Optional[T]` - Create an empty Optional
- `OfNullable[T comparable](value T, isZero func(T) bool) Optional[T]` - Create an Optional from a value that might be zero/null

### Operations
- `(o Optional[T]) IsPresent() bool` - Check if a value is present
- `(o Optional[T]) GetIfPresent() (T, bool)` - Return the value and a boolean indicating if it's present
- `(o Optional[T]) Get() (T, error)` - Return the value or an error if not present
- `(o Optional[T]) OrElse(defaultValue T) T` - Return the value or a default if absent
- `(o Optional[T]) OrElseGet(supplier func() T) T` - Return the value or get a default from a function
- `(o Optional[T]) OrElseThrow(err error) (T, error)` - Return the value or a specific error
- `(o Optional[T]) IfPresent(consumer func(T))` - Execute an action if the value is present

## Examples

For comprehensive examples, check the [examples directory](examples/main.go), which includes:

1. Basic Optional usage
2. Repository pattern with Optional
3. Chaining operations
4. Safe parsing with Optional
5. Handling nil and zero values

## Common Patterns

### Error handling

```go
// Traditional Go error handling
result, err := someFunction()
if err != nil {
    // handle error
}
// use result

// With Optional
resultOpt := tryGetResult()
if value, ok := resultOpt.GetIfPresent(); ok {
    // use value
} else {
    // handle missing value
}
```

### Default values

```go
// Traditional Go approach
var result string
if value, ok := someMap["key"]; ok {
    result = value
} else {
    result = "default"
}

// With Optional
result := findInMap("key").OrElse("default")
```

### Conditional execution

```go
// Traditional Go approach
if user, ok := getUser(id); ok {
    processUser(user)
}

// With Optional
getUserOptional(id).IfPresent(func(user User) {
    processUser(user)
})
```

## Error constants

- `ErrNoValuePresent` - Returned when trying to get a value from an empty Optional 