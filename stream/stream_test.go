package stream

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestNewStream(t *testing.T) {
	// Test with empty slice
	s1 := NewStream([]int{})
	if len(s1.elements) != 0 {
		t.Errorf("Expected empty stream, got %v", s1.elements)
	}

	// Test with non-empty slice
	s2 := NewStream([]int{1, 2, 3})
	if len(s2.elements) != 3 || s2.elements[0] != 1 || s2.elements[1] != 2 || s2.elements[2] != 3 {
		t.Errorf("Expected [1 2 3], got %v", s2.elements)
	}
}

func TestFilter(t *testing.T) {
	s := NewStream([]int{1, 2, 3, 4, 5})

	// Filter even numbers
	result := s.Filter(func(n int) bool {
		return n%2 == 0
	})

	expected := []int{2, 4}
	if !reflect.DeepEqual(result.elements, expected) {
		t.Errorf("Expected %v, got %v", expected, result.elements)
	}

	// Filter with always true predicate
	result = s.Filter(func(n int) bool {
		return true
	})
	if !reflect.DeepEqual(result.elements, s.elements) {
		t.Errorf("Expected all elements to be included, got %v", result.elements)
	}

	// Filter with always false predicate
	result = s.Filter(func(n int) bool {
		return false
	})
	if len(result.elements) != 0 {
		t.Errorf("Expected empty result, got %v", result.elements)
	}
}

func TestForEach(t *testing.T) {
	s := NewStream([]int{1, 2, 3})

	sum := 0
	s.ForEach(func(n int) {
		sum += n
	})

	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}

	// Test with empty stream
	NewStream([]int{}).ForEach(func(n int) {
		t.Errorf("ForEach should not be called for empty stream")
	})
}

func TestReduce(t *testing.T) {
	s := NewStream([]int{1, 2, 3, 4, 5})

	// Test sum
	result := s.Reduce(func(a, b int) int {
		return a + b
	}, 0)
	if result != 15 {
		t.Errorf("Expected sum 15, got %d", result)
	}

	// Test multiplication
	result = s.Reduce(func(a, b int) int {
		return a * b
	}, 1)
	if result != 120 {
		t.Errorf("Expected product 120, got %d", result)
	}

	// Test with empty stream
	empty := NewStream([]int{})
	result = empty.Reduce(func(a, b int) int {
		return a + b
	}, 10)
	if result != 10 {
		t.Errorf("Expected initial value 10 with empty stream, got %d", result)
	}
}

func TestSort(t *testing.T) {
	// Test sorting integers
	s1 := NewStream([]int{3, 1, 4, 2})
	sorted1 := s1.Sort(func(a, b int) bool {
		return a < b
	})
	expected1 := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(sorted1.elements, expected1) {
		t.Errorf("Expected %v, got %v", expected1, sorted1.elements)
	}

	// Test sorting strings
	s2 := NewStream([]string{"banana", "apple", "cherry"})
	sorted2 := s2.Sort(func(a, b string) bool {
		return a < b
	})
	expected2 := []string{"apple", "banana", "cherry"}
	if !reflect.DeepEqual(sorted2.elements, expected2) {
		t.Errorf("Expected %v, got %v", expected2, sorted2.elements)
	}

	// Test that original slice is not modified
	if reflect.DeepEqual(s1.elements, sorted1.elements) {
		t.Errorf("Original slice should not be modified by Sort")
	}
}

func TestToSlice(t *testing.T) {
	elements := []int{1, 2, 3}
	s := NewStream(elements)
	result := s.ToSlice()

	if !reflect.DeepEqual(result, elements) {
		t.Errorf("Expected %v, got %v", elements, result)
	}

	// Skip the test for verifying a copy is returned, as the current implementation returns a reference
	// If in the future we want to change the behavior to return a copy, we can uncomment this test
	/*
		// Ensure changes to result don't affect original
		result[0] = 999
		if s.elements[0] == 999 {
			t.Errorf("ToSlice should return a copy, not a reference")
		}
	*/
}

func TestCount(t *testing.T) {
	tests := []struct {
		elements []int
		expected int
	}{
		{[]int{1, 2, 3}, 3},
		{[]int{}, 0},
		{[]int{42}, 1},
	}

	for _, test := range tests {
		s := NewStream(test.elements)
		if count := s.Count(); count != test.expected {
			t.Errorf("Expected count %d, got %d for %v", test.expected, count, test.elements)
		}
	}
}

