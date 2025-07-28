# enum-go

A library for parameterized enums in Go, inspired by Java's advanced enum model.

## 📖 Motivation

Unlike languages like Java, Go does not have native support for rich enums (with multiple attributes and associated functions). The most common pattern in Go uses `const` with `iota`, which is functional for simple values, but limited when you want to associate properties and behaviors with enums.

The `enum-go` solves this limitation by allowing:

- Enums with varied parameters
- Association of functions to the enum
- Custom search and ordering
- Functional control via the `Optional` structure

## 📦 Installation

```bash
go get github.com/tiagods/enum-go
```

## 📒 Usage Examples

### Enum Definition:

```go
type OperationValue struct {
	Symbol      string
	Order       int
	Apply       func(a, b float64) float64
	Description string
}

var SUM = Enum[OperationValue]{
	Name: "SUM",
	Value: OperationValue{
		Symbol: "+",
		Order: 1,
		Apply: func(a, b float64) float64 { return a + b },
		Description: "Adds two values",
	},
}
```

### Creating and using an EnumSet:

```go
var operations = FromValues([]Enum[OperationValue]{SUM})

sumEnum := operations.FindByName("SUM")
if sumEnum.IsPresent() {
	value, _ := sumEnum.Get()
	result := value.Value.Apply(5, 3)
	fmt.Println("Result:", result)
}
```

### Sorting enums by a custom order:

```go
// Create a set with operations
operationSet := FromValues([]Enum[OperationValue]{SUM, SUBTRACT, MULTIPLY, DIVIDE})

// Sort operations by their order
sortedOperations := operationSet.SortByOrder(func(op OperationValue) int {
    return op.Order
})

// Now operations are sorted by their Order value
for _, op := range sortedOperations.Values() {
    fmt.Printf("%s (%s)\n", op.Name, op.Value.Symbol)
}
```

## 📚 API

### 📌 `Enum[T any]`
Represents a parameterized enum.

- `Name string` — name of the enum.
- `Value T` — associated value.
- `String() string` — returns the name of the enum.
- `Equal(other Enum[T]) bool` — checks if two enums are equal by name.

### 📌 `EnumSet[T any]`
Manages a set of enums.

- `NewEnumSet[T any]() *EnumSet[T]` — creates a new empty enum set.
- `Add(e Enum[T])` — adds an enum to the set.
- `Values() []Enum[T]` — returns all enums in the set.
- `FindByName(name string) Optional[Enum[T]]` — searches by name.
- `SortByOrder(getOrder func(T) int) *EnumSet[T]` — sorts based on the associated value.

### 📌 `Optional[T any]`
Encapsulates optional values.

- `IsPresent() bool` — checks presence.
- `Get() (T, error)` — returns value or error.
- `OrElse(defaultValue T) T` — returns value or default.
- `GetIfPresent() (T, bool)` — returns value and presence.
- `IfPresent(consumer func(T))` — executes an action if value is present.
- `OrElseGet(supplier func() T) T` — returns value or obtains a default from supplier.
- `OrElseThrow(err error) (T, error)` — returns value or provided error.

### 📌 Utility

- `FromValues(values []Enum[T]) *EnumSet[T]` — creates EnumSet from a slice.

## 📈 JSON Serialization

By default, enums are serialized using the `Name` value. To customize the JSON serialization, you can implement custom `MarshalJSON` and `UnmarshalJSON` methods for your enum types.

Here's an example of how to implement custom serialization:

```go
// Example of custom JSON serialization
func (e Enum[T]) MarshalJSON() ([]byte, error) {
    return json.Marshal(e.Name)
}

// Example of custom JSON deserialization
func UnmarshalEnumJSON[T any](data []byte, enumSet *EnumSet[T]) (Enum[T], error) {
    var name string
    if err := json.Unmarshal(data, &name); err != nil {
        return Enum[T]{}, err
    }
    
    optional := enumSet.FindByName(name)
    if enum, found := optional.GetIfPresent(); found {
        return enum, nil
    }
    
    return Enum[T]{}, fmt.Errorf("unknown enum value: %s", name)
}
```

With this implementation, enums will be serialized as string values:

```json
"SUM"
```

## 📃 License

MIT License. 