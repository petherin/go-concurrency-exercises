package main

import (
	"fmt"
	"sync"
)

// run the program and check that variable i
// was pinned for access from goroutine even after
// enclosing function returns.

func main() {
	var wg sync.WaitGroup

	incr := func(wg *sync.WaitGroup) {
		// Usually, when func returns, local variable i's
		// value would be lost as it goes out of scope.
		// But if the goroutine has a reference to it, it is
		// copied to the goroutine and persists even
		// after the enclosing function ends.
		// This is an enclosure.
		var i int
		wg.Add(1)

		// i is copied to goroutine and persists
		// inside it. This is done by copying the value of i
		// from the stack (used by function) to the heap
		// where the goroutine is running.
		go func() {
			defer wg.Done()
			i++
			fmt.Printf("value of i: %v\n", i)
		}()
		fmt.Println("return from function")

		// When we return, i goes out of scope for
		// for the function, but its copy lives on in the
		// goroutine.
		return
	}

	incr(&wg)
	wg.Wait()
	fmt.Println("done")
}