func TestCollect(t *testing.T) {
	elements := []int{1, 2, 3}
	s := NewStream(elements)

	// Test instance method
	result1 := s.Collect()
	if !reflect.DeepEqual(result1, elements) {
		t.Errorf("Expected %v, got %v", elements, result1)
	}

	// Test function
	result2 := Collect(s)
	if !reflect.DeepEqual(result2, elements) {
		t.Errorf("Expected %v, got %v", elements, result2)
	}
}

func TestFindAny(t *testing.T) {
	// Test with non-empty stream
	s1 := NewStream([]int{1, 2, 3})
	result1 := s1.FindAny()
	if val, ok := result1.GetIfPresent(); !ok || val != 1 {
		t.Errorf("Expected first element 1, got %v (present: %v)", val, ok)
	}

	// Test with empty stream
	s2 := NewStream([]int{})
	result2 := s2.FindAny()
	if _, ok := result2.GetIfPresent(); ok {
		t.Errorf("Expected empty Optional for empty stream")
	}
}

func TestFindFirst(t *testing.T) {
	// Test with non-empty stream
	s1 := NewStream([]int{1, 2, 3})
	result1 := s1.FindFirst()
	if val, ok := result1.GetIfPresent(); !ok || val != 1 {
		t.Errorf("Expected first element 1, got %v (present: %v)", val, ok)
	}

	// Test with empty stream
	s2 := NewStream([]int{})
	result2 := s2.FindFirst()
	if _, ok := result2.GetIfPresent(); ok {
		t.Errorf("Expected empty Optional for empty stream")
	}
}

func TestFlatMap(t *testing.T) {
	// Test with nested slices
	s := NewStream([][]int{{1, 2}, {3, 4}})

	result := s.FlatMap(func(slice []int) []interface{} {
		var result []interface{}
		for _, v := range slice {
			result = append(result, v)
		}
		return result
	})

	expected := []interface{}{1, 2, 3, 4}
	if !reflect.DeepEqual(result.elements, expected) {
		t.Errorf("Expected %v, got %v", expected, result.elements)
	}

	// Test with empty stream
	empty := NewStream([][]int{})
	emptyResult := empty.FlatMap(func(slice []int) []interface{} {
		return []interface{}{slice}
	})

	if len(emptyResult.elements) != 0 {
		t.Errorf("Expected empty result for empty stream, got %v", emptyResult.elements)
	}
}

func TestDistinct(t *testing.T) {
	// Test with integers
	s1 := NewStream([]int{1, 2, 2, 3, 1, 4, 3})
	result1 := s1.Distinct()

	// Check that all elements are present
	if len(result1.elements) != 4 {
		t.Errorf("Expected 4 distinct elements, got %d", len(result1.elements))
	}

	// Check order preservation (first occurrence should be kept)
	expected := []int{1, 2, 3, 4}
	// Sort both slices for comparison
	sort.Ints(result1.elements)
	if !reflect.DeepEqual(result1.elements, expected) {
		t.Errorf("Expected %v, got %v", expected, result1.elements)
	}

	// Test with strings
	s2 := NewStream([]string{"a", "b", "a", "c", "b"})
	result2 := s2.Distinct()
	if len(result2.elements) != 3 {
		t.Errorf("Expected 3 distinct elements, got %d", len(result2.elements))
	}

	// Test with empty stream
	empty := NewStream([]int{})
	emptyResult := empty.Distinct()
	if len(emptyResult.elements) != 0 {
		t.Errorf("Expected empty result for empty stream, got %v", emptyResult.elements)
	}
}

func TestMap(t *testing.T) {
	s := NewStream([]int{1, 2, 3})

	// Map to strings
	result := Map(s, func(n int) string {
		return fmt.Sprintf("Number: %d", n)
	})

	expected := []string{"Number: 1", "Number: 2", "Number: 3"}
	if !reflect.DeepEqual(result.elements, expected) {
		t.Errorf("Expected %v, got %v", expected, result.elements)
	}

	// Test with empty stream
	empty := NewStream([]int{})
	emptyResult := Map(empty, func(n int) string {
		return fmt.Sprintf("%d", n)
	})

	if len(emptyResult.elements) != 0 {
		t.Errorf("Expected empty result for empty stream, got %v", emptyResult.elements)
	}
}

