package stream

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tiagods/go-extras/optional"
)

// Stream represents a sequence of elements
type Stream[T any] struct {
	elements []T
}

// NewStream creates a new Stream from a slice
func NewStream[T any](elements []T) *Stream[T] {
	return &Stream[T]{elements: elements}
}

// Filter filters elements based on a predicate function
func (s *Stream[T]) Filter(predicate func(T) bool) *Stream[T] {
	var result []T
	for _, e := range s.elements {
		if predicate(e) {
			result = append(result, e)
		}
	}
	return NewStream(result)
}

// ForEach applies a function to each element in the stream
func (s *Stream[T]) ForEach(action func(T)) {
	for _, e := range s.elements {
		action(e)
	}
}

// Reduce reduces the elements to a single value using an aggregation function
func (s *Stream[T]) Reduce(reducer func(T, T) T, initialValue T) T {
	accumulator := initialValue
	for _, e := range s.elements {
		accumulator = reducer(accumulator, e)
	}
	return accumulator
}

// Sort sorts the elements based on a comparison function
func (s *Stream[T]) Sort(less func(T, T) bool) *Stream[T] {
	result := make([]T, len(s.elements))
	copy(result, s.elements)
	sort.Slice(result, func(i, j int) bool {
		return less(result[i], result[j])
	})
	return NewStream(result)
}

// ToSlice converts the Stream back to a slice
func (s *Stream[T]) ToSlice() []T {
	return s.elements
}

// Count returns the number of elements in the stream
func (s *Stream[T]) Count() int {
	return len(s.elements)
}

// Collect returns the stream elements as a slice
func (s *Stream[T]) Collect() []T {
	return s.elements
}

// FindAny returns an arbitrary element from the Stream
func (s *Stream[T]) FindAny() optional.Optional[T] {
	if len(s.elements) > 0 {
		return optional.Of(s.elements[0])
	}
	return optional.Empty[T]()
}

// FindFirst returns the first element from the Stream
func (s *Stream[T]) FindFirst() optional.Optional[T] {
	if len(s.elements) > 0 {
		return optional.Of(s.elements[0])
	}
	return optional.Empty[T]()
}

// FlatMap maps each element of the Stream to a new Stream and flattens the result
func (s *Stream[T]) FlatMap(mapper func(T) []interface{}) *Stream[interface{}] {
	var result []interface{}
	for _, e := range s.elements {
		mappedElements := mapper(e)                // Apply the mapping function
		result = append(result, mappedElements...) // Flatten the Stream
	}
	return NewStream(result)
}

// Distinct removes duplicate elements from the Stream
func (s *Stream[T]) Distinct() *Stream[T] {
	uniqueMap := make(map[interface{}]bool)
	var result []T
	for _, e := range s.elements {
		key := fmt.Sprintf("%v", e) // create a unique key based on the element's value
		if _, exists := uniqueMap[key]; !exists {
			uniqueMap[key] = true
			result = append(result, e)
		}
	}
	return NewStream(result)
}

// Map transforms elements from type T to type R
func Map[T any, R any](stream *Stream[T], mapper func(T) R) *Stream[R] {
	var result []R
	for _, e := range stream.elements {
		result = append(result, mapper(e))
	}
	return NewStream(result)
}

// FlatMap transforms elements from type T to []R and flattens the result
func FlatMap[T any, R any](stream *Stream[T], mapper func(T) []R) *Stream[R] {
	var result []R
	for _, e := range stream.elements {
		result = append(result, mapper(e)...)
	}
	return NewStream(result)
}

// Limit returns at most n elements from the stream
func Limit[T any](stream *Stream[T], n int) *Stream[T] {
	if n >= len(stream.elements) {
		return stream
	}
	return NewStream(stream.elements[:n])
}

// Collect returns a slice of elements from the stream
func Collect[T any](stream *Stream[T]) []T {
	return stream.elements
}

// Join concatenates the elements of the Stream into a single string
// If the Stream is empty, returns an empty string
// If the Stream has only one element, returns the string representation of that element
// For multiple elements, concatenates them using the provided separator
func (s *Stream[T]) Join(separator string) string {
	if len(s.elements) == 0 {
		return ""
	}

	if len(s.elements) == 1 {
		return fmt.Sprintf("%v", s.elements[0])
	}

	var result strings.Builder
	for i, e := range s.elements {
		if i > 0 {
			result.WriteString(separator)
		}
		result.WriteString(fmt.Sprintf("%v", e))
	}
	return result.String()
}

// GroupBy groups the Stream elements into a map using interface{} for keys
// This method accepts any key type but is less type-safe
// Use GroupByTyped for type-safe operations
func (s *Stream[T]) GroupBy(keyMapper func(T) interface{}) map[interface{}][]T {
	result := make(map[interface{}][]T)

	for _, e := range s.elements {
		key := keyMapper(e)
		result[key] = append(result[key], e)
	}

	return result
}

// GroupByAndTransform groups the Stream elements into a map and transforms the values
// This method accepts any key and value type but is less type-safe
// Use GroupByTyped and GroupByAndTransformTyped for type-safe operations
func (s *Stream[T]) GroupByAndTransform(keyMapper func(T) interface{}, valueMapper func(T) interface{}) map[interface{}][]interface{} {
	result := make(map[interface{}][]interface{})

	for _, e := range s.elements {
		key := keyMapper(e)
		value := valueMapper(e)
		result[key] = append(result[key], value)
	}

	return result
}

// GroupByString is a convenience method for grouping by string keys
func (s *Stream[T]) GroupByString(keyMapper func(T) string) map[string][]T {
	result := make(map[string][]T)

	for _, e := range s.elements {
		key := keyMapper(e)
		result[key] = append(result[key], e)
	}

	return result
}

// GroupByStringToString is a convenience method for grouping by string keys and transforming to string values
func (s *Stream[T]) GroupByStringToString(keyMapper func(T) string, valueMapper func(T) string) map[string][]string {
	result := make(map[string][]string)

	for _, e := range s.elements {
		key := keyMapper(e)
		value := valueMapper(e)
		result[key] = append(result[key], value)
	}

	return result
}

// GroupBy groups the Stream elements into a map using a key mapper function
// Keys are determined by the keyMapper function
// Elements with the same key are grouped into a slice
func GroupBy[T any, K comparable](stream *Stream[T], keyMapper func(T) K) map[K][]T {
	result := make(map[K][]T)

	for _, e := range stream.elements {
		key := keyMapper(e)
		result[key] = append(result[key], e)
	}

	return result
}

// GroupByWithValueMapper groups the Stream elements into a map using a key mapper function
// Keys are determined by the keyMapper function
// Values are transformed by the valueMapper function before being grouped
func GroupByWithValueMapper[T any, K comparable, V any](stream *Stream[T], keyMapper func(T) K, valueMapper func(T) V) map[K][]V {
	result := make(map[K][]V)

	for _, e := range stream.elements {
		key := keyMapper(e)
		value := valueMapper(e)
		result[key] = append(result[key], value)
	}

	return result
}
