package dynamicpool

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

type Task func()

type DynamicWorkerPool struct {
	taskCh             chan Task
	wg                 sync.WaitGroup
	scaleUpTimeout     time.Duration
	scaleDownTimeout   time.Duration
	initialWorkerCount int32
	minWorkerCount     int32
	maxWorkerCount     int32
	currentWorkerCount int32
	mu                 sync.Mutex
	closed             bool
}

func New(opts ...Options) (*DynamicWorkerPool, error) {
	pool := &DynamicWorkerPool{
		taskCh:             make(chan Task),
		wg:                 sync.WaitGroup{},
		scaleUpTimeout:     time.Second * 3,
		scaleDownTimeout:   time.Second * 10,
		initialWorkerCount: 3,
		minWorkerCount:     1,
		maxWorkerCount:     10,
		mu:                 sync.Mutex{},
		closed:             false,
	}

	for _, opt := range opts {
		if err := opt(pool); err != nil {
			return nil, err
		}
	}

	go pool.watchInterruptSignal()

	return pool, nil
}

func (p *DynamicWorkerPool) Start() {
	for i := 0; i < int(p.initialWorkerCount); i++ {
		p.spawnWorker()
	}
}

func (p *DynamicWorkerPool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return
	}
	p.closed = true
}

func (p *DynamicWorkerPool) watchInterruptSignal() {
	interruptCh := make(chan os.Signal)
	signal.Notify(interruptCh, os.Interrupt, os.Kill)
	sig := <-interruptCh
	signal.Reset(sig) // reset the interrupt signal in order to make the OS force close available
	p.Stop()
}

func (p *DynamicWorkerPool) AddTask(task Task) {
	for {
		if p.isPoolClosed() {
			return
		}
		select {
		case p.taskCh <- task:
			p.wg.Add(1)
			return
		case <-time.After(p.scaleUpTimeout):
			if p.canScaleUp() {
				p.spawnWorker()
				p.AddTask(task)
				return
			}
		}
	}
}

func (p *DynamicWorkerPool) canScaleUp() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.currentWorkerCount == p.maxWorkerCount {
		return false
	}

	return true
}

func (p *DynamicWorkerPool) canScaleDown() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.currentWorkerCount > p.minWorkerCount {
		return true
	}

	return false
}

func (p *DynamicWorkerPool) isPoolClosed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return true
	}

	return false
}

func (p *DynamicWorkerPool) spawnWorker() {
	if !p.canScaleUp() {
		return
	}
	atomic.AddInt32(&p.currentWorkerCount, 1)

	go p.worker()
}

func (p *DynamicWorkerPool) worker() {
	for {
		// if the pool is closed, don't spawn a new worker
		if p.isPoolClosed() {
			atomic.AddInt32(&p.currentWorkerCount, -1)
			return
		}
		select {
		case task, ok := <-p.taskCh:
			if !ok {
				// the task channel has been closed
				return
			}

			// if the pool is closed, don't start the task
			if p.isPoolClosed() {
				atomic.AddInt32(&p.currentWorkerCount, -1)
				p.wg.Done()
				return
			}
			task()
			p.wg.Done()
		case <-time.After(p.scaleDownTimeout):
			if p.canScaleDown() {
				atomic.AddInt32(&p.currentWorkerCount, -1)
				return
			}
		}
	}
}

func (p *DynamicWorkerPool) Wait() {
	p.wg.Wait()
}
