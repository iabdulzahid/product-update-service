package service

import (
	"context"
	"log"

	"github.com/iabdulzahid/product-update-service/internal/repository"
	"github.com/iabdulzahid/product-update-service/pkg/queue"
)

type ProductService struct {
	store   *repository.ProductStore
	queue   *queue.EventQueue
	workers int
}

func NewProductService(store *repository.ProductStore, queue *queue.EventQueue, workers int) *ProductService {
	return &ProductService{
		store:   store,
		queue:   queue,
		workers: workers,
	}
}

func (ps *ProductService) StartWorkers(ctx context.Context) {
	for i := 0; i < ps.workers; i++ {
		go func(id int) {
			for {
				select {
				case product, ok := <-ps.queue.Dequeue():
					if !ok {
						// queue closed
						return
					}
					ps.store.Update(product)
					log.Printf("Worker %d processed product %s", id, product.ProductID)
				case <-ctx.Done():
					log.Printf("Worker %d stopping due to ctx.Done()", id)
					return
				}
			}
		}(i)
	}
}
