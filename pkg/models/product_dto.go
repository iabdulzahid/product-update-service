package models

type ProductDTO struct {
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}
