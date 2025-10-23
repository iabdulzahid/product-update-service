package queue

import "github.com/iabdulzahid/product-update-service/internal/domain"

type EventQueue struct {
	queue chan *domain.Product
}

func NewEventQueue(size int) *EventQueue {
	return &EventQueue{
		queue: make(chan *domain.Product, size),
	}
}

func (eq *EventQueue) Enqueue(p *domain.Product) {
	eq.queue <- p
}

func (eq *EventQueue) Dequeue() <-chan *domain.Product {
	return eq.queue
}

func (eq *EventQueue) Close() {
	close(eq.queue)
}

// TryEnqueue attempts to add a product to the queue without blocking.
// Returns false if the queue is full.
func (q *EventQueue) TryEnqueue(product *domain.Product) bool {
	select {
	case q.queue <- product:
		return true
	default:
		return false // queue is full
	}
}
