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

	mu := sync.Mutex{}
	c := sync.NewCond(&mu)

	wg.Add(1)
	go func() {
		defer wg.Done()

		// suspend goroutine until sharedRsc is populated.

		log.Println("goroutine1: about to wait")
		c.L.Lock()
		for len(sharedRsc) < 1 {
			//time.Sleep(1 * time.Millisecond)
			c.Wait()
		}

		fmt.Printf("goroutine1: %s\n", sharedRsc["rsc1"])
		c.L.Unlock()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// suspend goroutine until sharedRsc is populated.

		log.Println("goroutine2: about to wait")
		c.L.Lock()
		for len(sharedRsc) < 2 {
			//time.Sleep(1 * time.Millisecond)
			c.Wait()
		}

		fmt.Printf("goroutine2: %s\n", sharedRsc["rsc2"])
		c.L.Unlock()
	}()

	c.L.Lock()
	// writes changes to sharedRsc
	time.Sleep(1 * time.Second)
	sharedRsc["rsc1"] = "foo"
	time.Sleep(1 * time.Second)
	sharedRsc["rsc2"] = "bar"
	log.Println("main routine: wake all goroutines with Broadcast")
	c.Broadcast() // Wakes up all goroutines
	c.L.Unlock()
	wg.Wait()
}
