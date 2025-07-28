package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tiagods/enum-go/optional"
)

// User represents a simple user model
type User struct {
	ID       int
	Name     string
	Email    string
	Age      int
	Role     string
	IsActive bool
}

// UserRepository simulates a database of users
type UserRepository struct {
	users map[int]User
}

// NewUserRepository creates a new user repository with sample data
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: map[int]User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com", Age: 28, Role: "Admin", IsActive: true},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com", Age: 34, Role: "User", IsActive: true},
			3: {ID: 3, Name: "Charlie", Email: "charlie@example.com", Age: 45, Role: "User", IsActive: false},
		},
	}
}

// FindUserByID returns an Optional with the user if found
func (r *UserRepository) FindUserByID(id int) optional.Optional[User] {
	if user, ok := r.users[id]; ok {
		return optional.Of(user)
	}
	return optional.Empty[User]()
}

// FindActiveUserByEmail simulates finding a user by email
func (r *UserRepository) FindActiveUserByEmail(email string) optional.Optional[User] {
	for _, user := range r.users {
		if user.Email == email && user.IsActive {
			return optional.Of(user)
		}
	}
	return optional.Empty[User]()
}

// ParseInt safely parses a string to int, returning an Optional
func ParseInt(s string) optional.Optional[int] {
	n, err := strconv.Atoi(s)
	if err != nil {
		return optional.Empty[int]()
	}
	return optional.Of(n)
}

// Example 1: Basic Optional Usage
func example1() {
	fmt.Println("\n=== Example 1: Basic Optional Usage ===")

	// Create a present optional
	presentOpt := optional.Of("Hello, Optional!")
	fmt.Println("Present optional:", presentOpt.IsPresent())

	// Create an empty optional
	emptyOpt := optional.Empty[string]()
	fmt.Println("Empty optional:", emptyOpt.IsPresent())

	// Get values safely
	if value, found := presentOpt.GetIfPresent(); found {
		fmt.Println("Value:", value)
	}

	// Use OrElse for default values
	value := emptyOpt.OrElse("Default Value")
	fmt.Println("OrElse result:", value)
}

// Example 2: Working with repositories and error handling
func example2() {
	fmt.Println("\n=== Example 2: Repository Pattern with Optional ===")

	repo := NewUserRepository()

	// Find user by ID
	userOpt := repo.FindUserByID(1)

	// Method 1: Using GetIfPresent
	if user, found := userOpt.GetIfPresent(); found {
		fmt.Printf("Found user: %s (ID: %d)\n", user.Name, user.ID)
	} else {
		fmt.Println("User not found")
	}

	// Method 2: Using IfPresent with a closure
	repo.FindUserByID(2).IfPresent(func(user User) {
		fmt.Printf("Using IfPresent: User %s has role %s\n", user.Name, user.Role)
	})

	// Method 3: Using Get with error handling
	notFoundOpt := repo.FindUserByID(999)
	user, err := notFoundOpt.Get()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("User:", user.Name) // This won't execute
	}
}

// Example 3: Chaining operations
func example3() {
	fmt.Println("\n=== Example 3: Chaining Operations ===")

	repo := NewUserRepository()

	// Find a user with GetIfPresent
	activeUser := repo.FindActiveUserByEmail("charlie@example.com")
	if user, found := activeUser.GetIfPresent(); found {
		fmt.Printf("Found active user: %s\n", user.Name)
	} else {
		fmt.Println("No active user found with that email")
	}

	// Using OrElse to provide default
	inactiveUser := repo.FindActiveUserByEmail("charlie@example.com")
	defaultUser := User{Name: "Default User"}
	user := inactiveUser.OrElse(defaultUser)
	fmt.Printf("User (with default): %s\n", user.Name)

	// Using OrElseGet with a function
	activeAdmin := repo.FindActiveUserByEmail("alice@example.com")
	result := activeAdmin.OrElseGet(func() User {
		return User{Name: "Fallback Admin", Role: "Admin"}
	})
	fmt.Printf("Admin result: %s (%s)\n", result.Name, result.Role)
}

// Example 4: Handling potential parsing errors
func example4() {
	fmt.Println("\n=== Example 4: Safe Parsing with Optional ===")

	inputs := []string{"42", "123", "not-a-number", "987"}

	for _, input := range inputs {
		// Parse string to int, returning an Optional
		numOpt := ParseInt(input)

		// Method 1: Check presence before using
		if numOpt.IsPresent() {
			num, _ := numOpt.Get()
			fmt.Printf("Successfully parsed '%s' to %d\n", input, num)
		} else {
			fmt.Printf("Failed to parse '%s'\n", input)
		}

		// Method 2: Using OrElseThrow with custom error
		num, err := numOpt.OrElseThrow(errors.New("invalid number format"))
		if err != nil {
			fmt.Printf("Error for '%s': %v\n", input, err)
		} else {
			fmt.Printf("Parsed value: %d\n", num)
		}
	}
}

// Example 5: Using OfNullable for zero value handling
func example5() {
	fmt.Println("\n=== Example 5: OfNullable for Zero Values ===")

	// Define what "zero" means for different types
	isZeroInt := func(i int) bool { return i == 0 }
	isZeroString := func(s string) bool { return s == "" }

	// Create optionals with OfNullable
	opt1 := optional.OfNullable(42, isZeroInt)
	opt2 := optional.OfNullable(0, isZeroInt)
	opt3 := optional.OfNullable("Hello", isZeroString)
	opt4 := optional.OfNullable("", isZeroString)

	fmt.Printf("OfNullable(42): present=%v\n", opt1.IsPresent())
	fmt.Printf("OfNullable(0): present=%v\n", opt2.IsPresent())
	fmt.Printf("OfNullable(\"Hello\"): present=%v\n", opt3.IsPresent())
	fmt.Printf("OfNullable(\"\"): present=%v\n", opt4.IsPresent())
}

func main() {
	fmt.Println("Optional Package Usage Examples")
	fmt.Println("===============================")

	example1()
	example2()
	example3()
	example4()
	example5()
}
