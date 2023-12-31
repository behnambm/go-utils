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

func NewWorkerPool(workerCount, bufferSize int) (*WorkerPool, error) {
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
			defer func() {
				//log.Printf("stopping worker #%d \n", workerIdx)
			}()

			for {
				select {
				case <-wp.done:
					return
				case task, ok := <-wp.taskCh:
					//log.Println("started to process")
					if !ok {
						//log.Printf("stopping worker %d with closed tasks channel\n", workerIdx)
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
