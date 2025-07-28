package enum

import (
	"encoding/json"
	"testing"
)

type ColorEnum struct {
	Hex string
	RGB [3]int
}

// Test enums with ColorEnum type
var (
	RED = Enum[ColorEnum]{
		Name: "RED",
		Value: ColorEnum{
			Hex: "#FF0000",
			RGB: [3]int{255, 0, 0},
		},
	}

	GREEN = Enum[ColorEnum]{
		Name: "GREEN",
		Value: ColorEnum{
			Hex: "#00FF00",
			RGB: [3]int{0, 255, 0},
		},
	}

	BLUE = Enum[ColorEnum]{
		Name: "BLUE",
		Value: ColorEnum{
			Hex: "#0000FF",
			RGB: [3]int{0, 0, 255},
		},
	}
)

// TestEnumCreation tests the creation and basic properties of Enum
func TestEnumCreation(t *testing.T) {
	tests := []struct {
		name      string
		enum      Enum[ColorEnum]
		expectHex string
		expectRGB [3]int
	}{
		{
			name:      "RED enum creation",
			enum:      RED,
			expectHex: "#FF0000",
			expectRGB: [3]int{255, 0, 0},
		},
		{
			name:      "GREEN enum creation",
			enum:      GREEN,
			expectHex: "#00FF00",
			expectRGB: [3]int{0, 255, 0},
		},
		{
			name:      "BLUE enum creation",
			enum:      BLUE,
			expectHex: "#0000FF",
			expectRGB: [3]int{0, 0, 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check name
			if tt.enum.Name != tt.enum.String() {
				t.Errorf("String() = %v, expected it to equal Name = %v", tt.enum.String(), tt.enum.Name)
			}

			// Check value properties
			if tt.enum.Value.Hex != tt.expectHex {
				t.Errorf("Value.Hex = %v, expected %v", tt.enum.Value.Hex, tt.expectHex)
			}

			for i, val := range tt.enum.Value.RGB {
				if val != tt.expectRGB[i] {
					t.Errorf("Value.RGB[%d] = %v, expected %v", i, val, tt.expectRGB[i])
				}
			}
		})
	}
}

// TestColorEnumString tests the String method of Enum with ColorEnum
func TestColorEnumString(t *testing.T) {
	tests := []struct {
		name     string
		enum     Enum[ColorEnum]
		expected string
	}{
		{"RED enum string representation", RED, "RED"},
		{"GREEN enum string representation", GREEN, "GREEN"},
		{"BLUE enum string representation", BLUE, "BLUE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enum.String(); got != tt.expected {
				t.Errorf("Enum.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestColorEnumEqual tests the Equal method of Enum with ColorEnum
func TestColorEnumEqual(t *testing.T) {
	// Create a copy of RED with same name but slightly different value
	redCopy := Enum[ColorEnum]{
		Name: "RED",
		Value: ColorEnum{
			Hex: "#FF0000",
			RGB: [3]int{254, 0, 0}, // Slightly different RGB
		},
	}

	tests := []struct {
		name     string
		enum1    Enum[ColorEnum]
		enum2    Enum[ColorEnum]
		expected bool
	}{
		{"Same enum instance", RED, RED, true},
		{"Different enum instances with same name", RED, redCopy, true},
		{"Different enum names", RED, GREEN, false},
		{"Different enum types", RED, GREEN, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enum1.Equal(tt.enum2); got != tt.expected {
				t.Errorf("Enum.Equal() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestEnumMarshalJSON tests the JSON serialization of enum values
func TestEnumMarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		enum         Enum[ColorEnum]
		expectedJSON string
	}{
		{"RED enum serialization", RED, `"RED"`},
		{"GREEN enum serialization", GREEN, `"GREEN"`},
		{"BLUE enum serialization", BLUE, `"BLUE"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test using the MarshalJSON method directly
			jsonBytes, err := tt.enum.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}

			if string(jsonBytes) != tt.expectedJSON {
				t.Errorf("MarshalJSON() = %v, want %v", string(jsonBytes), tt.expectedJSON)
			}

			// Test using the standard json.Marshal which should call the MarshalJSON method
			stdJsonBytes, err := json.Marshal(tt.enum)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
				return
			}

			if string(stdJsonBytes) != tt.expectedJSON {
				t.Errorf("json.Marshal() = %v, want %v", string(stdJsonBytes), tt.expectedJSON)
			}
		})
	}

	// Test serializing a slice of enums
	t.Run("Slice of enums", func(t *testing.T) {
		enums := []Enum[ColorEnum]{RED, GREEN, BLUE}
		expectedJSON := `["RED","GREEN","BLUE"]`

		jsonBytes, err := json.Marshal(enums)
		if err != nil {
			t.Errorf("json.Marshal(slice) error = %v", err)
			return
		}

		if string(jsonBytes) != expectedJSON {
			t.Errorf("json.Marshal(slice) = %v, want %v", string(jsonBytes), expectedJSON)
		}
	})

	// Test serializing an EnumSet
	t.Run("EnumSet serialization", func(t *testing.T) {
		type EnumSetJSON struct {
			Colors []Enum[ColorEnum] `json:"colors"`
		}

		set := EnumSetJSON{
			Colors: []Enum[ColorEnum]{RED, GREEN},
		}

		expectedJSON := `{"colors":["RED","GREEN"]}`

		jsonBytes, err := json.Marshal(set)
		if err != nil {
			t.Errorf("json.Marshal(EnumSet) error = %v", err)
			return
		}

		if string(jsonBytes) != expectedJSON {
			t.Errorf("json.Marshal(EnumSet) = %v, want %v", string(jsonBytes), expectedJSON)
		}
	})
}

// TestEnumWithComplexTypes tests Enum with complex types including functions
func TestEnumWithComplexTypes(t *testing.T) {
	type OperationValue struct {
		Symbol    string
		Operation func(a, b int) int
	}

	add := Enum[OperationValue]{
		Name: "ADD",
		Value: OperationValue{
			Symbol:    "+",
			Operation: func(a, b int) int { return a + b },
		},
	}

	subtract := Enum[OperationValue]{
		Name: "SUBTRACT",
		Value: OperationValue{
			Symbol:    "-",
			Operation: func(a, b int) int { return a - b },
		},
	}

	// Test the operation functions
	if result := add.Value.Operation(5, 3); result != 8 {
		t.Errorf("ADD operation = %v, expected 8", result)
	}

	if result := subtract.Value.Operation(5, 3); result != 2 {
		t.Errorf("SUBTRACT operation = %v, expected 2", result)
	}

	// Test string representation
	if add.String() != "ADD" {
		t.Errorf("String() = %v, expected ADD", add.String())
	}

	// Test equality
	if !add.Equal(add) {
		t.Errorf("Equal() returned false for the same enum")
	}

	if add.Equal(subtract) {
		t.Errorf("Equal() returned true for different enums")
	}
}
