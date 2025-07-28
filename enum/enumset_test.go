package enum

import (
	"testing"
)

// TestEnum is an enum type used for testing
type TestEnum int

// Constants for TestEnum
const (
	FIRST  TestEnum = 1
	SECOND TestEnum = 2
	THIRD  TestEnum = 3
)

// Test enums with TestEnum type
var (
	TestFirst  = Enum[TestEnum]{Name: "FIRST", Value: FIRST}
	TestSecond = Enum[TestEnum]{Name: "SECOND", Value: SECOND}
	TestThird  = Enum[TestEnum]{Name: "THIRD", Value: THIRD}
)

// TestEnumString tests the String method of Enum
func TestEnumString(t *testing.T) {
	tests := []struct {
		name     string
		enum     Enum[TestEnum]
		expected string
	}{
		{"First enum", TestFirst, "FIRST"},
		{"Second enum", TestSecond, "SECOND"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enum.String(); got != tt.expected {
				t.Errorf("Enum.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestEnumEqual tests the Equal method of Enum
func TestEnumEqual(t *testing.T) {
	tests := []struct {
		name     string
		enum1    Enum[TestEnum]
		enum2    Enum[TestEnum]
		expected bool
	}{
		{"Same enum", TestFirst, TestFirst, true},
		{"Different enums", TestFirst, TestSecond, false},
		{"Same name different value", Enum[TestEnum]{Name: "FIRST", Value: SECOND}, TestFirst, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enum1.Equal(tt.enum2); got != tt.expected {
				t.Errorf("Enum.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestNewEnumSet tests the NewEnumSet function
func TestNewEnumSet(t *testing.T) {
	set := NewEnumSet[TestEnum]()
	if set == nil {
		t.Errorf("NewEnumSet() returned nil")
	}
	if len(set.values) != 0 {
		t.Errorf("NewEnumSet() values length = %v, want 0", len(set.values))
	}
}

// TestEnumSetAdd tests the Add method of EnumSet
func TestEnumSetAdd(t *testing.T) {
	set := NewEnumSet[TestEnum]()
	set.Add(TestFirst)

	if len(set.values) != 1 {
		t.Errorf("EnumSet.Add() values length = %v, want 1", len(set.values))
	}

	if !set.values[0].Equal(TestFirst) {
		t.Errorf("EnumSet.Add() first value = %v, want %v", set.values[0], TestFirst)
	}

	set.Add(TestSecond)
	if len(set.values) != 2 {
		t.Errorf("EnumSet.Add() values length = %v, want 2", len(set.values))
	}
}

// TestEnumSetValues tests the Values method of EnumSet
func TestEnumSetValues(t *testing.T) {
	set := NewEnumSet[TestEnum]()
	set.Add(TestFirst)
	set.Add(TestSecond)

	values := set.Values()
	if len(values) != 2 {
		t.Errorf("EnumSet.Values() length = %v, want 2", len(values))
	}

	if !values[0].Equal(TestFirst) || !values[1].Equal(TestSecond) {
		t.Errorf("EnumSet.Values() returned incorrect values")
	}
}

// TestEnumSetFindByName tests the FindByName method of EnumSet
func TestEnumSetFindByName(t *testing.T) {
	set := NewEnumSet[TestEnum]()
	set.Add(TestFirst)
	set.Add(TestSecond)

	tests := []struct {
		name          string
		searchName    string
		expectedFound bool
		expectedEnum  Enum[TestEnum]
	}{
		{"Find existing enum", "FIRST", true, TestFirst},
		{"Find another existing enum", "SECOND", true, TestSecond},
		{"Find non-existent enum", "THIRD", false, Enum[TestEnum]{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := set.FindByName(tt.searchName)

			if tt.expectedFound {
				if enum, found := result.GetIfPresent(); !found {
					t.Errorf("EnumSet.FindByName(%v) expected found, got not found", tt.searchName)
				} else if !enum.Equal(tt.expectedEnum) {
					t.Errorf("EnumSet.FindByName(%v) = %v, want %v", tt.searchName, enum, tt.expectedEnum)
				}
			} else {
				if _, found := result.GetIfPresent(); found {
					t.Errorf("EnumSet.FindByName(%v) expected not found, got found", tt.searchName)
				}

				if result.IsPresent() {
					t.Errorf("EnumSet.FindByName(%v) expected empty Optional, got non-empty", tt.searchName)
				}
			}
		})
	}
}

// TestEnumSetSortByOrder tests the SortByOrder method of EnumSet
func TestEnumSetSortByOrder(t *testing.T) {
	// Create enums with different order values
	orderEnums := []Enum[int]{
		{Name: "THIRD", Value: 3},
		{Name: "FIRST", Value: 1},
		{Name: "SECOND", Value: 2},
	}

	set := FromValues(orderEnums)
	sortedSet := set.SortByOrder(func(i int) int { return i })

	// Check that sortedSet is the same instance as set (for method chaining)
	if sortedSet != set {
		t.Errorf("EnumSet.SortByOrder() didn't return the same instance for method chaining")
	}

	// Check that values are sorted
	values := sortedSet.Values()
	if len(values) != 3 {
		t.Errorf("SortByOrder() sorted values length = %v, want 3", len(values))
	}

	// Check the order of the sorted values
	expectedNames := []string{"FIRST", "SECOND", "THIRD"}
	for i, expectedName := range expectedNames {
		if values[i].Name != expectedName {
			t.Errorf("SortByOrder()[%d].Name = %v, want %v", i, values[i].Name, expectedName)
		}
	}
}

// TestFromValues tests the FromValues function
func TestFromValues(t *testing.T) {
	values := []Enum[TestEnum]{TestFirst, TestSecond}
	set := FromValues(values)

	if set == nil {
		t.Errorf("FromValues() returned nil")
	}

	if len(set.values) != 2 {
		t.Errorf("FromValues() values length = %v, want 2", len(set.values))
	}

	if !set.values[0].Equal(TestFirst) || !set.values[1].Equal(TestSecond) {
		t.Errorf("FromValues() returned incorrect values")
	}
}
