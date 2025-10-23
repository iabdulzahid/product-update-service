package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iabdulzahid/product-update-service/internal/domain"
	"github.com/iabdulzahid/product-update-service/internal/repository"
	"github.com/iabdulzahid/product-update-service/pkg/models"
	"github.com/iabdulzahid/product-update-service/pkg/queue"
)

type ProductHandler struct {
	Store *repository.ProductStore
	Queue *queue.EventQueue
}

func NewProductHandler(store *repository.ProductStore, queue *queue.EventQueue) *ProductHandler {
	return &ProductHandler{Store: store, Queue: queue}
}

func (h *ProductHandler) PostEventHandler(w http.ResponseWriter, r *http.Request) {
	var dto models.ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	product := &domain.Product{
		ProductID: dto.ProductID,
		Price:     dto.Price,
		Stock:     dto.Stock,
	}

	// Non-blocking enqueue with graceful handling if queue is full
	if !h.Queue.TryEnqueue(product) {
		http.Error(w, "Queue full, try again later", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	product, ok := h.Store.Get(id)
	if !ok {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
