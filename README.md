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

### Sync Package

#### Mutex

`Lock` and `Unlock` resources you want to be concurrent-safe.

To allow multiple reads, but with writes that hold a lock exclusively, use `RLock` and `RUnlock`.

#### Atomic

Atomic package can update variables in a concurrent-safe way.

#### Cond

Orchestrates goroutines that are waiting for a condition to be true.

`Wait` will suspend goroutine execution until a condition is met. You won't be able to stop on a breakpoint on a `Wait()` line because the entire goroutine is suspended, as opposed to simply blocking at that line.

`Signal` wakes one goroutine waiting on the `cond` variable.

`Broadcast` wakes all goroutiones waiting on the `cond` variable.

#### Once

`sync.Once` ensures only one call to `Do(funcValue)` ever calls the passed function, even on different goroutines.

This is good for things like Singletons, or resources that multiple goroutines need but that only need to be initialised once.

#### Pool

Used to constrain the creation of expensive resources like db connections, network connections, and memory.

Maintains a pool of a fixed number of resources that can be reused.

Code `Get`s resource from the pool and when finished, `Put`s it back in the pool for other code to use.

### Go Race Dectector

Race detector find race conditions in Go code.

`go test -race mypkg`

`go run -race mysrc.go`

`go build -race mycmd`. If you then execute the resulting build e.g. `./mycmd` any data races will be shown.

`go install -race mypkg`

Binary needs to be race-enabled so race detector can work on it.

Race-enabled binaries can be 10 times slower and use 10 times more memory, so don't release to production, use during testing phase.

### Concurrency Patterns
#### Pipeline

Used to process streams or batches of data. It's a series of stages connected by channels.

Each stage is represented by a goroutine.

![Pipeline Pattern](files/pipeline.png "Pipeline Pattern")

Goroutines can have the same input and output parameters, meaning we can chain them together however we want.

For example, we could do `square(decrement(square(ch)))`.

Separating stages out provides us with good separation of concerns.

If a stage is taking a long time we can increase the number of goroutines for that stage.




