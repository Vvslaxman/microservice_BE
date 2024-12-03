package utils

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	Workers int
	wg      sync.WaitGroup
	tasks   chan func()
}

func NewWorkerPool(workers int) *WorkerPool {
	pool := &WorkerPool{
		Workers: workers,
		tasks:   make(chan func(), 100),
	}
	return pool
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.Workers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for task := range wp.tasks {
				task()
			}
		}()
	}
}

func (wp *WorkerPool) AddTask(task func()) {
	wp.tasks <- task
}

func (wp *WorkerPool) Wait() {
	close(wp.tasks)
	wp.wg.Wait()
}

func ProcessImage(url string) {
	fmt.Printf("Processing image: %s\n", url)
	// Simulated processing logic
}
