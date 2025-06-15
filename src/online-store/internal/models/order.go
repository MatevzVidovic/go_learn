// internal/models/order.go
// Order represents a customer order

package models

import "time"

// Order represents a customer's order
type Order struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	ProductID  int       `json:"product_id" db:"product_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	TotalCents int       `json:"total_cents" db:"total_cents"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// OrderRequest represents data needed to create an order
type OrderRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

// OrderResponse includes product information with the order
type OrderResponse struct {
	ID          int       `json:"id"`
	ProductID   int       `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	TotalCents  int       `json:"total_cents"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// TotalInDollars returns the total price in dollars
func (o *Order) TotalInDollars() float64 {
	return float64(o.TotalCents) / 100.0
}

// MQTT Message Types
// These structs represent the data we send over MQTT

// UserRegisteredEvent is published when a new user registers
type UserRegisteredEvent struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	Timestamp int64  `json:"timestamp"`
}

// ProductCreatedEvent is published when a new product is created
type ProductCreatedEvent struct {
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

// OrderCreatedEvent is published when a new order is placed
type OrderCreatedEvent struct {
	OrderID    int   `json:"order_id"`
	UserID     int   `json:"user_id"`
	ProductID  int   `json:"product_id"`
	Quantity   int   `json:"quantity"`
	TotalCents int   `json:"total_cents"`
	Timestamp  int64 `json:"timestamp"`
}

// LowStockAlert is published when product stock is low
type LowStockAlert struct {
	ProductID    int    `json:"product_id"`
	ProductName  string `json:"product_name"`
	CurrentStock int    `json:"current_stock"`
	ReorderLevel int    `json:"reorder_level"`
	Timestamp    int64  `json:"timestamp"`
}
