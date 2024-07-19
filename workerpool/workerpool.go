package workerpool

import (
	"fmt"
	"log"
	"sync"
)

var (
	ErrNoWorkers          = fmt.Errorf("attempting to create worker pool with less than 1 worker")
	ErrNegativeBufferSize = fmt.Errorf("attempting to create worker pool with a negative buffer size")
)

type Task func()

type WorkerPool struct {
	workerCount int
	start       sync.Once
	stop        sync.Once
	taskCh      chan Task
	done        chan struct{}
	wg          sync.WaitGroup
}

func New(workerCount, bufferSize int) (*WorkerPool, error) {
	if workerCount < 1 {
		return nil, ErrNoWorkers
	}
	if bufferSize < 0 {
		return nil, ErrNegativeBufferSize
	}

	return &WorkerPool{
		workerCount: workerCount,
		taskCh:      make(chan Task, bufferSize),
		done:        make(chan struct{}),
		wg:          sync.WaitGroup{},
	}, nil
}

func (wp *WorkerPool) Start() {
	wp.start.Do(func() {
		log.Println("Starting worker pool...")
		wp.startWorkers()
	})
}

func (wp *WorkerPool) Stop() {
	wp.stop.Do(func() {
		log.Println("Stopping worker pool...")
		wp.done <- struct{}{}
		close(wp.taskCh)
		close(wp.done)
	})
}

func (wp *WorkerPool) startWorkers() {

	for i := 0; i < wp.workerCount; i++ {
		go func(workerIdx int) {
			for {
				select {
				case <-wp.done:
					return
				case task, ok := <-wp.taskCh:
					if !ok {
						return
					}
					task()
					wp.wg.Done()
				}
			}
		}(i + 1)
	}
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.wg.Add(1)
	wp.taskCh <- task
}

func (wp *WorkerPool) WaitUntilDone() {
	wp.wg.Wait()
}
