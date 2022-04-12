package main

import (
	"fmt"
	"math/rand"
	"time"
)

// identify the data race and fix the issue.

// broken main() func with data race
//func main() {
//	start := time.Now()
//	var t *time.Timer
//	t = time.AfterFunc(randomDuration(), func() {
//		fmt.Println(time.Now().Sub(start))
//		t.Reset(randomDuration())
//	})
//	time.Sleep(5 * time.Second)
//}

// no data race version
func main() {
	start := time.Now()
	var t *time.Timer
	ch := make(chan bool)

	t = time.AfterFunc(randomDuration(), func() {
		fmt.Println(time.Now().Sub(start))
		// Don't reset t in the goroutine.
		// Instead, send message on a channel.
		ch <- true
	})

	// Keep resetting t for 5 seconds.
	for time.Since(start) < 5*time.Second {
		// Wait for channel to have a value and then get
		// main routine to reset, rather than goroutine setting it.
		// This avoids concurrent access of t between main goroutine
		// and AfterFunc goroutine.
		<-ch
		t.Reset(randomDuration())
	}
}

func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}

//----------------------------------------------------
// (main goroutine) -> t <- (time.AfterFunc goroutine)
//----------------------------------------------------
// (working condition)
// main goroutine..
// t = time.AfterFunc()  // returns a timer..

// AfterFunc goroutine
// t.Reset()        // timer reset
//----------------------------------------------------
// (race condition- random duration is very small)
// AfterFunc goroutine might try and reset t before main goroutine
// has had a chance to set it on line 15
// t.Reset() // t = nil

// main goroutine..
// t = time.AfterFunc() // this line is run second but won't be
// hit if t.Reset above causes a panic
//----------------------------------------------------
