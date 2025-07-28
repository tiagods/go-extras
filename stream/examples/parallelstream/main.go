package main

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/tiagods/go-extras/stream"
)

// Simulate a computationally intensive task
func heavyComputation(n int) int {
	// Simulate intensive computation
	result := 0
	for i := 0; i < 10000000; i++ {
		result += i % n
	}
	return n * 2
}

func main() {
	fmt.Println("==== Parallel Stream Examples ====")

	// Generate a large dataset for processing
	data := make([]int, 100)
	for i := 0; i < len(data); i++ {
		data[i] = i + 1
	}

	// Example 1: Basic parallel processing
	fmt.Println("\n1. Basic parallel processing:")
	fmt.Printf("Number of available CPUs: %d\n", runtime.NumCPU())

	// Using default parallelism (using all available CPUs)
	start := time.Now()
	resultDefault := stream.NewStream(data).ParallelStream(func(n int) interface{} {
		return heavyComputation(n)
	}, 0) // 0 means use all available CPUs

	fmt.Printf("Processing with default parallelism took: %v\n", time.Since(start))
	fmt.Printf("First few results: %v\n", stream.Limit(resultDefault, 5).ToSlice())

	// Example 2: Limiting concurrency
	fmt.Println("\n2. Limited concurrency (2 goroutines):")
	start = time.Now()
	resultLimited := stream.NewStream(data).ParallelStream(func(n int) interface{} {
		return heavyComputation(n)
	}, 2) // Use only 2 goroutines

	fmt.Printf("Processing with 2 goroutines took: %v\n", time.Since(start))
	fmt.Printf("First few results: %v\n", stream.Limit(resultLimited, 5).ToSlice())

	// Example 3: More complex transformations
	fmt.Println("\n3. Complex parallel transformations:")

	// Create a sequence of prime numbers
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}

	// Parallel computation of prime factors
	result := stream.NewStream(primes).ParallelStream(func(p int) interface{} {
		// For each prime, find all numbers in 1-1000 that are divisible by it
		var divisible []int
		for i := 1; i <= 100; i++ {
			if i%p == 0 {
				divisible = append(divisible, i)
			}
		}
		return map[string]interface{}{
			"prime":          p,
			"divisibleCount": len(divisible),
			"firstFew":       divisible[:int(math.Min(float64(len(divisible)), 5))],
		}
	}, 0)

	fmt.Println("Prime factorization results:")
	result.ForEach(func(r interface{}) {
		data := r.(map[string]interface{})
		fmt.Printf("  Prime %d has %d divisible numbers, first few: %v\n",
			data["prime"], data["divisibleCount"], data["firstFew"])
	})

	// Example 4: Order of results in parallel processing
	fmt.Println("\n4. Note on result ordering:")
	fmt.Println("Results from parallel processing might not maintain original order.")
	fmt.Println("If order matters, you may need to add ordering information to results.")

	// Demonstrate potential ordering issues
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Process with varying sleep times to show order differences
	result = stream.NewStream(numbers).ParallelStream(func(n int) interface{} {
		// Sleep a random amount based on the number
		// (even numbers sleep longer to finish later)
		if n%2 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
		return n
	}, 0)

	fmt.Println("Original order:", numbers)
	fmt.Println("Result order:", result.ToSlice())
}
