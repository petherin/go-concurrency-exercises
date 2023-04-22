package main

import (
	"log"
	"sync"
)

// Based on https://go.dev/blog/pipelines
func main() {
	// basicPipeline()
	fanoutFanin()
}

// basicPipeline sets up the pipeline with 3 stages, a gen stage (source or producer).
// A second stage that does some work (worker).
// And a last stage that prints the results from channel `out` (sink or consumer).
func basicPipeline() {
	// gen returns a receive-only channel containing 2 and 3.
	c := gen(2, 3)

	// sq takes the channel and returns another receive-only channel with the results.
	out := sq(c)

	// Consume the output.
	log.Println(<-out) // 4
	log.Println(<-out) // 9

	// Set up another pipeline that chains stages together.
	// This can be done because gen returns <-chan int and
	// sq takes and returns the same type of <-chan int.
	for n := range sq(sq(gen(2, 3))) {
		log.Println(n) // 16 then 81
	}
}

// gen is a typical first stage in a pipeline.
// It takes in a variadic int. It creates an out channel.
// In a goroutine it ranges over the incoming numbers and send them to the out channel.
// After the range loop, it closes the out channel to signal there are no more numbers
// in the channel.
// While the goroutine runs, we return the channel.
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// sq is a typical second stage in a pipeline.
// It takes in a channel (prepared by `gen`). It creates an out channel
// to contain results.
// In a goroutine it ranges over the in channel and squares each number. It
// sends the result of each square to the out channel.
// When we have finished ranging over the in channel, we close the out channel
// to say there are no more results.
// While the goroutine runs, we return the channel.
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// fanoutFanin runs the gen stage once, as in the basic pipeline.
// But then it calls sq twice, giving us 2 channels. There is a goroutine running inside each sq.
// Send the 2 channels to merge(), and range over the return value of merge(). This is also
// a channel. Log out the values of the channel until the channel is closed (range will stop
// looping when this happens).
func fanoutFanin() {
	// gen some numbers into a channel
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)

	// Consume the merged output from c1 and c2.
	for n := range merge(c1, c2) {
		log.Println(n) // 4 then 9, or 9 then 4
	}
}

// merge takes in n channels. It merges them into a single out channel.
// It defines a function called `output` that takes in a channel, ranges over its values,
// sending each one to the single out channel. Then it calls wg.Done().
// Once the `output` channel is defined, add the number of incoming channels to the waitgroup.
// Range over the incoming channels and call the output function in a goroutine.
// While the output functions are running in goroutines, run another goroutine to wait for
// the waitgroup and then close the out channel.
// While this goroutine is running, return out.
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
