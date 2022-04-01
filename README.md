# Go concurrency exercises

Exercises and code walks included in this repository are part of Udemy course "concurrency in Go (Golang)".

https://www.udemy.com/course/concurrency-in-go-golang/?referralCode=5AE5A041D5793C048954

## Notes
My changes start at commit 35 https://github.com/petherin/go-concurrency-exercises/commit/dba9ef06ecf604a01481b7b771258f3e0795a56a.

### OS Threads and Goroutines
Goroutines run sequentially in OS threads.

Go will run as many OS threads as there are CPU cores on your machine. `runtime.GOMAXPROCS()` provides this number.

Parallelism is achieved by running goroutines sequentially on multiple os threads.

Each goroutine is given a time slice of 10ms.

### Channels

Default value for channel is `nil`. Reading/writing to a nil channel blocks forever. Create channels with `make` built-in function.

Closing a nil channel panics.

Owner of channel is a goroutine that instantiates, writes and closes a channel.

Channel users should only have a read view into the channel.

### Select

All `select` cases are considered simultaneously.

`select` waits until a case is ready to proceed.

If waiting for multiple channels in a `select` and more than one is ready, one will be picked at random.

`select`s useful for timeouts so long-running channels don't hold things up.

A `select` can be non-blocking if there's a `default` case. If a channel isn't ready in the `select` then a `default` case will exit the block without blocking.

Empty `select` blocks forever.

A nil channel in a `select` also blocks forever.
