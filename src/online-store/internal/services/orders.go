// internal/services/orders.go
// This file contains order-related business logic

package services

import (
	"database/sql"
	"fmt"
	"time"

	"online-store/internal/models"
	"online-store/internal/mqtt"
)

// OrderService handles order operations
type OrderService struct {
	db         *sql.DB
	mqttClient *mqtt.Client
}

// NewOrderService creates a new order service
func NewOrderService(db *sql.DB, mqttClient *mqtt.Client) *OrderService {
	return &OrderService{
		db:         db,
		mqttClient: mqttClient,
	}
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(userID int, req models.OrderRequest) (*models.OrderResponse, error) {
	// Start a database transaction
	// This ensures that if anything goes wrong, all changes are rolled back
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	
	// If something goes wrong, roll back the transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the product to check stock and calculate price
	var product models.Product
	err = tx.QueryRow(
		"SELECT id, name, price_cents, stock_quantity FROM products WHERE id = ?",
		req.ProductID,
	).Scan(&product.ID, &product.Name, &product.PriceCents, &product.StockQuantity)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Check if we have enough stock
	if product.StockQuantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: only %d items available", product.StockQuantity)
	}

	// Calculate total price
	totalCents := product.PriceCents * req.Quantity

	// Create the order
	result, err := tx.Exec(
		"INSERT INTO orders (user_id, product_id, quantity, total_cents, status) VALUES (?, ?, ?, ?, ?)",
		userID, req.ProductID, req.Quantity, totalCents, "pending",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get order ID: %w", err)
	}

	// Update product stock
	newStock := product.StockQuantity - req.Quantity
	_, err = tx.Exec(
		"UPDATE products SET stock_quantity = ? WHERE id = ?",
		newStock, req.ProductID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Create order response
	orderResponse := &models.OrderResponse{
		ID:          int(orderID),
		ProductID:   req.ProductID,
		ProductName: product.Name,
		Quantity:    req.Quantity,
		TotalCents:  totalCents,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Publish MQTT event that order was created
	event := models.OrderCreatedEvent{
		OrderID:    int(orderID),
		UserID:     userID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalCents: totalCents,
		Timestamp:  time.Now().Unix(),
	}
	
	if err := s.mqttClient.Publish("order/created", event); err != nil {
		fmt.Printf("Failed to publish order created event: %v", err)
	}

	// Check if stock is low after this order
	if newStock < 10 {
		alert := models.LowStockAlert{
			ProductID:    req.ProductID,
			ProductName:  product.Name,
			CurrentStock: newStock,
			ReorderLevel: 10,
			Timestamp:    time.Now().Unix(),
		}
		
		if err := s.mqttClient.Publish("inventory/low_stock", alert); err != nil {
			fmt.Printf("Failed to publish low stock alert: %v", err)
		}
	}

	return orderResponse, nil
}

// GetUserOrders returns all orders for a specific user
func (s *OrderService) GetUserOrders(userID int) ([]models.OrderResponse, error) {
	rows, err := s.db.Query(`
		SELECT o.id, o.product_id, p.name, o.quantity, o.total_cents, o.status, o.created_at
		FROM orders o
		JOIN products p ON o.product_id = p.id
		WHERE o.user_id = ?
		ORDER BY o.created_at DESC
	`, userID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []models.OrderResponse
	
	for rows.Next() {
		var order models.OrderResponse
		err := rows.Scan(
			&order.ID,
			&order.ProductID,
			&order.ProductName,
			&order.Quantity,
			&order.TotalCents,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// GetOrder returns a specific order (only if it belongs to the user)
func (s *OrderService) GetOrder(orderID, userID int) (*models.OrderResponse, error) {
	var order models.OrderResponse
	err := s.db.QueryRow(`
		SELECT o.id, o.product_id, p.name, o.quantity, o.total_cents, o.status, o.created_at
		FROM orders o
		JOIN products p ON o.product_id = p.id
		WHERE o.id = ? AND o.user_id = ?
	`, orderID, userID).Scan(
		&order.ID,
		&order.ProductID,
		&order.ProductName,
		&order.Quantity,
		&order.TotalCents,
		&order.Status,
		&order.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// UpdateOrderStatus updates the status of an order
// This method is called by MQTT handlers when payments are confirmed
func (s *OrderService) UpdateOrderStatus(orderID int, status string) error {
	result, err := s.db.Exec(
		"UPDATE orders SET status = ? WHERE id = ?",
		status, orderID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	// Publish MQTT event that order status changed
	event := struct {
		OrderID   int    `json:"order_id"`
		Status    string `json:"status"`
		Timestamp int64  `json:"timestamp"`
	}{
		OrderID:   orderID,
		Status:    status,
		Timestamp: time.Now().Unix(),
	}
	
	if err := s.mqttClient.Publish("order/status_changed", event); err != nil {
		fmt.Printf("Failed to publish order status changed event: %v", err)
	}

	return nil
}