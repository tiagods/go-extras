package stream

import (
	"runtime"
	"sync"
	"testing"
)

func TestParallelStream(t *testing.T) {
	// Test basic functionality
	s := NewStream([]int{1, 2, 3, 4, 5})

	// Double each number in parallel
	result := s.ParallelStream(func(i int) interface{} {
		return i * 2
	}, 2) // Use 2 goroutines

	// Check if all elements are processed correctly
	expectedValues := []int{2, 4, 6, 8, 10}
	for _, expected := range expectedValues {
		found := false
		for _, actual := range result.elements {
			if actual.(int) == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find %d in results, but didn't", expected)
		}
	}

	// Check that all results are present (order may vary)
	if len(result.elements) != len(s.elements) {
		t.Errorf("Expected %d results, got %d", len(s.elements), len(result.elements))
	}

	// Test with default number of goroutines (should use GOMAXPROCS)
	result2 := s.ParallelStream(func(i int) interface{} {
		return i * 3
	}, 0)

	// Check results again
	expectedValues2 := []int{3, 6, 9, 12, 15}
	for _, expected := range expectedValues2 {
		found := false
		for _, actual := range result2.elements {
			if actual.(int) == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find %d in results, but didn't", expected)
		}
	}
}

func TestParallelStreamWithEmptyStream(t *testing.T) {
	// Test with empty stream
	emptyStream := NewStream([]int{})
	result := emptyStream.ParallelStream(func(i int) interface{} {
		return i * 2
	}, 2)

	if len(result.elements) != 0 {
		t.Errorf("Expected empty result for empty stream, got %v", result.elements)
	}
}

func TestParallelStreamWithLargeData(t *testing.T) {
	// Skip this test if running in CI or short mode
	if testing.Short() {
		t.Skip("Skipping large data test in short mode")
	}

	// Create a large dataset
	data := make([]int, 1000)
	for i := 0; i < len(data); i++ {
		data[i] = i
	}

	s := NewStream(data)

	// Use a time-consuming operation
	result := s.ParallelStream(func(i int) interface{} {
		sum := 0
		for j := 0; j < 10000; j++ { // Artificial work
			sum += j % (i + 1)
		}
		return i * 2
	}, runtime.GOMAXPROCS(0))

	// Verify result count
	if len(result.elements) != len(data) {
		t.Errorf("Expected %d results, got %d", len(data), len(result.elements))
	}

	// Verify all elements were processed
	resultMap := make(map[int]bool)
	for _, v := range result.elements {
		resultMap[v.(int)] = true
	}

	for _, expected := range data {
		if !resultMap[expected*2] {
			t.Errorf("Missing result for input %d", expected)
			break // Just report one failure to avoid flooding output
		}
	}
}

func TestParallelStreamConcurrencyLimit(t *testing.T) {
	// This test verifies the goroutine limit is enforced
	s := NewStream([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// We'll use a shared counter and mutex to track maximum concurrency
	var (
		mu             sync.Mutex
		activeCount    int
		maxActiveCount int
	)

	// Process in parallel with artificial delays
	s.ParallelStream(func(i int) interface{} {
		// Increment active count
		mu.Lock()
		activeCount++
		if activeCount > maxActiveCount {
			maxActiveCount = activeCount
		}
		mu.Unlock()

		// Simulate work
		runtime.Gosched() // Force scheduler to potentially run other goroutines

		// Decrement active count
		mu.Lock()
		activeCount--
		mu.Unlock()

		return i
	}, 3) // Max 3 goroutines

	// Check that we never exceeded our limit
	if maxActiveCount > 3 {
		t.Errorf("Concurrency limit exceeded: wanted max 3 concurrent goroutines, got %d", maxActiveCount)
	}
}
