package main

import "log"

// Based on https://go.dev/blog/pipelines
func main() {
	basicPipeline()
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
