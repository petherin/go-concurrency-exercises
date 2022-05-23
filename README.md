# Go concurrency exercises

<!-- `make toc` to generate https://github.com/jonschlinkert/markdown-toc#cli -->

<!-- toc -->

- [OS Threads and Goroutines](#os-threads-and-goroutines)
- [Channels](#channels)
- [Select](#select)
- [Sync Package](#sync-package)
  * [Mutex](#mutex)
  * [Atomic](#atomic)
  * [Cond](#cond)
  * [Once](#once)
  * [Pool](#pool)
- [Go Race Dectector](#go-race-dectector)
- [Concurrency Patterns](#concurrency-patterns)
  * [Pipeline](#pipeline)
  * [Fan-out, Fan-in](#fan-out-fan-in)
  * [Cancellation of Goroutines](#cancellation-of-goroutines)
- [Context Package](#context-package)
  * [Cancellation Functions](#cancellation-functions)
  * [Data Functions](#data-functions)

<!-- tocstop -->

Exercises and code walks included in this repository are part of Udemy course "concurrency in Go (Golang)".

My changes start at commit 35 https://github.com/petherin/go-concurrency-exercises/commit/dba9ef06ecf604a01481b7b771258f3e0795a56a.

https://www.udemy.com/course/concurrency-in-go-golang/?referralCode=5AE5A041D5793C048954

## OS Threads and Goroutines
Goroutines run sequentially in OS threads.

Go will run as many OS threads as there are CPU cores on your machine. `runtime.GOMAXPROCS()` provides this number.

Parallelism is achieved by running goroutines sequentially on multiple os threads.

Each goroutine is given a time slice of 10ms.

## Channels

Default value for channel is `nil`. Reading/writing to a nil channel blocks forever. Create channels with `make` built-in function.

Closing a nil channel panics.

Owner of channel is a goroutine that instantiates, writes and closes a channel.

Channel users should only have a read view into the channel.

## Select

All `select` cases are considered simultaneously.

`select` waits until a case is ready to proceed.

If waiting for multiple channels in a `select` and more than one is ready, one will be picked at random.

`select`s useful for timeouts so long-running channels don't hold things up.

A `select` can be non-blocking if there's a `default` case. If a channel isn't ready in the `select` then a `default` case will exit the block without blocking.

Empty `select` blocks forever.

A nil channel in a `select` also blocks forever.

## Sync Package

### Mutex

`Lock` and `Unlock` resources you want to be concurrent-safe.

To allow multiple reads, but with writes that hold a lock exclusively, use `RLock` and `RUnlock`.

### Atomic

Atomic package can update variables in a concurrent-safe way.

### Cond

Orchestrates goroutines that are waiting for a condition to be true.

`Wait` will suspend goroutine execution until a condition is met. You won't be able to stop on a breakpoint on a `Wait()` line because the entire goroutine is suspended, as opposed to simply blocking at that line.

`Signal` wakes one goroutine waiting on the `cond` variable.

`Broadcast` wakes all goroutiones waiting on the `cond` variable.

### Once

`sync.Once` ensures only one call to `Do(funcValue)` ever calls the passed function, even on different goroutines.

This is good for things like Singletons, or resources that multiple goroutines need but that only need to be initialised once.

### Pool

Used to constrain the creation of expensive resources like db connections, network connections, and memory.

Maintains a pool of a fixed number of resources that can be reused.

Code `Get`s resource from the pool and when finished, `Put`s it back in the pool for other code to use.

## Go Race Dectector

Race detector find race conditions in Go code.

`go test -race mypkg`

`go run -race mysrc.go`

`go build -race mycmd`. If you then execute the resulting build e.g. `./mycmd` any data races will be shown.

`go install -race mypkg`

Binary needs to be race-enabled so race detector can work on it.

Race-enabled binaries can be 10 times slower and use 10 times more memory, so don't release to production, use during testing phase.

## Concurrency Patterns
### Pipeline

Used to process streams or batches of data. It's a series of stages connected by channels.

Each stage is represented by a goroutine.

![Pipeline Pattern](files/pipeline.png "Pipeline Pattern")

Goroutines can have the same input and output parameters, meaning we can chain them together however we want.

For example, we could do `square(decrement(square(ch)))`.

Separating stages out provides us with good separation of concerns.

If a stage is taking a long time we can increase the number of goroutines for that stage.

### Fan-out, Fan-in

If we have a stage in our pipeline that is taking too long and blocking subsequent stages, we can use fan-out, fan-in.

![Fan-Out Fan-In Pattern](files/pipeline-fanout-fanin.png "Fan-Out Fan-In Pattern")

Multiple goroutines for a stage are started. They take in items from the incoming channel and do work. This is fan-out.

They send the output on their own channels to `merge` goroutines. These merge the incoming multiple channels into a single output channel. This is fan-in.

![Fan-Out Fan-In Diagram](files/fanout-fanin-diagram.png "Fan-Out Fan-In Diagram")

### Cancellation of Goroutines

In the above pipeline, `main()` is waiting for values on the channel from `merge()`.

If there are multiple values on the channel from `merge()` but `main()` only reads one, it will block execution.

`merge()` goroutines will be blocked when trying to send more values.

`square()` and `generator()` will also be blocked on sending because `merge()` is blocked.

This is a goroutine leak.

We need a way to cancel goroutines if something goes wrong like this.

We can pass a read-only `done` channel to goroutines.

Then we can close the channel to send broadcast signal to all goroutines.

On receiving the signal on `done` channel, goroutines need to abandon work and terminate.

## Context Package

Serves two primary purposes.

* Provides API's for cancelling branches of the call graph (i.e. the calls made during a given request).
* Provides a data bag for transporting request-scoped data through the call graph.

### Cancellation Functions

`context.Background()` returns an empty context, it's the root of any context tree.

`context.TODO()` returns an empty context, intended to be a placeholder.

`ctx, cancel := context.WithCancel(context.Background())` returns a copy of parent context with a new Done channel. The returned `cancel` allows us to close the context's Done channel.

`cancel()` doesn't wait for anything, it just closes the Done channel and returns. It can be called by many goroutines simultaneously, but after the first call, it doesn't do anything.

`context.WithDeadline()` takes parent context and a time as input. It returns a new context that closes its Done channel when the machine's clock advances past the deadline.

`ctx.Deadline()` tells us if a deadline is associated with the context.

`context.WithTimeout()` takes parent context and a time duration as input. It returns a new context that closes its Done channel when the given duration expires.

`WithTimeout` is actually a wrapper around `WithDeadline`. What's the difference? 

`WithTimeout()`'s timer countdown begins from the moment the context is created.

`WithDeadine` sets an explicit time when timer will expire.

### Data Functions

`context.WithValue()` associates request-scoped values with a context.

`ctx.Value()` returns the value associated with the provided key.
