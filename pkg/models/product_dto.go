package models

type ProductDTO struct {
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
}

type Config struct {
	Port      int
	Workers   int
	QueueSize int
}
