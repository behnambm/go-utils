package workerpool

import (
	"runtime"
	"testing"
	"time"
)

func TestNewWorkerPool_NoError(t *testing.T) {
	_, err := NewWorkerPool(1, 1)
	if err != nil {
		t.Fatal("returned error - should not return errors")
	}
}

func TestNewWorkerPool_NoWorkerErr(t *testing.T) {
	_, err := NewWorkerPool(0, 1)
	if err != ErrNoWorkers {
		t.Fatal("returned invalid error - should return ErrNoWorkers")
	}
}

func TestNewWorkerPool_NegativeBufferSizeErr(t *testing.T) {
	_, err := NewWorkerPool(1, -1)
	if err != ErrNegativeBufferSize {
		t.Fatal("returned invalid error - should return ErrNegativeBufferSize")
	}
}

func TestNewWorkerPool_StartAndStopOnce(t *testing.T) {
	wp, _ := NewWorkerPool(1, 1)

	wp.Start()
	wp.Start()

	wp.Stop()
	wp.Stop()
}

func TestNewWorkerPool_GoroutineCount(t *testing.T) {
	wp, err := NewWorkerPool(20, 1)
	if err != nil {
		t.Fatal("returned error - should not return errors")
	}

	wp.Start()
	numGoroutine := runtime.NumGoroutine()
	if numGoroutine < 20 {
		t.Fatal("workers are not spawned")
	}
	wp.Stop()
}

func TestNewWorkerPool_AddTask(t *testing.T) {
	wp, _ := NewWorkerPool(1, 1)
	wp.Start()

	counter := 0
	increaseCounter := func() {
		counter++
	}

	wp.AddTask(increaseCounter)
	wp.WaitUntilDone()
	wp.Stop()

	if counter != 1 {
		t.Fatal("worker didn't change the value of counter")
	}
}

func TestNewWorkerPool_WaitUntilDone(t *testing.T) {
	wp, _ := NewWorkerPool(1, 1)
	wp.Start()

	increaseCounter := func() {
		time.Sleep(1 * time.Second)
	}

	wp.AddTask(increaseCounter)
	start := time.Now()
	wp.WaitUntilDone()
	totalTime := time.Since(start)
	wp.Stop()

	if totalTime.Seconds() < 1 {
		t.Fatal("worker didn't run the task")
	}
}

func TestNewWorkerPool_TaskRun_OneWorker(t *testing.T) {
	// this task should wait for 1 seconds
	wp, _ := NewWorkerPool(1, 0)
	wp.Start()

	increaseCounter := func() {
		time.Sleep(1 * time.Second)
	}

	start := time.Now()

	wp.AddTask(increaseCounter)
	wp.WaitUntilDone()

	totalTime := time.Since(start)

	wp.Stop()

	if totalTime.Seconds() < 1 {
		t.Fatal("worker didn't run the task in proper time")
	}
}

func TestNewWorkerPool_BlockingAddTask_OneWorker(t *testing.T) {
	// set 1 as worker count so the tasks will be processed one by one
	// this task should wait for 1 second to AddTask and 2 seconds in total
	wp, _ := NewWorkerPool(1, 0)
	wp.Start()

	increaseCounter := func() {
		time.Sleep(1 * time.Second)
	}
	start := time.Now()

	wp.AddTask(increaseCounter)
	wp.AddTask(increaseCounter)

	addWaitTime := time.Since(start)
	if addWaitTime.Seconds() < 1 {
		t.Fatal("AddTask didn't wait to send task to worker")
	}

	wp.WaitUntilDone()

	totalTime := time.Since(start)

	wp.Stop()

	if totalTime.Seconds() < 2 {
		t.Fatal("tasks didn't finish in proper time - should wait 2 seconds")
	}
}

func TestNewWorkerPool_NonBlockingAddTask_OneWorker_OneBuffer(t *testing.T) {
	// set 1 as worker count and 1 as buffer size
	// this task should not wait to AddTask and should wait 2 seconds in total
	wp, _ := NewWorkerPool(1, 1)
	wp.Start()

	increaseCounter := func() {
		time.Sleep(1 * time.Second)
	}
	start := time.Now()

	wp.AddTask(increaseCounter)
	wp.AddTask(increaseCounter)

	addWaitTime := time.Since(start)
	if int(addWaitTime.Seconds()) != 0 {
		t.Fatal("AddTask should not wait with buffered channel to enqueue tasks")
	}

	wp.WaitUntilDone()

	totalTime := time.Since(start)

	wp.Stop()

	if totalTime.Seconds() < 2 {
		t.Fatal("tasks didn't finish in proper time - should wait 2 seconds")
	}
}

func TestNewWorkerPool_NonBlockingAddTask_TwoWorker_OneBuffer(t *testing.T) {
	// set 2 as worker count and 1 as buffer size
	// this task should not wait to AddTask and should wait 1 seconds in total
	wp, _ := NewWorkerPool(2, 1)
	wp.Start()

	increaseCounter := func() {
		time.Sleep(1 * time.Second)
	}
	start := time.Now()

	wp.AddTask(increaseCounter)
	wp.AddTask(increaseCounter)

	addWaitTime := time.Since(start)
	if int(addWaitTime.Seconds()) != 0 {
		t.Fatal("AddTask should not wait with buffered channel to enqueue tasks")
	}

	wp.WaitUntilDone()

	totalTime := time.Since(start)

	wp.Stop()

	if totalTime.Seconds() < 1 || totalTime.Seconds() > 1.5 { // to make wait time is less than 2 seconds
		t.Fatal("tasks didn't finish in proper time - should wait 1 seconds")
	}
}
