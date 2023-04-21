package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	numberOfValues := 10
	wg := sync.WaitGroup{}

	values := make(chan int, numberOfValues)

	routinePool := numberOfValues / 2

	for i := 0; i < routinePool; i++ {
		wg.Add(1)

		go func(values chan int, i int, wg *sync.WaitGroup) {
			defer wg.Done()

			log.Printf("Goroutine %d waiting for value", i)
			for v := range values {
				log.Printf("Goroutine %d received value %d", i, v)
				time.Sleep(1 * time.Second)
			}
			log.Printf("Goroutine %d terminating", i)
		}(values, i, &wg)
	}

	for i := 0; i < numberOfValues; i++ {
		values <- i
	}

	close(values)
	wg.Wait()

	log.Println("main finished")
}
