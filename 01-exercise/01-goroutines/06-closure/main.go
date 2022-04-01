package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// What is the output?
	// Before fix, this prints out 4 three times.
	// This is because the loop increments i to 4
	// before any of the goroutines get a chance to run.
	// By the time they are running, i is already at 4.
	// Goroutines operate on the current value of their
	// variables at the time of their execution.

	// Broken code
	//for i := 1; i <= 3; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		fmt.Println(i)
	//	}()
	//}

	// Fixed code
	// If we want the goroutine to operate on a specific
	// value, we have to pass it to the goroutine.
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}
	wg.Wait()
}
