# Go concurrency exercises

Exercises and code walks included in this repository are part of Udemy course "concurrency in Go (Golang)".

https://www.udemy.com/course/concurrency-in-go-golang/?referralCode=5AE5A041D5793C048954

## Notes
My changes start at commit 35 https://github.com/petherin/go-concurrency-exercises/commit/dba9ef06ecf604a01481b7b771258f3e0795a56a.

Goroutines run sequentially in OS threads.

Go will run as many OS threads as there are CPU cores on your machine. `runtime.GOMAXPROCS()` provides this number.

Parallelism is achieved by running goroutines sequentially on multiple os threads.

Each goroutine is given a time slice of 10ms.