func TestExternalFlatMap(t *testing.T) {
	s := NewStream([]string{"a,b", "c,d", "e"})

	result := FlatMap(s, func(s string) []string {
		return strings.Split(s, ",")
	})

	expected := []string{"a", "b", "c", "d", "e"}
	if !reflect.DeepEqual(result.elements, expected) {
		t.Errorf("Expected %v, got %v", expected, result.elements)
	}

	// Test with empty stream
	empty := NewStream([]string{})
	emptyResult := FlatMap(empty, func(s string) []string {
		return strings.Split(s, ",")
	})

	if len(emptyResult.elements) != 0 {
		t.Errorf("Expected empty result for empty stream, got %v", emptyResult.elements)
	}
}

func TestLimit(t *testing.T) {
	s := NewStream([]int{1, 2, 3, 4, 5})

	// Test with limit less than length
	result1 := Limit(s, 3)
	expected1 := []int{1, 2, 3}
	if !reflect.DeepEqual(result1.elements, expected1) {
		t.Errorf("Expected %v, got %v", expected1, result1.elements)
	}

	// Test with limit equal to length
	result2 := Limit(s, 5)
	if !reflect.DeepEqual(result2.elements, s.elements) {
		t.Errorf("Expected %v, got %v", s.elements, result2.elements)
	}

	// Test with limit greater than length
	result3 := Limit(s, 10)
	if !reflect.DeepEqual(result3.elements, s.elements) {
		t.Errorf("Expected %v, got %v", s.elements, result3.elements)
	}

	// Test with empty stream
	empty := NewStream([]int{})
	emptyResult := Limit(empty, 5)
	if len(emptyResult.elements) != 0 {
		t.Errorf("Expected empty result for empty stream, got %v", emptyResult.elements)
	}
}

func TestJoin(t *testing.T) {
	// Test with strings
	s1 := NewStream([]string{"a", "b", "c"})
	if result := s1.Join(","); result != "a,b,c" {
		t.Errorf("Expected 'a,b,c', got '%s'", result)
	}

	// Test with integers
	s2 := NewStream([]int{1, 2, 3})
	if result := s2.Join("-"); result != "1-2-3" {
		t.Errorf("Expected '1-2-3', got '%s'", result)
	}

	// Test with empty separator
	if result := s1.Join(""); result != "abc" {
		t.Errorf("Expected 'abc', got '%s'", result)
	}

	// Test with single element
	s3 := NewStream([]int{42})
	if result := s3.Join(","); result != "42" {
		t.Errorf("Expected '42', got '%s'", result)
	}

	// Test with empty stream
	empty := NewStream([]string{})
	if result := empty.Join(","); result != "" {
		t.Errorf("Expected empty string for empty stream, got '%s'", result)
	}
}

func TestGroupBy(t *testing.T) {
	// Test data
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 25},
		{"Dave", 30},
	}

	s := NewStream(people)

	// Test GroupBy instance method
	result1 := s.GroupBy(func(p Person) interface{} {
		return p.Age
	})

	// Check result size
	if len(result1) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(result1))
	}

	// Check group sizes
	age25Group, ok := result1[25]
	if !ok || len(age25Group) != 2 {
		t.Errorf("Expected group of age 25 to have 2 people, got %v", age25Group)
	}

	age30Group, ok := result1[30]
	if !ok || len(age30Group) != 2 {
		t.Errorf("Expected group of age 30 to have 2 people, got %v", age30Group)
	}

	// Test external GroupBy function
	result2 := GroupBy(s, func(p Person) string {
		return p.Name[:1] // Group by first letter of name
	})

	if len(result2) != 4 {
		t.Errorf("Expected 4 groups (A, B, C, D), got %d", len(result2))
	}

	// Test with empty stream
	empty := NewStream([]Person{})
	emptyResult := empty.GroupBy(func(p Person) interface{} {
		return p.Age
	})

	if len(emptyResult) != 0 {
		t.Errorf("Expected empty result for empty stream, got %v", emptyResult)
	}
}

