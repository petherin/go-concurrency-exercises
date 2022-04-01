package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		for i := 0; i < 6; i++ {
			fmt.Printf("Sending: %d\n", i)
			// send iterator over channel
			ch <- i
		}

		// breaks out the loop below that's ranging over the channel
		close(ch)
	}()

	// range over channel to recv values
	for v := range ch {
		fmt.Printf("Received: %v\n", v)
	}

	// Output
	// When we send, we have to wait until the value is
	// received before we can send another value.
	// The output is a bit muddled up but the goroutine
	// is sending 0, then the main routine prints received 0.
	// The goroutine has to wait for the receive
	// before it can send the next value.
	//
	// Sending: 0
	// Sending: 1
	// Received: 0
	// Received: 1
	// Sending: 2
	// Sending: 3
	// Received: 2
	// Received: 3
	// Sending: 4
	// Sending: 5
	// Received: 4
	// Received: 5

}
