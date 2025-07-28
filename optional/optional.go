package optional

import "errors"

// Common errors returned by the package
var (
	ErrNoValuePresent = errors.New("no value present")
)

// Optional represents a value that may or may not be present.
// It is similar to Java's Optional or Rust's Option type.
type Optional[T any] struct {
	value T
	found bool
}

// Of creates an Optional with a present value
func Of[T any](value T) Optional[T] {
	return Optional[T]{value: value, found: true}
}

// Empty creates an empty Optional
func Empty[T any]() Optional[T] {
	return Optional[T]{found: false}
}

// OfNullable creates an Optional from a value that might be null (zero value in Go)
// If the value equals the zero value for its type and isZero returns true,
// it returns an empty Optional
func OfNullable[T comparable](value T, isZero func(T) bool) Optional[T] {
	if isZero(value) {
		return Empty[T]()
	}
	return Of(value)
}

// GetIfPresent returns the value and a boolean indicating if the value is present
func (o Optional[T]) GetIfPresent() (T, bool) {
	if o.found {
		return o.value, true
	}
	var empty T
	return empty, false
}

// IsPresent returns true if the value is present
func (o Optional[T]) IsPresent() bool {
	return o.found
}

// Get returns the value and an error if the value is not present
func (o Optional[T]) Get() (T, error) {
	if o.found {
		return o.value, nil
	}
	var empty T
	return empty, ErrNoValuePresent
}

// OrElse returns the value if present, or the provided default value
func (o Optional[T]) OrElse(defaultValue T) T {
	if o.found {
		return o.value
	}
	return defaultValue
}

// IfPresent executes an action if the value is present
func (o Optional[T]) IfPresent(consumer func(T)) {
	if o.found {
		consumer(o.value)
	}
}

// OrElseGet returns the value if present, or obtains a default value from a supplier function
func (o Optional[T]) OrElseGet(supplier func() T) T {
	if o.found {
		return o.value
	}
	return supplier()
}

// OrElseThrow returns the value if present, or returns the provided error
func (o Optional[T]) OrElseThrow(err error) (T, error) {
	if o.found {
		return o.value, nil
	}
	var empty T
	return empty, err
}
