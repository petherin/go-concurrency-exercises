package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// create pool of bytes.Buffers which can be reused.

var bufPool = sync.Pool{
	// New is run when there are no resources in pool and we need to create ome
	New: func() any {
		// This will only be called once and then bytes.Buffer will be reused on the sedonc call to log
		fmt.Println("allocating new bytes.Buffer")
		return new(bytes.Buffer)
	},
}

func log(w io.Writer, val string) {
	// We don't want to create a new bytes.Buffer every time this func is called.
	// Could be called by thousands of goroutines.
	// Better to put a few in a pool.
	//var b bytes.Buffer Instead of this we'll do this....
	b := bufPool.Get().(*bytes.Buffer) // Need to cast result of Get of the expected type

	b.Reset() // in case something else was using it

	b.WriteString(time.Now().Format("15:04:05"))
	b.WriteString(" : ")
	b.WriteString(val)
	b.WriteString("\n")

	w.Write(b.Bytes())

	bufPool.Put(b) // Put the buffer back in the pool
}

func main() {
	log(os.Stdout, "debug-string1")
	log(os.Stdout, "debug-string2")
}
