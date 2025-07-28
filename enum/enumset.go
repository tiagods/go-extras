package enum

import (
	"errors"
	"sort"

	"github.com/tiagods/enum-go/optional"
)

// ErrEnumNotFound is returned when an enum could not be found in the set
var ErrEnumNotFound = errors.New("enum not found")

// EnumSet is a collection of Enum values of the same type
type EnumSet[T any] struct {
	values []Enum[T]
}

// NewEnumSet creates a new empty EnumSet
func NewEnumSet[T any]() *EnumSet[T] {
	return &EnumSet[T]{values: []Enum[T]{}}
}

// Add appends an enum to the set
func (s *EnumSet[T]) Add(e Enum[T]) {
	s.values = append(s.values, e)
}

// Values returns all enums in the set
func (s *EnumSet[T]) Values() []Enum[T] {
	return s.values
}

// FindByName searches for an enum by its name and returns an Optional containing
// the enum if found, or an empty Optional if not found
func (s *EnumSet[T]) FindByName(name string) optional.Optional[Enum[T]] {
	for _, v := range s.values {
		if v.Name == name {
			return optional.Of(v)
		}
	}
	return optional.Empty[Enum[T]]()
}

// SortByOrder sorts the enums in the set using the provided ordering function
// and returns the same set for method chaining
func (s *EnumSet[T]) SortByOrder(getOrder func(T) int) *EnumSet[T] {
	sort.SliceStable(s.values, func(i, j int) bool {
		return getOrder(s.values[i].Value) < getOrder(s.values[j].Value)
	})
	return s
}

// FromValues creates a new EnumSet from a slice of Enum values
func FromValues[T any](values []Enum[T]) *EnumSet[T] {
	return &EnumSet[T]{values: values}
}