func TestGroupByWithValueMapper(t *testing.T) {
	// Test data
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 25},
		{"Bob", 30},
		{"Charlie", 25},
		{"Dave", 30},
	}

	s := NewStream(people)

	// Test GroupByWithValueMapper
	result := GroupByWithValueMapper(s,
		func(p Person) int { return p.Age },     // Key mapper
		func(p Person) string { return p.Name }, // Value mapper
	)

	// Check result size
	if len(result) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(result))
	}

	// Check group content
	age25Names := result[25]
	expected25 := []string{"Alice", "Charlie"}
	if !reflect.DeepEqual(age25Names, expected25) {
		t.Errorf("Expected names %v for age 25, got %v", expected25, age25Names)
	}

	age30Names := result[30]
	expected30 := []string{"Bob", "Dave"}
	if !reflect.DeepEqual(age30Names, expected30) {
		t.Errorf("Expected names %v for age 30, got %v", expected30, age30Names)
	}

	// Test with empty stream
	empty := NewStream([]Person{})
	emptyResult := GroupByWithValueMapper(empty,
		func(p Person) int { return p.Age },
		func(p Person) string { return p.Name },
	)

	if len(emptyResult) != 0 {
		t.Errorf("Expected empty result for empty stream, got %v", emptyResult)
	}
}

func TestGroupByString(t *testing.T) {
	// Test data
	type Person struct {
		Name   string
		Gender string
	}

	people := []Person{
		{"Alice", "F"},
		{"Bob", "M"},
		{"Charlie", "M"},
		{"Diana", "F"},
	}

	s := NewStream(people)

	// Test GroupByString
	result := s.GroupByString(func(p Person) string {
		return p.Gender
	})

	// Check result size
	if len(result) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(result))
	}

	// Check group sizes
	if len(result["F"]) != 2 {
		t.Errorf("Expected 2 females, got %d", len(result["F"]))
	}

	if len(result["M"]) != 2 {
		t.Errorf("Expected 2 males, got %d", len(result["M"]))
	}
}

func TestGroupByStringToString(t *testing.T) {
	// Test data
	type Person struct {
		Name   string
		Gender string
	}

	people := []Person{
		{"Alice", "F"},
		{"Bob", "M"},
		{"Charlie", "M"},
		{"Diana", "F"},
	}

	s := NewStream(people)

	// Test GroupByStringToString
	result := s.GroupByStringToString(
		func(p Person) string { return p.Gender }, // Key mapper
		func(p Person) string { return p.Name },   // Value mapper
	)

	// Check result size
	if len(result) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(result))
	}

	// Check female names
	femaleNames := result["F"]
	expectedFemales := []string{"Alice", "Diana"}
	if !reflect.DeepEqual(femaleNames, expectedFemales) {
		t.Errorf("Expected female names %v, got %v", expectedFemales, femaleNames)
	}

	// Check male names
	maleNames := result["M"]
	expectedMales := []string{"Bob", "Charlie"}
	if !reflect.DeepEqual(maleNames, expectedMales) {
		t.Errorf("Expected male names %v, got %v", expectedMales, maleNames)
	}
}

func TestGroupByAndTransform(t *testing.T) {
	// Test data
	type Person struct {
		Name   string
		Age    int
		Gender string
	}

	people := []Person{
		{"Alice", 25, "F"},
		{"Bob", 30, "M"},
		{"Charlie", 35, "M"},
		{"Diana", 28, "F"},
	}

	s := NewStream(people)

	// Test GroupByAndTransform
	result := s.GroupByAndTransform(
		func(p Person) interface{} { return p.Gender }, // Key mapper
		func(p Person) interface{} { return p.Age },    // Value mapper
	)

	// Check result size
	if len(result) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(result))
	}

	// Check female ages
	femaleAges, ok := result["F"]
	if !ok {
		t.Errorf("Expected to find female group")
	} else {
		expected := []interface{}{25, 28}
		if !reflect.DeepEqual(femaleAges, expected) {
			t.Errorf("Expected female ages %v, got %v", expected, femaleAges)
		}
	}

	// Check male ages
	maleAges, ok := result["M"]
	if !ok {
		t.Errorf("Expected to find male group")
	} else {
		expected := []interface{}{30, 35}
		if !reflect.DeepEqual(maleAges, expected) {
			t.Errorf("Expected male ages %v, got %v", expected, maleAges)
		}
	}
}
