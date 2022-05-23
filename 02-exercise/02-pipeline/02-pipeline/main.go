// generator() -> square() -> print

package main

import (
	"fmt"
	"sync"
)

func generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	// Implement fan-in
	// merge a list of channels to a single channel
	out := make(chan int)
	var wg sync.WaitGroup

	// anonymous function to merge values from multiple channels into a single channel
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))

	// range over the channels and run output to merge their values into a single channel
	for _, c := range cs {
		go output(c)
	}

	// goroutine to close the single channel when the multiple channels have been merged
	go func() {
		wg.Wait()
		close(out)
	}()

	// return the single channel so the caller can range over it as and when values are merged into it
	return out
}

func main() {
	in := generator(2, 3)

	// TODO: fan out square stage to run two instances.
	ch1 := square(in)
	ch2 := square(in)

	// TODO: fan in the results of square stages.
	for n := range merge(ch1, ch2) {
		fmt.Println(n)
	}
}
