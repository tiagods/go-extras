package enum

import "encoding/json"

// Enum is a generic enumeration type that associates a name with a value.
// T can be any type, allowing for flexible enum implementations.
type Enum[T any] struct {
	Name  string
	Value T
}

// String returns the name of the enum, implementing the Stringer interface.
func (e Enum[T]) String() string {
	return e.Name
}

// MarshalJSON implements the json.Marshaler interface.
func (e Enum[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Name)
}

// Equal checks if two enum instances are equal by comparing their names.
func (e Enum[T]) Equal(other Enum[T]) bool {
	return e.Name == other.Name
}
