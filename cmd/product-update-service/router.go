package main

import (
	"github.com/gorilla/mux"
	"github.com/iabdulzahid/product-update-service/internal/handler"
)

func NewRouter(h *handler.ProductHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/events", h.PostEventHandler).Methods("POST")
	r.HandleFunc("/products/{id}", h.GetProductHandler).Methods("GET")
	return r
}
