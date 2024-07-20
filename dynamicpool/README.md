# Dynamic WorkerPool 

This package offers a dynamic worker pool designed to handle varying workloads efficiently.
It automatically adjusts the number of goroutines based on current workload demands, 
scaling up or down as needed. By dynamically spawning and destroying goroutines, 
the pool optimizes resource usage while ensuring responsive performance. 
Ideal for tasks where workload fluctuates, this package simplifies concurrent
task management in Golang applications.

### Features
- **Configurable Worker Limits**: Set minimum and maximum number of workers to control resource usage.
- **Adaptive Scaling**: Define the required time to scale up or scale down based on workload changes.
- **Graceful Shutdown**: Monitor OS interruption signals and shut down the worker pool gracefully to ensure smooth termination.


### Example 

```go
package main

import (
	"fmt"
	"github.com/behnambm/go-utils/dynamicpool"
	"time"
)

func main() {
	dPool, _ := dynamicpool.New()
	dPool.Start()

	dPool.AddTask(func() {
		fmt.Println("task 1 - started")
		time.Sleep(20 * time.Second)
		fmt.Println("task 1 - done")
	})
	dPool.AddTask(func() {
		fmt.Println("task 2 - started")
		time.Sleep(5 * time.Second)
		fmt.Println("task 2 - done")
	})
	dPool.AddTask(func() {
		fmt.Println("task 3 - started")
		time.Sleep(12 * time.Second)
		fmt.Println("task 3 - done")
	})
	dPool.AddTask(func() {
		fmt.Println("task 4 - started")
		time.Sleep(1 * time.Second)
		fmt.Println("task 4 - done")
	})

	dPool.Wait()
	dPool.Stop()
}

```