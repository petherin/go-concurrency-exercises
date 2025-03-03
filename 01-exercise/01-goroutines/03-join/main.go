package main

import (
	"fmt"
	"sync"
)

func main() {
	// modify the program
	// to print the value as 1
	// deterministically.

	var data int
	var wg sync.WaitGroup

	wg.Add(1)

	// This never doesn't get the chance to run unless we use a WaitGroup
	go func() {
		defer wg.Done()
		data++
	}()

	wg.Wait()

	fmt.Printf("the value of data is %v\n", data)

	fmt.Println("Done")
}
