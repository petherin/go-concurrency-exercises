package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(1 * time.Second)
			msg := fmt.Sprintf("message%d", i)
			log.Printf("goroutine: sending %s on channel", msg)
			ch <- msg
		}
	}()

	// if there is no value on channel, do not block.
	// Before fix
	//for i := 0; i < 2; i++ {
	// This will block until there's a message, blocking
	// execution of the work below.
	//	m := <-ch
	//	log.Printf("main routine: %s received", m)
	//
	//	// Do some processing
	//	log.Println("main routine: processing for 1500 ms")
	//	time.Sleep(1500 * time.Millisecond)
	//}

	// After fix
	for i := 0; i < 2; i++ {
		// In this select, if there is no message then the default
		// case is hit, and the select does not block, so we can
		// proceed onto some other work. When the work is finished
		// the loop sends us back to the select.
		select {
		case m := <-ch:
			log.Printf("main routine: %s received", m)
		default:
			log.Println("main routine: no message received")
		}

		// Do some processing
		log.Println("main routine: processing for 1500 ms")
		time.Sleep(1500 * time.Millisecond)
	}
}
