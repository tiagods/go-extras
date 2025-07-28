package stream

import (
	"runtime"
	"sync"
)

// ParallelStream applies the `mapper` function to the Stream elements in parallel
// The number of goroutines can be specified by the user, or a default value will be used
func (s *Stream[T]) ParallelStream(mapper func(T) interface{}, maxGoroutines int) *Stream[interface{}] {
	// If the user didn't specify, use the number of available CPUs
	if maxGoroutines <= 0 {
		maxGoroutines = runtime.GOMAXPROCS(0)
	}

	var wg sync.WaitGroup
	resultChan := make(chan interface{}, len(s.elements))

	// Process elements in parallel with a limit on simultaneous goroutines
	sem := make(chan struct{}, maxGoroutines) // Semaphore to limit the number of simultaneous goroutines

	for _, e := range s.elements {
		wg.Add(1)
		sem <- struct{}{} // Acquire a "token" from the semaphore
		go func(el T) {
			defer wg.Done()
			defer func() { <-sem }() // Release the "token" after execution

			resultChan <- mapper(el) // Apply the mapper function and send to the channel
		}(e)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(resultChan)

	// Collect results from the channel and return a new Stream
	var result []interface{}
	for res := range resultChan {
		result = append(result, res)
	}

	return NewStream(result)
}
