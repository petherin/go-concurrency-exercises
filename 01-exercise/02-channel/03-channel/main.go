package main

import (
	"fmt"
)

func main() {
	// This provides a capacity of 6 in a buffered channel,
	// so up to 6 items can be added to the channel without blocking
	ch := make(chan int, 6)

	go func() {
		defer close(ch)

		// send all iterator values on channel without blocking
		for i := 0; i < 6; i++ {
			fmt.Printf("Sending: %d\n", i)
			ch <- i
		}
	}()

	for v := range ch {
		fmt.Printf("Received: %v\n", v)
	}

	// Output
	// The goroutine can send all 6 of its values without
	// waiting for them to be received. So it has time to send
	// all 6 and then execution reaches the range loop and
	// it prints out all the values on the channel without
	// having to wait for any.
	//
	// Sending: 0
	// Sending: 1
	// Sending: 2
	// Sending: 3
	// Sending: 4
	// Sending: 5
	// Received: 0
	// Received: 1
	// Received: 2
	// Received: 3
	// Received: 4
	// Received: 5

}
