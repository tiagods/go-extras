package main

import (
	"fmt"
	"math"

	"github.com/tiagods/go-extras/enum"
)

// OperationValue represents the properties and behavior of a mathematical operation
type OperationValue struct {
	Symbol      string
	Order       int
	Apply       func(a, b float64) float64
	Description string
}

var (
	// SUM represents the addition operation
	SUM = enum.Enum[OperationValue]{
		Name: "SUM",
		Value: OperationValue{
			Symbol:      "+",
			Order:       2,
			Apply:       func(a, b float64) float64 { return a + b },
			Description: "Adds two values",
		},
	}

	// SUBTRACT represents the subtraction operation
	SUBTRACT = enum.Enum[OperationValue]{
		Name: "SUBTRACT",
		Value: OperationValue{
			Symbol:      "-",
			Order:       3,
			Apply:       func(a, b float64) float64 { return a - b },
			Description: "Subtracts two values",
		},
	}

	// MULTIPLY represents the multiplication operation
	MULTIPLY = enum.Enum[OperationValue]{
		Name: "MULTIPLY",
		Value: OperationValue{
			Symbol:      "*",
			Order:       4,
			Apply:       func(a, b float64) float64 { return a * b },
			Description: "Multiplies two values",
		},
	}

	// DIVIDE represents the division operation
	DIVIDE = enum.Enum[OperationValue]{
		Name: "DIVIDE",
		Value: OperationValue{
			Symbol:      "/",
			Order:       5,
			Apply:       func(a, b float64) float64 { return a / b },
			Description: "Divides two values",
		},
	}

	// MODULUS represents the modulo operation
	MODULUS = enum.Enum[OperationValue]{
		Name: "MODULUS",
		Value: OperationValue{
			Symbol:      "%",
			Order:       1,
			Apply:       func(a, b float64) float64 { return math.Mod(a, b) },
			Description: "Returns the remainder of division between two values",
		},
	}
)

// OperationSet is a collection of all operations sorted by order
var OperationSet = enum.FromValues([]enum.Enum[OperationValue]{SUM, SUBTRACT, MULTIPLY, DIVIDE, MODULUS}).
	SortByOrder(func(op OperationValue) int { return op.Order })

func main() {
	fmt.Println("=== Enum and EnumSet Examples ===")

	// Example 1: Creating a new EnumSet
	fmt.Println("\n1. Creating an EnumSet manually:")
	customSet := enum.NewEnumSet[OperationValue]()
	customSet.Add(SUM)
	customSet.Add(MULTIPLY)
	fmt.Printf("Custom set contains %d operations\n", len(customSet.Values()))

	// Example 2: Using String() method
	fmt.Println("\n2. Using String() method on enum:")
	fmt.Printf("Operation name: %s\n", SUM.String())

	// Example 3: Using Equal method
	fmt.Println("\n3. Using Equal method to compare enums:")
	fmt.Printf("SUM equals SUBTRACT: %v\n", SUM.Equal(SUBTRACT))
	fmt.Printf("SUM equals SUM: %v\n", SUM.Equal(SUM))

	// Example 4: Finding by name
	fmt.Println("\n4. Finding an operation by name:")
	operation := OperationSet.FindByName("SUM")
	if op, found := operation.GetIfPresent(); found {
		fmt.Printf("Found operation: %s, result of 1+2=%v\n", op.Name, op.Value.Apply(1, 2))
	} else {
		fmt.Println("Operation not found")
	}

	// Example 5: Using SortByOrder
	fmt.Println("\n5. Demonstrating sorted order:")
	fmt.Println("Operations in sorted order:")
	for _, op := range OperationSet.Values() {
		fmt.Printf("- %s (order: %d, symbol: %s)\n", op.Name, op.Value.Order, op.Value.Symbol)
	}

	// Example 6: Using FromValues
	fmt.Println("\n6. Creating a set from values:")
	newSet := enum.FromValues([]enum.Enum[OperationValue]{DIVIDE, MODULUS})
	fmt.Printf("New set has %d operations: %s and %s\n",
		len(newSet.Values()), newSet.Values()[0].Name, newSet.Values()[1].Name)

	// Example 7: Functional example
	fmt.Println("\n7. Practical usage with operation application:")
	doCalculation := func(opName string, a, b float64) {
		opOptional := OperationSet.FindByName(opName)
		if op, found := opOptional.GetIfPresent(); found {
			result := op.Value.Apply(a, b)
			fmt.Printf("%v %s %v = %v\n", a, op.Value.Symbol, b, result)
		} else {
			fmt.Printf("Operation '%s' not found\n", opName)
		}
	}

	doCalculation("SUM", 10, 5)
	doCalculation("SUBTRACT", 10, 5)
	doCalculation("MULTIPLY", 10, 5)
	doCalculation("DIVIDE", 10, 5)
	doCalculation("MODULUS", 10, 3)
	doCalculation("POWER", 10, 2) // This should fail (operation not found)
}
