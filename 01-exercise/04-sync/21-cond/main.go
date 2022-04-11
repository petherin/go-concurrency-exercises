package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var sharedRsc = make(map[string]interface{})

func main() {
	var wg sync.WaitGroup

	mu := sync.Mutex{}     // Create mutex
	c := sync.NewCond(&mu) // Create condition and give it the mutex so we can lock and unlock

	wg.Add(1)
	go func() {
		defer wg.Done()

		// suspend goroutine until sharedRsc is populated.
		log.Println("goroutine: about to wait")
		c.L.Lock() // Lock the conditional variable
		for len(sharedRsc) == 0 {
			//time.Sleep(1 * time.Millisecond)
			// c.Wait() unlocks c.L and suspends goroutine execution.
			// Won't resume until a Signal or a Broadcast wakes it up.
			c.Wait()
		}

		log.Println("goroutine: finished waiting")
		fmt.Printf("goroutine: %s", sharedRsc["rsc1"])

		c.L.Unlock() // Release the lock
	}()

	c.L.Lock()                  // Acquire lock
	time.Sleep(2 * time.Second) // Wait around to prove the goroutine is suspended
	sharedRsc["rsc1"] = "foo"   // writes changes to sharedRsc
	c.Signal()                  // Wakes up one suspended goroutine
	c.L.Unlock()                // Release lock

	wg.Wait()
}
