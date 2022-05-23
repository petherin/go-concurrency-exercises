// generator() -> square() ->
//														-> merge -> print
//             -> square() ->
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func generator(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		// defer close(out) to ensure closure if goroutine terminated
		defer close(out)

		for _, n := range nums {
			select {
			case out <- n:
				// support termination of goroutine if signal sent on done channel
			case <-done:
				return
			}
		}
	}()

	return out
}

func square(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		// defer close(out) to ensure closure if goroutine terminated
		defer close(out)

		for n := range in {
			select {
			case out <- n * n:
				// support termination of goroutine if signal sent on done channel
			case <-done:
				return
			}
		}
	}()

	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	output := func(c <-chan int) {
		// defer wg.Done() to ensure it's run even if done
		// channel makes us return.
		defer wg.Done()

		for n := range c {
			// By adding a select we can catch signals on
			// the done channel and return if one is seen.
			// Otherwise we will perform the usual functionality
			// this code is doing
			// i.e. merging values on multiple channels to 'out'.
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	// Done channel to signal closures.
	// Use empty struct{} as we don't want to send any data,
	// just a signal.
	done := make(chan struct{})

	// Pass done channel to all stages of the pipeline.
	// It is idiomatic to pass the done channel as the first
	// paramater to the goroutine.
	in := generator(done, 2, 3)

	c1 := square(done, in)
	c2 := square(done, in)

	out := merge(done, c1, c2)

	// TODO: cancel goroutines after receiving one value.

	fmt.Println(<-out)

	// After reading one value from out channel, close the done channel to signal
	// goroutines to terminate.
	close(done)

	// To check if goroutines are being cancelled, allow 10ms
	// for goroutines to terminate, then print NumGoroutine.
	// There should only be 1, the main goroutine.
	time.Sleep(10)
	g := runtime.NumGoroutine()
	fmt.Printf("number of goroutines active = %d\n", g)
}
