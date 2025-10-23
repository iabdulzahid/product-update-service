package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/iabdulzahid/product-update-service/internal/handler"
	"github.com/iabdulzahid/product-update-service/internal/repository"
	"github.com/iabdulzahid/product-update-service/pkg/models"
	"github.com/iabdulzahid/product-update-service/pkg/queue"
)

// Test 1: Basic POST + GET flow
func TestPostAndGetProduct(t *testing.T) {
	store := repository.NewProductStore()
	eq := queue.NewEventQueue(10)

	done := make(chan struct{})
	go func() {
		for p := range eq.Dequeue() {
			store.Update(p)
			done <- struct{}{}
		}
	}()

	h := handler.NewProductHandler(store, eq)

	product := models.ProductDTO{ProductID: "p1", Price: 10.0, Stock: 5}
	payload, _ := json.Marshal(product)

	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	h.PostEventHandler(w, req)
	if w.Code != http.StatusAccepted {
		t.Fatalf("Expected 202, got %d", w.Code)
	}
	<-done

	retrieved, ok := store.Get("p1")
	if !ok || retrieved.Price != 10.0 || retrieved.Stock != 5 {
		t.Errorf("Expected %+v, got %+v", product, retrieved)
	}
}

// Test 2: Sequential updates for same product (later overrides earlier)
func TestSequentialUpdatesSameProduct(t *testing.T) {
	store := repository.NewProductStore()
	eq := queue.NewEventQueue(100)

	done := make(chan struct{})
	go func() {
		for p := range eq.Dequeue() {
			store.Update(p)
			done <- struct{}{}
		}
	}()

	h := handler.NewProductHandler(store, eq)

	updates := []models.ProductDTO{
		{ProductID: "p2", Price: 10, Stock: 5},
		{ProductID: "p2", Price: 20, Stock: 10},
		{ProductID: "p2", Price: 30, Stock: 15},
		{ProductID: "p2", Price: 40, Stock: 20},
		{ProductID: "p2", Price: 50, Stock: 25},
	}

	for _, u := range updates {
		payload, _ := json.Marshal(u)
		req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(payload))
		w := httptest.NewRecorder()
		h.PostEventHandler(w, req)
		if w.Code != http.StatusAccepted {
			t.Errorf("Expected 202, got %d", w.Code)
		}
		<-done
	}

	retrieved, ok := store.Get("p2")
	if !ok {
		t.Fatalf("Product not found")
	}
	if retrieved.Price != 50 || retrieved.Stock != 25 {
		t.Errorf("Expected last update to win, got %+v", retrieved)
	}
}

// Test 3: Concurrent updates for different products (concurrency safety)
func TestConcurrentUpdatesDifferentProducts(t *testing.T) {
	store := repository.NewProductStore()
	eq := queue.NewEventQueue(100)
	done := make(chan struct{})

	go func() {
		for p := range eq.Dequeue() {
			store.Update(p)
			done <- struct{}{}
		}
	}()

	h := handler.NewProductHandler(store, eq)
	var wg sync.WaitGroup

	products := []models.ProductDTO{
		{ProductID: "p3", Price: 15, Stock: 10},
		{ProductID: "p4", Price: 25, Stock: 20},
		{ProductID: "p5", Price: 35, Stock: 30},
		{ProductID: "p6", Price: 45, Stock: 40},
		{ProductID: "p7", Price: 55, Stock: 50},
	}

	for _, p := range products {
		wg.Add(1)
		go func(prod models.ProductDTO) {
			defer wg.Done()
			payload, _ := json.Marshal(prod)
			req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(payload))
			w := httptest.NewRecorder()
			h.PostEventHandler(w, req)
			if w.Code != http.StatusAccepted {
				t.Errorf("Expected 202, got %d", w.Code)
			}
			<-done
		}(p)
	}

	wg.Wait()

	for _, p := range products {
		retrieved, ok := store.Get(p.ProductID)
		if !ok {
			t.Errorf("Missing product %s", p.ProductID)
			continue
		}
		if retrieved.Price != p.Price || retrieved.Stock != p.Stock {
			t.Errorf("Incorrect data for %s: %+v", p.ProductID, retrieved)
		}
	}
}

// Test 4: GET on non-existent product
func TestGetNonExistentProduct(t *testing.T) {
	store := repository.NewProductStore()
	eq := queue.NewEventQueue(10)
	h := handler.NewProductHandler(store, eq)

	req := httptest.NewRequest("GET", "/products/unknown", nil)
	w := httptest.NewRecorder()
	h.GetProductHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

// Test 5: Queue full scenario (TryEnqueue returns false)
func TestQueueFullScenario(t *testing.T) {
	store := repository.NewProductStore()
	eq := queue.NewEventQueue(1) // intentionally small queue
	h := handler.NewProductHandler(store, eq)

	// Fill the queue manually without a consumer
	product1 := models.ProductDTO{ProductID: "p9", Price: 100, Stock: 10}
	payload1, _ := json.Marshal(product1)
	req1 := httptest.NewRequest("POST", "/events", bytes.NewBuffer(payload1))
	w1 := httptest.NewRecorder()
	h.PostEventHandler(w1, req1)
	if w1.Code != http.StatusAccepted {
		t.Fatalf("Expected 202, got %d", w1.Code)
	}

	// Now queue is full — next enqueue should fail
	product2 := models.ProductDTO{ProductID: "p10", Price: 200, Stock: 20}
	payload2, _ := json.Marshal(product2)
	req2 := httptest.NewRequest("POST", "/events", bytes.NewBuffer(payload2))
	w2 := httptest.NewRecorder()
	h.PostEventHandler(w2, req2)

	if w2.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 when queue full, got %d", w2.Code)
	}
}
