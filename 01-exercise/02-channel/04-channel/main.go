package main

import "fmt"

// Implement relaying of message with Channel Direction

// Accepts send-only channel
func genMsg(ch1 chan<- string) {
	// send message on ch1
	msg := "message"
	fmt.Printf("genMsg: sent '%s' to ch1\n", msg)
	ch1 <- "message"

	// If we do this we get a compile error:
	// Invalid operation: <-ch1 (receive from the send-only type chan<- string)
	// <-ch1
}

// ch1 is receive-only, ch2 is send-only
func relayMsg(ch1 <-chan string, ch2 chan<- string) {
	// recv message on ch1
	m := <-ch1
	fmt.Printf("relayMsg: received '%s' from ch1\n", m)

	// send it on ch2
	ch2 <- m
	fmt.Printf("relayMsg: sent '%s' to ch2\n", m)
}

func main() {
	// create ch1 and ch2
	ch1 := make(chan string)
	ch2 := make(chan string)

	// spin goroutine genMsg and relayMsg
	go genMsg(ch1)
	// recv message on ch2
	go relayMsg(ch1, ch2)

	v := <-ch2
	fmt.Printf("main routine received %s\n", v)
}
