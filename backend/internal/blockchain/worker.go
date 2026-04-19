package blockchain

import (
	"context"
	"log"
	"sync"
)

type EventHandler func(ctx context.Context, event AuctionCreatedEvent) error

type WorkerPool struct {
	numWorkers int
	events     <-chan AuctionCreatedEvent
	handler    EventHandler
}

func NewWorkerPool(numWorkers int, events <-chan AuctionCreatedEvent, handler EventHandler) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		events:     events,
		handler:    handler,
	}
}

// create the workers (go routines) based on numWorkers
func (wp *WorkerPool) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case event, ok := <-wp.events:
					if !ok {
						return
					}
					if err := wp.handler(ctx, event); err != nil {
						log.Printf("worker: failed to handle event: %v", err)
					}
				}
			}
		}()
	}

	wg.Wait()
	log.Println("all workers stopped")
}
