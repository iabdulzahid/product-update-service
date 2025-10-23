package repository

import (
	"sync"

	"github.com/iabdulzahid/product-update-service/internal/domain"
)

type ProductStore struct {
	mu       sync.RWMutex
	products map[string]*domain.Product
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		products: make(map[string]*domain.Product),
	}
}

func (ps *ProductStore) Update(product *domain.Product) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.products[product.ProductID] = product
}

func (ps *ProductStore) Get(id string) (*domain.Product, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	p, ok := ps.products[id]
	if !ok {
		return nil, false
	}
	// return a copy, Avoid exposing internal pointers to callers
	prod := *p
	return &prod, true
}
