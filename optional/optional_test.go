package optional

import (
	"errors"
	"fmt"
	"testing"
)

// CustomError for testing error interfaces
type CustomError struct {
	Code    int
	Message string
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func TestOptionalOf(t *testing.T) {
	opt := Of("test")
	if !opt.IsPresent() {
		t.Error("Optional should be present")
	}

	val, ok := opt.GetIfPresent()
	if !ok || val != "test" {
		t.Errorf("Expected value 'test', got %v, present: %v", val, ok)
	}
}

func TestOptionalEmpty(t *testing.T) {
	opt := Empty[string]()
	if opt.IsPresent() {
		t.Error("Optional should be empty")
	}

	val, ok := opt.GetIfPresent()
	if ok || val != "" {
		t.Errorf("Expected empty value, got %v, present: %v", val, ok)
	}
}

func TestOptionalOrElse(t *testing.T) {
	opt1 := Of("test")
	opt2 := Empty[string]()

	if opt1.OrElse("default") != "test" {
		t.Error("OrElse should return the value when present")
	}

	if opt2.OrElse("default") != "default" {
		t.Error("OrElse should return the default when empty")
	}
}

func TestOptionalOfNullable(t *testing.T) {
	isZero := func(s string) bool { return s == "" }

	opt1 := OfNullable("test", isZero)
	if !opt1.IsPresent() {
		t.Error("OfNullable should be present for non-zero value")
	}

	opt2 := OfNullable("", isZero)
	if opt2.IsPresent() {
		t.Error("OfNullable should be empty for zero value")
	}

	// Test with various types
	isZeroInt := func(i int) bool { return i == 0 }
	intOpt1 := OfNullable(42, isZeroInt)
	if !intOpt1.IsPresent() {
		t.Error("OfNullable should be present for non-zero int")
	}

	intOpt2 := OfNullable(0, isZeroInt)
	if intOpt2.IsPresent() {
		t.Error("OfNullable should be empty for zero int")
	}

	// Test with struct type
	type Person struct {
		Name string
		Age  int
	}
	isZeroPerson := func(p Person) bool { return p.Name == "" && p.Age == 0 }

	personOpt1 := OfNullable(Person{Name: "John", Age: 30}, isZeroPerson)
	if !personOpt1.IsPresent() {
		t.Error("OfNullable should be present for non-zero struct")
	}

	personOpt2 := OfNullable(Person{}, isZeroPerson)
	if personOpt2.IsPresent() {
		t.Error("OfNullable should be empty for zero struct")
	}
}

func TestOptionalGet(t *testing.T) {
	opt1 := Of("test")
	value, err := opt1.Get()
	if err != nil || value != "test" {
		t.Errorf("Get should return value without error, got value=%v, err=%v", value, err)
	}

	opt2 := Empty[string]()
	value, err = opt2.Get()
	if err != ErrNoValuePresent {
		t.Errorf("Get should return ErrNoValuePresent, got %v", err)
	}
	if value != "" {
		t.Errorf("Get should return empty value, got %v", value)
	}

	// Test with non-string type
	intOpt := Of(42)
	intVal, err := intOpt.Get()
	if err != nil || intVal != 42 {
		t.Errorf("Get should return int value without error, got value=%v, err=%v", intVal, err)
	}
}

func TestOptionalOrElseGet(t *testing.T) {
	supplierCalled := false
	supplier := func() string {
		supplierCalled = true
		return "supplied"
	}

	opt1 := Of("test")
	value := opt1.OrElseGet(supplier)
	if value != "test" {
		t.Errorf("OrElseGet should return the value when present, got %v", value)
	}
	if supplierCalled {
		t.Error("Supplier should not be called when value is present")
	}

	supplierCalled = false
	opt2 := Empty[string]()
	value = opt2.OrElseGet(supplier)
	if value != "supplied" {
		t.Errorf("OrElseGet should return supplier result when empty, got %v", value)
	}
	if !supplierCalled {
		t.Error("Supplier should be called when value is not present")
	}

	// Test with complex supplier logic
	complexSupplier := func() int {
		// Simulate some complex calculation
		result := 0
		for i := 1; i <= 5; i++ {
			result += i
		}
		return result // 15
	}

	intOpt := Empty[int]()
	intVal := intOpt.OrElseGet(complexSupplier)
	if intVal != 15 {
		t.Errorf("OrElseGet should return complex supplier result when empty, got %v", intVal)
	}
}

func TestOptionalOrElseThrow(t *testing.T) {
	customErr := errors.New("custom error")

	opt1 := Of("test")
	value, err := opt1.OrElseThrow(customErr)
	if err != nil || value != "test" {
		t.Errorf("OrElseThrow should return value without error, got value=%v, err=%v", value, err)
	}

	opt2 := Empty[string]()
	value, err = opt2.OrElseThrow(customErr)
	if err != customErr {
		t.Errorf("OrElseThrow should return custom error, got %v", err)
	}
	if value != "" {
		t.Errorf("OrElseThrow should return empty value, got %v", value)
	}

	// Test with custom error type
	customStructErr := &CustomError{Code: 404, Message: "Not found"}
	opt3 := Empty[int]()
	_, err = opt3.OrElseThrow(customStructErr)
	if err != customStructErr {
		t.Errorf("OrElseThrow should return custom struct error, got %v", err)
	}
}

func TestOptionalIfPresent(t *testing.T) {
	actionCalled := false
	action := func(s string) {
		if s != "test" {
			t.Errorf("Action received wrong value: %v", s)
		}
		actionCalled = true
	}

	opt1 := Of("test")
	opt1.IfPresent(action)
	if !actionCalled {
		t.Error("Action should be called when value is present")
	}

	actionCalled = false
	opt2 := Empty[string]()
	opt2.IfPresent(action)
	if actionCalled {
		t.Error("Action should not be called when value is not present")
	}

	// Test with multiple actions
	var capturedValue int
	intOpt := Of(42)
	intOpt.IfPresent(func(v int) {
		capturedValue = v
	})
	if capturedValue != 42 {
		t.Errorf("IfPresent should pass correct value to action, got %v", capturedValue)
	}

	// Test with more complex action
	var sum int
	Of(10).IfPresent(func(v int) {
		sum = v + 5
	})
	if sum != 15 {
		t.Errorf("IfPresent action should execute correctly, got sum=%v", sum)
	}
}

// Test multiple operations chained together
func TestOptionalChaining(t *testing.T) {
	// Create some test data
	type Result struct {
		Value string
	}

	process := func(input string) Optional[Result] {
		if input == "" {
			return Empty[Result]()
		}
		return Of(Result{Value: "Processed: " + input})
	}

	// Test chaining with present value
	result, found := process("test").GetIfPresent()
	if !found || result.Value != "Processed: test" {
		t.Errorf("Chained operation failed for present value, got %v, found=%v", result, found)
	}

	// Test chaining with empty value
	emptyResult := process("").OrElse(Result{Value: "Default"})
	if emptyResult.Value != "Default" {
		t.Errorf("Chained operation failed for empty value, got %v", emptyResult)
	}

	// Test with error handling
	_, err := process("").OrElseThrow(errors.New("processing failed"))
	if err == nil || err.Error() != "processing failed" {
		t.Errorf("Expected 'processing failed' error, got %v", err)
	}
}

// Test comparing with nil and zero values
func TestOptionalWithNilAndZeroValues(t *testing.T) {
	// Test with pointer types
	var ptr *string
	isNilPtr := func(p *string) bool { return p == nil }

	ptrOpt1 := OfNullable(ptr, isNilPtr)
	if ptrOpt1.IsPresent() {
		t.Error("Optional with nil pointer should not be present")
	}

	str := "hello"
	ptr = &str
	ptrOpt2 := OfNullable(ptr, isNilPtr)
	if !ptrOpt2.IsPresent() {
		t.Error("Optional with non-nil pointer should be present")
	}

	// Test with boolean type (which is comparable)
	isZeroBool := func(b bool) bool { return !b }

	boolOpt1 := OfNullable(false, isZeroBool)
	if boolOpt1.IsPresent() {
		t.Error("Optional with false should not be present")
	}

	boolOpt2 := OfNullable(true, isZeroBool)
	if !boolOpt2.IsPresent() {
		t.Error("Optional with true should be present")
	}

	// Test with fixed-size arrays (which are comparable)
	isEmptyArray := func(a [3]int) bool { return a == [3]int{0, 0, 0} }

	arrOpt1 := OfNullable([3]int{0, 0, 0}, isEmptyArray)
	if arrOpt1.IsPresent() {
		t.Error("Optional with zero array should not be present")
	}

	arrOpt2 := OfNullable([3]int{1, 2, 3}, isEmptyArray)
	if !arrOpt2.IsPresent() {
		t.Error("Optional with non-zero array should be present")
	}
}

func ExampleOptional() {
	// Criar um Optional com valor
	opt1 := Of("Hello, World!")
	fmt.Println("opt1 está presente:", opt1.IsPresent())

	// Obter valor se presente
	if val, ok := opt1.GetIfPresent(); ok {
		fmt.Println("Valor de opt1:", val)
	}

	// Criar um Optional vazio
	opt2 := Empty[string]()
	fmt.Println("opt2 está presente:", opt2.IsPresent())

	// Usar OrElse para fornecer um valor padrão
	fmt.Println("Valor de opt2 ou padrão:", opt2.OrElse("valor padrão"))

	// Output:
	// opt1 está presente: true
	// Valor de opt1: Hello, World!
	// opt2 está presente: false
	// Valor de opt2 ou padrão: valor padrão
}
