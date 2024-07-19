# WorkerPool 

The WorkerPool package for Go provides a convenient way to manage a pool of workers. It allows you to create multiple
workers and assign tasks to them for efficient processing.

### Example 

```go
package main

import (
	"time"

	"github.com/behnambm/go-utils/workerpool"
)

// Define your task
func myTask() {
	time.Sleep(1 * time.Second)
}

func main() {
	// Create a new worker pool with 1 worker and an unbuffered channel
	wp, _ := workerpool.New(1, 0) 
	wp.Start()

	// Add tasks to the worker pool
	wp.AddTask(myTask)
	wp.AddTask(myTask)

	// Wait for all tasks to be completed
	wp.WaitUntilDone()

	// Stop the worker pool
	wp.Stop()
}
```