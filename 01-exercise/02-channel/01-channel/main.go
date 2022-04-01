package main

import "fmt"

func main() {
	// Get the value computed from goroutine
	
	// Before, without a channel
	//go func(a, b int) {
	//	c := a + b
	//}(1, 2)
	//
	//fmt.Printf("computed value %v\n", r)

	// After, with a channel
	ch := make(chan int)
	go func(a, b int) {
		c := a + b
		ch <- c
	}(1, 2)

	r := <-ch
	fmt.Printf("computed value %v\n", r)
}
