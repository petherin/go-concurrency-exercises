package main

import "fmt"

func main() {
	// create channel owner goroutine which return channel and
	// writes data into channel and
	// closes the channel when done.

	// Idiomatically, the goroutine that opens, writes to,
	// and closes the channel, is the owner. The owner must be
	// the only one that opens, writes to and closes the channel.
	owner := func() <-chan int {
		ch := make(chan int)

		go func() {
			defer close(ch)
			for i := 0; i < 5; i++ {
				ch <- i
			}
		}()
		return ch
	}

	// All the consumer is allowed to do is read from the channel.
	consumer := func(ch <-chan int) {
		// read values from channel
		for v := range ch {
			fmt.Printf("Received: %d\n", v)
		}
		fmt.Println("Done receiving!")
	}

	ch := owner()
	consumer(ch)
}
